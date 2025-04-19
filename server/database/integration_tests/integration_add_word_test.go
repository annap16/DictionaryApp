 package integration_tests

import (
	"dictionary-app/server/graph/model"
)



func (s *IntegrationTestSuite) TestAddWord() {
	input := model.FullRecordInput{
		Word:        "apple",
		Translation: "jabłko",
		Examples:    []string{"An apple a day keeps the doctor away"},
	}

	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.True(ok)

	word, err := s.DB.ReceiveWordTranslation("apple")
	s.Require().NoError(err)
	s.Equal("apple", word.Word)
	s.Len(word.Translations, 1)
	s.Equal("jabłko", word.Translations[0].Translation)
	s.Len(word.Translations[0].ExampleSentences, 1)
}
