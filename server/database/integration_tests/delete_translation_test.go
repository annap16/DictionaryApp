package integration_tests

import (
	"dictionary-app/server/graph/model"
	"sync"
	"sync/atomic"
)

// --------------------------- SEQUENTIAL TESTS ---------------------------

// Deleting an existing translation
func (s *IntegrationTestSuite) TestDeleteTranslation() {
	word := "word"
	translation := "translation"

	input := model.FullRecordInput{
		Word: word,
		Translation: translation,
		Examples: []string{"example1"},
	}

	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	success, err := s.DB.DeleteTranslation(word, translation)
	s.Require().NoError(err)
	s.Require().True(success)

	receivedWord, err := s.DB.ReceiveWordTranslation(word)
	s.Require().NoError(err)
	s.Require().NotNil(receivedWord)
	s.Require().Len(receivedWord.Translations, 0)
}

// Deleting non existing translation
func (s *IntegrationTestSuite) TestDeleteNonExistingTranslation() {
	word := "nonexistingword"
	translation := "nonexistingtranslation"

	success, err := s.DB.DeleteTranslation(word, translation)
	s.Require().Error(err)
	s.EqualError(err, "Nie usunięto tłumaczenia - nie znaleziono związanego z nim słowa")	
	s.Require().False(success)
}

// --------------------------- PARALLEL TESTS ---------------------------

// Deleting different translations in parallel
func (s *IntegrationTestSuite) TestDeleteDifferentTranslationsInParallel() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		word := "word1"
		translation := "translation1"

		input := model.FullRecordInput{
			Word: word,
			Translation: translation,
			Examples: []string{"example1"},
		}
		ok, err := s.DB.AddWord(input)
		s.Require().NoError(err)
		s.Require().True(ok)

		success, err := s.DB.DeleteTranslation(word, translation)
		s.Require().NoError(err)
		s.Require().True(success)
	}()

	go func() {
		defer wg.Done()
		word := "word2"
		translation := "translation2"

		input := model.FullRecordInput{
			Word: word,
			Translation: translation,
			Examples: []string{"example2"},
		}
		ok, err := s.DB.AddWord(input)
		s.Require().NoError(err)
		s.Require().True(ok)

		success, err := s.DB.DeleteTranslation(word, translation)
		s.Require().NoError(err)
		s.Require().True(success)
	}()

	wg.Wait()
}

// Deleting the same translation concurrently
func (s *IntegrationTestSuite) TestDeleteSameTranslationConcurrently() {
	const concurrencyLevel = 5
	var wg sync.WaitGroup
	var successCount int32
	var errorCount int32

	word := "concurrentword"
	translation := "concurrenttranslation"

	input := model.FullRecordInput{
		Word: word,
		Translation: translation,
		Examples: []string{"example"},
	}

	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	wg.Add(concurrencyLevel)
	startBarrier := make(chan struct{})

	for i := 0; i < concurrencyLevel; i++ {
		go func() {
			defer wg.Done()
			<-startBarrier
			success, err := s.DB.DeleteTranslation(word, translation)

			if success==true && err==nil {
				atomic.AddInt32(&successCount, 1)
			} else if success==false && err!=nil{
				atomic.AddInt32(&errorCount, 1)
			}
		}()
	}

	close(startBarrier)
	wg.Wait()

	s.Require().Equal(int32(1), successCount)
	s.Require().Equal(int32(concurrencyLevel-1), errorCount)

	received, err := s.DB.ReceiveWordTranslation(word)
	s.Require().NoError(err)
	s.Require().NotNil(received)
	s.Require().Len(received.Translations, 0)
}

