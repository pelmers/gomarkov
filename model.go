package gomarkov

import (
    _ "fmt"
    "math/rand"
    "strings"
    "time"
)

var titles map[string]bool

type MarkovModel struct {
    Chain       map[string]map[string]float64
    Frequencies map[string]int
    Order       int
}

func init() {
    // seed random
    rand.Seed(time.Now().UnixNano())
    // some common titles
    titles = map[string]bool{
        "Mr.": true, "Mrs.": true, "Dr.": true, "M.": true, "Gov.": true,
        "Rep.": true, "Sen.": true, "Gen.": true, "St.": true, "Pres.": true,
    }
}

// Return a slice of tokens from some text.
func Tokenize(text string) []string {
    // maybe should add more logic to this?
    // Fields simply splits on whitespace
    return strings.Fields(text)
}

// Add an occurrence of state from key to the Markov chain
func addKey(chain map[string]map[string]float64, key, state string) {
    if _, exists := chain[key]; !exists {
        // allocate the map if it's not already there
        chain[key] = make(map[string]float64)
    }
    chain[key][state] += 1
}

// Normalize the frequencies of states into percentages
func normalizePct(chain map[string]map[string]float64,
    frequencies map[string]int) {
    for key := range chain {
        for state := range chain[key] {
            chain[key][state] /= float64(frequencies[key])
        }
    }
}

// Add the new tokens to the model
func (model *MarkovModel) Update(newtokens []string) {
    // turn the original percentages back into frequencies
    for key := range model.Chain {
        for state := range model.Chain[key] {
            model.Chain[key][state] *= float64(model.Frequencies[key])
        }
    }
    // add in the new tokens
    for index := range newtokens {
        if index < model.Order {
            continue
        }
        key := strings.Join(newtokens[index-model.Order:index], " ")
        addKey(model.Chain, key, newtokens[index])
        model.Frequencies[key] += 1
    }
    // re-normalize the frequencies into percentages
    normalizePct(model.Chain, model.Frequencies)
}

// Return a pointer to a MarkovModel of given order built from given tokens.
func BuildModel(tokens []string, order int) *MarkovModel {
    model := &MarkovModel{make(map[string]map[string]float64),
        make(map[string]int), order}
    model.Update(tokens)
    return model
}

// Return a random key from the Markov model.
func (model *MarkovModel) randKey() string {
    rkey := rand.Intn(len(model.Chain))
    for key := range model.Chain {
        if rkey == 0 {
            return key
        } else {
            rkey--
        }
    }
    panic("Unreachable area reached in model.randKey")
}

// Return the next word from the Markov model given the previous state.
func (model *MarkovModel) nextWord(previous string) string {
    if _, exists := model.Chain[previous]; !exists {
        // not in the chain, so randomly pick a key
        previous = model.randKey()
    }
    pct_threshold := rand.Float64()
    for next := range model.Chain[previous] {
        if model.Chain[previous][next] > pct_threshold {
            return next
        }
        pct_threshold -= model.Chain[previous][next]
    }
    panic("Unreachable code reached in model.nextWord")
}

// Return the key that occurs most frequently in the model
func (model *MarkovModel) MostFrequent() string {
    var max string
    for key := range model.Frequencies {
        if model.Frequencies[key] > model.Frequencies[max] {
            max = key
        }
    }
    return max
}

// Return a string of n_sentences sentences generated from the Markov model
// using the given prompt.
func (model *MarkovModel) Respond(prompt string, n_sentences int) string {
    // start off the response with a prompt
    response := Tokenize(prompt)
    if len(response) < model.Order {
        // the prompt wasn't long enough, so add some words to it
        for len(response) < model.Order {
            response = append(response, model.nextWord(prompt))
        }
    }
    for s := 0; s < n_sentences; s++ {
        word_count := 0
        // add more words until the last character is a period
        // but wait, what if it's a word like "Mr."? also check for titles
        for response[len(response)-1][len(response[len(response)-1])-1] != '.' ||
            titles[response[len(response)-1]] || word_count == 0 {
            // create a new prompt from the last Order words
            prompt = strings.Join(response[len(response)-model.Order:len(response)], " ")
            // call it a sentence if there's been 100 words already
            if word_count > 100 {
                break
            }
            // use the model to pick the next word
            response = append(response, model.nextWord(prompt))
            word_count++
        }
    }
    return strings.Join(response, " ")
}
