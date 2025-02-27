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
		createTranslation(input: $input)
	}
	`)

	request.Var("input", input)

	err := q.client.Run(context.Background(), request, nil)
	if err!=nil{
		log.Fatal("Error:", err)
	}

	return nil
 }

 func (q *QueriesHandler) SendReceiveMutation(ctx context.Context, input string) (string, error) {

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
	

	var response ReceiveResponse
	err := q.client.Run(context.Background(), request, &response)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return "", err
		}
		log.Fatal("Error:", err)
		return "", err
	}

	return q.ParseReceiveResponse(response), nil
}

func (q *QueriesHandler) ParseReceiveResponse(response ReceiveResponse) string{
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


func (q *QueriesHandler) SendRemoveMutation(ctx context.Context, input string) (bool, error) {
	request := graphql.NewRequest(`
		mutation ($word: String!) {
			deleteWord(word: $word)
		}
	`)
	request.Var("word", input)

	var response RemoveWordResponse
	
	if err := q.client.Run(context.Background(), request, &response); err != nil {
		
		log.Fatal("Error:", err)
		return false, err
	}

	return response.DeleteWord, nil

}

func(q *QueriesHandler) SendRemoveTranslationMutation(ctx context.Context, input string) (bool, error){
	request := graphql.NewRequest(`
		mutation ($translation: String!) {
    		deleteTranslation(translation: $translation)
		}
		`)
	request.Var("translation", input)

	var response RemoveTranslationResponse

    if err := q.client.Run(context.Background(), request, &response); err != nil {
        log.Fatal("Error:", err)
        return false, err
    }

	return response.DeleteTranslation, nil
}

func(q *QueriesHandler) SendRemoveExampleMutation(ctx context.Context, input string) (bool, error){
	request := graphql.NewRequest(`
	mutation($example: String!) {
		deleteExample(example: $example)
	}
	`)
	request.Var("example", input)

	var response RemoveExampleResponse

	if err := q.client.Run(context.Background(), request, &response); err != nil {
		log.Fatal("Error:", err)
	return false, err
	}	
	
	return response.DeleteExample, nil
}




