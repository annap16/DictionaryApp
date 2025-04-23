package main

import (
	"context"
	"dictionary-app/server/graph/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSendCreateMutation(t *testing.T) {
	input := model.FullRecordInput{
		Word:        "word1",
		Translation: "translation1",
		Examples:    []string{"example1"},
	}

	tests := []struct {
		name           string
		mockBehavior   func(mockClient *MockGraphQLClient)
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Success - CreateTranslation Returns True",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					response := args.Get(2).(*AddWordResponse)
					*response = AddWordResponse{CreateTranslation: true}
				}).Return(nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Failure - Error From GraphQL Client",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphQL Error"))
			},
			expectedResult: false,
			expectedError:  errors.New("graphQL Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			tt.mockBehavior(mockClient)

			queriesHandler := &QueriesHandlerQL{
				client: mockClient,
			}

			result, err := queriesHandler.SendCreateMutation(context.Background(), input)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestSendReceiveMutation(t *testing.T) {
	input := "word1"

	tests := []struct {
		name           string
		mockBehavior   func(mockClient *MockGraphQLClient)
		expectedResult string
		expectedError  error
	}{
		{
			name: "Success - Word Found",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					response := args.Get(2).(*ReceiveResponse)
					*response = ReceiveResponse{
						GetWordTranslation: struct {
							ID           string `json:"id"`
							Word         string `json:"word"`
							Translations []struct {
								ID               string `json:"id"`
								Translation      string `json:"translation"`
								ExampleSentences []struct {
									ID       string `json:"id"`
									Sentence string `json:"sentence"`
								} `json:"exampleSentences"`
							} `json:"translations"`
						}{
							Word: "word1",
							Translations: []struct {
								ID               string `json:"id"`
								Translation      string `json:"translation"`
								ExampleSentences []struct {
									ID       string `json:"id"`
									Sentence string `json:"sentence"`
								} `json:"exampleSentences"`
							}{
								{
									Translation: "translation1",
									ExampleSentences: []struct {
										ID       string `json:"id"`
										Sentence string `json:"sentence"`
									}{
										{Sentence: "example sentence 1"},
										{Sentence: "example sentence 2"},
									},
								},
							},
						},
					}
				}).Return(nil)
			},
			expectedResult: "Słowo: word1\nTłumaczenie: translation1\n\tPrzykład: example sentence 1\n\tPrzykład: example sentence 2",
			expectedError:  nil,
		},
		{
			name: "Failure - Error from GraphQL Client",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphQL Error"))
			},
			expectedResult: "",
			expectedError:  errors.New("graphQL Error"),
		},
		{
			name: "No Record Found",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Record not found"))
			},
			expectedResult: "",
			expectedError:  errors.New("Record not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			tt.mockBehavior(mockClient)

			queriesHandler := &QueriesHandlerQL{
				client: mockClient,
			}

			result, err := queriesHandler.SendReceiveMutation(context.Background(), input)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestSendRemoveMutation(t *testing.T) {
	input := "word1"
	tests := []struct {
		name           string
		mockBehavior   func(mockClient *MockGraphQLClient)
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Success - DeleteWord Returns True",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					response := args.Get(2).(*RemoveWordResponse)
					*response = RemoveWordResponse{DeleteWord: true}
				}).Return(nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Failure - Error From GraphQL Client",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphQL Error"))
			},
			expectedResult: false,
			expectedError:  errors.New("graphQL Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			tt.mockBehavior(mockClient)

			queriesHandler := &QueriesHandlerQL{
				client: mockClient,
			}

			result, err := queriesHandler.SendRemoveMutation(context.Background(), input)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestSendRemoveTranslationMutation(t *testing.T) {
	word := "word1"
	translation := "translation1"
	tests := []struct {
		name           string
		mockBehavior   func(mockClient *MockGraphQLClient)
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Success - DeleteTranslation Returns True",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					response := args.Get(2).(*RemoveTranslationResponse)
					*response = RemoveTranslationResponse{DeleteTranslation: true}
				}).Return(nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Failure - Error From GraphQL Client",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphQL Error"))
			},
			expectedResult: false,
			expectedError:  errors.New("graphQL Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			tt.mockBehavior(mockClient)

			queriesHandler := &QueriesHandlerQL{
				client: mockClient,
			}

			result, err := queriesHandler.SendRemoveTranslationMutation(context.Background(), word, translation)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestSendRemoveExampleMutation(t *testing.T) {
	input := model.FullRecordInput{
		Word:        "word1",
		Translation: "translation1",
		Examples:    []string{"example1"},
	}

	tests := []struct {
		name           string
		mockBehavior   func(mockClient *MockGraphQLClient)
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Success - DeleteExample Returns True",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					response := args.Get(2).(*RemoveExampleResponse)
					*response = RemoveExampleResponse{DeleteExample: true}
				}).Return(nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Failure - Error From GraphQL Client",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphQL Error"))
			},
			expectedResult: false,
			expectedError:  errors.New("graphQL Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			tt.mockBehavior(mockClient)

			queriesHandler := &QueriesHandlerQL{
				client: mockClient,
			}

			result, err := queriesHandler.SendRemoveExampleMutation(context.Background(), input)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestSendAddTranslationMutation(t *testing.T) {
	input := model.FullRecordInput{
		Word:        "word1",
		Translation: "translation1",
		Examples:    []string{"example1"},
	}

	tests := []struct {
		name           string
		mockBehavior   func(mockClient *MockGraphQLClient)
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Success - AddTranslation Returns True",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					response := args.Get(2).(*AddTranslationResponse)
					*response = AddTranslationResponse{AddTranslation: true}
				}).Return(nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Failure - Error From GraphQL Client",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphQL Error"))
			},
			expectedResult: false,
			expectedError:  errors.New("graphQL Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			tt.mockBehavior(mockClient)

			queriesHandler := &QueriesHandlerQL{
				client: mockClient,
			}

			result, err := queriesHandler.SendAddTranslationMutation(context.Background(), input)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestSendAddExampleMutation(t *testing.T) {
	input := model.FullRecordInput{
		Word:        "word1",
		Translation: "translation1",
		Examples:    []string{"example1"},
	}

	tests := []struct {
		name           string
		mockBehavior   func(mockClient *MockGraphQLClient)
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Success - AddExample Returns True",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					response := args.Get(2).(*AddExampleResponse)
					*response = AddExampleResponse{AddExample: true}
				}).Return(nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Failure - Error From GraphQL Client",
			mockBehavior: func(mockClient *MockGraphQLClient) {
				mockClient.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphQL Error"))
			},
			expectedResult: false,
			expectedError:  errors.New("graphQL Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGraphQLClient)
			tt.mockBehavior(mockClient)

			queriesHandler := &QueriesHandlerQL{
				client: mockClient,
			}

			result, err := queriesHandler.SendAddExampleMutation(context.Background(), input)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestParseReceiveResponse(t *testing.T) {

	mockResponse := ReceiveResponse{
		GetWordTranslation: struct {
			ID           string `json:"id"`
			Word         string `json:"word"`
			Translations []struct {
				ID               string `json:"id"`
				Translation      string `json:"translation"`
				ExampleSentences []struct {
					ID       string `json:"id"`
					Sentence string `json:"sentence"`
				} `json:"exampleSentences"`
			} `json:"translations"`
		}{
			Word: "word1",
			Translations: []struct {
				ID               string `json:"id"`
				Translation      string `json:"translation"`
				ExampleSentences []struct {
					ID       string `json:"id"`
					Sentence string `json:"sentence"`
				} `json:"exampleSentences"`
			}{
				{
					Translation: "translation1",
					ExampleSentences: []struct {
						ID       string `json:"id"`
						Sentence string `json:"sentence"`
					}{
						{Sentence: "example sentence 1"},
						{Sentence: "example sentence 2"},
					},
				},
			},
		},
	}

	expectedOutput := "Słowo: word1\nTłumaczenie: translation1\n\tPrzykład: example sentence 1\n\tPrzykład: example sentence 2"

	q := &QueriesHandlerQL{}

	result := q.ParseReceiveResponse(mockResponse)

	assert.Equal(t, expectedOutput, result)
}
