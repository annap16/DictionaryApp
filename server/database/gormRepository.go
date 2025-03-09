package database


import (
	"gorm.io/gorm"
	"log"
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
	return tx.Create(word).Error
}

func (r *GormRepository)CreateTranslation(translation *Translation, tx *gorm.DB) error{
	return tx.Create(translation).Error
}

func (r *GormRepository)CreateExample(example *ExampleSentence, tx *gorm.DB) error{
	return tx.Create(example).Error
}

func (r *GormRepository)DeleteWord(word string, tx *gorm.DB) (bool, error){
	result := tx.Where("word = LOWER(?)", word).Delete(&Word{})
	if result.Error != nil {
		log.Println("Error while deleting word from a DB:", result.Error)
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
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
		return false, nil
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
		return false, nil
	}
	return true, nil
}
func (r *GormRepository)GetWord(word string, result *Word, tx *gorm.DB) error{
	return tx.Preload("Translations.ExampleSentences").Where("word = ?", word).Find(&result).Error    
}
func (r *GormRepository)GetTranslation(wordID uint, translation string, result *Translation, tx *gorm.DB) error{
	return tx.Where("translation = LOWER(?) AND word_id=(?)", translation, wordID).Find(&result).Error    
}







