package initialize

import (
	"github.com/10gen/realm-cli/internal/cli"
	"github.com/10gen/realm-cli/internal/cloud/realm"
	"github.com/10gen/realm-cli/internal/terminal"

	"github.com/AlecAivazis/survey/v2"
)

const (
	fromApp      = "an existing app"
	fromTemplate = "a template"
)

type from struct {
	Type       string
	GroupID    string
	AppID      string
	TemplateID string
}

type inputs struct {
	AppData         cli.AppData
	Project         string
	From            string
	FromType        string
	Name            string
	DeploymentModel cli.AppDeploymentModel
	Location        cli.AppLocation
}

func (i *inputs) Resolve(profile *cli.Profile, ui terminal.UI) error {
	appData, appDataErr := cli.ResolveAppData(profile.WorkingDirectory)
	if appDataErr != nil {
		return appDataErr
	}
	if appData.Name != "" {
		return errProjectExists{}
	}

	if i.From != "" {
		i.FromType = determineFromType(i.From)

		// TODO(REALMC-XXXX): implement Realm app templates
		// if err := ui.AskOne(&i.FromType, &survey.Select{
		// 	Message: "Initialize app from what?",
		// 	Options: []string{fromApp, fromTemplate},
		// }); err != nil {
		// 	return err
		// }
	} else {
		if i.Name == "" {
			if err := ui.AskOne(&i.Name, &survey.Input{Message: "App Name"}); err != nil {
				return err
			}
		}
		if i.DeploymentModel == cli.AppDeploymentModelNil {
			i.DeploymentModel = flagDeploymentModelDefault
		}
		if i.Location == cli.AppLocationNil {
			i.Location = flagLocationDefault
		}
	}
	return nil
}

func (i *inputs) resolveFrom(ui terminal.UI, client realm.Client) (from, error) {
	f := from{Type: i.FromType}

	switch i.FromType {
	case fromApp:
		app, err := cli.ResolveApp(ui, client, realm.AppFilter{GroupID: i.Project, App: i.From})
		if err != nil {
			return from{}, err
		}
		f.GroupID = app.GroupID
		f.AppID = app.ID
	case fromTemplate:
		// do nothing

		// TODO(REALMC-XXXX): implement Realm app templates
		// template, err := cli.ResolveTemplate(ui, client, realm.TemplateFilter{GroupID: i.Project, Template: i.From})
		// if err != nil {
		// 	return from{}, err
		// }
		// f.TemplateID = template.ID
	}
	return f, nil
}

func determineFromType(from string) string {
	return fromApp

	// TODO(REALMC-XXXX): implement Realm app templates
	// if strings.HasPrefix(from, "__some_template_prefix?") {
	// 	return fromTemplate
	// }
	// return fromApp
}
