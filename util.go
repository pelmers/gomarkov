package gomarkov

// utilities for making markov models

import (
    "net/http"
    "net/http/httputil"
    "strings"
    "io/ioutil"
)

func TokenizeFromGutenberg(url string) ([]string, error) {
    var tokens []string
    resp, err := http.Get(url)
    if err != nil {
        return tokens, err
    }
    body, err := httputil.DumpResponse(resp, true)
    if err != nil {
        return tokens, err
    }
    text := string(body)
    start := strings.Index(text, "START")
    end := strings.LastIndex(text, " END OF")
    if start == -1 || end == -1 {
        tokens = Tokenize(text)
    } else {
        tokens = Tokenize(text[start:end])
    }
    return tokens, err
}

func TokenizeFromFile(path string) ([]string, error) {
    var tokens []string
    body, err := ioutil.ReadFile(path)
    if err != nil {
        return tokens, err
    }
    tokens = Tokenize(string(body))
    return tokens, err
}
