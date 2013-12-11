package gomarkov

import (
    _ "fmt"
    "testing"
)

func TestTokenizeFromGutenberg(t *testing.T) {
    url := "http://www.gutenberg.org/ebooks/7700.txt.utf-8"
    tokens, err := TokenizeFromGutenberg(url)
    if err != nil {
        t.Errorf("Tokenization failed with error: %v", err)
    } else {
        test_model := BuildModel(tokens, 2)
        test_model.Respond("By the", 2)
        if test_model.MostFrequent() != "of the" {
            t.Errorf("Most common chain in 7700 is 'of the', not '%v'",
                test_model.MostFrequent())
        }
    }
}
