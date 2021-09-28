package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/yosa12978/WikiMD/internal/pkg/dto"
	"github.com/yosa12978/WikiMD/internal/pkg/models"
	mongodb "github.com/yosa12978/WikiMD/internal/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICommitRepository interface {
	GetCommitByID(id_hex string) (*models.Commit, error)
	CreateCommit(commit_dto dto.CreateCommitDTO, username string) error
	GetCommitsByPageID(id_hex string) ([]models.Commit, error)
	ChangeCommit(id_hex string) (string, error)
}

type CommitRepository struct {
	db *mongo.Database
}

func NewCommitRepository() ICommitRepository {
	return &CommitRepository{db: mongodb.GetClient()}
}

func (cr *CommitRepository) GetCommitByID(id_hex string) (*models.Commit, error) {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return nil, errors.New("commit not found")
	}
	filter := bson.M{"_id": id}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var commit models.Commit
	cr.db.Collection("commits").FindOne(ctx, filter).Decode(&commit)
	return &commit, nil
}

func (cr *CommitRepository) CreateCommit(commit_dto dto.CreateCommitDTO, username string) error {
	page, err := NewPageRepository().GetPageObj(commit_dto.PageID)
	id := primitive.NewObjectID()
	if err != nil {
		return err
	}
	commit := models.Commit{
		ID:   id,
		Name: commit_dto.Name,
		Body: commit_dto.Body,
		Page: page.ID.Hex(),
		User: username,
		Time: time.Now().Unix(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = cr.db.Collection("commits").InsertOne(ctx, commit)
	if err != nil {
		return err
	}
	page.Commits = append(page.Commits, commit)
	page.LastCommitID = id.Hex()
	page.Name = commit.Name
	pageid, err := primitive.ObjectIDFromHex(commit_dto.PageID)
	if err != nil {
		return errors.New("page not found")
	}
	filter := bson.M{"_id": pageid}
	_, err = cr.db.Collection("pages").ReplaceOne(ctx, filter, page)
	return err
}

func (cr *CommitRepository) GetCommitsByPageID(id string) ([]models.Commit, error) {
	f_opts := options.Find().SetSort(bson.M{"_id": -1})
	filter := bson.M{"page": id}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cursor, err := cr.db.Collection("commits").Find(ctx, filter, f_opts)
	if err != nil {
		return nil, err
	}
	var commits []models.Commit
	if err = cursor.All(ctx, &commits); err != nil {
		return nil, err
	}
	return commits, nil
}

func (cr *CommitRepository) ChangeCommit(id_hex string) (string, error) {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return "", errors.New("commit not found")
	}

	cfilter := bson.M{"_id": id}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var commit models.Commit
	err = cr.db.Collection("commits").FindOne(context.TODO(), cfilter).Decode(&commit)
	if err != nil {
		log.Println(id.Hex())
		return "", errors.New("commit not found")
	}

	page_id, err := primitive.ObjectIDFromHex(commit.Page)
	if err != nil {
		return "", errors.New("page not found")
	}
	pagefilter := bson.M{"_id": page_id}

	var page models.Page
	err = cr.db.Collection("pages").FindOne(ctx, pagefilter).Decode(&page)
	if err != nil {
		return "", errors.New("page not found")
	}

	page.LastCommitID = id.Hex()
	_, err = cr.db.Collection("pages").ReplaceOne(context.TODO(), pagefilter, page)
	return page.ID.Hex(), err
}
