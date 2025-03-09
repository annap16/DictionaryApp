package main

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
	"dictionary-app/server/graph/model"
)

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
			expectedOutput: "Słowo zostało dodane do słownika\n",
			expectedResult: true,
		},
		{
			name: "Failure - Word Already Exists",
			command: "create word translation [example]",
			mockResponse: false,
			mockError: nil,
			expectedOutput: "Podane słowo istnieje już w słowniku\n",
			expectedResult: true, 
		},
		{
			name: "Failure - Mutation Error",
			command: "create word translation [example]",
			mockResponse: false,
			mockError: fmt.Errorf("mutation failed"),
			expectedOutput: "Wystąpił błąd podczas dodawania do bazy danych: mutation failed\n",
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
		expectedInput model.FullRecordInput
		expectedOutput string
		expectedResult bool
	}{
		{
			name: "Success - Input Check With One Example",
			command: "create ksiazka book [I love my new book]",
			expectedInput: model.FullRecordInput{
				Word: "ksiazka",
				Translation: "book",
				Examples: []string{"I love my new book"},
			},
			expectedOutput:"Słowo zostało dodane do słownika\n",
			expectedResult: true,
		},
		{
			name: "Success - Input Check With Two Examples",
			command: "create samochod car [A fast car] [My car is blue]",
			expectedInput: model.FullRecordInput{
				Word: "samochod",
				Translation: "car",
				Examples: []string{"A fast car", "My car is blue"},
			},
			expectedOutput: "Słowo zostało dodane do słownika\n",
			expectedResult: true,
		},
		{
			name: "Success - Input Check Without Any Examples",
			command: "create kot cat",
			expectedInput: model.FullRecordInput{
				Word: "kot",
				Translation: "cat",
				Examples: []string{},
			},
			expectedOutput: "Słowo zostało dodane do słownika\n",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			mockQueriesHandler := &MockQueriesHandler{client: mockClient}
			
			mockQueriesHandler.On("SendCreateMutation", mock.Anything, mock.MatchedBy(func(input model.FullRecordInput) bool {
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
			expectedOutput: "Podane słowo nie istnieje w słowniku\n",
			expectedResult: true,
		},
		{
			name: "Error - Mutation Error",
			command: "receive ksiazka",
			mockResponse: "",
			mockError: fmt.Errorf("mutation failed"),
			expectedOutput: "Wystąpił błąd podczas wyszukiwania w słowniku: mutation failed\n",
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

// --------------------------MODIFY COMMAND------------------------------------

type MockModifyFactory struct{
	mock.Mock
}

// func (m *MockModifyFactory) CreateAction(handler QueriesHandler, command string) (ModifyAction, error){
// 	args := m.Called(handler, command)
// 	if args.Get(0) != nil {
// 		return args.Get(0).(ModifyAction),
// 	}
// 	return nil
// }
func (m *MockModifyFactory) CreateAction(handler QueriesHandler, command string) (ModifyAction, error) {
	args := m.Called(handler, command)
	return args.Get(0).(ModifyAction), args.Error(1)
}






type MockModifyAction struct{
	mock.Mock
}

func (m *MockModifyAction) Execute() (bool, error){
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func TestHandleCommandModify(t *testing.T) {
	tests := []struct {
		name string
		command string
		factoryResponse ModifyAction
		execResponse bool
		expectedOutput string
		expectedResult bool
	}{
		{
			name: "Success - Example Added",
			command: "modify add example word translation [example]",
			factoryResponse: new(MockModifyAction),
			execResponse: true,
			expectedOutput: "Słowo zostało zmodyfikowane poprawnie\n",
			expectedResult: true,
		},
		{
			name: "Success - Translation Added",
			command: "modify add translation word translation",
			factoryResponse: new(MockModifyAction),
			execResponse: true,
			expectedOutput: "Słowo zostało zmodyfikowane poprawnie\n",
			expectedResult: true,
		},
		{
			name: "Failure - Add Example",
			command: "modify add example word translation [example]",
			factoryResponse: new(MockModifyAction),
			execResponse: false,
			expectedOutput: "Wprowadzono niepoprawne polecenie\n",
			expectedResult: true,
		},
		{
			name: "Failure - Add Translation",
			command: "modify add translation word translation",
			factoryResponse: new(MockModifyAction),
			execResponse: false,
			expectedOutput: "Wprowadzono niepoprawne polecenie\n",
			expectedResult: true,
		},
		{
			name: "Success - Example Deleted",
			command: "modify delete example word translation [example]",
			factoryResponse: new(MockModifyAction),
			execResponse: true,
			expectedOutput: "Słowo zostało zmodyfikowane poprawnie\n",
			expectedResult: true,
		},
		{
			name: "Success - Translation Deleted",
			command: "modify delete translation word translation",
			factoryResponse: new(MockModifyAction),
			execResponse: true,
			expectedOutput: "Słowo zostało zmodyfikowane poprawnie\n",
			expectedResult: true,
		},
		{
			name: "Failure - Delete Example",
			command: "modify delete example word translation [example]",
			factoryResponse: new(MockModifyAction),
			execResponse: false,
			expectedOutput: "Wprowadzono niepoprawne polecenie\n",
			expectedResult: true,
		},
		{
			name: "Failure - Delete Translation",
			command: "modify delete translation word translation",
			factoryResponse: new(MockModifyAction),
			execResponse: false,
			expectedOutput: "Wprowadzono niepoprawne polecenie\n",
			expectedResult: true,
		},


	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFactory := new(MockModifyFactory)
			mockModifyAction := tt.factoryResponse.(*MockModifyAction) 

			mockFactory.On("CreateAction", mock.Anything, mock.Anything).Return(mockModifyAction, nil)
			mockModifyAction.On("Execute").Return(tt.execResponse, nil)

			handler := &ModifyCommandHandler{
				modifyFactory: mockFactory, 
			}

			output := captureStdout(func() {
				result := handler.HandleCommand(tt.command)
				assert.Equal(t, tt.expectedResult, result)
			})

			assert.Equal(t, tt.expectedOutput, output)

			mockFactory.AssertExpectations(t)
			mockModifyAction.AssertExpectations(t)
		})
	}
}

func TestHandleCommandModify_OtherCommand(t *testing.T){
	mockClient := new(MockGraphQLClient)
    mockQueriesHandler := &MockQueriesHandler{client: mockClient}

    handler := &ModifyCommandHandler{handler: mockQueriesHandler}

    command := "create something"

    success := handler.HandleCommand(command)

    assert.False(t, success) 
}

// --------------------------REMOVE COMMAND------------------------------------


func TestHandleCommandRemove(t *testing.T){
	tests := []struct {
		name string
		command string
		mockResponse bool
		mockError error
		expectedOutput string
		expectedResult bool
	}{
		{
			name: "Success - Word Removed",
			command: "remove slowo",
			mockResponse: true,
			mockError: nil,
			expectedOutput: "Podane słowo i powiązane z nim dane zostały usunięte\n",
			expectedResult: true,
		},
		{
			name: "Failure - Word doesn't exist in the dictionary",
			command: "remove slowo",
			mockResponse: false,
			mockError: nil,
			expectedOutput: "Podane słowo nie istnieje w słowniku\n",
			expectedResult: true, 
		},
		{
			name: "Failure - Mutation Error",
			command: "remove slowo",
			mockResponse: false,
			mockError: fmt.Errorf("Mutation Failed"),
			expectedOutput: "Wystąpił błąd podczas usuwania słowa ze słownika: Mutation Failed\n",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			mockQueriesHandler := &MockQueriesHandler{client: mockClient}

			mockQueriesHandler.On("SendRemoveMutation", mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			handler := &RemoveCommandHandler{handler: mockQueriesHandler}

			output := captureStdout(func() {
				result := handler.HandleCommand(tt.command)
				assert.Equal(t, tt.expectedResult, result)
			})

			assert.Equal(t, tt.expectedOutput, output)

			mockQueriesHandler.AssertExpectations(t)
		})
	}
}

func TestHandleCommandRemoveSendModifyMutation_OtherCommand(t *testing.T){
	mockClient := new(MockGraphQLClient)
    mockQueriesHandler := &MockQueriesHandler{client: mockClient}

    handler := &RemoveCommandHandler{handler: mockQueriesHandler}

    command := "create something"

    success := handler.HandleCommand(command)

    mockQueriesHandler.AssertNotCalled(t, "SendRemoveMutation")
    assert.False(t, success) 
}

