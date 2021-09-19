package dto

import "github.com/yosa12978/WikiMD/internal/pkg/models"

type CreatePageDTO struct {
	Name string
	Body string
}

type AddCommitDTO struct {
	PageID string
	Commit *models.Commit
}

type UserSessionDTO struct {
	Username string
	Role     models.Role
}

type CreateCommitDTO struct {
	Name   string
	Body   string
	PageID string
	User   string
}
