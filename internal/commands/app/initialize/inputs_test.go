package initialize

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/10gen/realm-cli/internal/cli"
	"github.com/10gen/realm-cli/internal/cloud/realm"
	"github.com/10gen/realm-cli/internal/utils/test/assert"
	"github.com/10gen/realm-cli/internal/utils/test/mock"

	"github.com/Netflix/go-expect"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAppInitInputsResolve(t *testing.T) {
	t.Run("Should return an error if ran from a directory that already has a project", func(t *testing.T) {
		profile, teardown := mock.NewProfileFromTmpDir(t, "app_init_input_test")
		defer teardown()

		assert.Nil(t, cli.WriteFile(filepath.Join(profile.WorkingDirectory, realm.FileAppConfig), 0666, strings.NewReader(`{"name":"eggcorn"}`)))

		var i inputs
		assert.Equal(t, errProjectExists{}, i.Resolve(profile, nil))
	})

	for _, tc := range []struct {
		description string
		inputs      inputs
		procedure   func(c *expect.Console)
		test        func(t *testing.T, i inputs)
	}{
		{
			description: "With no flags set should prompt for just name and set cli.AppLocation and deployment model to defaults",
			procedure: func(c *expect.Console) {
				c.ExpectString("App Name")
				c.SendLine("test-app")
				c.ExpectEOF()
			},
			test: func(t *testing.T, i inputs) {
				assert.Equal(t, "test-app", i.Name)
				assert.Equal(t, flagDeploymentModelDefault, i.DeploymentModel)
				assert.Equal(t, flagLocationDefault, i.Location)
			},
		},
		{
			description: "With a name flag set should prompt for nothing else and set cli.AppLocation and deployment model to defaults",
			inputs:      inputs{Name: "test-app"},
			procedure:   func(c *expect.Console) {},
			test: func(t *testing.T, i inputs) {
				assert.Equal(t, "test-app", i.Name)
				assert.Equal(t, flagDeploymentModelDefault, i.DeploymentModel)
				assert.Equal(t, flagLocationDefault, i.Location)
			},
		},
		{
			description: "With name cli.AppLocation and deployment model flags set should prompt for nothing else",
			inputs: inputs{
				Name:            "test-app",
				DeploymentModel: cli.AppDeploymentModelLocal,
				Location:        cli.AppLocationOregon,
			},
			procedure: func(c *expect.Console) {},
			test: func(t *testing.T, i inputs) {
				assert.Equal(t, "test-app", i.Name)
				assert.Equal(t, cli.AppDeploymentModelLocal, i.DeploymentModel)
				assert.Equal(t, cli.AppLocationOregon, i.Location)
			},
		},
		{
			description: "With a from flag set should set from type to app",
			inputs:      inputs{From: "test-app"},
			procedure:   func(c *expect.Console) {},
			test: func(t *testing.T, i inputs) {
				assert.Equal(t, fromApp, i.FromType)
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			profile := mock.NewProfile(t)

			_, console, _, ui, consoleErr := mock.NewVT10XConsole()
			assert.Nil(t, consoleErr)
			defer console.Close()

			doneCh := make(chan (struct{}))
			go func() {
				defer close(doneCh)
				tc.procedure(console)
			}()

			assert.Nil(t, tc.inputs.Resolve(profile, ui))

			console.Tty().Close() // flush the writers
			<-doneCh              // wait for procedure to complete

			tc.test(t, tc.inputs)
		})
	}
}

func TestAppInitInputsResolveApp(t *testing.T) {
	t.Run("Should do nothing if from type is not set", func(t *testing.T) {
		var i inputs
		f, err := i.resolveFrom(nil, nil)
		assert.Nil(t, err)
		assert.Equal(t, from{}, f)
	})

	t.Run("Should return the app id and group id of specified app if from type is set to app", func(t *testing.T) {
		var appFilter realm.AppFilter
		app := realm.App{
			ID:          primitive.NewObjectID().Hex(),
			GroupID:     primitive.NewObjectID().Hex(),
			ClientAppID: "test-app-abcde",
			Name:        "test-app",
		}

		client := mock.RealmClient{}
		client.FindAppsFn = func(filter realm.AppFilter) ([]realm.App, error) {
			appFilter = filter
			return []realm.App{app}, nil
		}

		i := inputs{
			FromType: fromApp,
			Project:  app.GroupID,
			From:     app.ClientAppID,
		}

		f, err := i.resolveFrom(nil, client)
		assert.Nil(t, err)

		assert.Equal(t, from{Type: fromApp, GroupID: app.GroupID, AppID: app.ID}, f)
		assert.Equal(t, realm.AppFilter{GroupID: app.GroupID, App: app.ClientAppID}, appFilter)
	})

	t.Run("Should do nothing if from type is set to template", func(t *testing.T) {
		i := inputs{FromType: fromTemplate}
		f, err := i.resolveFrom(nil, nil)
		assert.Nil(t, err)
		assert.Equal(t, from{Type: fromTemplate}, f)
	})
}
