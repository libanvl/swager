package stoker

import (
	"strings"
)

type FlagHandler[Context any] struct {
	flag   Flag[Context]
	tokens TokenList
}

func (fh FlagHandler[Context]) Handle(context Context) error {
	return fh.flag.HandleTokens(context, fh.tokens)
}

type FlagHandlerList[Context any] []FlagHandler[Context]

func (fhl FlagHandlerList[Context]) HandleAll(context Context) error {
	for _, fh := range fhl {
		if err := fh.Handle(context); err != nil {
			return err
		}
	}

	return nil
}

type Parser[Context any] interface {
	Parse(args ...string) FlagHandlerList[Context]
	Present(flag string) bool
}

func NewParser[Context any](flags ...Flag[Context]) *parser[Context] {
	return &parser[Context]{flags: flags, present: make([]string, 0)}
}

type parser[Context any] struct {
	flags   FlagList[Context]
	present []string
}

func (p parser[Context]) Parse(args ...string) FlagHandlerList[Context] {
	result := make(FlagHandlerList[Context], 0)

	var currflag Flag[Context] = nil
	var currlist TokenList = nil

	for _, token := range args {
		// the current token is a flag
		if nextflag, ok := p.flags.FindByName(strings.ToLower(token)); ok {
			// - complete the work for the current flag
			if currflag != nil {
				// -- add the current tokenlist to the tokenset for the current flag
				result = append(result, FlagHandler[Context]{flag: currflag, tokens: currlist})
			}

			// - reset the current tokenlist
			currlist = nil

			// - add the next flag token to the list of present flags
			if !contains(p.present, nextflag.Name()) {
				p.present = append(p.present, nextflag.Name())
			}

			// - starting new tokenlist
			currlist = make(TokenList, 0)

			currflag = nextflag
			continue
		}

		// the current token is an argument
		// - append the token to the current tokenlist
		if currlist != nil {
			currlist = append(currlist, token)
		}
	}

	if currflag != nil {
		// - add the current tokenlist to the tokenset for the last flag in the list
		result = append(result, FlagHandler[Context]{flag: currflag, tokens: currlist})
	}

	return result
}

func (p *parser[Context]) Present(flag string) bool {
	return contains(p.present, strings.ToLower(flag))
}

func contains[T ~string](ts []T, arg T) bool {
	for _, s := range ts {
		if arg == s {
			return true
		}
	}

	return false
}
