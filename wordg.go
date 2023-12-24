// wordg.go - Implement the game "Wordle".  This is a typical
// word guessing game.  During the game, Side 1 choose a 5-letter word.
// Side 2 makes guesses about the word.
// Each guess must be a valid five-letter word.
// The letters of the guess entered by Side 2 will be highlighted by Side 1 to indicate
// the accuracy of the guess.
// Below we describe the highlighting done by Wordle, plus in parentheses the encoding
// we use:
// A letter will be shown in green (y) if the letter is in the word, and it is in the correct spot.
// A letter will be shown in yellow (p) if the letter is in the word, but it is not in the correct spot.
// A letter will be shown in gray (n) if the letter is not in the word.
//
// Mark Riordan  2023-12-06

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type ParseResult int

// This is a weird way of creating kind of an enum.
// The constants are for the different program modes.
const (
	BAD ParseResult = iota
	RUN
	GUESS
)

func usage() {
	var usageMsg = []string{
		"wordg: Program to play Wordle.",
		"Usage: wordg {--run | --guess }",
		"where:",
		"--run   specifies that the program should think of a word and let you guess it.",
		"--guess specifies that the program should makes guesses about a word some",
		"        other entity is thinking of.",
	}
	for _, line := range usageMsg {
		fmt.Println(line)
	}
}

func parseCmdLine() ParseResult {
	result := BAD
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 1 {
		arg1 := argsWithoutProg[0]
		if arg1 == "--run" {
			result = RUN
		} else if arg1 == "--guess" {
			result = GUESS
		}
	}
	return result
}

func isKnownWord(word string) bool {
	found := false
	for _, knownWord := range AllWords {
		if knownWord == word {
			found = true
			break
		}
	}
	return found
}

func runGame() {
	word := AllWords[rand.Intn(len(AllWords))]
	//fmt.Println("The word is " + word)
	scanner := bufio.NewScanner(os.Stdin)
	for running := true; running; {
		fmt.Print(" Guess: ")
		scanner.Scan()
		guess := scanner.Text()
		if "q" == guess {
			fmt.Println("The word was " + word)
			break
		} else if len(guess) != 5 {
			fmt.Println("Guesses must be exactly 5 lowercase letters")
		} else {
			// The guess must be a known word
			if !isKnownWord(guess) {
				fmt.Println(guess + " is not a valid word")
			} else {
				response := [5]string{" ", " ", " ", " ", " "}
				// First, scan for the correct letters in the correct places.
				// We need to have this information to later determine whether
				// a given letter that matches a letter in a different position
				// is a "p" or "n".
				for j := 0; j < len(guess); j++ {
					guessCh := guess[j : j+1]
					//fmt.Println("Looking at char " + ch + " " + response)
					wordCh := word[j : j+1]
					if guessCh == wordCh {
						response[j] = "y"
					}
				}
				for j := 0; j < len(guess); j++ {
					guessCh := guess[j : j+1]
					//fmt.Println("Looking at char " + ch + " " + response)
					wordCh := word[j : j+1]
					if guessCh != wordCh {
						// Iterate through the correct word, to see if this char
						// is found elsewhere in the word.
						found := false
						for k := 0; k < len(word); k++ {
							if k != j {
								if guessCh == word[k:k+1] && response[k] != "y" {
									// The guessed char is in the word, and not at
									// a position that is a correct guess.
									found = true
								}
							}
						}
						if found {
							response[j] = "p"
						} else {
							response[j] = "n"
						}
					}
				}
				responseStr := strings.Join(response[:], "")
				fmt.Println("Result: " + responseStr)
				if responseStr == "yyyyy" {
					fmt.Println("Congratulations!")
					running = false
				}
			}
		}
	}
	// for _, element := range allWords {
	// 	fmt.Println(element)
	// }
}

func doGuesses() {
	//validLetters := [5][26]string
}

func main() {
	ParseResult := parseCmdLine()
	if ParseResult == BAD {
		usage()
	} else if ParseResult == GUESS {
		fmt.Println("Not yet implemented.")
		//doGuesses()
	} else if ParseResult == RUN {
		runGame()
	}
}
