package pull

import (
	"errors"

	"github.com/10gen/realm-cli/internal/app"
	"github.com/10gen/realm-cli/internal/cli"
	"github.com/10gen/realm-cli/internal/cloud/realm"
	"github.com/10gen/realm-cli/internal/terminal"

	"github.com/mitchellh/go-homedir"
)

const (
	flagFrom      = "from"
	flagFromShort = "a"
	flagFromUsage = "specify the app to pull changes down from"

	flagProject      = "project"
	flagProjectUsage = "the MongoDB cloud project id"

	flagAppVersion      = "app-version"
	flagAppVersionUsage = "specify the app config version to pull changes down as"

	flagTarget      = "target"
	flagTargetShort = "t"
	flagTargetUsage = "provide the path to export a Realm app to"

	flagIncludeDependencies      = "include-dependencies"
	flagIncludeDependenciesShort = "d"
	flagIncludeDependenciesUsage = "include to to push Realm app dependencies changes as well"

	flagIncludeHosting      = "include-hosting"
	flagIncludeHostingShort = "s"
	flagIncludeHostingUsage = "include to push Realm app hosting changes as well"

	flagDryRun      = "dry-run"
	flagDryRunShort = "x"
	flagDryRunUsage = "include to run without writing any changes to the file system"
)

var (
	errConfigVersionMismatch = errors.New("must export an app with the same config version as found in the current project directory")
)

type inputs struct {
	Project             string
	From                string
	Target              string
	AppVersion          realm.AppConfigVersion
	IncludeDependencies bool
	IncludeHosting      bool
	DryRun              bool
}

func (i *inputs) Resolve(profile *cli.Profile, ui terminal.UI) error {
	wd := i.Target
	if wd == "" {
		wd = profile.WorkingDirectory
	}

	appDir, appConfig, appConfigErr := app.ResolveConfig(wd)
	if appConfigErr != nil {
		return appConfigErr
	}

	var target string
	if i.Target == "" {
		if appDir == "" {
			return errProjectNotFound{}
		}
		target = appDir
	} else {
		t, err := homedir.Expand(i.Target)
		if err != nil {
			return err
		}
		target = t
	}
	i.Target = target

	if appDir != "" {
		if i.AppVersion == realm.AppConfigVersionZero {
			i.AppVersion = appConfig.ConfigVersion
		} else if i.AppVersion != appConfig.ConfigVersion && appConfig.ConfigVersion != realm.AppConfigVersionZero {
			return errConfigVersionMismatch
		}

		if i.From == "" {
			i.From = appConfig.String()
		}
	}
	return nil
}

type from struct {
	GroupID string
	AppID   string
}

func (i inputs) resolveFrom(ui terminal.UI, client realm.Client) (from, error) {
	f := from{GroupID: i.Project}

	if i.From == "" {
		return f, nil
	}

	a, err := app.Resolve(ui, client, realm.AppFilter{GroupID: i.Project, App: i.From})
	if err != nil {
		if _, ok := err.(app.ErrAppNotFound); !ok {
			return from{}, err
		}
	}

	f.GroupID = a.GroupID
	f.AppID = a.ID
	return f, nil
}