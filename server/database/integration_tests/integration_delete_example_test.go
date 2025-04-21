package integration_tests

import (
	"dictionary-app/server/graph/model"
	"sync"
	"sync/atomic"
)

// --------------------------- SEQUENTIAL TESTS ---------------------------

// Deleting an existing example
func (s *IntegrationTestSuite) TestDeleteExample() {
	word := "exampleword"
	translation := "exampletranslation"
	example := "exampletodelete"

	input := model.FullRecordInput{
		Word: word,
		Translation: translation,
		Examples: []string{example},
	}

	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	success, err := s.DB.DeleteExample(input)
	s.Require().NoError(err)
	s.Require().True(success)

	received, err := s.DB.ReceiveWordTranslation(word)
	s.Require().NoError(err)
	s.Require().NotNil(received)
	s.Require().Len(received.Translations[0].ExampleSentences, 0)
}

// Deleting an example from a non existing word
func (s *IntegrationTestSuite) TestDeleteExample_NonExistingWord() {
	input := model.FullRecordInput{
		Word: "nonexistentword",
		Translation: "translation",
		Examples: []string{"example"},
	}

	success, err := s.DB.DeleteExample(input)
	s.Require().NoError(err)
	s.Require().False(success)
}

// Deleting an example from a non existing translation
func (s *IntegrationTestSuite) TestDeleteExample_NonExistingTranslation() {
	word := "word"
	inputWord := model.FullRecordInput{
		Word: word,
		Translation: "translation",
		Examples: []string{"example"},
	}
	ok, err := s.DB.AddWord(inputWord)
	s.Require().NoError(err)
	s.Require().True(ok)

	input := model.FullRecordInput{
		Word:  word,
		Translation: "nonexistingtranslation",
		Examples: []string{"example"},
	}

	success, err := s.DB.DeleteExample(input)
	s.Require().NoError(err)
	s.Require().False(success)
}

// --------------------------- PARALLEL TESTS ---------------------------

// Deleting different examples in parallel
func (s *IntegrationTestSuite) TestDeleteDifferentExamplesInParallel() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		input := model.FullRecordInput{
			Word: "word1",
			Translation: "translation1",
			Examples: []string{"example1"},
		}
		ok, err := s.DB.AddWord(input)
		s.Require().NoError(err)
		s.Require().True(ok)

		success, err := s.DB.DeleteExample(input)
		s.Require().NoError(err)
		s.Require().True(success)
	}()

	go func() {
		defer wg.Done()
		input := model.FullRecordInput{
			Word: "word2",
			Translation: "translation2",
			Examples: []string{"example2"},
		}
		ok, err := s.DB.AddWord(input)
		s.Require().NoError(err)
		s.Require().True(ok)

		success, err := s.DB.DeleteExample(input)
		s.Require().NoError(err)
		s.Require().True(success)
	}()

	wg.Wait()
}

// Deleting the same example concurrently
func (s *IntegrationTestSuite) TestDeleteSameExampleConcurrently() {
	const concurrencyLevel = 5
	var wg sync.WaitGroup
	var successCount int32
	var errorCount int32

	word := "word"
	translation := "translation"
	example := "example"

	input := model.FullRecordInput{
		Word: word,
		Translation: translation,
		Examples: []string{example},
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
			success, err := s.DB.DeleteExample(input)

			if success && err == nil {
				atomic.AddInt32(&successCount, 1)
			} else if !success && err != nil {
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
	s.Require().Len(received.Translations[0].ExampleSentences, 0)
}
