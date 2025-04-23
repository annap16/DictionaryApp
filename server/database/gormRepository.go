package database


import (
	"gorm.io/gorm"
	"fmt"
	"log"
	customerrors "dictionary-app/server/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"errors"
)

type Repository interface {
	TransactionWrapper(DB Database, fn func(tx *gorm.DB) (bool, error)) (bool, error)
    CreateWord(word *Word, tx *gorm.DB) error
	CreateTranslation(translation *Translation, tx *gorm.DB) error
	CreateExample(example *ExampleSentence, tx *gorm.DB) error
	DeleteWord(word string, tx *gorm.DB) (bool, error)
	DeleteTranslation(wordID uint, translation string, tx *gorm.DB) (bool, error)
	DeleteExample(translationID uint, example string, tx *gorm.DB) (bool, error)
	GetWord(word string, result *Word, tx *gorm.DB) error
	GetTranslation(wordID uint, translation string, result *Translation, tx *gorm.DB) error
}

type GormRepository struct {
}

func (r *GormRepository) TransactionWrapper(DB Database, fn func(tx *gorm.DB) (bool, error)) (bool, error) {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	success, err := fn(tx)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return success, nil
}

func (r *GormRepository)CreateWord(word *Word, tx *gorm.DB) error{
	err := tx.Create(word).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return customerrors.NewDuplicateKeyError("Podane słowo istnieje już w bazie danych")
			} else {
				fmt.Println("Postgres error:", pgErr.Message)
			}
		} else {
			fmt.Println("Other error:", err)
		}
	}
	return err
}

func (r *GormRepository)CreateTranslation(translation *Translation, tx *gorm.DB) error{
	err := tx.Create(translation).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return customerrors.NewDuplicateKeyError("Nie można dodać przykładu – narusza on unikalność rekordów")
			case "23503":
				return customerrors.NewForeignKeyError("Nie można dodać przykładu – powiązane tłumaczenie nie istnieje")
			default:
				fmt.Println("Postgres error:", pgErr.Message)
			}
		} else {
			fmt.Println("Other error:", err)
		}
		return err
	}
	return nil

}

func (r *GormRepository)CreateExample(example *ExampleSentence, tx *gorm.DB) error{
	err := tx.Create(example).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return customerrors.NewDuplicateKeyError("Nie można dodać przykładu – narusza on unikalność rekordów")
			case "23503":
				return customerrors.NewForeignKeyError("Nie można dodać przykładu – powiązane tłumaczenie nie istnieje")
			default:
				fmt.Println("Postgres error:", pgErr.Message)
			}
		} else {
			fmt.Println("Other error:", err)
		}
		return err
	}
	return nil
}

func (r *GormRepository)DeleteWord(word string, tx *gorm.DB) (bool, error){
	result := tx.Where("word = LOWER(?)", word).Delete(&Word{})
	if result.Error != nil {
		log.Println("Error while deleting word from a DB:", result.Error)
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, customerrors.NewNotFoundError("Nie usunięto podanego słowa - słowa nie znaleziono")
	}
	return true, nil
}

func (r *GormRepository)DeleteTranslation(wordID uint, translation string, tx *gorm.DB) (bool, error) {
	result :=  tx.Where("word_id=(?) AND translation=(?)", wordID, translation).Delete(&Translation{}) 
	if result.Error != nil {
		log.Println("Error while deleting translation :", result.Error)
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, customerrors.NewNotFoundError("Nie usunięto podanego tłumaczenia - tłumaczenia nie znaleziono")
	}
	return true, nil
}

func (r *GormRepository)DeleteExample(translationID uint, sentence string, tx *gorm.DB) (bool, error) {
	result := tx.Where("translation_id=(?) AND LOWER(sentence)=LOWER(?)", translationID, sentence).Delete(&ExampleSentence{}) 
	if result.Error != nil {
		log.Println("Error while deleting example sentence:", result.Error)
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false,  customerrors.NewNotFoundError("Nie usunięto podanego przykładu - przykładu nie znaleziono")
	}
	return true, nil
}
func (r *GormRepository)GetWord(word string, result *Word, tx *gorm.DB) error{
	return tx.Preload("Translations.ExampleSentences").Where("word = ?", word).Find(&result).Error    
}
func (r *GormRepository)GetTranslation(wordID uint, translation string, result *Translation, tx *gorm.DB) error{
	return tx.Where("translation = LOWER(?) AND word_id=(?)", translation, wordID).Find(&result).Error    
}







