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
	var existingWord Word
	if err := dbI.DB.Where("LOWER(word) = LOWER(?)", input.Word).Find(&existingWord).Error; err != nil {
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

	if err := dbI.DB.Create(&word).Error; err != nil {
		fmt.Println("Error while adding a word to a DB", err)
		return false, err
	}

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
		Where("LOWER(word) = LOWER(?)", input).                
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
	var word Word
	if err := dbI.DB.Where("LOWER(word) = LOWER(?)", input).Find(&word).Error; err != nil {
		log.Println("Error while checking word existance in a DB:", err)
		return false, err
	}
	if word.ID == 0 {
		return false, nil
	}
	if err := dbI.DB.Where("LOWER(word) = LOWER(?)", input).Delete(&Word{}).Error; err != nil {
		log.Println("Error while deleting word from a DB:", err)
		return false, err
	}

	return true, nil
}


func (dbI *DBInterface) DeleteTranslation(input string) (bool, error) {
	var translation Translation
	if err := dbI.DB.Where("LOWER(translation) = LOWER(?)", input).Find(&translation).Error; err != nil {
		log.Println("Error while searching for translation existance in a DB:", err)
		return false, err
	}
	if translation.ID == 0 {
		return false, nil
	}
	if err := dbI.DB.Where("LOWER(translation) = LOWER(?)", input).Delete(&Translation{}).Error; err != nil {
		log.Println("Error while deleting translation:", err)
		return false, err
	}

	return true, nil
}

func (dbI *DBInterface) DeleteExample(input string)(bool, error){
	var example ExampleSentence
	if err := dbI.DB.Where("LOWER(sentence) = LOWER(?)", input).Find(&example).Error; err != nil {
		log.Println("Error while searching for example existance in a DB:", err)
		return false, err
	}
	if example.ID == 0 {
		return false, nil
	}
	if err := dbI.DB.Where("LOWER(sentence) = LOWER(?)", input).Delete(&ExampleSentence{}).Error; err != nil {
		log.Println("Error while deleting example:", err)
		return false, err
	}

	return true, nil
}

func (dbI *DBInterface) AddTranslation(input model.CreateTranslationInput) (bool, error) {
	var existingWord Word

	err := dbI.DB.Where("LOWER(word) = LOWER(?)", input.Word).Find(&existingWord).Error    
	if err != nil{
		log.Println("Error while searching for word existance in a DB:", err)
		return false, err
	}
	if existingWord.ID==0{
		return false, err
	}

	var existingTranslation Translation
	err = dbI.DB.Where("LOWER(translation) = LOWER(?)", input.Translation).Find(&existingTranslation).Error    
	if err != nil{
		log.Println("Error while searching for translation existance in a DB:", err)
		return false, err
	}
	if existingTranslation.ID!=0{
		return false, err
	}


	exampleSentences := createExampleSentences(input.Examples)

	newTranslation := Translation{
		Translation:      input.Translation,
		WordID:           existingWord.ID,
		ExampleSentences: exampleSentences,
	}

	if err := dbI.DB.Create(&newTranslation).Error; err != nil {
		log.Println("Error while adding translation:", err)
		return false, err
	}

	return true, nil
}


func (dbI *DBInterface) AddExample( translation string, examples []string) (bool, error) {
	var existingTranslation Translation

	err := dbI.DB.Where("LOWER(translation) = LOWER(?)", translation).Find(&existingTranslation).Error    
	if err != nil{
		log.Println("Error while searching for translation existance in a DB:", err)
		return false, err
	}
	if existingTranslation.ID==0{
		return false, err
	}

	for _, example := range examples {
		var existingExample ExampleSentence
		err := dbI.DB.Where("LOWER(sentence) = LOWER(?)", example).Find(&existingExample).Error
		if err != nil {
			log.Println("Error while searching for example existence in DB:", err)
			return false, err
		}
		
		if existingExample.ID != 0 && existingTranslation.ID == existingExample.TranslationID {
			return false, nil 
		}
	}
	

	exampleSentences := createExampleSentences(examples)
	
	for _, example := range exampleSentences {
		example.TranslationID = existingTranslation.ID 
		if err := dbI.DB.Create(&example).Error; err != nil {
			log.Println("Error while adding example sentence:", err)
			return false, err
		}
	}

	return true, nil
}





