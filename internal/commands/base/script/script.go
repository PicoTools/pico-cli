package script

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/scripts"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

func loadCmd(*console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "load",
		Short:                 "Load script by path on FS",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
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

func listCmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:                   "list",
		Short:                 "List registred scripts",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			registeredScripts := scripts.GetScripts()
			if len(registeredScripts) == 0 {
				notificator.PrintWarning("no scripts registered")
				return
			}
			for _, v := range registeredScripts {
				notificator.Print("[%s] %s", v.GetAddedAt().Format("01/02 15:04:05"), v.GetPath())
			}
		},
	}
}

func removeCmd(*console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove",
		Short:                 "Remove registred scripts",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := scripts.RemoveExternalByPath(args[0]); err != nil {
				notificator.PrintError("%s", err.Error())
				return
			}
			notificator.PrintInfo("script %s removed", args[0])
		},
	}
	carapace.Gen(cmd).PositionalCompletion(externalScriptsCompleter())
	return cmd
}

func reloadCmd(*console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "reload",
		Short:                 "Reload script/all scripts",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				// reload all scripts
				if err := scripts.Rebuild(); err != nil {
					notificator.PrintError("%s", err.Error())
					return
				}
				notificator.PrintInfo("all scripts reloaded")
				return
			}
			if err := scripts.ReloadExternalByPath(args[0]); err != nil {
				notificator.PrintInfo("%s", err.Error())
				return
			}
			notificator.PrintInfo("script %s reloaded", args[0])
		},
	}
	carapace.Gen(cmd).PositionalCompletion(externalScriptsCompleter())
	return cmd
}

func Cmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "scripts",
		Short:                 "Manage scripts",
		DisableFlagsInUseLine: true,
		GroupID:               constants.BaseGroupId,
	}
	cmd.AddCommand(
		loadCmd(c),
		listCmd(c),
		reloadCmd(c),
		removeCmd(c),
	)
	return cmd
}

func externalScriptsCompleter() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		var suggestions []string
		for _, v := range scripts.GetScripts() {
			suggestions = append(suggestions, v.GetPath())
		}
		return carapace.ActionValues(suggestions...)
	})
}
