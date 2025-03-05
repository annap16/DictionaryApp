package main

import(
	"strings"
	"fmt"
	"context"
	"regexp"
	"github.com/annap16/DictionaryApp/graph/model"
)

func ParseQuery(query string) ([]string) {
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(query, -1)

	var sentences []string
	for _, match := range matches {
		sentences = append(sentences, match[1]) 
	}

	return sentences
}


//Implementing responsibility chain pattern for commands handling

type CommandHandler interface{
	HandleCommand(command string) bool
	SetNext(handler CommandHandler)
}

type CreateCommandHandler struct{
	next CommandHandler
	handler QueriesHandler
}

func (c *CreateCommandHandler) HandleCommand(command string) bool{
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="create"){
		if !CheckCreateSyntax(command){
			fmt.Println("Incorrect command")
			return true
		}

		word := strings.ToLower(commandSplitted[1])
		translation := strings.ToLower(commandSplitted[2]) 
		sentences := ParseQuery(command)

		input := model.FullRecordInput{
			Word:        word,
			Translation: translation,
			Examples:    sentences,
		}

		ctx := context.Background()
		success, err := c.handler.SendCreateMutation(ctx, input)
		if err != nil {
			fmt.Println("Error occured while adding word to the dictionary")
		}else if !success{
			fmt.Println("Given word already exists in the dictionary")
		}else{
			fmt.Println("Word successfully added to dictionary")
		}

		return true
	}
	if(c.next!=nil){
		return c.next.HandleCommand(command)
	}
	return false
}


func (c *CreateCommandHandler) SetNext(handler CommandHandler) {
	c.next = handler
}

type ReceiveCommandHandler struct{
	next CommandHandler
	handler QueriesHandler
}

func (r *ReceiveCommandHandler) HandleCommand(command string) bool{
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="receive"){
		if(len(commandSplitted)!=2){
			fmt.Println("Wrong command")
			return true
		}
		ctx := context.Background()
		received, err := r.handler.SendReceiveMutation(ctx, strings.ToLower(commandSplitted[1]))
		if err!=nil{
			fmt.Println("Error while receiving a word")
		}else if received!=""{
			fmt.Println(received)
		}else{
			fmt.Println("The word doesn't exist in the dictionary")
		}
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
	handler QueriesHandler
}

func (m *ModifyCommandHandler) HandleCommand(command string) bool{	
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="modify"){
		if(len(commandSplitted)<3){
			fmt.Println("Wrong command")
		}
		modifyAction := ModifyCommandFactory(m.handler, command)
		var success bool
		if modifyAction!= nil{
			success = modifyAction.Execute()
		}
		if success{
			fmt.Println("Word modified successfully")
		}else{
			fmt.Println("Something went wrong")
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
	handler QueriesHandler
}

func (r *RemoveCommandHandler) HandleCommand(command string) bool{
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="remove"){
		if(len(commandSplitted)!=2){
			fmt.Println("Wrong command")
			return true
		}
		ctx := context.Background()
		success, err := r.handler.SendRemoveMutation(ctx, commandSplitted[1])
		if err!=nil{
			fmt.Println("An error occured while removing word from the dictionary")
		}else if success{
			fmt.Println("Word and related data deleted successfully")
		}else{
			fmt.Println("The word dosen't exist in the dictionary")
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
