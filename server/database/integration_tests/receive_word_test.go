package integration_tests

import (
	"dictionary-app/server/graph/model"
	"sync"
)

// Receiving single word from database
func (s *IntegrationTestSuite) TestReceiveWordTranslation_Success() {
	input := model.FullRecordInput{
		Word:        "existingword",
		Translation: "existingtranslation",
		Examples:    []string{"Example 1", "Example 2"},
	}

	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	word, err := s.DB.ReceiveWordTranslation("existingword")
	s.Require().NoError(err)
	s.Require().Equal("existingword", word.Word)
	s.Require().Len(word.Translations, 1)
	s.Require().Equal("existingtranslation", word.Translations[0].Translation)
	s.Require().Len(word.Translations[0].ExampleSentences, 2)
}

// Receiving non existent word from database
func (s *IntegrationTestSuite) TestReceiveWordTranslation_NotFound() {
	_, err := s.DB.ReceiveWordTranslation("nonexistentword")
	s.Require().Error(err)
	s.EqualError(err, "Nie znaleziono podanego słowa w słowniku")
}

// Receiving word after deleting it from database
func (s *IntegrationTestSuite) TestReceiveWordTranslation_AfterDeletion() {
	input := model.FullRecordInput{
		Word:        "tobedeleted",
		Translation: "temptranslation",
		Examples:    []string{"temp example"},
	}

	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	_, err = s.DB.DeleteWord("tobedeleted")
	s.Require().NoError(err)

	_, err = s.DB.ReceiveWordTranslation("tobedeleted")
	s.Require().Error(err)
	s.EqualError(err, "Nie znaleziono podanego słowa w słowniku")
}

// Receiving the same word many times sequentially
func (s *IntegrationTestSuite) TestReceiveWordTranslation_RepeatedCalls() {
	input := model.FullRecordInput{
		Word:        "repeatword",
		Translation: "repeattranslation",
		Examples:    []string{"example1"},
	}

	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	for i := 0; i < 10; i++ {
		word, err := s.DB.ReceiveWordTranslation("repeatword")
		s.Require().NoError(err)
		s.Require().Equal("repeatword", word.Word)
		s.Require().Equal("repeattranslation", word.Translations[0].Translation)
	}
}

// Receving the same word many times in parallel
func (s *IntegrationTestSuite) TestReceiveWordTranslation_ConcurrentReads() {
	input := model.FullRecordInput{
		Word:        "concurrentreadword",
		Translation: "readtranslation",
		Examples:    []string{"Example A", "Example B"},
	}

	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	const readers = 10
	var wg sync.WaitGroup
	wg.Add(readers)

	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			word, err := s.DB.ReceiveWordTranslation("concurrentreadword")
			s.Require().NoError(err)
			s.Require().Equal("concurrentreadword", word.Word)
			s.Require().Equal("readtranslation", word.Translations[0].Translation)
			s.Require().Len(word.Translations[0].ExampleSentences, 2)
		}()
	}

	wg.Wait()
}
