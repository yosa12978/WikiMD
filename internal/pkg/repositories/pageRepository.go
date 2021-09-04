package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/yosa12978/WikiMD/internal/pkg/dto"
	"github.com/yosa12978/WikiMD/internal/pkg/models"
	mongodb "github.com/yosa12978/WikiMD/internal/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IPageRepository interface {
	CreatePage(page_dto *dto.CreatePageDTO, username string) error
	ReadPage(id_hex string) (*models.Page, error)
	GetPages() []models.Page
	SearchPages(query string) []models.Page
	AddCommit(pageid_hex string, commit *models.Commit) error
	DeletePage(id_hex string) error
}

type PageRepository struct {
	db *mongo.Database
}

func NewPageRepository() IPageRepository {
	return &PageRepository{db: mongodb.GetClient()}
}

func (pr *PageRepository) CreatePage(page_dto *dto.CreatePageDTO, username string) error {
	commit := models.Commit{
		ID:   primitive.NewObjectID(),
		Name: page_dto.Name,
		Body: page_dto.Body,
		Page: models.Page{},
		User: username,
	}
	page := models.Page{
		ID:           primitive.NewObjectID(),
		LastCommitID: commit.ID.Hex(),
		Commits:      []models.Commit{commit},
	}
	commit.Page = page
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	session, err := pr.db.Client().StartSession()
	if err != nil {
		return err
	}
	err = pr.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		_, err = pr.db.Collection("commits").InsertOne(sessionContext, commit)
		if err != nil {
			return err
		}

		_, err = pr.db.Collection("pages").InsertOne(sessionContext, page)
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		if err = session.CommitTransaction(sessionContext); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (pr *PageRepository) ReadPage(id_hex string) (*models.Page, error) {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return nil, errors.New("page not found")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{
		"_id": id,
	}
	var page models.Page
	err = pr.db.Collection("pages").FindOne(ctx, filter).Decode(&page)
	if err != nil {
		return nil, errors.New("page not found")
	}
	return &page, nil
}

func (pr *PageRepository) GetPages() []models.Page {
	f_opt := options.Find().SetSort(bson.M{"_id": -1})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var pages []models.Page
	cursor, err := pr.db.Collection("pages").Find(ctx, bson.M{}, f_opt)
	if err != nil {
		return pages
	}
	var wg sync.WaitGroup
	for cursor.Next(ctx) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var page models.Page
			err := cursor.Decode(&page)
			if err != nil {
				return
			}
			pages = append(pages, page)
		}()
	}
	wg.Wait()
	return pages
}

func (pr *PageRepository) SearchPages(query string) []models.Page {
	filter := bson.M{
		"$regex": bson.M{"$or": []bson.M{
			{"name": query},
			{"body": query},
		}},
	}
	f_opt := options.Find().SetSort(bson.M{"_id": -1})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var pages []models.Page
	cursor, err := pr.db.Collection("pages").Find(ctx, filter, f_opt)
	if err != nil {
		return pages
	}
	var wg sync.WaitGroup
	for cursor.Next(ctx) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var page models.Page
			err := cursor.Decode(page)
			if err != nil {
				return
			}
			pages = append(pages, page)
		}()
	}
	wg.Wait()
	return pages
}

func (pr *PageRepository) AddCommit(pageid_hex string, commit *models.Commit) error {
	id, err := primitive.ObjectIDFromHex(pageid_hex)
	if err != nil {
		return errors.New("page not found")
	}
	page, err := pr.ReadPage(pageid_hex)
	if err != nil {
		return err
	}
	page.Commits = append(page.Commits, *commit)
	page.LastCommitID = commit.ID.Hex()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = pr.db.Collection("pages").ReplaceOne(ctx, bson.M{"_id": id}, *page)
	if err != nil {
		return errors.New("replacement error")
	}
	return nil
}

func (pr *PageRepository) DeletePage(id_hex string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return errors.New("page not found")
	}
	filter := bson.M{"_id": id}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = pr.db.Collection("pages").DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("page not found")
	}
	return nil
}
