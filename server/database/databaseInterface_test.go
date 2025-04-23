package database

import (
	"testing"
	"os"
	"bytes"
	"dictionary-app/server/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"errors"	
	customerrors "dictionary-app/server/errors"
)

func captureStdout(f func()) string {
	originalStdout := os.Stdout
	readPipe, writePipe, _ := os.Pipe()
	defer readPipe.Close()

	os.Stdout = writePipe
	defer func() { os.Stdout = originalStdout }()

	f()

	writePipe.Close() 
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(readPipe)
	return buf.String()
}


func TestAddWord(t *testing.T) {
	tests := []struct {
		name string
		input model.FullRecordInput
		mockCreateWord func(mockRepo *MockRepository, mockTx *gorm.DB)
		expectedResult bool
		expectedError error
		expectedOutput string
	}{
		{
			name: "Success - Word Added",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockCreateWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateWord", mock.Anything, mockTx).Return(nil)
			},
			expectedResult: true,
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Error - General Database Error",
			input: model.FullRecordInput{
				Word: "samochod",
				Translation: "car",
				Examples: []string{"A fast car", "My car is blue"},
			},
			mockCreateWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateWord", mock.Anything, mockTx).Return(errors.New("DB internal error"))
			},
			expectedResult: false,
			expectedError: errors.New("DB internal error"),
			expectedOutput: "Error while adding a word to a DB DB internal error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockDB := new(MockDatabase)
			dbInterface := &DBInterface{
				DB: mockDB,
				repo: mockRepo,
			}
			mockTx := &gorm.DB{} 

			var capturedSuccess bool
			var capturedError error
			mockRepo.On("TransactionWrapper", mockDB, mock.AnythingOfType("func(*gorm.DB) (bool, error)")).
				Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) (bool, error))
					capturedSuccess, capturedError = fn(mockTx)
				}).
				Return(func(mock.Arguments) bool {
					return capturedSuccess
				}, func(mock.Arguments) error {
					return capturedError
				})

			tt.mockCreateWord(mockRepo, mockTx)

			var success bool
			var err error

			output := captureStdout(func() {
				success, err = dbInterface.AddWord(tt.input)

			})

			assert.Equal(t, tt.expectedOutput, output)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else if tt.expectedResult{
				assert.NoError(t, err)
				assert.True(t, success)
			}else{
				assert.NoError(t, err)
				assert.False(t, success)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestReceiveWord(t *testing.T){
	tests := []struct {
		name string
		word string
		mockGetWord func(mockRepo *MockRepository, mockTx *gorm.DB)
		expectedResult *model.Word
		expectedError error
		expectedOutput string
	}{
		{
			name: "Success - Word Received",
			word: "ksiazka",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 1
					word.Word = "ksiazka"
					word.Translations = []Translation{
						{ID: 1, Translation: "book", WordID: 1},
					}
				}).Return(nil)
			},
			expectedResult: &model.Word{
				ID: "1", 
				Word: "ksiazka",
				Translations: []*model.Translation{
					{ID: "1", Translation: "book"},
				},
			},
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Word Not Found",
			word: "ksiazka",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 0
				}).Return(nil)
			},
			expectedResult: nil,
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Word",
			word: "ksiazka",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 0
				}).Return(errors.New("Word error"))
			},
			expectedResult: nil,
			expectedError: errors.New("Word error"),
			expectedOutput: "Error while loading word from a DB: Word error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockDB := new(MockDatabase)
			dbInterface := &DBInterface{
				DB: mockDB,
				repo: mockRepo,
			}
			mockTx := &gorm.DB{} 

			var capturedSuccess bool
			var capturedError error
			mockRepo.On("TransactionWrapper", mockDB, mock.AnythingOfType("func(*gorm.DB) (bool, error)")).
				Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) (bool, error))
					capturedSuccess, capturedError = fn(mockTx)
				}).
				Return(func(mock.Arguments) bool {
					return capturedSuccess
				}, func(mock.Arguments) error {
					return capturedError
				})

			tt.mockGetWord(mockRepo, mockTx)

			var result *model.Word
			var err error

			output := captureStdout(func() {
				result, err = dbInterface.ReceiveWordTranslation(tt.word)
			})

			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectedResult, result)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteWord(t *testing.T){
	tests := []struct {
		name string
		input string
		mockDeleteWord func(mockRepo *MockRepository, mockTx *gorm.DB)
		expectedResult bool
		expectedError error
	}{
		{
			name: "Success - Word Deleted",
			input: "ksiazka",
			mockDeleteWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("DeleteWord", mock.Anything, mockTx).Return(true, nil)
			},
			expectedResult: true,
			expectedError: nil,
		},
		{
			name: "Failure - Word Not Found",
			input: "ksiazka",
			mockDeleteWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("DeleteWord", mock.Anything, mockTx).Return(false, nil)
			},
			expectedResult: false,
			expectedError: nil,
		},
		{
			name: "Failure - General Database Error",
			input: "ksiazka",
			mockDeleteWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("DeleteWord", mock.Anything, mockTx).Return(false, errors.New("DB internal error"))
			},
			expectedResult: false,
			expectedError: errors.New("DB internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockDB := new(MockDatabase)
			dbInterface := &DBInterface{
				DB: mockDB,
				repo: mockRepo,
			}
			mockTx := &gorm.DB{} 

			var capturedSuccess bool
			var capturedError error
			mockRepo.On("TransactionWrapper", mockDB, mock.AnythingOfType("func(*gorm.DB) (bool, error)")).
				Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) (bool, error))
					capturedSuccess, capturedError = fn(mockTx)
				}).
				Return(func(mock.Arguments) bool {
					return capturedSuccess
				}, func(mock.Arguments) error {
					return capturedError
				})

			tt.mockDeleteWord(mockRepo, mockTx)

			success, err := dbInterface.DeleteWord(tt.input)


			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else if tt.expectedResult{
				assert.NoError(t, err)
				assert.True(t, success)
			}else{
				assert.NoError(t, err)
				assert.False(t, success)
			}

			mockRepo.AssertExpectations(t)
		})
	}

}

func TestDeleteTranslation(t *testing.T){
	tests := []struct {
		name string
		word string
		translation string
		mockDeleteTranslation func(mockRepo *MockRepository, mockTx *gorm.DB)
		mockGetWord func(mockRepo *MockRepository, mockTx *gorm.DB)
		expectedResult bool
		expectedError error
		expectedOutput string
	}{
		{
			name: "Success - Translation Deleted",
			word: "ksiazka",
			translation: "book",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 1
				}).Return(nil)
			},
			mockDeleteTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("DeleteTranslation", uint(1), "book", mockTx).Return(true, nil)
			},
			expectedResult: true,
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Translation Not Found",
			word: "ksiazka",
			translation: "book",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 1
				}).Return(nil)
			},
			mockDeleteTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("DeleteTranslation", uint(1), "book", mockTx).Return(false, customerrors.NewNotFoundError("Nie usunięto podanego tłumaczenia - tłumaczenia nie znaleziono"))
			},
			expectedResult: false,
			expectedError: customerrors.NewNotFoundError("Nie usunięto podanego tłumaczenia - tłumaczenia nie znaleziono"),
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Translation",
			word: "ksiazka",
			translation: "book",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 1
				}).Return(nil)
			},
			mockDeleteTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("DeleteTranslation", uint(1), "book", mockTx).Return(false, errors.New("Translation error"))
			},
			expectedResult: false,
			expectedError: errors.New("Translation error"),
			expectedOutput: "",
		},
		{
			name: "Failure - Word Not Found",
			word: "ksiazka",
			translation: "book",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 0
				}).Return(nil)
			},
			mockDeleteTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			expectedResult: false,
			expectedError: customerrors.NewNotFoundError("Nie usunięto tłumaczenia - nie znaleziono związanego z nim słowa"),
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Word",
			word: "ksiazka",
			translation: "book",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 0
				}).Return(errors.New("Word error"))
			},
			mockDeleteTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			expectedResult: false,
			expectedError: errors.New("Word error"),
			expectedOutput: "Error while searching for word existance in a DB: Word error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockDB := new(MockDatabase)
			dbInterface := &DBInterface{
				DB: mockDB,
				repo: mockRepo,
			}
			mockTx := &gorm.DB{} 

			var capturedSuccess bool
			var capturedError error
			mockRepo.On("TransactionWrapper", mockDB, mock.AnythingOfType("func(*gorm.DB) (bool, error)")).
				Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) (bool, error))
					capturedSuccess, capturedError = fn(mockTx)
				}).
				Return(func(mock.Arguments) bool {
					return capturedSuccess
				}, func(mock.Arguments) error {
					return capturedError
				})

			tt.mockDeleteTranslation(mockRepo, mockTx)
			tt.mockGetWord(mockRepo, mockTx)

			var success bool
			var err error

			output := captureStdout(func() {
				success, err = dbInterface.DeleteTranslation(tt.word, tt.translation)
			})

			assert.Equal(t, tt.expectedOutput, output)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else if tt.expectedResult{
				assert.NoError(t, err)
				assert.True(t, success)
			}else{
				assert.NoError(t, err)
				assert.False(t, success)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteExample(t *testing.T){
	tests := []struct {
		name string
		input model.FullRecordInput
		mockDeleteExample func(mockRepo *MockRepository, mockTx *gorm.DB)
		mockGetWord func(mockRepo *MockRepository, mockTx *gorm.DB)
		mockGetTranslation func(mockRepo *MockRepository, mockTx *gorm.DB)
		expectedResult bool
		expectedError error
		expectedOutput string
	}{
		{
			name: "Success - Example Deleted",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockDeleteExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("DeleteExample", uint(1), "I love my new book", mockTx).Return(true, nil)
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 1
				}).Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 1
					}).
					Return(nil)
			},			
			expectedResult: true,
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Example Not Found",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockDeleteExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("DeleteExample", uint(1), "I love my new book", mockTx).Return(false, nil)
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 1
				}).Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 1
					}).
					Return(nil)
			},			
			expectedResult: false,
			expectedError: errors.New("Nie można było usunąć przykładu ze słownika. Anulowano całą operację"),
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Example",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockDeleteExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("DeleteExample", uint(1), "I love my new book", mockTx).Return(false, errors.New("Example error"))
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 1
				}).Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 1
					}).
					Return(nil)
			},			
			expectedResult: false,
			expectedError: errors.New("Nie można było usunąć przykładu ze słownika. Anulowano całą operację"),
			expectedOutput: "",
		},
		{
			name: "Failure - Translation Not Found",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockDeleteExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {

			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 1
				}).Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 0
					}).
					Return(nil)
			},			
			expectedResult: false,
			expectedError: customerrors.NewNotFoundError("Nie usunięto przykładu - nie znaleziono związanego z nim tłumaczenia"),
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Translation",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockDeleteExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {

			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 1
				}).Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 0
					}).
					Return(errors.New("Translation Error"))
			},			
			expectedResult: false,
			expectedError: errors.New("Translation Error"),
			expectedOutput: "Error while searching for translation existance in a DB: Translation Error\n",
		},
		{
			name: "Failure - Word Not Found",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockDeleteExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {

			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 0
				}).Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},			
			expectedResult: false,
			expectedError:  customerrors.NewNotFoundError("Nie usunięto przykładu - nie znaleziono związanego z nim słowa"),
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Word",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockDeleteExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {

			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).Run(func(args mock.Arguments) {
					word := args.Get(1).(*Word)
					word.ID = 0
				}).Return(errors.New("Word error"))
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},			
			expectedResult: false,
			expectedError: errors.New("Word error"),
			expectedOutput: "Error while searching for word existance in a DB: Word error\n",
		},
		

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockDB := new(MockDatabase)
			dbInterface := &DBInterface{
				DB: mockDB,
				repo: mockRepo,
			}
			mockTx := &gorm.DB{} 

			var capturedSuccess bool
			var capturedError error
			mockRepo.On("TransactionWrapper", mockDB, mock.AnythingOfType("func(*gorm.DB) (bool, error)")).
				Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) (bool, error))
					capturedSuccess, capturedError = fn(mockTx)
				}).
				Return(func(mock.Arguments) bool {
					return capturedSuccess
				}, func(mock.Arguments) error {
					return capturedError
				})

			tt.mockDeleteExample(mockRepo, mockTx)
			tt.mockGetWord(mockRepo, mockTx)
			tt.mockGetTranslation(mockRepo, mockTx)

			var success bool
			var err error

			output := captureStdout(func() {
				success, err = dbInterface.DeleteExample(tt.input)
			})

			assert.Equal(t, tt.expectedOutput, output)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else if tt.expectedResult{
				assert.NoError(t, err)
				assert.True(t, success)
			}else{
				assert.NoError(t, err)
				assert.False(t, success)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAddTranslation(t *testing.T){
	tests := []struct {
		name string
		input model.FullRecordInput
		mockCreateTranslation func(mockRepo *MockRepository, mockTx *gorm.DB)
		mockGetWord func(mockRepo *MockRepository, mockTx *gorm.DB)
		expectedResult bool
		expectedError error
		expectedOutput string
	}{
		{
			name: "Success - Translation Created",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1 
					}).
					Return(nil)
			},
			mockCreateTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateTranslation", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(0).(*Translation)
						translation.ID = 1 
					}).
					Return(nil)
			},
			expectedResult: true,
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Duplicated Translation",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1 
					}).
					Return(nil)
			},
			mockCreateTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateTranslation", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(0).(*Translation)
						translation.ID = 0
					}).
					Return(customerrors.NewDuplicateKeyError("Nie można dodać przykładu – narusza on unikalność rekordów"))
			},
			expectedResult: false,
			expectedError: errors.New("Nie można dodać tłumaczenia – narusza ono unikalność rekordów"),
			expectedOutput: "",
		},
		{
			name: "Failure - Foreign Key Violation - Race Conditions",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1 
					}).
					Return(nil)
			},
			mockCreateTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateTranslation", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(0).(*Translation)
						translation.ID = 0
					}).
					Return(customerrors.NewForeignKeyError("Nie można dodać przykładu – powiązane tłumaczenie nie istnieje"))
			},
			expectedResult: false,
			expectedError: errors.New("Wystąpił błąd podczas dodawania tłumaczenia"),
			expectedOutput: "",
		},
		{
			name: "Failure - Random Database Error With Translation",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1 
					}).
					Return(nil)
			},
			mockCreateTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateTranslation", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(0).(*Translation)
						translation.ID = 0
					}).
					Return(errors.New("Translation Error"))
			},
			expectedResult: false,
			expectedError: errors.New("Translation Error"),
			expectedOutput: "",
		},
		{
			name: "Failure - Word Not Found",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 0
					}).
					Return(nil)
			},
			mockCreateTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			expectedResult: false,
			expectedError: customerrors.NewNotFoundError("Nie dodano tłumaczenia - nie znaleziono związanego z nim słowa"),
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Word",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 0
					}).
					Return(errors.New("Word error"))
			},
			mockCreateTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			expectedResult: false,
			expectedError: errors.New("Word error"),
			expectedOutput: "Error while searching for word existance in a DB: Word error\n",
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockDB := new(MockDatabase)
			dbInterface := &DBInterface{
				DB: mockDB,
				repo: mockRepo,
			}
			mockTx := &gorm.DB{} 

			var capturedSuccess bool
			var capturedError error
			mockRepo.On("TransactionWrapper", mockDB, mock.AnythingOfType("func(*gorm.DB) (bool, error)")).
				Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) (bool, error))
					capturedSuccess, capturedError = fn(mockTx)
				}).
				Return(func(mock.Arguments) bool {
					return capturedSuccess
				}, func(mock.Arguments) error {
					return capturedError
				})

			tt.mockCreateTranslation(mockRepo, mockTx)
			tt.mockGetWord(mockRepo, mockTx)

			var success bool
			var err error

			output := captureStdout(func() {
				success, err = dbInterface.AddTranslation(tt.input)
			})

			assert.Equal(t, tt.expectedOutput, output)

			
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else if tt.expectedResult{
				assert.NoError(t, err)
				assert.True(t, success)
			}else{
				assert.NoError(t, err)
				assert.False(t, success)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAddExample(t *testing.T){
	tests := []struct {
		name string
		input model.FullRecordInput
		mockCreateExample func(mockRepo *MockRepository, mockTx *gorm.DB)
		mockGetWord func(mockRepo *MockRepository, mockTx *gorm.DB)
		mockGetTranslation func(mockRepo *MockRepository, mockTx *gorm.DB)
		expectedResult bool
		expectedError error
		expectedOutput string
	}{
		{
			name: "Success - Example Added",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockCreateExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateExample",  mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						example := args.Get(0).(*ExampleSentence)
						example.ID = 1 
					}).
					Return(nil)
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1
					}).
					Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", uint(1), "book", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 1 
					}).
					Return(nil)
			},
			expectedResult: true,
			expectedError:  nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Example Duplicate",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockCreateExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateExample",  mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						example := args.Get(0).(*ExampleSentence)
						example.ID = 0 
					}).
					Return(customerrors.NewDuplicateKeyError("Nie można dodać przykładu – narusza on unikalność rekordów"))
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1
					}).
					Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", uint(1), "book", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 1 
					}).
					Return(nil)
			},
			expectedResult: false,
			expectedError:  errors.New("Nie można dodać przykładu - narusza ono unikalność rekordów"),
			expectedOutput: "",
		},
		{
			name: "Failure - Foriegn Key Violation - Race Conditions",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockCreateExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateExample",  mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						example := args.Get(0).(*ExampleSentence)
						example.ID = 0 
					}).
					Return(customerrors.NewForeignKeyError("Nie można dodać przykładu – powiązane tłumaczenie nie istnieje"))
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1
					}).
					Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", uint(1), "book", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 1 
					}).
					Return(nil)
			},
			expectedResult: false,
			expectedError:  errors.New("Wystąpił błąd podczas dodawania przykładu"),
			expectedOutput: "",
		},
		{
			name: "Failure - Random Error Database Example",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockCreateExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("CreateExample",  mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						example := args.Get(0).(*ExampleSentence)
						example.ID = 0 
					}).
					Return(errors.New("Example error"))
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1
					}).
					Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", uint(1), "book", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 1 
					}).
					Return(nil)
			},
			expectedResult: false,
			expectedError:  errors.New("Example error"),
			expectedOutput: "",
		},
		{
			name: "Failure - Translation Not Found",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockCreateExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1
					}).
					Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", uint(1), "book", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 0
					}).
					Return(nil)
			},
			expectedResult: false,
			expectedError:  customerrors.NewNotFoundError("Nie dodano tłumaczenia - nie znaleziono związanego z nim tłumaczenia"),
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Translation",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockCreateExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1
					}).
					Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", uint(1), "book", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 0
					}).
					Return(errors.New("Translation Error"))
			},
			expectedResult: false,
			expectedError:  errors.New("Translation Error"),
			expectedOutput: "Error while searching for translation existance in a DB: Translation Error\n",
		},
		{
			name: "Failure - Word Not Found",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockCreateExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 0
					}).
					Return(nil)
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			expectedResult: false,
			expectedError: customerrors.NewNotFoundError("Nie dodano przykładu - nie znaleziono związanego z nim słowa"),
			expectedOutput: "",
		},
		{
			name: "Failure - Word Not Found Error",
			input: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			mockCreateExample: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 0
					}).
					Return(errors.New("Word Error"))
			},
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
			},
			expectedResult: false,
			expectedError: errors.New("Word Error"),
			expectedOutput: "Error while searching for word existance in a DB: Word Error\n",
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockDB := new(MockDatabase)
			dbInterface := &DBInterface{
				DB: mockDB,
				repo: mockRepo,
			}
			mockTx := &gorm.DB{} 

			var capturedSuccess bool
			var capturedError error
			mockRepo.On("TransactionWrapper", mockDB, mock.AnythingOfType("func(*gorm.DB) (bool, error)")).
				Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) (bool, error))
					capturedSuccess, capturedError = fn(mockTx)
				}).
				Return(func(mock.Arguments) bool {
					return capturedSuccess
				}, func(mock.Arguments) error {
					return capturedError
				})

			tt.mockCreateExample(mockRepo, mockTx)
			tt.mockGetWord(mockRepo, mockTx)
			tt.mockGetTranslation(mockRepo, mockTx)

			var success bool
			var err error

			output := captureStdout(func() {
				success, err = dbInterface.AddExample(tt.input)
			})

			assert.Equal(t, tt.expectedOutput, output)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else if tt.expectedResult{
				assert.NoError(t, err)
				assert.True(t, success)
			}else{
				assert.NoError(t, err)
				assert.False(t, success)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetWordID(t *testing.T){
	tests := []struct {
		name string
		word string
		mockGetWord func(mockRepo *MockRepository, mockTx *gorm.DB)
		expectedResult uint
		expectedError error
		expectedOutput string
	}{
		{
			name: "Success - Got Word ID",
			word: "ksiazka",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 1
					}).
					Return(nil)
			},
			expectedResult: 1,
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Word Not Found",
			word: "ksiazka",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 0
					}).
					Return(nil)
			},
			expectedResult: 0,
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Word",
			word: "ksiazka",
			mockGetWord: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetWord", "ksiazka", mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						word := args.Get(1).(*Word)
						word.ID = 0
					}).
					Return(errors.New("Word Error"))
			},
			expectedResult: 0,
			expectedError: errors.New("Word Error"),
			expectedOutput: "Error while searching for word existance in a DB: Word Error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockDB := new(MockDatabase)
			dbInterface := &DBInterface{
				DB: mockDB,
				repo: mockRepo,
			}
			mockTx := &gorm.DB{} 
			
			tt.mockGetWord(mockRepo, mockTx)

			var result uint
			var err error

			output := captureStdout(func() {
				result, err = dbInterface.GetWordID(mockTx, tt.word)
			})

			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectedResult, result)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} 

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetTranslationID(t *testing.T){
	tests := []struct {
		name string
		wordID uint
		translation string
		mockGetTranslation func(mockRepo *MockRepository, mockTx *gorm.DB)
		expectedResult uint
		expectedError error
		expectedOutput string
	}{
		{
			name: "Success - Got Translation ID",
			wordID: 1,
			translation: "book",
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 1
					}).
					Return(nil)
			},
			expectedResult: 1,
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Translation Not Found",
			wordID: 1,
			translation: "book",
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 0
					}).
					Return(nil)
			},
			expectedResult: 0,
			expectedError: nil,
			expectedOutput: "",
		},
		{
			name: "Failure - Database Error With Translation",
			wordID: 1,
			translation: "book",
			mockGetTranslation: func(mockRepo *MockRepository, mockTx *gorm.DB) {
				mockRepo.On("GetTranslation", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.Anything, mockTx).
					Run(func(args mock.Arguments) {
						translation := args.Get(2).(*Translation)
						translation.ID = 0
					}).
					Return(errors.New("Translation Error"))
			},
			expectedResult: 0,
			expectedError: errors.New("Translation Error"),
			expectedOutput: "Error while searching for translation existance in a DB: Translation Error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockDB := new(MockDatabase)
			dbInterface := &DBInterface{
				DB: mockDB,
				repo: mockRepo,
			}
			mockTx := &gorm.DB{} 
			
			tt.mockGetTranslation(mockRepo, mockTx)

			var result uint
			var err error

			output := captureStdout(func() {
				result, err = dbInterface.GetTranslationID(mockTx, tt.wordID, tt.translation)
			})

			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectedResult, result)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} 

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateExampleSentences(t *testing.T) {
	tests := []struct {
		name string
		examples []string
		expected []ExampleSentence
	}{
		{
			name:     "Valid input",
			examples: []string{"This is an example", "Another example sentence"},
			expected: []ExampleSentence{
				{Sentence: "This is an example"},
				{Sentence: "Another example sentence"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createExampleSentences(tt.examples)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertTranslations(t *testing.T) {
	tests := []struct {
		name string
		translations []Translation
		expected []*model.Translation
	}{
		{
			name: "Valid translations",
			translations: []Translation{
				{ID: 1, Translation: "Book", ExampleSentences: []ExampleSentence{{Sentence: "This is a book"}}},
				{ID: 2, Translation: "Pen", ExampleSentences: []ExampleSentence{{Sentence: "This is a pen"}}},
			},
			expected: []*model.Translation{
				{ID: "1", Translation: "Book", ExampleSentences: []*model.ExampleSentence{{ID: "0", Sentence: "This is a book"}}},
				{ID: "2", Translation: "Pen", ExampleSentences: []*model.ExampleSentence{{ID: "0", Sentence: "This is a pen"}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertTranslations(tt.translations)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertExampleSentences(t *testing.T) {
	tests := []struct {
		name string
		sentences []ExampleSentence
		expected []*model.ExampleSentence
	}{
		{
			name: "Valid sentences",
			sentences: []ExampleSentence{
				{ID: 1, Sentence: "This is a test"},
				{ID: 2, Sentence: "Another test sentence"},
			},
			expected: []*model.ExampleSentence{
				{ID: "1", Sentence: "This is a test"},
				{ID: "2", Sentence: "Another test sentence"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertExampleSentences(tt.sentences)
			assert.Equal(t, tt.expected, result)
		})
	}
}






