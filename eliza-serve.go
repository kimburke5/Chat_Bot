package main

import (
	"fmt"
	"net/http"
	"strings"
    "math/rand"
    "regexp"
    "bufio"
    "os"
	"./goeliza"
	
)
func main() {

	//serve the files from the /public folder
	fs := http.Dir("./public")
	fileServer := http.FileServer(fs)
	//handle requests to /
	http.Handle("/", fileServer)
	//handle request to /userinput
	http.HandleFunc("/userinput", userInputHandler)
	//listen on tcp and serve requests on port :8080
	http.ListenAndServe(":8080", nil)
}

func userInputHandler(w http.ResponseWriter, r *http.Request) {
	//executed when a request is made to the userInput
	userinput := r.URL.Query().Get("userinput")
	reply := startEliza(userinput)
	fmt.Fprintf(w, reply)
} //inputHandler

func startEliza(input string) string {

	fmt.Println("Eliza: " + ElizaHi())

    for {
        statement := getInput()
        fmt.Println("Eliza: " + ReplyTo(statement))

       if IsQuitStatement(statement) {
            break
        }
    }
	return getInput()
}
// ElizaHi will return a random introductory sentence for ELIZA.
func ElizaHi() string {
    return randChoice(goeliza.Introductions)
}

// ElizaHi will return a random goodbye sentence for ELIZA.
func ElizaBye() string {
    return randChoice(goeliza.Goodbyes)
}

// ReplyTo will construct a reply for a given statement using ELIZA's rules.
func ReplyTo(statement string) string {
    // First, preprocess the statement for more effective matching
    statement = preprocess(statement)

    // Then, we check if this is a quit statement
    if IsQuitStatement(statement) {
        return ElizaBye()
    }

    // Next, we try to match the statement to a statement that ELIZA can 
    // recognize, and construct a pre-determined, appropriate response.
    for pattern, responses := range goeliza.Psychobabble {
        re := regexp.MustCompile(pattern)
        matches := re.FindStringSubmatch(statement)

        // If the statement matched any recognizable statements.
        if len(matches) > 0 {
            // If we matched a regex group in parentheses, get the first match.
            // The matched regex group will match a "fragment" that will form 
            // part of the response, for added realism.
            var fragment string
            if len(matches) > 1 {
                fragment = reflect(matches[1])
            }

            // Choose a random appropriate response, and format it with the 
            // fragment, if needed.
            response := randChoice(responses)
            if strings.Contains(response, "%s") {
                response = fmt.Sprintf(response, fragment)
            }
            return response
        }
    }

    // If no patterns were matched, return a default response.
    return randChoice(goeliza.DefaultResponses)
}

// IsQuitStatement returns if the statement is a quit statement
func IsQuitStatement(statement string) bool {
    statement = preprocess(statement)
    for _, quitStatement := range goeliza.QuitStatements {
        if statement == quitStatement {
            return true
        }
    }
    return false
}

// preprocess will do some normalization on a statement for better regex matching
func preprocess(statement string) string {
    statement = strings.TrimRight(statement, "\n.!")
    statement = strings.ToLower(statement)
    return statement
}

// reflect flips a few words in an input fragment (such as "I" -> "you").
func reflect(fragment string) string {
    words := strings.Split(fragment, " ")
    for i, word := range words {
        if reflectedWord, ok := goeliza.ReflectedWords[word]; ok {
            words[i] = reflectedWord
        }
    }
    return strings.Join(words, " ")
}

// randChoice returns a random element in an (string) array.
func randChoice(list []string) string {
    randIndex := rand.Intn(len(list))
    return list[randIndex]
}

func getInput() string {
    fmt.Print("You: ")
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    return input
}