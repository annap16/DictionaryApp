package database

import (
	"database/sql"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) TransactionWrapper(DB Database, fn func(tx *gorm.DB) (bool, error)) (bool, error) {
	args := m.Called(DB, fn)
	if returnFn, ok := args.Get(0).(func(mock.Arguments) bool); ok {
		if errFn, ok := args.Get(1).(func(mock.Arguments) error); ok {
			return returnFn(args), errFn(args)
		}
	}
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CreateWord(word *Word, tx *gorm.DB) error {
	args := m.Called(word, tx)
	return args.Error(0)
}

func (m *MockRepository) CreateTranslation(translation *Translation, tx *gorm.DB) error {
	args := m.Called(translation, tx)
	return args.Error(0)
}

func (m *MockRepository) CreateExample(example *ExampleSentence, tx *gorm.DB) error {
	args := m.Called(example, tx)
	return args.Error(0)
}

func (m *MockRepository) DeleteWord(word string, tx *gorm.DB) (bool, error) {
	args := m.Called(word, tx)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) DeleteTranslation(wordID uint, translation string, tx *gorm.DB) (bool, error) {
	args := m.Called(wordID, translation, tx)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) DeleteExample(translationID uint, example string, tx *gorm.DB) (bool, error) {
	args := m.Called(translationID, example, tx)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetWord(word string, result *Word, tx *gorm.DB) error {
	args := m.Called(word, result, tx)
	return args.Error(0)
}

func (m *MockRepository) GetTranslation(wordID uint, translation string, result *Translation, tx *gorm.DB) error {
	args := m.Called(wordID, translation, result, tx)
	return args.Error(0)
}

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Begin(opts ...*sql.TxOptions) *gorm.DB {
	args := m.Called(opts)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDatabase) AutoMigrate(dst ...interface{}) error {
	args := m.Called(dst)
	return args.Error(0)
}
