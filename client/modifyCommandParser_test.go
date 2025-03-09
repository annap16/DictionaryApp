package main

import(
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"

)

func TestParseModifyCommand(t *testing.T) {
	tests := []struct {
		name string
		command string
		expectedError error
		expectedModel *ModifyCommandParams
	}{
		{
			name:    "Success - Add Translation",
			command: "modify add translation word translation [example1] [example2]",
			expectedError: nil,
			expectedModel: &ModifyCommandParams{
				Action:     "add",
				TargetType: "translation",
				Word:       "word",
				Translation: "translation",
				Examples:   []string{"example1", "example2"},
			},
		},
		{
			name:    "Success - Add Example",
			command: "modify add example word translation [example1] [example2]",
			expectedError: nil,
			expectedModel: &ModifyCommandParams{
				Action:     "add",
				TargetType: "example",
				Word:       "word",
				Translation: "translation",
				Examples:   []string{"example1", "example2"},
			},
		},
		{
			name:    "Success - Delete Example",
			command: "modify delete example word translation [example1]",
			expectedError: nil,
			expectedModel: &ModifyCommandParams{
				Action:     "delete",
				TargetType: "example",
				Word:       "word",
				Translation: "translation",
				Examples:   []string{"example1"},
			},
		},
		{
			name:    "Success - Delete Translation",
			command: "modify delete translation word translation",
			expectedError: nil,
			expectedModel: &ModifyCommandParams{
				Action:     "delete",
				TargetType: "translation",
				Word:       "word",
				Translation: "translation",
				Examples:  nil,
			},
		},
		{
			name:    "Failure - Wrong Key Word",
			command: "modify create translation word translation",
			expectedError: errors.New("Niepoprawna składnia dla polecenia modyfikacji słowa"),
			expectedModel: nil,
		},
		{
			name:    "Failure - Add - Wrong Command",
			command: "modify add translation word translation [][Wrong]]",
			expectedError: errors.New("Niepoprawna składnia dla polecenia dodawania"),
			expectedModel: nil,
		},
		{
			name:    "Failure - Delete - Wrong Command",
			command: "modify delete translation word translation [][Wrong]]",
			expectedError: errors.New("Niepoprawna składnia dla polecenia usuwania"),
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
