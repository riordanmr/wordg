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
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type RunType int

// This is a weird way of creating kind of an enum.
// The constants are for the different program modes.
const (
	BAD RunType = iota
	RUN
	GUESS
)

const LETTERS_IN_WORD = 5

var MyScanner bufio.Scanner

// Map: index is a letter, value is the minimum number of occurrences of that
// letter in the word we are trying to guess.  We don't populate with letters
// that we don't yet know are required.
var requiredLetters = make(map[string]int)

type Settings struct {
	runType RunType
	word    string
	errMsg  string
}

func usage() {
	var usageMsg = []string{
		"wordg: Program to play Wordle.",
		"Usage: wordg {--run | --guess } [--word=word]",
		"where:",
		"--run   specifies that the program should think of a word and let you guess it.",
		"--guess specifies that the program should makes guesses about a word some",
		"        other entity is thinking of.",
		"word    applies only to --run mode, and specifies the word the program should",
		"        think of. Optional; the default is for wordg to select aa word randomly.",
	}
	for _, line := range usageMsg {
		fmt.Println(line)
	}
}

func parseCmdLine() Settings {
	var settings Settings
	var run bool
	var guess bool
	flag.BoolVar(&run, "run", false, "Have the program think of a word and make you guess")
	flag.BoolVar(&guess, "guess", false, "Have the program try to guess the word")
	flag.StringVar(&settings.word, "word", "", "The word the program is thinking of in run mode. If not supplied, the program will chose a word at random.")

	flag.Parse()

	if (run && guess) || (!run && !guess) {
		settings.errMsg = "You must specify either --guess or --run"
	} else {
		if run {
			settings.runType = RUN
		} else {
			settings.runType = GUESS
		}
	}
	return settings
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

func readGuessResult() string {
	MyScanner.Scan()
	response := MyScanner.Text()
	return response
}

func runGame(word string) {
	if len(word) == 0 {
		word = AllWords[rand.Intn(len(AllWords))]
	}
	//fmt.Println("The word is " + word)
	for running := true; running; {
		fmt.Print(" Guess: ")
		MyScanner.Scan()
		guess := MyScanner.Text()
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
}

// Define a Set type as a map with a boolean value
type StringSet map[string]bool

// Function to add an element to the set
func (set StringSet) Add(element string) {
	set[element] = true
}

// Function to remove an element from the set
func (set StringSet) Remove(element string) {
	delete(set, element)
}

// Function to remove all elements from the set
func (set StringSet) RemoveAll() {
	clear(set)
}

// Function to check if an element exists in the set
func (set StringSet) Contains(element string) bool {
	return set[element]
}

// Return true if we found the correct word.
func processResponse(validLetters *[LETTERS_IN_WORD]StringSet, myGuess string,
	response string) bool {
	foundAnswer := false
	if response == "yyyyy" {
		foundAnswer = true
	} else if len(response) != LETTERS_IN_WORD {
		fmt.Printf("Response must be of length %v\n", LETTERS_IN_WORD)
	} else {
		// Loop through the letters in the response.
		var charToCountThisGuess map[string]int = make(map[string]int)
		for ipos := 0; ipos < LETTERS_IN_WORD; ipos++ {
			respCh := response[ipos : ipos+1]
			guessCh := myGuess[ipos : ipos+1]
			if respCh == "n" {
				for j := 0; j < LETTERS_IN_WORD; j++ {
					validLetters[j].Remove(guessCh)
				}
			} else if respCh == "y" {
				validLetters[ipos].RemoveAll()
				validLetters[ipos].Add(guessCh)

				_, present := charToCountThisGuess[guessCh]
				if present {
					charToCountThisGuess[guessCh]++
				} else {
					charToCountThisGuess[guessCh] = 1
				}
			} else if respCh == "p" {
				validLetters[ipos].Remove(guessCh)

				_, present := charToCountThisGuess[guessCh]
				if present {
					charToCountThisGuess[guessCh]++
				} else {
					charToCountThisGuess[guessCh] = 1
				}
			} else {
				fmt.Println("Unexpected response char: " + respCh)
			}
		}
		// Now we have accumulated in charToCountThisGuess the info from the response
		// for this guessed word. Apply this knowledge to requiredLetters, which will
		// reflect required letters info from all responses so far.
		for requiredCh, count := range charToCountThisGuess {
			oldCount, present := requiredLetters[requiredCh]
			if present {
				if count > oldCount {
					requiredLetters[requiredCh] = count
				}
			} else {
				requiredLetters[requiredCh] = count
			}
		}
	}
	return foundAnswer
}

func printSetOfValidLetters(validLetters *[LETTERS_IN_WORD]StringSet) {
	// Debug print the set of valid letters for each position.
	for k := 0; k < len(validLetters); k++ {
		fmt.Print(k, " ")
		msg := ""
		for idx := 0; idx < 26; idx++ {
			ch := "abcdefghijklmnopqrstuvwxyz"[idx : idx+1]
			if validLetters[k][ch] {
				msg += ch
			}
		}
		fmt.Println(msg)
	}
}

func makeMapFromWord(word string) map[string]int {
	mapLetterToCount := make(map[string]int)

	for j := 0; j < len(word); j++ {
		ch := word[j : j+1]
		count, present := mapLetterToCount[ch]
		if present {
			mapLetterToCount[ch] = count + 1
		} else {
			mapLetterToCount[ch] = 1
		}
	}

	return mapLetterToCount
}

func doGuesses() {
	// Define an array of sets, one for each position in the word being guessed.
	// Initially populate each set with all possible letters.
	fmt.Println(("doGuesses here"))
	var validLetters [LETTERS_IN_WORD]StringSet
	for i, _ := range validLetters {
		validLetters[i] = make(map[string]bool)
	}
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	for idx := 0; idx < len(validLetters); idx++ {
		for ia := 0; ia < len(alphabet); ia++ {
			validLetters[idx].Add(alphabet[ia : ia+1])
		}
	}

	var response string = ""
	for {
		//printSetOfValidLetters(&validLetters)
		var myGuess string
		// Loop through the list of words, finding the first one
		// that matches the clues we have so far.
		for _, guess := range AllWords {
			matches := true
			// Loop through the letters of this guess.
			for ilet := 0; ilet < len(guess); ilet++ {
				if !validLetters[ilet].Contains(guess[ilet : ilet+1]) {
					// This potential guess is incompatible with the clues so far,
					// so stop analyzing this potential guess.
					matches = false
					break
				}
			}
			if matches {
				// The word matches according to validLetters, but does it have
				// all the letters we know are in the word?
				mapLetterToCountThisWord := makeMapFromWord(guess)
				for letter, numRequired := range requiredLetters {
					countThisGuess, present := mapLetterToCountThisWord[letter]
					if !present {
						matches = false
					} else if countThisGuess < numRequired {
						matches = false
					}
				}
				if matches {
					myGuess = guess
					fmt.Println(myGuess)
					break
				}
			}
		}
		if len(myGuess) == 0 {
			fmt.Println("I could not find a matching word")
		}

		fmt.Print("Resp: ")
		response = readGuessResult()
		if response == "q" {
			break
		}
		if processResponse(&validLetters, myGuess, response) {
			break
		}
	}
}

func main() {
	settings := parseCmdLine()
	if len(settings.errMsg) != 0 {
		fmt.Println(settings.errMsg)
		usage()
	} else {
		MyScanner = *bufio.NewScanner(os.Stdin)
		if settings.runType == GUESS {
			doGuesses()
		} else if settings.runType == RUN {
			runGame(settings.word)
		}
	}
}
