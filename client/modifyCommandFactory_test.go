package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAction(t *testing.T) {
	mockHandler := new(MockQueriesHandler)

	tests := []struct {
		name           string
		command        string
		expectedAction ModifyAction
		expectedError  error
	}{
		{
			name:    "Success - Create Add Command",
			command: "modyfikuj dodaj tłumaczenie word translation",
			expectedAction: &ModifyAddCommand{
				handler:     mockHandler,
				targetType:  "tłumaczenie",
				word:        "word",
				translation: "translation",
				examples:    nil,
			},
			expectedError: nil,
		},
		{
			name:    "Success - Create Delete Command",
			command: "modyfikuj usuń tłumaczenie word translation",
			expectedAction: &ModifyDeleteCommand{
				handler:     mockHandler,
				targetType:  "tłumaczenie",
				word:        "word",
				translation: "translation",
				examples:    nil,
			},
			expectedError: nil,
		},
		{
			name:           "Failure - Wrong Key word",
			command:        "modyfikuj create translation word translation",
			expectedAction: nil,
			expectedError:  errors.New("Niepoprawna składnia dla polecenia modyfikacji słowa"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := &ModifyCommandFactory{}

			action, err := factory.CreateAction(mockHandler, tt.command)

			assert.Equal(t, tt.expectedAction, action)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
