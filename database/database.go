package database

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Connection *gorm.DB
}

func NewDatabase() (*Database, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return nil, nil
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %w", err)
		return nil, nil
	}

	fmt.Println("Connected to PostgreSQL!")

	database := &Database{Connection: db}
	database.migrate()
	return database, nil
}

func (d *Database) migrate() {
	err := d.Connection.AutoMigrate(&Word{}, &Translation{}, &ExampleSentence{})
	if err != nil {
		log.Fatal("Error migrating database: %v", err)
		return
	}
	fmt.Println("Database migration completed.")
}
