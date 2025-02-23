package main

import(
	"strings"
	"fmt"
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
	commandSplitted := strings.Split(command, " ")
	if(commandSplitted[0]=="create"){
		fmt.Println("create")
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

