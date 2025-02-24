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
			log.Fatal("Error:", err)
		}

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

type ReciveCommandHandler struct{
	next CommandHandler
}

func (r *ReciveCommandHandler) HandleCommand(command string) bool{
	commandSplitted := strings.Split(command, " ")
	if(commandSplitted[0]=="recive"){
		fmt.Println("recive")
		return true
	}
	if(r.next!=nil){
		return r.next.HandleCommand(command)
	}
	return false
}

func (r *ReciveCommandHandler) SetNext(handler CommandHandler) {
	r.next = handler
}

type ModifyCommandHandler struct{
	next CommandHandler
}

func (m *ModifyCommandHandler) HandleCommand(command string) bool{
	commandSplitted := strings.Split(command, " ")
	if(commandSplitted[0]=="modify"){
		fmt.Println("modify")
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
	commandSplitted := strings.Split(command, " ")
	if(commandSplitted[0]=="remove"){
		fmt.Println("remove")
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

