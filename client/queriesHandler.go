package main

import(
	"context"
	"strings"
	"fmt"
	"github.com/machinebox/graphql"
	"dictionary-app/server/graph/model"
)

type QueriesHandlerQL struct{
	client GraphQLClient
}


func (q *QueriesHandlerQL) SendCreateMutation(ctx context.Context, input model.FullRecordInput) (bool, error) {
	request := graphql.NewRequest(`
	mutation ($input: FullRecordInput!) {
		createTranslation(input: $input)
	}
	`)

	request.Var("input", input)

	var response AddWordResponse

	err := q.client.Run(context.Background(), request, &response)
	if err!=nil{
		fmt.Println("Error while adding a word:", err)
		return false, err
	}

	return response.CreateTranslation, nil
 }



 func (q *QueriesHandlerQL) SendReceiveMutation(ctx context.Context, input string) (string, error) {

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
		if strings.Contains(err.Error(), "Record not found") {
			return "", nil
		}
		fmt.Println("Error while searching for a word:", err)
		return "", err
	}

	return q.ParseReceiveResponse(response), nil
}

func (q *QueriesHandlerQL) ParseReceiveResponse(response ReceiveResponse) string{
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

func (q *QueriesHandlerQL) SendRemoveMutation(ctx context.Context, input string) (bool, error) {
	request := graphql.NewRequest(`
		mutation ($word: String!) {
			deleteWord(word: $word)
		}
	`)
	request.Var("word", input)

	var response RemoveWordResponse
	
	if err := q.client.Run(context.Background(), request, &response); err != nil {
		
		fmt.Println("Error while removing word from a dictionary:", err)
		return false, err
	}

	return response.DeleteWord, nil

}

func(q *QueriesHandlerQL) SendRemoveTranslationMutation(ctx context.Context, word string, input string) (bool, error){
	request := graphql.NewRequest(`
		mutation ($word: String! $translation: String!) {
    		deleteTranslation(word: $word translation: $translation)
		}
		`)
	request.Var("translation", input)
	request.Var("word", word)

	var response RemoveTranslationResponse

    if err := q.client.Run(context.Background(), request, &response); err != nil {
        fmt.Println("Error while deleting translation:", err)
        return false, err
    }

	return response.DeleteTranslation, nil
}

func(q *QueriesHandlerQL) SendRemoveExampleMutation(ctx context.Context, input model.FullRecordInput) (bool, error){
	request := graphql.NewRequest(`
	mutation($input: FullRecordInput!) {
		deleteExample(input: $input)
	}
	`)

	request.Var("input", input)

	var response RemoveExampleResponse

	if err := q.client.Run(context.Background(), request, &response); err != nil {
		fmt.Println("Error while deleting example:", err)
	return false, err
	}	
	
	return response.DeleteExample, nil
}

func (q *QueriesHandlerQL) SendAddTranslationMutation(ctx context.Context, input model.FullRecordInput) (bool, error) {
    request := graphql.NewRequest(`
        mutation ($input: FullRecordInput!) {
            addTranslation(input: $input)
        }
    `)

    request.Var("input", input)

	var response AddTranslationResponse

    err := q.client.Run(context.Background(), request, &response)
	if err != nil {
        fmt.Println("Error while adding translation:", err)
        return false, err
    }

    return response.AddTranslation, nil
}

func (q *QueriesHandlerQL) SendAddExampleMutation(ctx context.Context, input model.FullRecordInput) (bool, error) {
    request := graphql.NewRequest(`
        mutation ($input: FullRecordInput!) {
            addExample(input: $input)
        }
    `)

    request.Var("input", input)

	var response AddExampleResponse

    err := q.client.Run(context.Background(), request, &response)
	if err != nil {
        fmt.Println("Error while adding example sentences:", err)
        return false, err
    }

    return response.AddExample, nil
}




