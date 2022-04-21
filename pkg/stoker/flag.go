package stoker

type TokenList []string
type TokenListHandler[Context any] func(context Context, tokens TokenList) error

type Flag[Context any] interface {
	Name() string
	HandleTokens(context Context, tokens TokenList) error
}

func NewFlag[Context any](name string, handler TokenListHandler[Context]) *flag[Context] {
	if len(name) < 1 {
		return nil
	}

	return &flag[Context]{name: name, handler: handler}
}

type FlagList[Context any] []Flag[Context]

func (fl FlagList[Context]) FindByName(name string) (Flag[Context], bool) {
	for _, d := range fl {
		if name == d.Name() {
			return d, true
		}
	}

	return nil, false
}

type flag[Context any] struct {
	name    string
	handler TokenListHandler[Context]
}

func (f flag[Context]) Name() string {
	return f.name
}

func (f flag[Context]) HandleTokens(context Context, tokens TokenList) error {
	return f.handler(context, tokens)
}
