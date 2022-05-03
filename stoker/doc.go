/*
Package stocker is for specialized argument parsing.

A list of string tokens are processed in order. When a token matches a defined Flag,
the following tokens are captured into a TokenList until a new Flag token
is found, or the end of input.

Each Flag has an associated TokenListHandler which is called once for each TokenList after parsing is
complete and the HandleAll method is called. A TokenListHandler additionally takes a context object passed with HandleAll.

Example

Given the flags: --foo and --bar
And the input: one two --foo three four --bar --five

  * Token one and two are ignored.
  * Token three and four are placed in a TokenList for the --foo Flag
  * Token --five is placed in a TokenList for the --bar Flag

Parser

A set of flags is used to create a Parser instance. The flags in a single Parser must have the same type for the
TokenListHandler context object. Calling Parse returns a FlagHandlerList. HandleAll can be called on the list to
call each handler in the order the flags where found during parsing.

Example

```
context := new(ExampleContext)

common_handler := func(context *ExampleContext, tl stoker.TokenList) error {
  return nil
}

parser := stoker.NewParser[*ExampleContext](
  stoker.NewFlag("--foo", common_handler),
  stoker.NewFlag("--bar", common_handler)
)

handlers := parser.Parse(os.Args...)

if err := handlers.HandleAll(context); err := nil {
  log.Fatal(err)
}
```
*/
package stoker
