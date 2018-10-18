package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)


// Configuration
var leet_dict string
var wordlist_file string
var verbose bool
var show_help bool

// In-Memory Storage
var dictionary map [string] []string
var regex_wordlist []regexp.Regexp

// Parallelism
var wg sync.WaitGroup

// Convert word to match using dictionary
func wordToRegex(word string) string {
	var result bytes.Buffer

	word = strings.TrimSpace(word)

	lower := strings.ToLower(word)
	upper := strings.ToUpper(word)

	result.WriteString("(")
	result.WriteString(lower)
	result.WriteString("|")
	result.WriteString(upper)

	// Add additional parts here!
	result.WriteString("|")
	result.WriteString("(")
	for _,c := range word {
		var possible_matches_char []string
		possible_matches_char = append(possible_matches_char, string(c))
		if string(c) == strings.ToUpper(string(c)){
			possible_matches_char = append(possible_matches_char,
					strings.ToLower(string(c)))
		} else {
			possible_matches_char = append(possible_matches_char,
					strings.ToUpper(string(c)))
		}
		for _,item := range dictionary[string(c)] {
			if len(item) < 1 {
				continue
			}
			if verbose {
				fmt.Printf("\t%s\n", item)
			}
			r := strings.NewReplacer(
					"[", "\\[",
					"]", "\\]",
					"\\", "[\\\\]",
					"|", "\\|",
					"(", "\\(",
					")", "\\)",
					"}", "\\}",
					"{", "\\{",
					"+", "\\+",
					"*", "\\*",
			)
			item = r.Replace(item)
			possible_matches_char = append(possible_matches_char, item)
		}

		result.WriteString("(")
		result.WriteString(strings.Join(possible_matches_char, "|"))
		result.WriteString(")")
	}
	result.WriteString(")")

	result.WriteString(")")
	return result.String()
}

/*
 * Returns true if successful, false with error if not.
 */
func loadDictionary() (bool, error){
	dictionary = make(map[string][]string)

	f, err := os.Open(leet_dict)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		key_val := strings.Split(line, ":")

		if len(key_val) <= 1 {
			log.Println(key_val)
			log.Printf("Can't split `%s` on semicolon. Skipping.\n", line)
			continue
		}

		possible_repls := strings.Split(key_val[1], ",")

		// Trim all words
		for i,v := range possible_repls {
			possible_repls[i] = strings.TrimSpace(v)
		}

		if verbose {
			fmt.Printf("Adding %s:%s\n", key_val[0], strings.Join(possible_repls,","))
		}

		// Store upper and lowercase, just in case. :p
		dictionary[strings.ToUpper(key_val[0])] = possible_repls
		dictionary[strings.ToLower(key_val[0])] = possible_repls
	}

	return true, nil
}

/*
 * Returns true if successful, false with error if not.
 */
func loadWordlist() (bool, error){
	regex_wordlist = []regexp.Regexp{}

	// Use dictionary to generate list of regexes for matching replacements
	f, err := os.Open(wordlist_file)
	if err != nil {
		return false, err
	}
	defer f.Close()

	wordscanner := bufio.NewScanner(f)
	for wordscanner.Scan() {
		word := wordscanner.Text()

		// Skip empty lines
		if len(word) == 0 {
			continue
		}
		word = wordToRegex(word)
		if verbose {
			fmt.Printf("Adding regex:%s\n", word)
		}

		word_regex := regexp.MustCompile(word)
		regex_wordlist = append(regex_wordlist, *word_regex)
	}
	return true, nil
}

func printHighlight(line string, subword string) {
	loc := strings.Index(line, subword)
	if loc == -1 {
		log.Printf("Error, %s not in %s\n", subword, line)
		return
	}
	fmt.Printf("%s", line[:loc])
	fmt.Printf("\x1B[31m%s\x1B[0m", subword)
	fmt.Printf("%s", line[loc+len(subword):])
	fmt.Println()
}

func readLoop(){
	scanIn := bufio.NewScanner(os.Stdin)
	onlyalpha := regexp.MustCompile("[A-Za-z]+")

	for scanIn.Scan() {
		line := scanIn.Text()
		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		for _,regex := range regex_wordlist {
			matches := regex.FindAll([]byte(line), -1)
			if len(matches) > 0 {
				for _,match := range matches {
					alpha_text := onlyalpha.Find(match)
					if alpha_text != nil && string(alpha_text) == string(match) {
						// Skip if not 1337 speak.
						continue
					}
					printHighlight(line, string(match))
				}
			}
		}
	}
}

func init(){
	flag.StringVar(&leet_dict, "d", "list.dict", "Detection dictionary.")
	flag.StringVar(&wordlist_file, "w", "words.txt", "Wordlist to detect on.")
	flag.BoolVar(&verbose, "v", false, "Verbose mode.")
	flag.BoolVar(&show_help, "h", false, "Show help menu.")
}

func main(){
	flag.Parse()

	if show_help {
		flag.PrintDefaults()
		return
	}

	loadDictionary()
	fmt.Printf("Dict size:%d\n", len(dictionary))

	loadWordlist()
	fmt.Printf("Wordlist size:%d\n", len(regex_wordlist))

	readLoop()
}
