package database

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"database/sql"
)
type Database interface {
	Begin(opts ...*sql.TxOptions) *gorm.DB         
    AutoMigrate(dst ...interface{}) error
}

type GormDatabase struct {
	Connection *gorm.DB
}

func (g *GormDatabase) Begin(opts ...*sql.TxOptions) *gorm.DB {
	return g.Connection.Begin() 
}

func (g *GormDatabase) AutoMigrate(dst ...interface{}) error {
	return g.Connection.AutoMigrate(dst...)
}

func NewGormDatabase() (*GormDatabase, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return nil, nil
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to GormDatabase: %w", err)
		return nil, nil
	}

	fmt.Println("Connected to PostgreSQL!")

	GormDatabase := &GormDatabase{Connection: db}
	GormDatabase.migrate()
	return GormDatabase, nil
}

func (d *GormDatabase) migrate() {
	err := d.Connection.AutoMigrate(&Word{}, &Translation{}, &ExampleSentence{})
	if err != nil {
		log.Fatal("Error migrating GormDatabase: %v", err)
		return
	}
	fmt.Println("GormDatabase migration completed.")
}

