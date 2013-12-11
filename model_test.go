package gomarkov

import (
    _ "fmt"
    "strings"
    "testing"
)

func TestTokenize(t *testing.T) {
    str := "1 2 3 4  5      abcd \n def"
    correct := []string{"1", "2", "3", "4", "5", "abcd", "def"}
    pass := true
    for i, s := range Tokenize(str) {
        if s != correct[i] {
            pass = false
        }
    }
    if !pass {
        t.Error("Expected " + strings.Join(correct, " ") + ", got " + strings.Join(Tokenize(str), " "))
    }
    str2 := "John"
    correct2 := []string{"John"}
    if Tokenize(str2)[0] != correct2[0] {
        t.Error("Expected " + strings.Join(correct2, " ") + ", got " + strings.Join(Tokenize(str2), " "))
    }
}

func TestBuildModel(t *testing.T) {
    test_data := []string{"John", "and", "Bob", "went", "to", "a", "store", "and",
        "Bob", "went", "to", "school", "with", "Mr.", "John."}
    test_model1 := BuildModel(test_data, 1)
    if test_model1.Chain["to"]["a"] != 0.5 {
        t.Errorf("Probability of [a] from [to] should have been 0.5, not %v", test_model1.Chain["to"]["a"])
        t.Error(test_model1)
    }
    test_model2 := BuildModel(test_data, 2)
    if test_model2.Chain["went to"]["a"] != 0.5 {
        t.Errorf("Probability of [a] from [went to] should have been 0.5, not %v", test_model2.Chain["went to"]["school"])
        t.Error(test_model2)
    }
}

func TestMostFrequent(t *testing.T) {
    test_data := []string{"John", "and", "Bob", "went", "to", "a", "store", "but",
        "Boey", "went", "to", "school", "with", "Mr.", "John."}
    test_model1 := BuildModel(test_data, 2)
    if test_model1.MostFrequent() != "went to" {
        t.Errorf("The most frequent key should be 'went to', not '%v'", test_model1.MostFrequent())
    }
}

func TestRespond(t *testing.T) {
    test_data := []string{"John", "and", "Bob", "went", "to", "a", "store", "and",
        "Bob", "went", "to", "school", "with", "Mr.", "John."}
    test_model1 := BuildModel(test_data, 1)
    if test_model1.Respond("with", 1) != "with Mr. John." {
        t.Errorf("Expected response of '%v', got '%v'", "with Mr. John.", test_model1.Respond("with", 1))
    }
    test_model1.Respond("John ", 1)
    test_model2 := BuildModel(test_data, 2)
    test_model2.Respond("John and Bob", 10)
    test_model3 := BuildModel(test_data, 4)
    test_model3.Respond("to school with", 5)
}
