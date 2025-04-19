package integration_tests

import (
	"dictionary-app/server/graph/model"
	"sync"
	"sync/atomic"
	"time"
	"fmt"
)

// --------------------------- SEQUENTIAL TESTS ---------------------------

// Adding single word to database
func (s *IntegrationTestSuite) TestAddWord() {
    input := model.FullRecordInput{
        Word: "slowo1",
        Translation: "translation1",
        Examples: []string{"Example number 1", "Example number 2"},
    }
    ok, err := s.DB.AddWord(input)
    s.Require().NoError(err)
    s.True(ok)

    word, err := s.DB.ReceiveWordTranslation("slowo1")
    s.Require().NoError(err)
    s.Equal("slowo1", word.Word)
    s.Len(word.Translations, 1)
    s.Equal("translation1", word.Translations[0].Translation) 
    s.Len(word.Translations[0].ExampleSentences, 2)
}

// Adding two different words sequentially
func (s *IntegrationTestSuite) TestAddMultipleWords() {
    input1 := model.FullRecordInput{
        Word: "slowo2",
        Translation: "translation2",
        Examples: []string{"Example for word 2"},
    }
    ok, err := s.DB.AddWord(input1)
    s.Require().NoError(err)
    s.True(ok)

    input2 := model.FullRecordInput{
        Word: "slowo3",
        Translation: "translation3",
        Examples: []string{"Example for word 3"},
    }
    ok, err = s.DB.AddWord(input2)
    s.Require().NoError(err)
    s.True(ok)

    word1, err := s.DB.ReceiveWordTranslation("slowo2")
    s.Require().NoError(err)
    s.Equal("slowo2", word1.Word)
    s.Equal("translation2", word1.Translations[0].Translation)

    word2, err := s.DB.ReceiveWordTranslation("slowo3")
    s.Require().NoError(err)
    s.Equal("slowo3", word2.Word)
    s.Equal("translation3", word2.Translations[0].Translation)
}

// Adding the same word twice sequentially
func (s *IntegrationTestSuite) TestAddDuplicateWord() {
    input1 := model.FullRecordInput{
        Word: "duplicateWord",
        Translation: "firstTranslation",
        Examples: []string{"First example"},
    }
    ok, err := s.DB.AddWord(input1)
    s.Require().NoError(err)
    s.True(ok)

    input2 := model.FullRecordInput{
        Word: "duplicateWord",
        Translation: "secondTranslation",
        Examples: []string{"Second example"},
    }
    ok, err = s.DB.AddWord(input2)
    s.Require().NoError(err)
    s.True(ok)

    // Verify word has both translations
    word, err := s.DB.ReceiveWordTranslation("duplicateWord")
    s.Require().NoError(err)
    s.Equal("duplicateWord", word.Word)
    s.Len(word.Translations, 2)
    
    translations := []string{
        word.Translations[0].Translation,
        word.Translations[1].Translation,
    }
    s.Contains(translations, "firstTranslation")
    s.Contains(translations, "secondTranslation")
}

// --------------------------- PARALLEL TESTS ---------------------------

// Adding different words in parallel
func (s *IntegrationTestSuite) TestAddDifferentWordsInParallel() {
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        input := model.FullRecordInput{
            Word: "parallelWord1",
            Translation: "parallelTranslation1",
            Examples: []string{"Parallel example 1"},
        }
        ok, err := s.DB.AddWord(input)
        s.Require().NoError(err)
        s.True(ok)
    }()

    go func() {
        defer wg.Done()
        input := model.FullRecordInput{
            Word: "parallelWord2",
            Translation: "parallelTranslation2",
            Examples: []string{"Parallel example 2"},
        }
        ok, err := s.DB.AddWord(input)
        s.Require().NoError(err)
        s.True(ok)
    }()

    wg.Wait()

    word1, err := s.DB.ReceiveWordTranslation("parallelWord1")
    s.Require().NoError(err)
    s.Equal("parallelWord1", word1.Word)

    word2, err := s.DB.ReceiveWordTranslation("parallelWord2")
    s.Require().NoError(err)
    s.Equal("parallelWord2", word2.Word)
}

// Adding the same word with the same translation in parallel
func (s *IntegrationTestSuite) addSameWordSameTranslationConcurrently(word, translation string, examples []string, concurrencyLevel int) {
    var wg sync.WaitGroup
    startBarrier := make(chan struct{})
    var successCount int32
    var errorCount int32

    wg.Add(concurrencyLevel)

    for i := 0; i < concurrencyLevel; i++ {
        go func() {
            defer wg.Done()
            input := model.FullRecordInput{
                Word: word,
                Translation: translation,
                Examples: examples,
            }
            <-startBarrier 
            _, err := s.DB.AddWord(input) 

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

// Repeatedly runs the test for adding the same word with the same translation to improve reliability
func (s *IntegrationTestSuite) TestAddSameWordSameTranslationConcurrentlyRepeated() {
	const runs = 50
	const concurrencyLevel = 5
	const baseWord = "concurrentWord"
	const translation = "concurrentTranslation"
	examples := []string{"Same example 1", "Same example 2"}

	for i := 0; i < runs; i++ {
		testWord := fmt.Sprintf("%s_%d", baseWord, i)
		s.cleanupTestWord(testWord)

		s.Run(fmt.Sprintf("Run_%d", i), func() {
			s.addSameWordSameTranslationConcurrently(testWord, translation, examples, concurrencyLevel)
		})
	}
}

// Adding the same word with the different translation in parallel
func (s *IntegrationTestSuite) addSameWordWithDetailedConcurrencyCheck(word string, concurrencyLevel int) {
    var wg sync.WaitGroup
    startBarrier := make(chan struct{})
    var successCount int32

    expectedTranslations := make([]string, concurrencyLevel)

    for i := 1; i <= concurrencyLevel; i++ {
        translationNum := i
        expectedTranslations[i-1] = fmt.Sprintf("concurrentTranslation%d", translationNum)

        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            input := model.FullRecordInput{
                Word: word,
                Translation: fmt.Sprintf("concurrentTranslation%d", n),
                Examples: []string{fmt.Sprintf("Concurrent example %d", n)},
            }
            <-startBarrier
            ok, err := s.DB.AddWord(input)
            if ok && err == nil {
                atomic.AddInt32(&successCount, 1)
            }
        }(translationNum)
    }

    time.Sleep(100 * time.Millisecond)
    close(startBarrier)
    wg.Wait()

    record, err := s.DB.ReceiveWordTranslation(word)
    s.Require().NoError(err)
    s.Require().Equal(word, record.Word)

    s.Require().Equal(int32(concurrencyLevel), successCount)

    s.Require().Len(record.Translations, concurrencyLevel)

    actualTranslations := make(map[string]bool)
    for _, t := range record.Translations {
        actualTranslations[t.Translation] = true
    }

    for _, expected := range expectedTranslations {
        s.True(actualTranslations[expected], "Missing expected translation: %s", expected)
    }
}

// Repeatedly runs the test for adding the same word with different translation to improve reliability
func (s *IntegrationTestSuite) TestAddSameWordWithHighConcurrencyRepeatedly() {
    const runs = 50
    const concurrencyLevel = 5

    for i := 0; i < runs; i++ {
        testWord := fmt.Sprintf("concurrentTestWord_%d", i)
        s.cleanupTestWord(testWord)

        s.Run(fmt.Sprintf("Run_%d", i), func() {
            s.addSameWordWithDetailedConcurrencyCheck(testWord, concurrencyLevel)
        })
    }
}

// Helper function for deleting words between runs
func (s *IntegrationTestSuite) cleanupTestWord(word string) {
    _, _ = s.DB.DeleteWord(word)
    
    _, err := s.DB.ReceiveWordTranslation(word)
    if err == nil {
        s.Fail("Cleanup failed - word still exists: %s", word)
    }
}
