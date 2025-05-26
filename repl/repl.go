package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fliptv97/monkey-interpreter/lexer"
	"github.com/fliptv97/monkey-interpreter/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	for {
		_, err := fmt.Fprint(out, PROMPT)
		if err != nil {
			return err
		}

		scanned := scanner.Scan()
		if !scanned {
			return nil
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			_, err = fmt.Fprintf(out, "%+v\n", tok)
			if err != nil {
				return err
			}
		}
	}
}
