package database

import (
	"fmt"
	"log"
	"errors"
	"strings"
	"github.com/annap16/DictionaryApp/graph/model"
	"gorm.io/gorm"
)

type DBInterface struct {
	DB *gorm.DB
}

func NewDBInterface(db *gorm.DB) *DBInterface {
	return &DBInterface{DB: db}
}

func (dbI *DBInterface) TransactionWrapper(fn func(tx *gorm.DB) (bool, error)) (bool, error) {
	tx := dbI.DB.Begin()
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

func (dbI *DBInterface) AddWord(input model.FullRecordInput) (bool, error) {	
	exampleSentences := createExampleSentences(input.Examples)

	translation := Translation{ 
		Translation:      input.Translation,
		ExampleSentences: exampleSentences, 
	}

	word := Word{ 
		Word:         input.Word,
		Translations: []Translation{translation}, 
	}

	return dbI.TransactionWrapper(func(tx *gorm.DB) (bool, error){
		err := tx.Create(&word).Error
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				return false, nil
			}
			fmt.Println("Error while adding a word to a DB", err)
			return false, err
	}
	return true, nil
	})

}

func createExampleSentences(examples []string) []ExampleSentence {
	var exampleSentences []ExampleSentence
	for _, sentence := range examples {
		exampleSentences = append(exampleSentences, ExampleSentence{
			Sentence: sentence,
		})
	}
	return exampleSentences
}

func (dbI *DBInterface) ReceiveWordTranslation(input string) (*model.Word, error) {
	var result Word

	err := dbI.DB.
		Preload("Translations.ExampleSentences"). 
		Where("word = ?", input).                
		Find(&result).Error    

	if err != nil {
		log.Println("Error while loading word from a DB:", err)
		return nil, err
	}
	if result.ID == 0 {
		err = errors.New("Record not found")
		return nil, err
	}
						
	modelResult := &model.Word{
		ID:   fmt.Sprintf("%d", result.ID),
		Word: result.Word,
		Translations: convertTranslations(result.Translations), 
	}

	return modelResult, nil
}

func convertTranslations(translations []Translation) []*model.Translation {
	var modelTranslations []*model.Translation
	for _, t := range translations {
		modelTranslations = append(modelTranslations, &model.Translation{
			ID:              fmt.Sprintf("%d", t.ID),
			Translation:     t.Translation,
			ExampleSentences: convertExampleSentences(t.ExampleSentences),
		})
	}
	return modelTranslations
}

func convertExampleSentences(sentences []ExampleSentence) []*model.ExampleSentence {
	var modelSentences []*model.ExampleSentence
	for _, s := range sentences {
		modelSentences = append(modelSentences, &model.ExampleSentence{
			ID:       fmt.Sprintf("%d", s.ID),
			Sentence: s.Sentence,
		})
	}
	return modelSentences
}

func (dbI *DBInterface) DeleteWord(input string) (bool, error) {

	return dbI.TransactionWrapper(func(tx *gorm.DB) (bool, error){
		result := tx.Where("word = LOWER(?)", input).Delete(&Word{})
		if result.Error != nil {
			log.Println("Error while deleting word from a DB:", result.Error)
			return false, result.Error
		}
		if result.RowsAffected == 0 {
			return false, nil
		}

		return true, nil
	})
}

func (dbI *DBInterface) GetWordID(tx *gorm.DB, word string) (uint, error) {
	var existingWord Word

	err := tx.Where("word = ?", word).Find(&existingWord).Error    
	if err != nil{
		log.Println("Error while searching for word existance in a DB:", err)
		return 0, err
	}

	return uint(existingWord.ID), nil
}

func (dbI *DBInterface) GetTranslationID(tx *gorm.DB, wordID uint, translation string) (uint, error) {
	var existingTranslation Translation

	err := tx.Where("translation = LOWER(?) AND word_id=(?)", translation, wordID).Find(&existingTranslation).Error    
	if err != nil{
		log.Println("Error while searching for translation existance in a DB:", err)
		return 0, err
	}
	return uint(existingTranslation.ID), nil
}


func (dbI *DBInterface) DeleteTranslation(word string, translation string) (bool, error) {
	return dbI.TransactionWrapper(func(tx *gorm.DB) (bool, error){
		wordID, err := dbI.GetWordID(tx, word)
		if err!=nil{
			return false, err
		} else if wordID==0{
			return false, nil
		}

		result := tx.Where("word_id=(?) AND translation=(?)", wordID, translation).Delete(&Translation{}) 
		if result.Error != nil {
			log.Println("Error while deleting translation :", result.Error)
			return false, err
		}
		if result.RowsAffected == 0 {
			return false, nil
		}

		return true, nil
	})
}

func (dbI *DBInterface) DeleteExample(input model.FullRecordInput)(bool, error){
	return dbI.TransactionWrapper(func(tx *gorm.DB) (bool, error){
		wordID, err := dbI.GetWordID(tx, input.Word)
		if err!=nil{
			return false, err
		} else if wordID==0{
			return false, nil
		}

		translationID, err := dbI.GetTranslationID(tx, wordID, input.Translation)
		if err != nil{
			return false, err
		}else if translationID==0{
			return false, nil
		}
		
		for _, example := range input.Examples {
			result := tx.Where("translation_id=(?) AND LOWER(sentence)=LOWER(?)", translationID, example).Delete(&ExampleSentence{}) 
			if result.Error != nil {
				log.Println("Error while deleting example sentence:", result.Error)
				return false, err
			}
			if result.RowsAffected == 0 {
				return false, nil
			}
		}

		return true, nil
	})

}

func (dbI *DBInterface) AddTranslation(input model.FullRecordInput) (bool, error) {
	return dbI.TransactionWrapper(func(tx *gorm.DB) (bool, error){
		wordID, err := dbI.GetWordID(tx, input.Word)
		if err!=nil{
			return false, err
		} else if wordID==0{
			return false, nil
		}

		exampleSentences := createExampleSentences(input.Examples)

		newTranslation := Translation{
			Translation:      input.Translation,
			WordID:           wordID,
			ExampleSentences: exampleSentences,
		}

		if err := tx.Create(&newTranslation).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				return false, nil
			}
			if strings.Contains(err.Error(), "violates foreign key constraint") {
				log.Println("Error while searching for word existance in a DB:", err)
				return false, errors.New("Word does not exist")
			}
			log.Println("Error while adding translation:", err)
			return false, err
		}

		return true, nil
	})
}


func (dbI *DBInterface) AddExample(input model.FullRecordInput) (bool, error) {
	return dbI.TransactionWrapper(func(tx *gorm.DB) (bool, error){
		wordID, err := dbI.GetWordID(tx, input.Word)
		if err!=nil{
			return false, err
		} else if wordID==0{
			return false, nil
		}

		translationID, err := dbI.GetTranslationID(tx, wordID, input.Translation)
		if err != nil{
			return false, err
		}else if translationID==0{
			return false, nil
		}

		exampleSentences := createExampleSentences(input.Examples)
		
		for _, example := range exampleSentences {
			example.TranslationID = translationID
			if err := tx.Create(&example).Error; err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return false, nil
				}
				if strings.Contains(err.Error(), "violates foreign key constraint") {
					log.Println("Error while searching for translation existance in a DB:", err)
					return false, errors.New("Translation does not exist")
				}
				log.Println("Error while adding example sentence:", err)
				return false, err
			}
		}

		return true, nil
	})
}





