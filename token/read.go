package token

import (
	"bufio"
	"io"
	"iter"
	"regexp"
	"strconv"
)

var (
	IntRegexp = regexp.MustCompile(`^\d+$`)
)

func Read(text io.Reader) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		scanner := bufio.NewScanner(text)
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			switch raw := scanner.Text(); raw {
			case "def":
				yield(&KeywordDef{})

			case "=":
				yield(&OperatorAssign{})

			case "+":
				yield(&OperatorPlus{})

			default:
				switch {
				case IntRegexp.MatchString(raw):
					value, _ := strconv.Atoi(raw)
					yield(&LiteralInt{Value: value})
				default:
					yield(&Identifier{Name: raw})
				}
			}
		}
	}
}
