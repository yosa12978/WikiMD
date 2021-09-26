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
	ChangeCommit(id_hex string, page_id_hex string) error
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return errors.New("commit not found")
	}

	commit_filter := bson.M{"_id": id}
	var commit models.Commit
	err = cr.db.Collection("commits").FindOne(ctx, commit_filter).Decode(&commit)
	if err != nil {
		return errors.New("commit not found")
	}

	page_id, err := primitive.ObjectIDFromHex(commit.Page)
	if err != nil {
		return errors.New("page not found")
	}

	page_filter := bson.M{"_id": page_id}
	var page models.Page
	err = cr.db.Collection("pages").FindOne(ctx, page_filter).Decode(&page)
	if err != nil {
		return errors.New("page not found")
	}

	for i := 0; i < len(page.Commits)-1; i++ {
		for j := i; j < len(page.Commits)-1-i; j++ {
			if page.Commits[j].ID.Hex() > page.Commits[j+1].ID.Hex() {
				temp := page.Commits[j]
				page.Commits[j] = page.Commits[j+1]
				page.Commits[j+1] = temp
			}
		}
	}

	var s int
	var low int = 0
	var high int = len(page.Commits) - 1
	for low <= high {
		var middle int = (low + high) / 2
		var guess string = page.Commits[middle].ID.Hex()
		if guess == id_hex {
			s = middle
			break
		} else if guess > id_hex {
			high = middle - 1
		} else {
			low = middle + 1
		}
	}
	page.Commits = append(page.Commits[:s], page.Commits[s+1:]...)
	cr.ChangeCommit(page.Commits[len(page.Commits)-1].ID.Hex(), page_id.Hex())

	_, err = cr.db.Collection("commits").DeleteOne(ctx, commit_filter)
	if err != nil {
		return errors.New("commit not found")
	}

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

func (cr *CommitRepository) ChangeCommit(id_hex string, page_id_hex string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return errors.New("commit not found")
	}
	page_id, err := primitive.ObjectIDFromHex(page_id_hex)
	if err != nil {
		return errors.New("page not found")
	}

	pagefilter := bson.M{"_id": page_id}
	cfilter := bson.M{"_id": id}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var page models.Page
	err = cr.db.Collection("pages").FindOne(ctx, pagefilter).Decode(&page)
	if err != nil {
		return errors.New("page not found")
	}

	var commit models.Commit
	err = cr.db.Collection("commits").FindOne(ctx, cfilter).Decode(&commit)
	if err != nil {
		return errors.New("commit not found")
	}

	page.LastCommitID = id.Hex()
	_, err = cr.db.Collection("pages").ReplaceOne(context.TODO(), pagefilter, page)
	return err
}
