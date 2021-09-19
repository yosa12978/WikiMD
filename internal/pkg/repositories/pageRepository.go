package repositories

import (
	"context"
	"errors"
	"time"

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
	ReadPage(id_hex string) (*models.Commit, error)
	GetPages() []models.Page
	SearchPages(query string) []models.Page
	AddCommit(pageid_hex string, commit *models.Commit) error
	DeletePage(id_hex string, username string) error
	GetPageObj(id_hex string) (*models.Page, error)
}

type PageRepository struct {
	db *mongo.Database
}

func NewPageRepository() IPageRepository {
	return &PageRepository{db: mongodb.GetClient()}
}

func (pr *PageRepository) CreatePage(page_dto *dto.CreatePageDTO, username string) error {
	pageid := primitive.NewObjectID()
	commit := models.Commit{
		ID:   primitive.NewObjectID(),
		Name: page_dto.Name,
		Body: page_dto.Body,
		Page: pageid.Hex(),
		User: username,
		Time: time.Now().Unix(),
	}
	page := models.Page{
		ID:           pageid,
		Name:         page_dto.Name,
		LastCommitID: commit.ID.Hex(),
		Commits:      []models.Commit{commit},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := pr.db.Collection("commits").InsertOne(ctx, commit)
	if err != nil {
		return err
	}
	_, err = pr.db.Collection("pages").InsertOne(ctx, page)
	return err
}

func (pr *PageRepository) ReadPage(id_hex string) (*models.Commit, error) {
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

	commit, err := NewCommitRepository().GetCommitByID(page.LastCommitID)
	return commit, err
}

func (pr *PageRepository) GetPages() []models.Page {
	f_opt := options.Find().SetSort(bson.M{"name": 1})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var pages []models.Page
	cursor, err := pr.db.Collection("pages").Find(ctx, bson.M{}, f_opt)
	if err != nil {
		return pages
	}
	// var wg sync.WaitGroup
	for cursor.Next(ctx) {
		// wg.Add(1)
		// go func(c *mongo.Cursor) {
		// 	defer wg.Done()
		var page models.Page
		err := cursor.Decode(&page)
		if err != nil {
			continue
			// return
		}
		pages = append(pages, page)
		// 	}(cursor)
	}
	// wg.Wait()
	return pages
}

func (pr *PageRepository) SearchPages(query string) []models.Page {
	filter := bson.M{"name": primitive.Regex{Pattern: query, Options: "i"}}
	f_opt := options.Find().SetSort(bson.M{"name": -1})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var pages []models.Page
	cursor, err := pr.db.Collection("pages").Find(ctx, filter, f_opt)
	if err != nil {
		return pages
	}
	for cursor.Next(ctx) {
		var page models.Page
		err := cursor.Decode(&page)
		if err != nil {
			continue
		}
		pages = append(pages, page)
	}
	return pages
}
func (pr *PageRepository) GetPageObj(id_hex string) (*models.Page, error) {
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
	return &page, err
}

func (pr *PageRepository) AddCommit(pageid_hex string, commit *models.Commit) error {
	id, err := primitive.ObjectIDFromHex(pageid_hex)
	if err != nil {
		return errors.New("page not found")
	}
	page, err := pr.GetPageObj(pageid_hex)
	if err != nil {
		return err
	}
	page.Commits = append(page.Commits, *commit)
	page.LastCommitID = commit.ID.Hex()
	page.Name = commit.Name
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = pr.db.Collection("pages").ReplaceOne(ctx, bson.M{"_id": id}, *page)
	if err != nil {
		return errors.New("replacement error")
	}
	return nil
}

func (pr *PageRepository) DeletePage(id_hex string, username string) error {
	urole, err := NewUserRepository().GetUserRole(username)
	if err != nil {
		return err
	}
	if urole == models.USER_ROLE {
		return errors.New("forbidden")
	}
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
