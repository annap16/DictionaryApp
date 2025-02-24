package main

import(
	"context"
	"log"
	"strings"
	"github.com/machinebox/graphql"
	"github.com/annap16/DictionaryApp/graph/model"
)

type QueriesHandler struct{
	client *graphql.Client
}


func (q *QueriesHandler) SendCreateMutation(ctx context.Context, input model.CreateTranslationInput) error {
	request := graphql.NewRequest(`
	mutation ($input: CreateTranslationInput!) {
		createTranslation(input: $input) {
			id
			word
			translations {
				id
				translation
				exampleSentences {
					id
					sentence
				}
			}
		}
	}
`)	
	request.Var("input", input)

	err := q.client.Run(context.Background(), request, nil)
	if err!=nil{
		log.Fatal("Error:", err)
	}

	return nil
 }

 // TODO move to more accurate place
 type ReciveResponse struct {
	GetWordTranslation struct { 
		ID           string `json:"id"`
		Word         string `json:"word"`
		Translations []struct {
			ID              string `json:"id"`
			Translation     string `json:"translation"`
			ExampleSentences []struct {
				ID       string `json:"id"`
				Sentence string `json:"sentence"`
			} `json:"exampleSentences"`
		} `json:"translations"`
	} `json:"getWordTranslation"`
}


 func (q *QueriesHandler) SendReciveMutation(ctx context.Context, input string) (string, error) {

	request := graphql.NewRequest(`
		query ($word: String!) {
			getWordTranslation(word: $word) {
				word 
				translations {
					translation
					exampleSentences {
						sentence
					}
				}
			}
		}
	`)

	request.Var("word", input)

	var response ReciveResponse
	if err := q.client.Run(context.Background(), request, &response); err != nil {
		log.Fatal(err)
	}
	
	return q.ParseReceiveResponse(response), nil
}

func (q *QueriesHandler) ParseReceiveResponse(response ReciveResponse) string{
	var result string
	result += "Word: " + response.GetWordTranslation.Word + "\n"

	for _, translation := range response.GetWordTranslation.Translations {
		result += "Translation: " + translation.Translation + "\n"
		for _, sentence := range translation.ExampleSentences {
			result += "Examples: " + sentence.Sentence + "\n"
		}
	}
	result = strings.TrimSuffix(result, "\n")
	return result
}


