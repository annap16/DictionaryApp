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
	client *graphql.Client
}

func (c *CreateCommandHandler) HandleCommand(command string) bool{
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="create"){
		if !CheckCreateSyntax(command){
			fmt.Println("Incorrect command")
			return true
		}

		handler := &QueriesHandler{client: c.client}
		word := strings.ToLower(commandSplitted[1])
		translation := strings.ToLower(commandSplitted[2]) 
		sentences := ParseQuery(command)

		input := model.CreateTranslationInput{
			Word:        word,
			Translation: translation,
			Examples:    sentences,
		}

		ctx := context.Background()
		success, err := handler.SendCreateMutation(ctx, input)
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

func CheckCreateSyntax(command string) bool{
	pattern := `(?i)^create\s+(\S+)\s+(\S+)(\s+\[.*?\])*?$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(command)
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
	client *graphql.Client
}

func (r *ReceiveCommandHandler) HandleCommand(command string) bool{
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="receive"){
		if(len(commandSplitted)!=2){
			fmt.Println("Wrong command")
			return true
		}
		handler := &QueriesHandler{client: r.client}
		ctx := context.Background()
		received, err := handler.SendReceiveMutation(ctx, strings.ToLower(commandSplitted[1]))
		if err!=nil{
			fmt.Println("Error while receving a word")
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
	client *graphql.Client
}

func (m *ModifyCommandHandler) HandleCommand(command string) bool{	
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="modify"){
		if(len(commandSplitted)<3){
			fmt.Println("Wrong command")
		}
		handler := &QueriesHandler{client: m.client}
		modifyType := strings.ToLower(commandSplitted[1])
		if modifyType=="delete"{
			m.handleUpdateDeleteCommand(command, commandSplitted, handler)
		} else if modifyType=="add"{
			m.handleUpdateAddCommand(command, commandSplitted, handler)
		} else{
			fmt.Println("Wrong command")
		}
		return true
	}

	if(m.next!=nil){
		return m.next.HandleCommand(command)
	}
	return false
}

func (m *ModifyCommandHandler) handleUpdateDeleteCommand(command string, commandSplitted []string, handler *QueriesHandler){
	var success bool
	var err error

	if(strings.ToLower(commandSplitted[2])=="translation"){
		if len(commandSplitted)!=4{
			fmt.Println("Wrong command")
		}
		success, err = handler.SendRemoveTranslationMutation(context.Background(), commandSplitted[3])
	}else if (strings.ToLower(commandSplitted[2])=="example"){
		sentence := ParseQuery(command)
		if len(sentence) !=1 {
			fmt.Println("Wrong command. You should specify only one example and use brackets")
			return
		}
		success, err = handler.SendRemoveExampleMutation(context.Background(), strings.ToLower(commandSplitted[3]), sentence[0])
	}else{
		fmt.Println("Wrong command")
	}

	if err != nil {
		log.Println("Error while modyfing word:", err)
	} else if success{
		fmt.Println("Word modified successfully")
	}else{
		fmt.Println("Example/Translation dosen't exist in the dictionary")
	}

}

func (m *ModifyCommandHandler) handleUpdateAddCommand(command string, commandSplitted []string, handler *QueriesHandler){
	var success bool
	var err error
	if strings.ToLower(commandSplitted[2])=="translation"{
		if !CheckAddTranslationSyntax(command){
			fmt.Println("Wrong command")
			return
		}
		sentences := ParseQuery(command)
		input := model.CreateTranslationInput{
			Word:        strings.ToLower(commandSplitted[3]),
			Translation: strings.ToLower(commandSplitted[4]),
			Examples:    sentences,
		}
		success, err = handler.SendAddTranslationMutation(context.Background(), input)
	} else if strings.ToLower(commandSplitted[2])=="example"{
		if !CheckAddExampleSyntax(command){
			fmt.Println("Wrong command syntax")
			return 
		}
		sentences := ParseQuery(command)
		success, err = handler.SendAddExampleMutation(context.Background(), commandSplitted[3], sentences)
	}else {
		fmt.Println("Wrong command. You can only delete example or translation while modyfing word")
		return
	}

	if err != nil {
		log.Fatal("Error:", err)
	} else if success{
		fmt.Println("Word modified successfully")
	}else{
		fmt.Println("Couldn't add examples/translation because it already exists or couldnt be found")
	}
}

func CheckAddExampleSyntax(command string) bool{
	pattern := `(?i)^modify\s+add\s+example\s+(\S+)(\s+\[.*?\])+$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(command)
}

func CheckAddTranslationSyntax(command string) bool{
	pattern := `(?i)^modify\s+add\s+translation\s+(\S+)\s+(\S+)(\s+\[.*?\])*$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(command)
}


func (m *ModifyCommandHandler) SetNext(handler CommandHandler) {
	m.next = handler
}


type RemoveCommandHandler struct{
	next CommandHandler
	client *graphql.Client
}

func (r *RemoveCommandHandler) HandleCommand(command string) bool{
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="remove"){
		if(len(commandSplitted)!=2){
			fmt.Println("Wrong command")
			return true
		}
		handler := &QueriesHandler{client: r.client}
		ctx := context.Background()
		success, err := handler.SendRemoveMutation(ctx, commandSplitted[1])
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

