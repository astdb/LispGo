// LispGo is a toy lisp interpreter writtn in Go, in line with the spirit of http://norvig.com/lispy.html,
// https://maryrosecook.com/blog/post/little-lisp-interpreter and http://steve-yegge.blogspot.com.au/2007/06/rich-programmer-food.html

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	breadr := bufio.NewReader(os.Stdin)		// buffered reader to read console input
	prompt := "glisop>>"	// console prompt

	// run REPL-style until keyboard interrupt
	for {
		// show prompt
		fmt.Print(prompt)
		in, _ := breadr.ReadString('\n')

		if strings.TrimSpace(in) == "" {
			// no input - continue showing prompt

		} else {
			// input!! - evaluate
			fmt.Printf("%s\n", parse(in)) // echo parsed input
		}
	}
}

// split input statement into tokens (names, numbers and parentheses)
func tokenize(input string) []string {
	input = strings.TrimSpace(input)

	// prepare regexps to match open/close parentheses
	r_spc_op, err := regexp.Compile(`\(`) // replace all opening parentheses with [opening parentheses with spaces around them]
	r_spc_cl, err := regexp.Compile(`\)`) // replace all closing parentheses with [closing parentheses with spaces around them]

	// error in regexp?
	if err != nil {
		fmt.Printf("Regexp error in tokenizer: %v\n", err)
		return nil
	}

	// perform replacement using prepared regexps and return result of splitting on whitespace
	input = r_spc_op.ReplaceAllString(r_spc_cl.ReplaceAllString(input, " ) "), " ( ")
	return strings.Split(input, " ")
}

// parenthesize() takes the tokens produced by tokenize() and produces a nested array that mimics the structure of the Lisp code
func parenthesize(tokens []string, list []*category) []*category {
	if list == nil {
		return parenthesize(tokens, []*category{})
	} else {
		token := strings.TrimSpace(tokens[0])
		tokens = tokens[1:]

		if token == "" {
			// return list.pop();
			// return list[len(list)-1]
			if len(list) > 0 {
				return list[len(list)-1:]
			} else {
				return nil
			}
			
		} else if token == "(" {
			// list.push(parenthesize(input, []));
			// return parenthesize(input, list);
			list = append(list, parenthesize(tokens, []*category{})...)
			return parenthesize(tokens, list)
		} else if token == ")" {
			return list
		} /*else {
			return parenthesize(tokens, append(list, categorize(token)))
		} */
		return parenthesize(tokens, append(list, categorize(token)))
	}
}

func categorize(token string) *category {
	_, err := strconv.ParseFloat(strings.TrimSpace(token), 64)

	if err == nil {
		// number
		// return { type:'literal', value: parseFloat(input) };
		return &category{"literal", token}

	} else {
		tok_runes := []rune(token)
		if tok_runes[0] == '"' && tok_runes[len(tok_runes)-1] == '"' {
			// return { type:'literal', value: input.slice(1, -1) };
			return &category{"literal", string(tok_runes[1 : len(tok_runes)-1])}

		} else {
			// return { type:'identifier', value: input };
			return &category{"identifier", token}
		}
	}
}

func parse(input string) []*category {
	return parenthesize(tokenize(input), nil)
	/*cats := parenthesize(tokenize(input), nil)

	for _, c := range cats {
		c.toString()
	}
	fmt.Print("\n")*/
}

type category struct {
	typ string // type
	val string // value
}

func (cat *category) toString() string {
	return fmt.Sprintf("[type: %s, value: %s] ", cat.typ, cat.val)
}
