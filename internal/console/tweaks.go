package console

import (
	"unicode"

	"github.com/reeflective/readline"
	"golang.org/x/text/unicode/rangetable"
)

// applyTweaks adds some tweaks for console
func applyTweaks(shell *readline.Shell) {
	// basic tweaks
	shell.Config.Set("history-autosuggest", true)
	shell.Config.Set("colored-completion-prefix", true)

	// not supported yet by reeflective/readline
	//shell.Config.Set("meta-flag", true)
	//shell.Config.Set("input-meta", true)
	//shell.Config.Set("output-meta", true)
	//shell.Config.Set("convert-meta", false)

	// add support of UTF-8 characters
	rangetable.Visit(unicode.Cyrillic, func(r rune) {
		shell.Config.Bind("emacs", string(r), "self-insert", false)
		shell.Config.Bind("vi", string(r), "self-insert", false)
	})
}
