package main

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
	"dictionary-app/server/graph/model"
)


func TestModifyAddCommand_Execute(t *testing.T) {
	mockHandler := new(MockQueriesHandler)

	tests := []struct {
		name string
		input ModifyAddCommand
		expectedResult bool
		expectedOutput string
		mockBehavior func()

	}{
		{
			name: "Success - Execute With Translation Type",
			input: ModifyAddCommand{
				handler: mockHandler,
				targetType: "translation",
				word: "word1",
				translation: "translation1",
				examples: []string{"example1"},
			},
			expectedResult: true,
			expectedOutput: "",
			mockBehavior: func() {
				mockHandler.On("SendAddTranslationMutation", mock.Anything, model.FullRecordInput{
					Word: "word1",
					Translation: "translation1",
					Examples: []string{"example1"},
				}).Return(true, nil)
			},
		},
		{
			name: "Success - Execute With Example Type",
			input: ModifyAddCommand{
				handler: mockHandler,
				targetType: "example",
				word: "word2",
				translation: "translation2",
				examples: []string{"example2"},
			},
			expectedResult: true,
			expectedOutput: "", 
			mockBehavior: func() {
				mockHandler.On("SendAddExampleMutation", mock.Anything, model.FullRecordInput{
					Word: "word2",
					Translation: "translation2",
					Examples: []string{"example2"},
				}).Return(true, nil)
			},
		},
		{
			name: "Failure - Execute With Invalid TargetType",
			input: ModifyAddCommand{
				handler: mockHandler,
				targetType: "invalid",
				word: "word3",
				translation: "translation3",
				examples: []string{"example3"},
			},
			expectedResult: false,
			expectedOutput: "Invalid modify add command\n",
			mockBehavior: func() {
			},
		},
		{
			name: "Failure - Execute With Error From SendAddTranslationMutation",
			input: ModifyAddCommand{
				handler: mockHandler,
				targetType: "translation",
				word: "word4",
				translation: "translation4",
				examples: []string{"example4"},
			},
			expectedResult: false,
			expectedOutput: "Error: Add Error\n", 
			mockBehavior: func() {
				mockHandler.On("SendAddTranslationMutation", mock.Anything, model.FullRecordInput{
					Word: "word4",
					Translation: "translation4",
					Examples: []string{"example4"},
				}).Return(false, errors.New("Add Error")) 
			},
		},
		{
			name: "Failure - Execute With Error From SendAddExampleMutation",
			input: ModifyAddCommand{
				handler: mockHandler,
				targetType: "example",
				word: "word5",
				translation: "translation5",
				examples: []string{"example5"},
			},
			expectedResult: false,
			expectedOutput: "Error: Add Error\n", 
			mockBehavior: func() {
				mockHandler.On("SendAddExampleMutation", mock.Anything, model.FullRecordInput{
					Word: "word5",
					Translation: "translation5",
					Examples: []string{"example5"},
				}).Return(false, errors.New("Add Error"))
			},
		},
		{
			name: "Failure - Execute With False From SendAddTranslationMutation",
			input: ModifyAddCommand{
				handler: mockHandler,
				targetType: "translation",
				word: "word6",
				translation: "translation6",
				examples: []string{"example6"},
			},
			expectedResult: false,
			expectedOutput: "", 
			mockBehavior: func() {
				mockHandler.On("SendAddTranslationMutation", mock.Anything, model.FullRecordInput{
					Word: "word6",
					Translation: "translation6",
					Examples: []string{"example6"},
				}).Return(false, nil) 
			},
		},
		{
			name: "Failure - Execute With False From SendAddExampleMutation",
			input: ModifyAddCommand{
				handler: mockHandler,
				targetType: "example",
				word: "word7",
				translation: "translation7",
				examples: []string{"example7"},
			},
			expectedResult: false,
			expectedOutput: "", 
			mockBehavior: func() {
				mockHandler.On("SendAddExampleMutation", mock.Anything, model.FullRecordInput{
					Word: "word7",
					Translation: "translation7",
					Examples: []string{"example7"},
				}).Return(false, nil)
			},
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			var result bool
			output := captureStdout(func() {
				result = tt.input.Execute()
			})

			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectedResult, result)
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestModifyDeleteCommand_Execute(t *testing.T) {
	mockHandler := new(MockQueriesHandler)

	tests := []struct {
		name string
		input ModifyDeleteCommand
		expectedResult bool
		expectedOutput string
		mockBehavior func()
	}{
		{
			name: "Success - Execute With Translation Type",
			input: ModifyDeleteCommand{
				handler: mockHandler,
				targetType: "translation",
				word: "word1",
				translation: "translation1",
			},
			expectedResult: true,
			expectedOutput: "",
			mockBehavior: func() {
				mockHandler.On("SendRemoveTranslationMutation", mock.Anything, "word1", "translation1").Return(true, nil)
			},
		},
		{
			name: "Success - Execute With Example Type",
			input: ModifyDeleteCommand{
				handler: mockHandler,
				targetType: "example",
				word: "word2",
				translation: "translation2",
				examples: []string{"example2"},
			},
			expectedResult: true,
			expectedOutput: "",
			mockBehavior: func() {
				mockHandler.On("SendRemoveExampleMutation", mock.Anything, model.FullRecordInput{
					Word: "word2",
					Translation: "translation2",
					Examples: []string{"example2"},
				}).Return(true, nil)
			},
		},
		{
			name: "Failure - Execute With Invalid TargetType",
			input: ModifyDeleteCommand{
				handler: mockHandler,
				targetType: "invalid",
				word: "word3",
				translation: "translation3",
				examples: []string{"example3"},
			},
			expectedResult: false,
			expectedOutput: "Invalid modify delete command\n",
			mockBehavior: func() {},
		},
		{
			name: "Failure - Execute With Error From SendRemoveTranslationMutation",
			input: ModifyDeleteCommand{
				handler: mockHandler,
				targetType: "translation",
				word: "word4",
				translation: "translation4",
			},
			expectedResult: false,
			expectedOutput: "Error: Remove Error\n", 
			mockBehavior: func() {
				mockHandler.On("SendRemoveTranslationMutation", mock.Anything, "word4", "translation4").Return(false, errors.New("Remove Error"))
			},
		},
		{
			name: "Failure - Execute With Error From SendRemoveExampleMutation",
			input: ModifyDeleteCommand{
				handler: mockHandler,
				targetType: "example",
				word: "word5",
				translation: "translation5",
				examples: []string{"example5"},
			},
			expectedResult: false,
			expectedOutput: "Error: Remove Error\n", 
			mockBehavior: func() {
				mockHandler.On("SendRemoveExampleMutation", mock.Anything, model.FullRecordInput{
					Word: "word5",
					Translation: "translation5",
					Examples: []string{"example5"},
				}).Return(false, errors.New("Remove Error"))
			},
		},
		{
			name: "Failure - Execute With False From SendRemoveTranslationMutation",
			input: ModifyDeleteCommand{
				handler: mockHandler,
				targetType: "translation",
				word: "word6",
				translation: "translation6",
			},
			expectedResult: false,
			expectedOutput: "", 
			mockBehavior: func() {
				mockHandler.On("SendRemoveTranslationMutation", mock.Anything, "word6", "translation6").Return(false, nil) 
			},
		},
		{
			name: "Failure - Execute With False From SendRemoveExampleMutation",
			input: ModifyDeleteCommand{
				handler: mockHandler,
				targetType: "example",
				word: "word7",
				translation: "translation7",
				examples: []string{"example7"},
			},
			expectedResult: false,
			expectedOutput: "", 
			mockBehavior: func() {
				mockHandler.On("SendRemoveExampleMutation", mock.Anything, model.FullRecordInput{
					Word: "word7",
					Translation: "translation7",
					Examples: []string{"example7"},
				}).Return(false, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			var result bool
			output := captureStdout(func() {
				result = tt.input.Execute()
			})

			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectedResult, result)
			mockHandler.AssertExpectations(t)
		})
	}
}
