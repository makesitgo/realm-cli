package cmd

import (
	"fmt"

	"github.com/10gen/realm-cli/internal/cli"
	"github.com/10gen/realm-cli/internal/flags"

	"github.com/spf13/cobra"
	"honnef.co/go/tools/version"
)

const (
	cliName = "realm-cli"
)

var (
	rootCmd = &cobra.Command{
		Version: version.Version,
		Use:     cliName,
		Short:   "CLI tool to manage your MongoDB Realm application",
		Long:    fmt.Sprintf("Use %s command help for information on a specific command", cliName),
	}
)

// Execute runs the CLI
func Execute() {
	cli.HandleErr(rootCmd.Execute())
}

func init() {
	profile, profileErr := cli.NewDefaultProfile()
	cli.HandleErr(profileErr)

	var config cli.Config

	rootCmd.PersistentFlags().StringVarP(&profile.Name, flags.Profile, flags.ProfileShort, cli.DefaultProfile, flags.ProfileUsage)
	rootCmd.PersistentFlags().BoolVar(&config.UI.DisableColors, flags.DisableColors, false, flags.DisableColorsUsage)
	rootCmd.PersistentFlags().StringVar(&config.UI.OutputFormat, flags.OutputFormat, flags.DefaultOutputFormat, flags.OutputFormatUsage)
	rootCmd.PersistentFlags().StringVar(&config.UI.OutputTarget, flags.OutputTarget, "", flags.OutputTargetUsage)
	rootCmd.PersistentFlags().StringVar(&config.Command.RealmBaseURL, flags.RealmBaseURL, flags.DefaultRealmBaseURL, flags.RealmBaseURLUsage)

	factory := cli.CommandFactory{profile, config}

	rootCmd.AddCommand(factory.Build(cli.LoginCommand))
	rootCmd.AddCommand(factory.Build(cli.LogoutCommand))
	rootCmd.AddCommand(factory.Build(cli.WhoamiCommand))

	// pass a list of functions to be run before each command's Execute function
	cobra.OnInitialize(func() { cli.HandleErr(profile.Load()) })
}