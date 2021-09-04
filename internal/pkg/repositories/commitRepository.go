package repositories

type ICommitRepository interface {
}

type CommitRepository struct{}

func NewCommitRepository() ICommitRepository {
	return &CommitRepository{}
}
