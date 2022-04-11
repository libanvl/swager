/*
Package stocker is for specialized argument parsing.
*/
package stoker

import (
	"errors"
	"strings"
)

type TokenSetProcessor func(TokenSet) error
type TokenList []string

type def struct {
	flag    string
	hasArgs bool
}

type defList []*def

type parser struct {
	Defs    defList
	Present TokenList
}

type TokenSet []TokenList
type TokenMap map[string]TokenSet

func Parser(heads ...*def) *parser {
	return &parser{Defs: heads, Present: make([]string, 0)}
}

func Def(flag string) *def {
	return &def{flag: flag, hasArgs: true}
}

func DefArgs(flag string, hasArgs bool) *def {
	return &def{flag: flag, hasArgs: hasArgs}
}

func (p parser) Parse(args ...string) TokenMap {
	result := make(TokenMap)

	var currdef *def = nil
	var prevdef *def = nil
	var currlist TokenList = nil
	for _, token := range args {
		tokenlower := strings.ToLower(token)
		var ok bool = false
		if currdef, ok = p.Defs.findByFlag(tokenlower); ok {
			if !p.Present.Contains(currdef.flag) {
				p.Present = append(p.Present, currdef.flag)
			}
			// starting new tokenlist

			if !currdef.hasArgs {
				continue
			}

			// move currlist to map
			// - ensure the map has a tokenset for the flag
			if _, ok := result[currdef.flag]; !ok {
				result[currdef.flag] = make(TokenSet, 0)
			}

			if prevdef != nil && len(currlist) > 0 {
				// - add the current tokenlist to the tokenset for the last flag
				result[prevdef.flag] = append(result[prevdef.flag], currlist)
			}

			// reset the current tokenlist
			prevdef = currdef
			currlist = nil
			continue
		}

		// append the token to the current tokenlist
		if currlist == nil {
			currlist = TokenList{token}
		} else {
			currlist = append(currlist, token)
		}
	}

	if prevdef != nil {
		result[prevdef.flag] = append(result[prevdef.flag], currlist)
	}

	return result
}

func (m TokenMap) ProcessSet(flag string, p TokenSetProcessor) error {
	if p == nil {
		return errors.New("TokenSetProcessor cannot be nil")
	}

	if ts, ok := m[flag]; ok {
		return p(ts)
	}

	return nil
}

func (ss TokenList) Contains(arg string) bool {
	for _, s := range ss {
		if arg == s {
			return true
		}
	}

	return false
}

func (dl defList) findByFlag(flag string) (*def, bool) {
	for _, d := range dl {
		if flag == d.flag {
			return d, true
		}
	}

	return nil, false
}
