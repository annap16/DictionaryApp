package integration_tests

import (
	"dictionary-app/server/graph/model"
	"sync"
	"sync/atomic"
	"time"
	"fmt"
)

// --------------------------- SEQUENTIAL TESTS ---------------------------

// Adding a new translation to an existing word
func (s *IntegrationTestSuite) TestAddTranslation() {
    wordInput := model.FullRecordInput{
        Word: "word",
        Translation: "translation",
        Examples: []string{"Example1", "Example2"},
    }
    ok, err := s.DB.AddWord(wordInput)
    s.Require().NoError(err)
    s.True(ok)

    translationInput := model.FullRecordInput{
        Word: "word",
        Translation: "translation2",
        Examples: []string{"Example3", "Example4"},
    }

    ok, err = s.DB.AddTranslation(translationInput)
    s.Require().NoError(err)
    s.True(ok)

    word, err := s.DB.ReceiveWordTranslation("word")
    s.Require().NoError(err)
    s.Equal("word", word.Word)
    s.Len(word.Translations, 2)
    
    s.Equal(word.Translations[0].Translation, "translation")
    s.Equal(word.Translations[1].Translation, "translation2")
}

// Adding a translation to a non-existent word
func (s *IntegrationTestSuite) TestAddTranslationNonExistentWord() {
    input := model.FullRecordInput{
        Word: "word",
        Translation: "translation",
        Examples: []string{"Example"},
    }
    ok, err := s.DB.AddTranslation(input)
    s.Require().Error(err)
    s.EqualError(err, "Nie dodano tłumaczenia - nie znaleziono związanego z nim słowa")
    s.False(ok)
}

// Adding a duplicate translation
func (s *IntegrationTestSuite) TestAddDuplicateTranslation() {
    input := model.FullRecordInput{
        Word: "word",
        Translation: "translation",
        Examples: []string{"Example"},
    }
    ok, err := s.DB.AddWord(input)
    s.Require().NoError(err)
    s.True(ok)

    ok, err = s.DB.AddTranslation(input)
    s.Error(err)
    s.Equal("Nie można dodać tłumaczenia – narusza ono unikalność rekordów", err.Error())
    s.False(ok)

    word, err := s.DB.ReceiveWordTranslation("word")
    s.Require().NoError(err)
    s.Equal("word", word.Word)
    s.Len(word.Translations, 1)
}

// Adding two different translation sequentially
func (s *IntegrationTestSuite) TestAddMultipleTranslationsSequentially() {
	word := "word"

	ok, err := s.DB.AddWord(model.FullRecordInput{
		Word: word,
		Translation: "translation1",
		Examples: []string{"Example1"},
	})
	s.Require().NoError(err)
	s.Require().True(ok)

	input1 := model.FullRecordInput{
		Word:  word,
		Translation: "translation2",
		Examples: []string{"Example2"},
	}
	ok, err = s.DB.AddTranslation(input1)
	s.Require().NoError(err)
	s.True(ok)

	input2 := model.FullRecordInput{
		Word: word,
		Translation: "translation3",
		Examples: []string{"Example3"},
	}
	ok, err = s.DB.AddTranslation(input2)
	s.Require().NoError(err)
	s.True(ok)

	wordRecord, err := s.DB.ReceiveWordTranslation(word)
	s.Require().NoError(err)
	s.Equal(word, wordRecord.Word)
	s.Len(wordRecord.Translations, 3) 

    s.Equal("translation1", wordRecord.Translations[0].Translation)
	s.Equal("translation2", wordRecord.Translations[1].Translation)
	s.Equal("translation3", wordRecord.Translations[2].Translation)
}

// Adding the same translation twice sequentially
func (s *IntegrationTestSuite) TestAddDuplicateTranslationSequentially() {
	word := "word"
	translation := "translation2"
	examples1 := []string{"Example1"}
	examples2 := []string{"Example2"}

	ok, err := s.DB.AddWord(model.FullRecordInput{
		Word: word,
		Translation: "translation",
		Examples: []string{"Example"},
	})
	s.Require().NoError(err)
	s.Require().True(ok)

	ok, err = s.DB.AddTranslation(model.FullRecordInput{
		Word: word,
		Translation: translation,
		Examples: examples1,
	})
	s.Require().NoError(err)
	s.True(ok)

	ok, err = s.DB.AddTranslation(model.FullRecordInput{
		Word: word,
		Translation: translation,
		Examples: examples2,
	})
	s.Require().False(ok)
	s.Require().Error(err)
	s.Contains(err.Error(), "narusza ono unikalność rekordów")

	wordRecord, err := s.DB.ReceiveWordTranslation(word)
	s.Require().NoError(err)

	var matchCount int
	for _, t := range wordRecord.Translations {
		if t.Translation == translation {
			matchCount++
		}
	}
	s.Equal(1, matchCount, "Expected only one instance of the duplicate translation")
}

// --------------------------- PARALLEL TESTS ---------------------------

// Adding the same translations in parallel
func (s *IntegrationTestSuite) addSameTranslationConcurrently(word, translation string, examples []string, concurrencyLevel int) {
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
                Word: word,
                Translation: translation,
                Examples: examples,
            }
            _, err := s.DB.AddTranslation(input)
            if err != nil {
                if err.Error() == "Nie można dodać tłumaczenia – narusza ono unikalność rekordów" {
                    atomic.AddInt32(&errorCount, 1)
                } else {
                    s.Require().NoError(err)
                }
            } else {
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
    s.Equal(translation, wordRecord.Translations[0].Translation)
    s.Require().Len(wordRecord.Translations[0].ExampleSentences, 2)
    s.Require().Equal(int32(concurrencyLevel), successCount+errorCount)
}

// Repeatedly runs the test for adding the same translation to improve reliability
func (s *IntegrationTestSuite) TestAddSameTranslationConcurrentlyRepeatedly() {
    const runs = 50
    const concurrencyLevel = 5
    const baseWord = "word"
    const translation = "translation"
    examples := []string{"Example"}

    for i := 0; i < runs; i++ {
        testWord := fmt.Sprintf("%s_%d", baseWord, i)

        ok, err := s.DB.AddWord(model.FullRecordInput{
            Word: testWord,
            Translation: translation,
            Examples: []string{"Example1", "Example2"},
        })
        s.Require().NoError(err)
        s.Require().True(ok)

        s.Run(fmt.Sprintf("Run_%d", i), func() {
            s.addSameTranslationConcurrently(testWord, translation, examples, concurrencyLevel)
        })
    }
}

// Adding different translations in parallel
func (s *IntegrationTestSuite) addDifferentTranslationsConcurrently(word string, concurrencyLevel int) {
    var wg sync.WaitGroup
    startBarrier := make(chan struct{})
    var successCount int32
    expectedTranslations := make([]string, concurrencyLevel)

    for i := 0; i < concurrencyLevel; i++ {
        translation := fmt.Sprintf("translation_%d", i)
        expectedTranslations[i] = translation

        wg.Add(1)
        go func(tr string) {
            defer wg.Done()
            <-startBarrier
            input := model.FullRecordInput{
                Word: word,
                Translation: tr,
                Examples: []string{"Example for " + tr},
            }
            ok, err := s.DB.AddTranslation(input)
            if ok && err == nil {
                atomic.AddInt32(&successCount, 1)
            }
        }(translation)
    }

    time.Sleep(100 * time.Millisecond)
    close(startBarrier)
    wg.Wait()

    wordRecord, err := s.DB.ReceiveWordTranslation(word)
    s.Require().NoError(err)
    s.Require().Equal(word, wordRecord.Word)
    s.Require().Len(wordRecord.Translations, concurrencyLevel+1)

    actualSet := make(map[string]bool)
    for _, t := range wordRecord.Translations {
        actualSet[t.Translation] = true
    }

    for _, expected := range expectedTranslations {
        s.True(actualSet[expected], "Missing translation: %s", expected)
    }

    s.Require().Equal(int32(concurrencyLevel), successCount)
}

// Repeatedly runs the test for adding different translations to improve reliability
func (s *IntegrationTestSuite) TestAddDifferentTranslationsConcurrentlyRepeatedly() {
    const runs = 50
    const concurrencyLevel = 5
    const baseWord = "word"

    for i := 0; i < runs; i++ {
        testWord := fmt.Sprintf("%s_%d", baseWord, i)

        ok, err := s.DB.AddWord(model.FullRecordInput{
            Word: testWord,
            Translation: "translation",
            Examples: []string{"Example"},
        })
        s.Require().NoError(err)
        s.Require().True(ok)

        s.Run(fmt.Sprintf("Run_%d", i), func() {
            s.addDifferentTranslationsConcurrently(testWord, concurrencyLevel)
        })
    }
}

