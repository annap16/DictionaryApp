
package integration_tests

import (
	"context"
	"fmt"
	"testing"
	"log"
	"os"

    "io"
    "gorm.io/gorm/logger"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"dictionary-app/server/database"
)

type IntegrationTestSuite struct {
	suite.Suite
	DB        *database.DBInterface
	container tc.Container
}

func (s *IntegrationTestSuite) SetupSuite() {
	if err := godotenv.Load("./.env.test"); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	req := tc.ContainerRequest{
		Image:        "postgres:14-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": os.Getenv("POSTGRES_PASSWORD"),
			"POSTGRES_USER": os.Getenv("POSTGRES_USER"),
			"POSTGRES_DB": os.Getenv("POSTGRES_DB"),
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	s.container = container

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("host=%s port=%s user=user password=password dbname=testdb sslmode=disable TimeZone=UTC", host, port.Port())

	gormLogger := logger.New(
		log.New(io.Discard, "", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Silent,
		},
	)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	s.Require().NoError(err)

	err = db.AutoMigrate(&database.Word{}, &database.Translation{}, &database.ExampleSentence{})
	s.Require().NoError(err)

	repo := &database.GormRepository{}
	dbWrapper := &database.GormDatabase{Connection: db}
	s.DB = database.NewDBInterface(dbWrapper, repo)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.container != nil {
		s.container.Terminate(ctx)
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
