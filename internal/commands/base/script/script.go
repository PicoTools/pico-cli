package script

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/scripts"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

// Cmd returns command "scripts"
func Cmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scripts",
		Short:   "Manage external scripts (extender for agent interaction)",
		GroupID: constants.BaseGroupId,
	}

	cmd.AddCommand(
		// load external script to storage
		loadCmd(c),
		// list loaded external scripts
		listCmd(c),
		// reload external script from FS
		reloadCmd(c),
		// unload external script from storage
		unloadCmd(c),
	)

	return cmd
}

// loadCmd returns command "load" for "script"
func loadCmd(*console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load <file>",
		Short: "Load script by path on FS",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := scripts.RegisterExternalByPath(args[0]); err != nil {
				notificator.PrintError("%s", err.Error())
				return
			}
			notificator.PrintInfo("script successfully registered")
		},
	}
	carapace.Gen(cmd).PositionalCompletion(carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		return carapace.ActionFiles()
	}))
	return cmd
}

// listCmd returns command "list" for "script"
func listCmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List loaded external scripts",
		Run: func(cmd *cobra.Command, args []string) {
			registeredScripts := scripts.GetScripts()

			if len(registeredScripts) == 0 {
				notificator.PrintWarning("no scripts loaded yet")
				return
			}

			for _, v := range registeredScripts {
				notificator.PrintInfo("(%s) %s", v.GetAddedAt().Format("01/02 15:04:05"), v.GetPath())
			}
		},
	}
}

// unloadCmd returns command "unload" for "script"
func unloadCmd(*console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unload <file>",
		Short: "Unload already loaded script by file path",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := scripts.RemoveExternalByPath(args[0]); err != nil {
				notificator.PrintError("%s", err.Error())
				return
			}
			notificator.PrintInfo("script %s unloaded", args[0])
		},
	}

	// FS autocomplete
	carapace.Gen(cmd).PositionalCompletion(externalScriptsCompleter())

	return cmd
}

// reloadCmd returns command "reload" for "script"
func reloadCmd(*console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload [file]",
		Short: "Reload exact or all scripts",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				if err := scripts.Rebuild(); err != nil {
					notificator.PrintError("%s", err.Error())
					return
				}
				notificator.PrintInfo("storage with scripts reloaded")
			default:
				if err := scripts.ReloadExternalByPath(args[0]); err != nil {
					notificator.PrintError("%s", err.Error())
					return
				}
				notificator.PrintInfo("script %s reloaded", args[0])
			}
		},
	}

	// FS autocomplete
	carapace.Gen(cmd).PositionalCompletion(externalScriptsCompleter())

	return cmd
}

// externalScriptsCompleter returns FS autocompleter
func externalScriptsCompleter() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		var suggestions []string
		for _, v := range scripts.GetScripts() {
			suggestions = append(suggestions, v.GetPath())
		}
		return carapace.ActionValues(suggestions...)
	})
}
