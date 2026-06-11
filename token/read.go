package token

import (
	"bufio"
	"io"
	"iter"
	"regexp"
	"strconv"
)

var (
	IntRegexp   = regexp.MustCompile(`^-?\d+$`)
	FloatRegexp = regexp.MustCompile(`^-?\d+\.\d+$`)
)

func Read(text io.Reader) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		scanner := bufio.NewScanner(text)
		scanner.Split(readSplitter)

		for scanner.Scan() {
			switch raw := scanner.Text(); raw {
			case "def":
				yield(NewKeywordDefine())

			case "free":
				yield(NewKeywordFree())

			case "ret":
				yield(NewKeywordReturn())

			case "defer":
				yield(NewKeywordDefer())

			case "=":
				yield(NewOperatorAssign())

			case "+":
				yield(NewOperatorPlus())

			case "-":
				yield(NewOperatorMinus())

			case "{":
				yield(NewCurlyBracketOpen())

			case "}":
				yield(NewCurlyBracketClose())

			case "(":
				yield(NewParenthesesOpen())

			case ")":
				yield(NewParenthesesClose())

			default:
				switch {
				case FloatRegexp.MatchString(raw):
					value, _ := strconv.ParseFloat(raw, 64)
					yield(NewLiteralFloat(value))

				case IntRegexp.MatchString(raw):
					value, _ := strconv.Atoi(raw)
					yield(NewLiteralInt(value))

				default:
					yield(NewIdentifier(raw))
				}
			}
		}
	}
}
