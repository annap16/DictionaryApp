# Dictionary Application
## Table of Contents
- [Application Overview](#application-overview)
- [Setup and Usage Instructions](#setup-and-usage-instructions)
- [Command Schema in Command-Line Application](#command-schema-in-command-line-appliction)
  - [Available Commands](#available-commands)
- [Database Design](#database-design)
- [GraphQL Queries and Mutation](#graphql-queries-and-mutation)
- [Example Client Commands](#example-client-commands)
## Application Overview
The application consists of a client and a server. The server is responsible for handling communication with the database, processing translations concurrently, and implementing other features as described in the recruitment assignment. The client is a simple command-line application that parses and interprets user input. It allows users to interact with the application easily by validating input and converting it into GraphQL queries.  

All commands in the client are written in Polish, as this is a Polish-to-English dictionary rather than an English-to-Polish one. The application includes unit tests covering all crucial functions.  

## Setup and Usage Instructions
The application runs using Docker for both the client and server. Follow these steps to set it up:
1. Download the project and create a `.env` file in the main project directory using provided `env.example` file with example account.
2. Start Docker by running: <pre> ```sudo docker-compose up --build ``` </pre>
3.  In a separate terminal, start the client by executing:<pre> ```sudo docker-compose exec client /bin/sh ``` </pre> Then, once inside, run:  <pre> ```/dictionary-client ``` </pre>
4.  After successfully launching the client, you can interact with the application using example commands and leave the application at any time using `exit` command.
> [!NOTE]   
> Even without launching the client, you can still test and interact with the backend directly using tools like GraphQL Playground, Postman, or any other API testing tool that supports GraphQL.

## Command Schema in Command-Line Application
Commands such as `dodaj`, `usuń`, and `modyfikuj` are **not case-sensitive** in the client application.  
### Available Commands
- **Add a new word with a translation and examples:** <pre>```dodaj {word_in_polish} {word_in_english} [Example one] ... [Example n]``` </pre>  The list of examples can be empty.
- **Retrieve translations of an existing word:*** <pre>```sprawdź {word_in_polish}``` </pre>
- **Modify an existing translation by adding new example sentences:** <pre> ```modyfikuj dodaj przykład {word_in_polish} {word_in_english} [New example one] [New example two] ... [New example n] ```</pre> At least one example must be provided.
- **Add a new translation for an existing word (examples are optional):** <pre>```modyfikuj dodaj tłumaczenie {word_in_polish} {word_in_english} [Example one] [Example two] ... [Example n]```</pre> The list of examples can be empty.
- **Remove specific example sentences from a translation:** <pre>```modyfikuj usuń przykład {word_in_polish} {word_in_english} [Example one] [Example two] ... [Example n] ```</pre> At least one example must be specified.
- **Delete a specific translation of a word:** <pre>```modyfikuj usuń tłumaczenie {word_in_polish} {translation_to_delete} ```</pre>
- **Delete an entire word and all its translations:** <pre>```usuń {word_in_polish}```</pre>

## Database Design
Below is the entity-relationship diagram (ERD) representing the database tables:
![Database](https://github.com/user-attachments/assets/29584b36-660e-42da-a6ce-2ceb69535a8a)
The following is the database schema, including indexes omitted in image.
```sql
-- Database Schema

CREATE TABLE Word (
  ID INTEGER PRIMARY KEY,
  Word TEXT UNIQUE NOT NULL
);

CREATE TABLE Translation (

  ID INTEGER PRIMARY KEY,
  Translation TEXT NOT NULL,
  WordID INTEGER NOT NULL,
  UNIQUE (Translation, WordID),
  FOREIGN KEY (WordID) REFERENCES Word(ID)
);

CREATE TABLE ExampleSentence (
  ID INTEGER PRIMARY KEY,
  Sentence TEXT NOT NULL,


  TranslationID INTEGER NOT NULL,
  UNIQUE (Sentence, TranslationID),
  FOREIGN KEY (TranslationID) REFERENCES Translation(ID)
);
```
## GraphQL Queries and Mutation
Examples of GraphQL queries and mutations are provided in the file: [GraphQLexamples.txt](https://github.com/user-attachments/files/19872124/GraphQLexamples.txt)

## Example Client Commands
- **Adding a new word:** <pre>```dodaj książka book [I love my new book] [This is a book]```</pre>
- **Retrieving an existing word:** <pre>```sprawdź książka```</pre>
- **Modifying an existing word - adding an example sentence:** <pre>```modyfikuj dodaj przykład książka book [New example] [New example nr 2]```</pre>
- **Modifying an existing word - adding a translation (with examples):** <pre>```modyfikuj dodaj tłumaczenie książka tome [This is a tome]```</pre>
- **Modifying an existing word - removing an example sentence:** <pre>```modyfikuj usuń przykład książka book [New example] [New example nr 2]```</pre>
- **Modifying an existing word - removing a translation:** <pre>```modyfikuj usuń tłumaczenie książka tome```</pre>
- **Deleting an existing word:** <pre>```usuń książka```</pre>

