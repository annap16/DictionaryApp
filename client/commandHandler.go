package main

import(
	"strings"
	"fmt"
	"context"
	"dictionary-app/server/graph/model"
)

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
	if(strings.ToLower(commandSplitted[0])=="dodaj"){
		if !CheckCreateSyntax(command){
			fmt.Println("Niepoprawne polecenie")
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
			fmt.Println("Wystąpił błąd podczas dodawania do bazy danych:", err)
		}else if !success{
			fmt.Println("Podane słowo istnieje już w słowniku")
		}else{
			fmt.Println("Słowo zostało dodane do słownika")
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
	if(strings.ToLower(commandSplitted[0])=="sprawdź"){
		if(len(commandSplitted)!=2){
			fmt.Println("Niepoprawne polecenie")
			return true
		}
		ctx := context.Background()
		received, err := r.handler.SendReceiveMutation(ctx, strings.ToLower(commandSplitted[1]))
		if err!=nil{
			fmt.Println("Wystąpił błąd podczas wyszukiwania w słowniku:", err)
		}else if received!=""{
			fmt.Println(received)
		}else{
			fmt.Println("Podane słowo nie istnieje w słowniku")
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
	modifyFactory CommandFactory
}

func (m *ModifyCommandHandler) HandleCommand(command string) bool{	
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="modyfikuj"){
		var success bool
		modifyAction, err := m.modifyFactory.CreateAction(m.handler, command)

		if modifyAction!= nil{
			success, err = modifyAction.Execute()
		}else if err!=nil{
			fmt.Println(err)
			return true
		}

		if success{
			fmt.Println("Słowo zostało zmodyfikowane poprawnie")
		}else if err==nil{
			fmt.Println("Wprowadzono niepoprawne polecenie")
		}else{
			fmt.Println(err)
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
	if(strings.ToLower(commandSplitted[0])=="usuń"){
		if(len(commandSplitted)!=2){
			fmt.Println("Niepoprawne polecenie")
			return true
		}
		ctx := context.Background()
		success, err := r.handler.SendRemoveMutation(ctx, commandSplitted[1])
		if err!=nil{
			fmt.Println("Wystąpił błąd podczas usuwania słowa ze słownika:", err)
		}else if success{
			fmt.Println("Podane słowo i powiązane z nim dane zostały usunięte")
		}else{
			fmt.Println("Podane słowo nie istnieje w słowniku")
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
