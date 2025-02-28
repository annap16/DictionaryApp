package main

import(

)

type ReceiveResponse struct {
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

type AddWordResponse struct {
	CreateTranslation bool `json:"createTranslation"`
}

type RemoveWordResponse struct {
	DeleteWord bool `json:"deleteWord"`
}

type RemoveTranslationResponse struct {
	DeleteTranslation bool `json:"deleteTranslation"`
}

type RemoveExampleResponse struct {
	DeleteExample bool `json:"deleteExample"`
}