package repository

type RepositorySetup struct {
	HistoryRepository *HistoryRepository
}

func Setup() *RepositorySetup {
	return &RepositorySetup{
		HistoryRepository: NewHistoryRepository(),
	}
}
