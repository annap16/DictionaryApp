package integration_tests

import (
	"dictionary-app/server/graph/model"
	"sync"
	"sync/atomic"
)

// --------------------------- SEQUENTIAL TESTS ---------------------------

// Deleting an existing word
func (s *IntegrationTestSuite) TestDeleteWord() {
	input := model.FullRecordInput{
		Word: "wordtodelete",
		Translation: "translationtodelete",
		Examples: []string{"example1", "example2"},
	}
	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	word, err := s.DB.ReceiveWordTranslation("wordtodelete")
	s.Require().NoError(err)
	s.Equal("wordtodelete", word.Word)
	s.Len(word.Translations, 1)
	s.Equal("translationtodelete", word.Translations[0].Translation)

	ok, err = s.DB.DeleteWord("wordtodelete")
	s.Require().NoError(err)
	s.Require().True(ok)

	_, err = s.DB.ReceiveWordTranslation("wordtodelete")
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "Record not found")
}

// Deleting non existing word
func (s *IntegrationTestSuite) TestDeleteNonExistingWord() {
	ok, err := s.DB.DeleteWord("nonexistingword")
	s.Require().NoError(err)
	s.Require().False(ok) 

	_, err = s.DB.ReceiveWordTranslation("nonexistingword")
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "Record not found")
}

// --------------------------- PARALLEL TESTS ---------------------------

// Deleting different words in parallel
func (s *IntegrationTestSuite) TestDeleteWordsInParallel() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		input := model.FullRecordInput{
			Word: "word1todelete",
			Translation: "translation1todelete",
			Examples: []string{"example1", "example2"},
		}
		ok, err := s.DB.AddWord(input)
		s.Require().NoError(err)
		s.Require().True(ok)

		ok, err = s.DB.DeleteWord("word1todelete")
		s.Require().NoError(err)
		s.Require().True(ok)

		_, err = s.DB.ReceiveWordTranslation("word1todelete")
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "Record not found")
	}()

	go func() {
		defer wg.Done()
		input := model.FullRecordInput{
			Word: "word2todelete",
			Translation: "translation2todelete",
			Examples: []string{"example1", "example2"},
		}
		ok, err := s.DB.AddWord(input)
		s.Require().NoError(err)
		s.Require().True(ok)

		ok, err = s.DB.DeleteWord("word2todelete")
		s.Require().NoError(err)
		s.Require().True(ok)

		_, err = s.DB.ReceiveWordTranslation("word2todelete")
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "Record not found")
	}()

	wg.Wait()
}

// Deleting the same word concurrently
func (s *IntegrationTestSuite) TestDeleteSameWordConcurrently() {
	const concurrencyLevel = 5
	var wg sync.WaitGroup
	var successCount int32
	var errorCount int32

	word := "concurrentdeleteword"

	input := model.FullRecordInput{
		Word: word,
		Translation: "translationtodelete",
		Examples: []string{"example1", "example2"},
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
			success, err := s.DB.DeleteWord(word)
			if err != nil {
					s.Require().NoError(err)
			} else if success {
				atomic.AddInt32(&successCount, 1)
			} else if !success{
				atomic.AddInt32(&errorCount, 1)	
			} 
		}()
	}

	close(startBarrier)

	wg.Wait()

	s.Require().Equal(int32(concurrencyLevel-1), errorCount)

	_, err = s.DB.ReceiveWordTranslation(word)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "Record not found")
}


