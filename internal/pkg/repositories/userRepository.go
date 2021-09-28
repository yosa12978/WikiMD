package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/yosa12978/WikiMD/internal/pkg/dto"
	"github.com/yosa12978/WikiMD/internal/pkg/models"
	mongodb "github.com/yosa12978/WikiMD/internal/pkg/mongo"
	"github.com/yosa12978/WikiMD/pkg/crypto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
	LogInUser(username string, password string) (*dto.UserSessionDTO, error)
	CreateUser(username string, password string) error
	ReadUser(username string) (*models.User, error)
	UpdateUser(id_hex string, username string) error
	DeleteUser(id_hex string, username string) error
	GetUserRole(username string) (models.Role, error)
	ReadByNameAndPass(username string, password string) (models.User, error)
	ChangeRole(role models.Role, username string) error
}

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository() IUserRepository {
	return &UserRepository{db: mongodb.GetClient()}
}

func (ur *UserRepository) LogInUser(username string, password string) (*dto.UserSessionDTO, error) {
	user, err := ur.ReadByNameAndPass(username, password)
	if err != nil {
		return nil, err
	}
	compl := dto.UserSessionDTO{Username: user.Username, Role: user.Role}
	return &compl, nil
}

func (ur *UserRepository) CreateUser(username string, password string) error {
	_, err := ur.ReadUser(username)
	if err == nil {
		return errors.New("username is already in use")
	}
	new_user := models.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: crypto.GetMD5(password),
		//Email:    email,
		Regdate: time.Now().Unix(),
		Token:   crypto.GetToken32(),
		Role:    models.USER_ROLE,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = ur.db.Collection("users").InsertOne(ctx, new_user)
	return err
}

func (ur *UserRepository) ReadUser(username string) (*models.User, error) {
	filter := bson.M{"username": username}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var user models.User
	err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (ur *UserRepository) UpdateUser(id_hex string, username string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return errors.New("user not found")
	}
	set_list := bson.M{"$set": []bson.M{{"username": username}}}
	filter := bson.M{"_id": id}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = ur.db.Collection("users").UpdateOne(ctx, filter, set_list)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
	}
	return err
}

func (ur *UserRepository) DeleteUser(id_hex string, username string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return errors.New("user not found")
	}
	filter := bson.M{"$and": []bson.M{{"_id": id}, {"username": username}}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = ur.db.Collection("users").DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
	}
	return err
}

func (ur *UserRepository) GetUserRole(username string) (models.Role, error) {
	filter := bson.M{"username": username}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var user models.User
	err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return "", errors.New("user not found")
	}
	return user.Role, nil
}

func (ur *UserRepository) ReadByNameAndPass(username string, password string) (models.User, error) {
	filter := bson.M{"$and": []bson.M{{"username": username}, {"password": crypto.GetMD5(password)}}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var user models.User
	err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user not found")
		}
	}
	return user, err
}

func (ur *UserRepository) ChangeRole(role models.Role, username string) error {
	filter := bson.M{"username": username}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var user models.User
	err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return errors.New("user not found")
	}
	user.Role = role
	_, err = ur.db.Collection("users").ReplaceOne(ctx, filter, user)
	return err
}
