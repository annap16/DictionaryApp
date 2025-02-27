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


func (dbI *DBInterface) AddWord(input model.CreateTranslationInput) error {
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
		fmt.Println("Error while adding new word to a database", err)
		return err
	}

	return nil
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
		log.Fatal("Error:", err)
		return nil, err
	}
	if result.ID == 0 {
		err = errors.New("record not found")
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
	if err := dbI.DB.Where("word = ?", input).Find(&word).Error; err != nil {
		log.Fatal("Error:", err)
		return false, err
	}
	if word.ID == 0 {
		return false, nil
	}
	if err := dbI.DB.Where("word = ?", input).Delete(&Word{}).Error; err != nil {
		log.Fatal("Error deleting word:", err)
		return false, err
	}

	return true, nil
}


func (dbI *DBInterface) DeleteTranslation(input string) (bool, error) {
	var translation Translation
	if err := dbI.DB.Where("translation = ?", input).Find(&translation).Error; err != nil {
		log.Fatal("Error:", err)
		return false, err
	}
	if translation.ID == 0 {
		return false, nil
	}
	if err := dbI.DB.Where("translation = ?", input).Delete(&Translation{}).Error; err != nil {
		log.Fatal("Error deleting translation:", err)
		return false, err
	}


	return true, nil

}

func (dbI *DBInterface) DeleteExample(input string)(bool, error){
	var example ExampleSentence
	if err := dbI.DB.Where("sentence = ?", input).Find(&example).Error; err != nil {
		log.Fatal("Error:", err)
		return false, err
	}
	if example.ID == 0 {
		return false, nil
	}
	if err := dbI.DB.Where("sentence = ?", input).Delete(&ExampleSentence{}).Error; err != nil {
		log.Fatal("Error deleting example:", err)
		return false, err
	}

	return true, nil
}

func (dbI *DBInterface) AddTranslation(input model.CreateTranslationInput) (bool, error) {
	var existingWord Word

	err := dbI.DB.Preload("Translations.ExampleSentences"). Where("word = ?", input.Word).Find(&existingWord).Error    
	if err != nil{
		fmt.Println("Error while finding the word:", err)
		return false, err
	}
	if existingWord.ID==0{
		return false, err
	}

	exampleSentences := createExampleSentences(input.Examples)

	newTranslation := Translation{
		Translation:      input.Translation,
		WordID:           existingWord.ID,
		ExampleSentences: exampleSentences,
	}

	if err := dbI.DB.Create(&newTranslation).Error; err != nil {
		fmt.Println("Error while saving translation:", err)
		return false, err
	}

	return true, nil
}


func (dbI *DBInterface) AddExample( translation string, examples []string) (bool, error) {
	var existingTranslation Translation

	err := dbI.DB.Preload("ExampleSentences"). Where("translation = ?", translation).Find(&existingTranslation).Error    
	if err != nil{
		fmt.Println("Error while finding the translation:", err)
		return false, err
	}
	if existingTranslation.ID==0{
		return false, err
	}

	exampleSentences := createExampleSentences(examples)
	
	
	for _, example := range exampleSentences {
		example.TranslationID = existingTranslation.ID 
		if err := dbI.DB.Create(&example).Error; err != nil {
			fmt.Println("Error while saving example sentence:", err)
			return false, err
		}
	}


	return true, nil
}





