package main

import (
    model "github.com/pelmers/gomarkov"
    "bufio"
    "os"
    "fmt"
)

// Given a string slice of options, return the index of the user's choice.
func GetChoice(options []string) int {
    var choice int
    for i, opt := range options {
        fmt.Printf("%d) %v\n", (i+1), opt)
    }
    fmt.Print("Please choose an option: ")
    _, err := fmt.Scan(&choice)
    if err != nil || choice > len(options) || choice < 1 {
        fmt.Println("An error occurred. Please retry.")
        return GetChoice(options)
    }
    return choice
}

func MakeTokens(tokenizer func(string) ([]string, error), prompt string) []string {
    var tokens []string
    fmt.Printf("%s: ", prompt)
    _, err := fmt.Scan(&prompt)
    if err != nil {
        fmt.Println("An error occurred reading the input URL.")
        return tokens
    }
    tokens, err = tokenizer(prompt)
    if err != nil {
        fmt.Println(err)
    }
    return tokens
}

func main() {
    menu := []string{"Parse from URL",
    "Parse from file",
    "Generate sentences",
    "Clear books",
    "Exit"}
    choice := 0
    generator := model.BuildModel(make([]string, 0), 2)
    reader := bufio.NewReader(os.Stdin)
    for choice != 5 {
        choice = GetChoice(menu)
        switch choice {
        case 1:
            generator.Update(MakeTokens(model.TokenizeFromGutenberg, "URL"))
        case 2:
            generator.Update(MakeTokens(model.TokenizeFromFile, "Path"))
        case 3:
            fmt.Print("Prompt: ")
            prompt, _ := reader.ReadString('\n')
            fmt.Println(generator.Respond(prompt, 4))
        case 4:
            generator = model.BuildModel(make([]string, 0), 2)
        }
    }
}
