package migrator

type Repository interface{}

type Migrator struct {
	repo Repository
}

func New(repo Repository) *Migrator {
	return &Migrator{
		repo: repo,
	}
}

func (m *Migrator) Create() error {}
