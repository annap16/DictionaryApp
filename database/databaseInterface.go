package database

import (
	"fmt"
	"log"
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
		First(&result).Error    

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, err
		}
		log.Fatal("Error:", err)
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
	if err := dbI.DB.Where("word = ?", input).First(&word).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return false, nil
		}
		log.Fatal("Error:", err)
		return false, err
	}
	if err := dbI.DB.Where("word = ?", input).Delete(&Word{}).Error; err != nil {
		log.Fatal("Error deleting word:", err)
		return false, err
	}

	return true, nil
}





