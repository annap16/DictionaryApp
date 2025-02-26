package database

import (
	"fmt"
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

