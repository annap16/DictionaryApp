package database

import (
	"fmt"
	"errors"
	"strings"
	"dictionary-app/server/graph/model"
	"gorm.io/gorm"
)

type DBInterface struct {
	DB Database
	repo Repository
}

func NewDBInterface(db Database, repo Repository) *DBInterface {
	return &DBInterface{DB: db, repo: repo}
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

	return dbI.repo.TransactionWrapper(dbI.DB,func(tx *gorm.DB) (bool, error){
		err := dbI.repo.CreateWord(&word, tx)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				return dbI.AddTranslation(input)
			}
			fmt.Println("Error while adding a word to a DB", err)
			return false, err
		}
		return true, nil
	})
}

func (dbI *DBInterface) ReceiveWordTranslation(input string) (*model.Word, error) {
	var result Word

	success, err := dbI.repo.TransactionWrapper(dbI.DB,func(tx *gorm.DB) (bool, error){
		err := dbI.repo.GetWord(input, &result, tx)
		if err != nil {
			fmt.Println("Error while loading word from a DB:", err)
			return false, err
		}
		if result.ID == 0 {
			err := errors.New("Record not found")
			return false, err
		}
		return true, nil
	})

	if !success{
		return nil, err
	}
						
	modelResult := &model.Word{
		ID:   fmt.Sprintf("%d", result.ID),
		Word: result.Word,
		Translations: convertTranslations(result.Translations), 
	}

	return modelResult, nil
}


func (dbI *DBInterface) DeleteWord(input string) (bool, error) {
	return dbI.repo.TransactionWrapper(dbI.DB,func(tx *gorm.DB) (bool, error){
		return dbI.repo.DeleteWord(input, tx)
	})
}

func (dbI *DBInterface) GetWordID(tx *gorm.DB, word string) (uint, error) {
	var existingWord Word

	err := dbI.repo.GetWord(word, &existingWord, tx)
	if err != nil{
		fmt.Println("Error while searching for word existance in a DB:", err)
		return 0, err
	}

	return uint(existingWord.ID), nil
}

func (dbI *DBInterface) GetTranslationID(tx *gorm.DB, wordID uint, translation string) (uint, error) {
	var existingTranslation Translation

	err := dbI.repo.GetTranslation(wordID, translation, &existingTranslation, tx)
	if err != nil{
		fmt.Println("Error while searching for translation existance in a DB:", err)
		return 0, err
	}
	return uint(existingTranslation.ID), nil
}


func (dbI *DBInterface) DeleteTranslation(word string, translation string) (bool, error) {
	return dbI.repo.TransactionWrapper(dbI.DB,func(tx *gorm.DB) (bool, error){
		wordID, err := dbI.GetWordID(tx, word)
		if err!=nil{
			return false, err
		} else if wordID==0{
			return false, nil
		}

		success, err := dbI.repo.DeleteTranslation(wordID, translation, tx)
		if !success && err==nil{
			return success, errors.New("Nie znaleziono tłumaczenia")
		}
		
		return success, err
	})
}

func (dbI *DBInterface) DeleteExample(input model.FullRecordInput)(bool, error){
	return dbI.repo.TransactionWrapper(dbI.DB,func(tx *gorm.DB) (bool, error){
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
			success, err := dbI.repo.DeleteExample(translationID, example, tx)
			if !success || err!=nil{
				return success, errors.New("Nie można było usunąć przykładu ze słownika. Anulowano całą operację")
			}
		}

		return true, nil
	})

}

func (dbI *DBInterface) AddTranslation(input model.FullRecordInput) (bool, error) {
	return dbI.repo.TransactionWrapper(dbI.DB,func(tx *gorm.DB) (bool, error){
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

		err = dbI.repo.CreateTranslation(&newTranslation, tx)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				return false, errors.New("Nie można dodać tłumaczenia – narusza ono unikalność rekordów")
			}
			if strings.Contains(err.Error(), "violates foreign key constraint") {
				return false, errors.New("Wystąpił błąd podczas dodawania tłumaczenia")
			}
			return false, err
		}
		return true, nil
		})
}

func (dbI *DBInterface) AddExample(input model.FullRecordInput) (bool, error) {
	return dbI.repo.TransactionWrapper(dbI.DB,func(tx *gorm.DB) (bool, error){
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
			err:= dbI.repo.CreateExample(&example, tx)
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					return false, errors.New("Nie można dodać przykładu - narusza ono unikalność rekordów")
				}
				if strings.Contains(err.Error(), "violates foreign key constraint") {
					return false, errors.New("Wystąpił błąd podczas dodawania przykładu")
				}
				return false, err
			}
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





