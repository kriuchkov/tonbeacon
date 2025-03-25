package containers

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	// DefaultPostgresUser is the default user for the Postgres container.
	defaultPostgresUser = "test"

	// DefaultPostgresPassword is the default password for the Postgres container.
	defaultPostgresPassword = "test"

	// DefaultPostgresDatabase is the default database for the Postgres container.
	defaultPostgresDatabase = "testdb"
)

type PostgresOptions struct {
	User     string
	Password string
	Database string
}

func (o *PostgresOptions) SetDefaults() {
	if o.User == "" {
		o.User = defaultPostgresUser
	}
	if o.Password == "" {
		o.Password = defaultPostgresPassword
	}
	if o.Database == "" {
		o.Database = defaultPostgresDatabase
	}
}

type Postgres struct {
	testcontainers.Container
	user     string
	password string
	database string
}

func NewPostgres(ctx context.Context, opt *PostgresOptions) (*Postgres, error) {
	opt.SetDefaults()

	container := Postgres{
		user:     opt.User,
		password: opt.Password,
		database: opt.Database,
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     container.user,
			"POSTGRES_PASSWORD": container.password,
			"POSTGRES_DB":       container.database,
		},
	}

	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "start container")
	}

	container.Container = pgC
	return &container, nil
}

func (p *Postgres) GetDSN(ctx context.Context) (string, error) {
	host, err := p.Host(ctx)
	if err != nil {
		return "", errors.Wrap(err, "get host")
	}

	port, err := p.MappedPort(context.Background(), "5432")
	if err != nil {
		return "", errors.Wrap(err, "get port")
	}

	dsn := "postgres://" + p.user + ":" + p.password + "@" + host + ":" + port.Port() + "/" + p.database + "?sslmode=disable"
	return dsn, nil
}

func (p *Postgres) GetUser() string {
	return p.user
}

func (p *Postgres) GetPassword() string {
	return p.password
}

func (p *Postgres) GetDatabase() string {
	return p.database
}
