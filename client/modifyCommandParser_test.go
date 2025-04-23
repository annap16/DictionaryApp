package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseModifyCommand(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		expectedError error
		expectedModel *ModifyCommandParams
	}{
		{
			name:          "Success - Add Translation",
			command:       "modyfikuj dodaj tłumaczenie word translation [example1] [example2]",
			expectedError: nil,
			expectedModel: &ModifyCommandParams{
				Action:      "dodaj",
				TargetType:  "tłumaczenie",
				Word:        "word",
				Translation: "translation",
				Examples:    []string{"example1", "example2"},
			},
		},
		{
			name:          "Success - Add Example",
			command:       "modyfikuj dodaj przykład word translation [example1] [example2]",
			expectedError: nil,
			expectedModel: &ModifyCommandParams{
				Action:      "dodaj",
				TargetType:  "przykład",
				Word:        "word",
				Translation: "translation",
				Examples:    []string{"example1", "example2"},
			},
		},
		{
			name:          "Success - Delete Example",
			command:       "modyfikuj usuń przykład word translation [example1]",
			expectedError: nil,
			expectedModel: &ModifyCommandParams{
				Action:      "usuń",
				TargetType:  "przykład",
				Word:        "word",
				Translation: "translation",
				Examples:    []string{"example1"},
			},
		},
		{
			name:          "Success - Delete Translation",
			command:       "modyfikuj usuń tłumaczenie word translation",
			expectedError: nil,
			expectedModel: &ModifyCommandParams{
				Action:      "usuń",
				TargetType:  "tłumaczenie",
				Word:        "word",
				Translation: "translation",
				Examples:    nil,
			},
		},
		{
			name:          "Failure - Wrong Key Word",
			command:       "modify create translation word translation",
			expectedError: errors.New("Niepoprawna składnia dla polecenia modyfikacji słowa"),
			expectedModel: nil,
		},
		{
			name:          "Failure - Add - Wrong Command",
			command:       "modyfikuj dodaj tłumaczenie word translation [][Wrong]]",
			expectedError: errors.New("Niepoprawna składnia dla polecenia dodawania"),
			expectedModel: nil,
		},
		{
			name:          "Failure - Delete - Wrong Command",
			command:       "modyfikuj usuń tłumaczenie word translation [][Wrong]]",
			expectedError: errors.New("Niepoprawna składnia dla polecenia usuwania"),
			expectedModel: nil,
		},
		{
			name:          "Failure - Command To Short For Processing",
			command:       "modify add",
			expectedError: errors.New("Wprowadzono niepoprawne polecenie"),
			expectedModel: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseModifyCommand(tt.command)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedModel, result)
		})
	}

}
