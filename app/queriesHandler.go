package main

import(
	"context"
	"log"
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