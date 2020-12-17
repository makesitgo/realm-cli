package initialize

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/10gen/realm-cli/internal/cli"
	"github.com/10gen/realm-cli/internal/cloud/realm"
	"github.com/10gen/realm-cli/internal/terminal"

	"github.com/spf13/pflag"
)

// Command is the `app init` command
var Command = cli.CommandDefinition{
	Use:         "init",
	Aliases:     []string{"initialize"},
	Display:     "app init",
	Description: "Initialize a Realm app in your current local directory",
	Help:        "",
	Command:     &command{},
}

type command struct {
	inputs      inputs
	realmClient realm.Client
}

func (cmd *command) Flags(fs *pflag.FlagSet) {
	fs.StringVar(&cmd.inputs.Project, flagProject, "", flagProjectUsage)
	fs.StringVarP(&cmd.inputs.From, flagFrom, flagFromShort, "", flagFromUsage)
	fs.StringVarP(&cmd.inputs.Name, flagName, flagNameShort, "", flagNameUsage)
	fs.VarP(&cmd.inputs.DeploymentModel, flagDeploymentModel, flagDeploymentModelShort, flagDeploymentModelUsage)
	fs.VarP(&cmd.inputs.Location, flagLocation, flagLocationShort, flagLocationUsage)
}

func (cmd *command) Inputs() cli.InputResolver {
	return &cmd.inputs
}

func (cmd *command) Setup(profile *cli.Profile, ui terminal.UI) error {
	cmd.realmClient = realm.NewAuthClient(profile.RealmBaseURL(), profile.Session())
	return nil
}

func (cmd *command) Handler(profile *cli.Profile, ui terminal.UI) error {
	from, fromErr := cmd.inputs.resolveFrom(ui, cmd.realmClient)
	if fromErr != nil {
		return fromErr
	}

	switch from.Type {
	case fromApp:
		return cmd.initializeFromApp(profile.WorkingDirectory, from.GroupID, from.AppID)
	case fromTemplate:
		return cmd.initializeFromTemplate(profile.WorkingDirectory, from.TemplateID)
	}
	return cmd.initialize(profile.WorkingDirectory)
}

func (cmd *command) Feedback(profile *cli.Profile, ui terminal.UI) error {
	return ui.Print(terminal.NewTextLog("Successfully initialized app"))
}

func (cmd *command) initialize(wd string) error {
	data := []byte(fmt.Sprintf(`{
    "config_version": %s,
    "name": %q,
    "location": %q,
    "deployment_model": %q,
    "security": {},
    "custom_user_data_config": {
        "enabled": false
    },
    "sync": {
        "development_mode_enabled": false
    }
}
`, realm.DefaultAppConfigVersion, cmd.inputs.Name, cmd.inputs.Location, cmd.inputs.DeploymentModel))

	return cli.WriteFile(filepath.Join(wd, realm.FileAppConfig), 0666, bytes.NewReader(data))
}

func (cmd *command) initializeFromApp(wd, groupID, appID string) error {
	_, zipPkg, exportErr := cmd.realmClient.Export(groupID, appID, realm.ExportRequest{IsTemplated: true})
	if exportErr != nil {
		return exportErr
	}
	return cli.WriteZip(wd, zipPkg)
}

func (cmd *command) initializeFromTemplate(wd, templateID string) error {
	// TODO(REALMC-XXXX): implement Realm app templates
	return errors.New("initializing from templates is not yet supported")
}
