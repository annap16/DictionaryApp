package integration_tests

import (
	"dictionary-app/server/graph/model"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// --------------------------- SEQUENTIAL TESTS ---------------------------

// Adding a new example to an existing word
func (s *IntegrationTestSuite) TestAddExample() {
	input := model.FullRecordInput{
		Word:        "word",
		Translation: "translation",
		Examples:    []string{"Example1"},
	}
	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	exampleInput := model.FullRecordInput{
		Word:        "word",
		Translation: "translation",
		Examples:    []string{"Example2"},
	}
	ok, err = s.DB.AddExample(exampleInput)
	s.Require().NoError(err)
	s.Require().True(ok)

	wordRecord, err := s.DB.ReceiveWordTranslation("word")
	s.Require().NoError(err)
	s.Equal("word", wordRecord.Word)
	s.Len(wordRecord.Translations, 1)

	s.Len(wordRecord.Translations[0].ExampleSentences, 2)
	s.Equal("Example1", wordRecord.Translations[0].ExampleSentences[0].Sentence)
	s.Equal("Example2", wordRecord.Translations[0].ExampleSentences[1].Sentence)
}

// Adding an example to a non existent word
func (s *IntegrationTestSuite) TestAddExampleNonExistentWord() {
	input := model.FullRecordInput{
		Word:        "word",
		Translation: "translation",
		Examples:    []string{"Example"},
	}
	ok, err := s.DB.AddExample(input)
	s.Require().Error(err)
	s.EqualError(err, "Nie dodano przykładu - nie znaleziono związanego z nim słowa")
	s.False(ok)
}

// Adding an example to a non existent translation
func (s *IntegrationTestSuite) TestAddExampleNonExistentTranslation() {
	input := model.FullRecordInput{
		Word:        "word",
		Translation: "translation",
		Examples:    []string{"Example1"},
	}
	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	input2 := model.FullRecordInput{
		Word:        "word",
		Translation: "translation2",
		Examples:    []string{"Example"},
	}
	ok, err = s.DB.AddExample(input2)
	s.Require().Error(err)
	s.EqualError(err, "Nie dodano tłumaczenia - nie znaleziono związanego z nim tłumaczenia")
	s.False(ok)
}

// Adding a duplicate example
func (s *IntegrationTestSuite) TestAddDuplicateExample() {
	input := model.FullRecordInput{
		Word:        "word",
		Translation: "translation",
		Examples:    []string{"Example"},
	}
	ok, err := s.DB.AddWord(input)
	s.Require().NoError(err)
	s.Require().True(ok)

	ok, err = s.DB.AddExample(input)
	s.Require().False(ok)
	s.Require().Error(err)
	s.Equal("Nie można dodać przykładu - narusza ono unikalność rekordów", err.Error())

	wordRecord, err := s.DB.ReceiveWordTranslation("word")
	s.Require().NoError(err)
	s.Len(wordRecord.Translations[0].ExampleSentences, 1)
	s.Equal("Example", wordRecord.Translations[0].ExampleSentences[0].Sentence)
}

// Adding different examples sequentially
func (s *IntegrationTestSuite) TestAddMultipleExamplesSequentially() {
	word := "word"
	translation := "translation"

	ok, err := s.DB.AddWord(model.FullRecordInput{
		Word:        word,
		Translation: translation,
		Examples:    []string{"Example1"},
	})
	s.Require().NoError(err)
	s.Require().True(ok)

	example1 := model.FullRecordInput{
		Word:        word,
		Translation: translation,
		Examples:    []string{"Example2"},
	}
	ok, err = s.DB.AddExample(example1)
	s.Require().NoError(err)
	s.True(ok)

	example2 := model.FullRecordInput{
		Word:        word,
		Translation: translation,
		Examples:    []string{"Example3"},
	}
	ok, err = s.DB.AddExample(example2)
	s.Require().NoError(err)
	s.True(ok)

	wordRecord, err := s.DB.ReceiveWordTranslation(word)
	s.Require().NoError(err)
	s.Equal(word, wordRecord.Word)

	s.Len(wordRecord.Translations[0].ExampleSentences, 3)
	s.Equal("Example1", wordRecord.Translations[0].ExampleSentences[0].Sentence)
	s.Equal("Example2", wordRecord.Translations[0].ExampleSentences[1].Sentence)
	s.Equal("Example3", wordRecord.Translations[0].ExampleSentences[2].Sentence)

}

// --------------------------- PARALLEL TESTS ---------------------------

// Adding the same examples in parallel
func (s *IntegrationTestSuite) addSameExampleConcurrently(word, translation string, example string, concurrencyLevel int) {
	var wg sync.WaitGroup
	startBarrier := make(chan struct{})
	var successCount int32
	var errorCount int32

	wg.Add(concurrencyLevel)

	for i := 0; i < concurrencyLevel; i++ {
		go func() {
			defer wg.Done()
			<-startBarrier

			input := model.FullRecordInput{
				Word:        word,
				Translation: translation,
				Examples:    []string{example},
			}
			ok, err := s.DB.AddExample(input)
			if err != nil {
				if err.Error() == "Nie można dodać przykładu - narusza ono unikalność rekordów" {
					atomic.AddInt32(&errorCount, 1)
				} else {
					s.Require().NoError(err)
				}
			} else if ok {
				atomic.AddInt32(&successCount, 1)
			}
		}()
	}

	time.Sleep(100 * time.Millisecond)
	close(startBarrier)
	wg.Wait()

	wordRecord, err := s.DB.ReceiveWordTranslation(word)
	s.Require().NoError(err)
	s.Require().Equal(word, wordRecord.Word)
	s.Require().Len(wordRecord.Translations, 1)
	s.Require().Len(wordRecord.Translations[0].ExampleSentences, 1)
	s.Equal("Example", wordRecord.Translations[0].ExampleSentences[0].Sentence)
	s.Require().Equal(int32(concurrencyLevel), successCount+errorCount)
}

// Repeatedly runs the test for adding the same example to improve reliability
func (s *IntegrationTestSuite) TestAddSameExampleConcurrentlyRepeatedly() {
	const runs = 50
	const concurrencyLevel = 5
	const baseWord = "word"
	const translation = "translation"
	example := "Example"

	for i := 0; i < runs; i++ {
		testWord := fmt.Sprintf("%s_%d", baseWord, i)

		ok, err := s.DB.AddWord(model.FullRecordInput{
			Word:        testWord,
			Translation: translation,
			Examples:    []string{example},
		})
		s.Require().NoError(err)
		s.Require().True(ok)

		s.Run(fmt.Sprintf("Run_%d", i), func() {
			s.addSameExampleConcurrently(testWord, translation, example, concurrencyLevel)
		})
	}
}

// Adding different examples in parallel
func (s *IntegrationTestSuite) addDifferentExamplesConcurrently(word, translation string, concurrencyLevel int) {
	var wg sync.WaitGroup
	startBarrier := make(chan struct{})
	var successCount int32
	expectedExamples := make([]string, concurrencyLevel)

	for i := 0; i < concurrencyLevel; i++ {
		ex := fmt.Sprintf("example_%d", i)
		expectedExamples[i] = ex

		wg.Add(1)
		go func(example string) {
			defer wg.Done()
			<-startBarrier

			input := model.FullRecordInput{
				Word:        word,
				Translation: translation,
				Examples:    []string{example},
			}
			ok, err := s.DB.AddExample(input)
			if ok && err == nil {
				atomic.AddInt32(&successCount, 1)
			} else if err != nil {
				s.T().Logf("Error adding example %q: %v", example, err)
			}
		}(ex)
	}

	time.Sleep(100 * time.Millisecond)
	close(startBarrier)
	wg.Wait()

	wordRecord, err := s.DB.ReceiveWordTranslation(word)
	s.Require().NoError(err)
	s.Require().Equal(word, wordRecord.Word)
	s.Require().Len(wordRecord.Translations, 1)

	exampleSet := make(map[string]bool)
	for _, ex := range wordRecord.Translations[0].ExampleSentences {
		exampleSet[ex.Sentence] = true
	}

	for _, expected := range expectedExamples {
		s.True(exampleSet[expected], "Missing example: %s", expected)
	}

	s.Require().Equal(int32(concurrencyLevel), successCount)
}

// Repeatedly runs the test for adding different exmaples to improve reliability
func (s *IntegrationTestSuite) TestAddDifferentExamplesConcurrentlyRepeatedly() {
	const runs = 50
	const concurrencyLevel = 5
	const baseWord = "word"
	const translation = "translation"

	for i := 0; i < runs; i++ {
		testWord := fmt.Sprintf("%s_diff_%d", baseWord, i)

		ok, err := s.DB.AddWord(model.FullRecordInput{
			Word:        testWord,
			Translation: translation,
			Examples:    []string{"base"},
		})
		s.Require().NoError(err)
		s.Require().True(ok)

		s.Run(fmt.Sprintf("Run_%d", i), func() {
			s.addDifferentExamplesConcurrently(testWord, translation, concurrencyLevel)
		})
	}
}
