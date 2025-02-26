package main

import(
	"strings"
	"fmt"
	"context"
	"log"
	"regexp"
	"github.com/annap16/DictionaryApp/graph/model"
	"github.com/machinebox/graphql"

)

//Implementing responsibility chain pattern for commands handling

type CommandHandler interface{
	HandleCommand(command string) bool
	SetNext(handler CommandHandler)
}

type CreateCommandHandler struct{
	next CommandHandler
}

func (c *CreateCommandHandler) HandleCommand(command string) bool{

	// TODO make NewClient request one time in other func
	client := graphql.NewClient("http://localhost:8080/query")
	handler := &QueriesHandler{client: client}

	commandSplitted := strings.Split(command, " ")
	if(commandSplitted[0]=="create"){
		if len(commandSplitted)<3 {
			log.Fatal("Incorrect command")
			return false
		}
		word := commandSplitted[1]
		translation := commandSplitted[2] 
		sentences := ParseQuery(command)

		input := model.CreateTranslationInput{
			Word:        word,
			Translation: translation,
			Examples:    sentences,
		}

		ctx := context.Background()
		err := handler.SendCreateMutation(ctx, input)
		if err != nil {
			log.Fatal("Error after send create mutation:", err)
		}

		fmt.Println("Word successfully added to dictionary")

		return true
	}
	if(c.next!=nil){
		return c.next.HandleCommand(command)
	}
	return false
}

func ParseQuery(query string) ([]string) {
	re := regexp.MustCompile(`\[(.*?)\]`)

	matches := re.FindAllStringSubmatch(query, -1)

	var sentences []string
	for _, match := range matches {
		sentences = append(sentences, match[1]) 
	}

	return sentences
}

func (c *CreateCommandHandler) SetNext(handler CommandHandler) {
	c.next = handler
}

type ReceiveCommandHandler struct{
	next CommandHandler
}

func (r *ReceiveCommandHandler) HandleCommand(command string) bool{
	// TODO make NewClient request one time in other func
	client := graphql.NewClient("http://localhost:8080/query")
	handler := &QueriesHandler{client: client}

	commandSplitted := strings.Split(command, " ")
	if(commandSplitted[0]=="receive"){
		if(len(commandSplitted)!=2){
			log.Fatal("Wrong command")
		}
		ctx := context.Background()
		received, err := handler.SendReceiveMutation(ctx, commandSplitted[1])
		if err != nil {
			log.Fatal("Error:", err)
		}

		fmt.Println(received)

		return true
	}
	if(r.next!=nil){
		return r.next.HandleCommand(command)
	}
	return false
}

func (r *ReceiveCommandHandler) SetNext(handler CommandHandler) {
	r.next = handler
}

type ModifyCommandHandler struct{
	next CommandHandler
}

func (m *ModifyCommandHandler) HandleCommand(command string) bool{
	commandSplitted := strings.Split(command, " ")
	if(commandSplitted[0]=="modify"){
		if(len(commandSplitted)<3){
			log.Fatal("Wrong command")
		}
		modifyType := commandSplitted[1]
		if modifyType=="delete"{
			// TODO
		} else if modifyType=="add"{
			//TODO
		} else{
			log.Fatal("Wrong command")
		}
		return true
	}
	if(m.next!=nil){
		return m.next.HandleCommand(command)
	}
	return false
}

func (m *ModifyCommandHandler) SetNext(handler CommandHandler) {
	m.next = handler
}


type RemoveCommandHandler struct{
	next CommandHandler
}

func (r *RemoveCommandHandler) HandleCommand(command string) bool{
	// TODO make NewClient request one time in other func
	client := graphql.NewClient("http://localhost:8080/query")
	handler := &QueriesHandler{client: client}

	commandSplitted := strings.Split(command, " ")
	if(commandSplitted[0]=="remove"){
		if(len(commandSplitted)!=2){
			log.Fatal("Wrong command")
		}

		ctx := context.Background()
		err := handler.SendRemoveMutation(ctx, commandSplitted[1])
		if err != nil {
			log.Fatal("Error:", err)
		} else{
			fmt.Println("Successfully deleted translation with word:", commandSplitted[1])
		}
		return true
	}
	if(r.next!=nil){
		return r.next.HandleCommand(command)
	}
	return false
}

func (r *RemoveCommandHandler) SetNext(handler CommandHandler) {
	r.next = handler
}

