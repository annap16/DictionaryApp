package main

import (
	"testing"
	"context"
	"bytes"
	"os"
	"fmt"
	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
	"github.com/annap16/DictionaryApp/graph/model"
	"github.com/machinebox/graphql"
)

// --------------------------HELPER FUNCTIONS TESTS-----------------------------

func TestCheckCreateSyntax(t *testing.T) {
	tests := []struct {
		name string
		command string
		want bool
	}{
		{
			name : "Success - Command Without Examples",
			command: "create xyz abc",
			want: true,
		},
		{
			name: "Success - Command With Capital Letters", 
			command: "CREATe something else", 
			want: true,
		},
		{
			name: "Success - Command With 2 Examples",
			command: "create word translation [Sentence one] [Sentence two]", 
			want: true,
		},
		{
			name: "Failure - Too Many Arguments Before Examples",
			command: "Create word translation something [else]", 
			want: false,
		},
		{
			name: "Failure - Not Enough Arguments Before Examples",
			command: "create missingTranslation [value1] [value2]", 
			want: false,
		},
		{
			name: "Failure - Empty Create Command With Whitespace Character",
			command: "create ", 
			want: false,
		},
		{
			name: "Failure - Command Different Than Create",
			command: "random command", 
			want: false,
		},
		{
			name: "Failure - Wrong Brackets",
			command: "create word translation [sentence[sentence2]]]", 
			want: false,
		},
		{
			name: "Failure - Empty Brackets",
			command: "create word translation []", 
			want: false,
		},
		{
			name: "Failure - Correct And Wrong Brackets Combined",
			command: "create word translation ][sentence]][]", 
			want: false,
		},
		{
			name: "Failure - Correct And Empty Brackets Combined",
			command: "create word [][ here] translation", 
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckCreateSyntax(tt.command)
			if got != tt.want {
				t.Errorf("CheckCreateSyntax(%q) = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name string
		query string
		want []string
	}{
		{
			name: "Success - Create Query With 2 Examples",
			query: "create word translation [value1] [value2 value3]", 
			want: []string{"value1", "value2 value3"},
		},
		{
			name: "Success - Create Query Without Any Examples",
			query: "create word translation", 
			want: []string{},
		},
		{
			name: "Success - Create Query With 1 Example",
			query: "create word translation [single]", 
			want: []string{"single"},
		},
		{
			name: "Success - Create Query 3 Examples",
			query: "create word translation [one] [two] [three]", 
			want: []string{"one", "two", "three"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseQuery(tt.query)
			if !equalSlices(got, tt.want) {
				t.Errorf("ParseQuery(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

// --------------------------HELPER FUNCTION-----------------------------------

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// --------------------------MOCK STRUCTURE------------------------------------

type MockGraphQLClient struct {
	mock.Mock
}

func (m *MockGraphQLClient) Run(ctx context.Context, req *graphql.Request, res interface{}) error {
	args := m.Called(ctx, req, res)
	return args.Error(0)
}

type MockQueriesHandler struct {
	client *MockGraphQLClient
	mock.Mock
}

func (m *MockQueriesHandler) SendCreateMutation(ctx context.Context, input model.CreateTranslationInput) (bool, error) {
	args := m.Called(ctx, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendReceiveMutation(ctx context.Context, input string) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}

func (m *MockQueriesHandler) SendRemoveMutation(ctx context.Context, input string) (bool, error) {
	args := m.Called(ctx, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendRemoveTranslationMutation(ctx context.Context, input string) (bool, error) {
	args := m.Called(ctx, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendRemoveExampleMutation(ctx context.Context, translation string, input string) (bool, error) {
	args := m.Called(ctx, translation, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendAddTranslationMutation(ctx context.Context, input model.CreateTranslationInput) (bool, error) {
	args := m.Called(ctx, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendAddExampleMutation(ctx context.Context, translation string, sentences []string) (bool, error) {
	args := m.Called(ctx, translation, sentences)
	return args.Bool(0), args.Error(1)
}

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

// --------------------------CREATE COMMAND-------------------------------------

func TestHandleCommandCreate(t *testing.T) {
	tests := []struct {
		name string
		command string
		mockResponse bool
		mockError error
		expectedOutput string
		expectedResult bool
	}{
		{
			name: "Success - Word Created",
			command: "create word translation [example]",
			mockResponse: true,
			mockError: nil,
			expectedOutput: "Word successfully added to dictionary\n",
			expectedResult: true,
		},
		{
			name: "Failure - Word Already Exists",
			command: "create word translation [example]",
			mockResponse: false,
			mockError: nil,
			expectedOutput: "Given word already exists in the dictionary\n",
			expectedResult: true, 
		},
		{
			name: "Failure - Mutation Error",
			command: "create word translation [example]",
			mockResponse: false,
			mockError: fmt.Errorf("mutation failed"),
			expectedOutput: "Error occured while adding word to the dictionary\n",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			mockQueriesHandler := &MockQueriesHandler{client: mockClient}

			mockQueriesHandler.On("SendCreateMutation", mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			handler := &CreateCommandHandler{handler: mockQueriesHandler}

			output := captureStdout(func() {
				result := handler.HandleCommand(tt.command)
				assert.Equal(t, tt.expectedResult, result)
			})

			assert.Equal(t, tt.expectedOutput, output)

			mockQueriesHandler.AssertExpectations(t)
		})
	}
}

func TestHandleCommandCreate_OtherCommand(t *testing.T) {
    mockClient := new(MockGraphQLClient)
    mockQueriesHandler := &MockQueriesHandler{client: mockClient}

    handler := &CreateCommandHandler{handler: mockQueriesHandler}

    command := "add something"

    success := handler.HandleCommand(command)

    mockQueriesHandler.AssertNotCalled(t, "SendCreateMutation")
    assert.False(t, success) 
}

func TestHandleCommandCreate_CorrectInput(t *testing.T) {
	tests := []struct {
		name string
		command string
		expectedInput model.CreateTranslationInput
		expectedOutput string
		expectedResult bool
	}{
		{
			name: "Success - Input Check With One Example",
			command: "create ksiazka book [I love my new book]",
			expectedInput: model.CreateTranslationInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			expectedOutput: "Word successfully added to dictionary\n",
			expectedResult: true,
		},
		{
			name: "Success - Input Check With Two Examples",
			command: "create samochod car [A fast car] [My car is blue]",
			expectedInput: model.CreateTranslationInput{
				Word: "samochod",
				Translation: "car",
				Examples: []string{"A fast car", "My car is blue"},
			},
			expectedOutput: "Word successfully added to dictionary\n",
			expectedResult: true,
		},
		{
			name: "Success - Input Check Without Any Examples",
			command: "create kot cat",
			expectedInput: model.CreateTranslationInput{
				Word: "kot",
				Translation: "cat",
				Examples: []string{},
			},
			expectedOutput: "Word successfully added to dictionary\n",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			mockQueriesHandler := &MockQueriesHandler{client: mockClient}
			
			mockQueriesHandler.On("SendCreateMutation", mock.Anything, mock.MatchedBy(func(input model.CreateTranslationInput) bool {
				return input.Word == tt.expectedInput.Word &&
					input.Translation == tt.expectedInput.Translation &&
					len(input.Examples) == len(tt.expectedInput.Examples) &&
					equalSlices(input.Examples, tt.expectedInput.Examples)
			})).Return(true, nil)

			handler := &CreateCommandHandler{handler: mockQueriesHandler}


			output := captureStdout(func() {
				result := handler.HandleCommand(tt.command)
				assert.Equal(t, tt.expectedResult, result)
			})

			assert.Equal(t, tt.expectedOutput, output)

			mockQueriesHandler.AssertExpectations(t)
		})
	}
}

// --------------------------RECEIVE COMMAND------------------------------------

func TestHandleCommandReceive(t *testing.T){
	tests := []struct {
		name string
		command string
		mockResponse string
		mockError error
		expectedOutput string
		expectedResult bool
	}{
		{
			name: "Success - Word Received",
			command: "receive word",
			mockResponse: "result",
			mockError: nil,
			expectedOutput: "result\n",
			expectedResult: true,
		},
		{
			name: "Failure - Word Not Found",
			command: "receive word",
			mockResponse: "",
			mockError: nil,
			expectedOutput: "The word doesn't exist in the dictionary\n",
			expectedResult: true,
		},
		{
			name: "Error - Mutation Error",
			command: "receive ksiazka",
			mockResponse: "",
			mockError: fmt.Errorf("mutation failed"),
			expectedOutput: "Error while receiving a word\n",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			mockQueriesHandler := &MockQueriesHandler{client: mockClient}

			mockQueriesHandler.On("SendReceiveMutation", mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			handler := &ReceiveCommandHandler{handler: mockQueriesHandler}

			output := captureStdout(func() {
				result := handler.HandleCommand(tt.command)
				assert.Equal(t, tt.expectedResult, result)
			})

			assert.Equal(t, tt.expectedOutput, output)

			mockQueriesHandler.AssertExpectations(t)
		})
	}
	
}

func TestHandleCommandReceive_OtherCommand(t *testing.T) {
    mockClient := new(MockGraphQLClient)
    mockQueriesHandler := &MockQueriesHandler{client: mockClient}

    handler := &ReceiveCommandHandler{handler: mockQueriesHandler}

    command := "create something"

    success := handler.HandleCommand(command)

    mockQueriesHandler.AssertNotCalled(t, "SendReceiveMutation")
    assert.False(t, success) 
}

func TestHandleCommandReceive_CorrectInput(t *testing.T) {
	tests := []struct {
		name string
		command string
		expectedInput string
		expectedOutput string
		expectedResult bool
	}{
		{
			name: "Success - Word Received (ksiazka)",
			command: "receive ksiazka",
			expectedInput: "ksiazka",
			expectedOutput: "result\n",
			expectedResult: true,
		},
		{
			name: "Success - Word Received (samochod)",
			command: "receive samochod",
			expectedInput: "samochod",
			expectedOutput: "result\n",
			expectedResult: true,
		},
		{
			name: "Success - Word Received (kot)",
			command: "receive kot",
			expectedInput: "kot",
			expectedOutput: "result\n",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			mockQueriesHandler := &MockQueriesHandler{client: mockClient}
			
			mockQueriesHandler.On("SendReceiveMutation", mock.Anything, mock.MatchedBy(func(input string) bool {
				return input == tt.expectedInput})).Return("result", nil)

			handler := &ReceiveCommandHandler{handler: mockQueriesHandler}

			output := captureStdout(func() {
				result := handler.HandleCommand(tt.command)
				assert.Equal(t, tt.expectedResult, result)
			})

			assert.Equal(t, tt.expectedOutput, output)

			mockQueriesHandler.AssertExpectations(t)
		})
	}
}

