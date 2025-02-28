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
	if(strings.ToLower(commandSplitted[0])=="create"){

		
		if !CheckCreateSyntax(command){
			fmt.Println("Incorrect command")
			return true
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
		success := handler.SendCreateMutation(ctx, input)
		
		if !success{
			fmt.Println("Given word already exists in the dictionary")
			return true
		}

		fmt.Println("Word successfully added to dictionary")

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
}

func (r *ReceiveCommandHandler) HandleCommand(command string) bool{
	// TODO make NewClient request one time in other func
	client := graphql.NewClient("http://localhost:8080/query")
	handler := &QueriesHandler{client: client}

	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="receive"){
		if(len(commandSplitted)!=2){
			log.Fatal("Wrong command")
		}
		ctx := context.Background()
		received, err := handler.SendReceiveMutation(ctx, commandSplitted[1])
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				fmt.Println("The word dosen't exist in the dictionary")
				return true			}
			log.Fatal("Error:", err)
		}else if received!=""{
			fmt.Println(received)
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
}

func (m *ModifyCommandHandler) HandleCommand(command string) bool{
	// TODO make NewClient request one time in other func
	client := graphql.NewClient("http://localhost:8080/query")
	handler := &QueriesHandler{client: client}
	
	commandSplitted := strings.Split(command, " ")
	if(strings.ToLower(commandSplitted[0])=="modify"){
		if(len(commandSplitted)<4){
			fmt.Println("Wrong command")
		}
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
		success, err = handler.SendRemoveTranslationMutation(context.Background(), commandSplitted[3])
	}else if (strings.ToLower(commandSplitted[2])=="example"){
		sentence := ParseQuery(command)
		if len(sentence) !=1{
			fmt.Println("Wrong command. You should specify only one example and use brackets")
			return
		}
		success, err = handler.SendRemoveExampleMutation(context.Background(), sentence[0])

	}else{
		fmt.Println("Wrong command. You can on delete example or translation while modyfing word")
		return
	}

	if err != nil {
		log.Fatal("Error:", err)
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
		if len(commandSplitted)<5 {
			fmt.Println("Wrong command. Not enugh input arguments for adding new translation")
			return
		}
		sentences := ParseQuery(command)
		input := model.CreateTranslationInput{
			Word:        commandSplitted[3],
			Translation: commandSplitted[4],
			Examples:    sentences,
		}
		success, err = handler.SendAddTranslationMutation(context.Background(), input)
	} else if strings.ToLower(commandSplitted[2])=="example"{
		if len(commandSplitted)<4{
			fmt.Println("Wrong command. Not enugh input arguments for adding new examples")
			return
		}
		sentences := ParseQuery(command)
		success, err = handler.SendAddExampleMutation(context.Background(), commandSplitted[3], sentences)
	}else {
		fmt.Println("Wrong command. You can on delete example or translation while modyfing word")
		return
	}

	if err != nil {
		log.Fatal("Error:", err)
	} else if success{
		fmt.Println("Word modified successfully")
	}else{
		fmt.Println("Example/Translation dosen't exist in the dictionary")
	}
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
	if(strings.ToLower(commandSplitted[0])=="remove"){
		if(len(commandSplitted)!=2){
			log.Fatal("Wrong command")
		}

		ctx := context.Background()
		success, err := handler.SendRemoveMutation(ctx, commandSplitted[1])
		if err != nil {
			log.Fatal("Error:", err)
		} else if success{
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

