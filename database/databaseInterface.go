package database

import (
	"fmt"
	"log"
	"errors"
	"github.com/annap16/DictionaryApp/graph/model"
	"gorm.io/gorm"
)

type DBInterface struct {
	DB *gorm.DB
}

func NewDBInterface(db *gorm.DB) *DBInterface {
	return &DBInterface{DB: db}
}


func (dbI *DBInterface) AddWord(input model.CreateTranslationInput) (bool, error) {
	tx := dbI.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() 
		}
	}()

	var existingWord Word
	if err := tx.Where("word = ?", input.Word).Find(&existingWord).Error; err != nil {
		tx.Rollback() 
		log.Println("Error while checking word existance in a DB:", err)
		return false, err
	}
	if existingWord.ID!=0{
		return false, nil
	}
	
	exampleSentences := createExampleSentences(input.Examples)

	translation := Translation{ 
		Translation:      input.Translation,
		ExampleSentences: exampleSentences, 
	}

	word := Word{ 
		Word:         input.Word,
		Translations: []Translation{translation}, 
	}

	err := tx.Create(&word).Error
	if err != nil {
		tx.Rollback() 
		fmt.Println("Error while adding a word to a DB", err)
		return false, err
	}
	
	tx.Commit()

	return true, nil
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
	tx := dbI.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() 
		}
	}()

	var word Word
	if err := tx.Where("word = LOWER(?)", input).Find(&word).Error; err != nil {
		tx.Rollback() 
		log.Println("Error while checking word existance in a DB:", err)
		return false, err
	}
	if word.ID == 0 {
		tx.Rollback() 
		return false, nil
	}
	if err := tx.Where("word = LOWER(?)", input).Delete(&Word{}).Error; err != nil {
		tx.Rollback() 
		log.Println("Error while deleting word from a DB:", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}


func (dbI *DBInterface) DeleteTranslation(input string) (bool, error) {
	tx := dbI.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() 
		}
	}()

	var translation Translation
	if err := tx.Where("translation = LOWER(?)", input).Find(&translation).Error; err != nil {
		tx.Rollback() 
		log.Println("Error while searching for translation existance in a DB:", err)
		return false, err
	}
	if translation.ID == 0 {
		tx.Rollback() 
		return false, nil
	}


	if err := tx.Where("translation = LOWER(?)", input).Delete(&Translation{}).Error; err != nil {
		tx.Rollback() 

		log.Println("Error while deleting translation:", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (dbI *DBInterface) DeleteExample(translation string, input string)(bool, error){
	tx := dbI.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() 
		}
	}()
	
	var existingTranslation Translation
	if err := tx.Where("translation = ?", translation).Find(&existingTranslation).Error; err!=nil{
		tx.Rollback() 
		log.Println("Error while searching for translation existance in a DB:", err)
		return false, err
	}
	if existingTranslation.ID==0{
		tx.Rollback() 
		return false, nil
	}


	var example ExampleSentence
	if err := tx.Where("LOWER(sentence) = LOWER(?) AND translation_id=?", input, existingTranslation.ID).Find(&example).Error; err != nil {
		tx.Rollback() 
		log.Println("Error while searching for example existance in a DB:", err)
		return false, err
	}
	if example.ID == 0 {
		tx.Rollback() 
		return false, nil
	}
	if err := tx.Where("LOWER(sentence) = LOWER(?) AND translation_id=?", input, existingTranslation.ID).Delete(&ExampleSentence{}).Error; err != nil {
		tx.Rollback() 
		log.Println("Error while deleting example:", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (dbI *DBInterface) AddTranslation(input model.CreateTranslationInput) (bool, error) {
	tx := dbI.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() 
		}
	}()
	
	var existingWord Word

	err := tx.Where("word = ?", input.Word).Find(&existingWord).Error    
	if err != nil{
		tx.Rollback() 
		log.Println("Error while searching for word existance in a DB:", err)
		return false, err
	}
	if existingWord.ID==0{
		tx.Rollback() 
		return false, err
	}

	var existingTranslation Translation
	err = tx.Where("translation = ?", input.Translation).Find(&existingTranslation).Error    
	if err != nil{
		tx.Rollback() 
		log.Println("Error while searching for translation existance in a DB:", err)
		return false, err
	}
	if existingTranslation.ID!=0{
		tx.Rollback() 
		return false, err
	}

	exampleSentences := createExampleSentences(input.Examples)

	newTranslation := Translation{
		Translation:      input.Translation,
		WordID:           existingWord.ID,
		ExampleSentences: exampleSentences,
	}

	if err := tx.Create(&newTranslation).Error; err != nil {
		tx.Rollback() 
		log.Println("Error while adding translation:", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}


func (dbI *DBInterface) AddExample(translation string, examples []string) (bool, error) {
	tx := dbI.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() 
		}
	}()

	var existingTranslation Translation

	err := tx.Where("translation = LOWER(?)", translation).Find(&existingTranslation).Error    
	if err != nil{
		tx.Rollback() 
		log.Println("Error while searching for translation existance in a DB:", err)
		return false, err
	}
	if existingTranslation.ID==0{
		tx.Rollback() 
		return false, err
	}	

	exampleSentences := createExampleSentences(examples)
	
	for _, example := range exampleSentences {
		example.TranslationID = existingTranslation.ID 
		if err := tx.Create(&example).Error; err != nil {
			tx.Rollback() 
			log.Println("Error while adding example sentence:", err)
			return false, err
		}
	}

	tx.Commit()

	return true, nil
}






