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

type ICommitRepository interface {
	GetCommitByID(id_hex string) (*models.Commit, error)
	CreateCommit(commit_dto dto.CreateCommitDTO, username string) error
	DeleteCommit(id_hex string) error
	GetCommitsByPageID(id_hex string) ([]models.Commit, error)
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

func (cr *CommitRepository) DeleteCommit(id_hex string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return errors.New("commit not found")
	}
	filter := bson.M{"_id": id}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = cr.db.Collection("commits").DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("commit not found")
	}
	return nil
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
