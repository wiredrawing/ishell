package main

import (
    "fmt"

    "github.com/c-bata/go-prompt"
)

func completer(in prompt.Document) []prompt.Suggest {
    s := []prompt.Suggest{
        {Text: "users", Description: "Store the username and age"},
        {Text: "articles", Description: "Store the article text posted by user"},
        {Text: "comments", Description: "Store the text commented to articles"},
        {Text: "groups", Description: "Combine users with specific rules"},
    }
    return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	for {
		in := prompt.Input(">>> ", completer, prompt.OptionTitle("sql-prompt"))
		fmt.Println("Your input: " + in);
	}

}