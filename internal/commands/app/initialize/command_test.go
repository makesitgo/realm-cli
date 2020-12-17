package initialize

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/10gen/realm-cli/internal/cloud/realm"
	"github.com/10gen/realm-cli/internal/utils/test/assert"
	"github.com/10gen/realm-cli/internal/utils/test/mock"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAppInitSetup(t *testing.T) {
	t.Run("Should construct a Realm client with the configured base url", func(t *testing.T) {
		profile := mock.NewProfile(t)
		profile.SetRealmBaseURL("http://localhost:8080")

		cmd := &command{inputs: inputs{
			Name: "test-app",
		}}
		assert.Nil(t, cmd.realmClient)

		assert.Nil(t, cmd.Setup(profile, nil))
		assert.NotNil(t, cmd.realmClient)
	})
}

func TestAppInitHandler(t *testing.T) {
	t.Run("Should initialize an empty project when no from type is specified", func(t *testing.T) {
		profile, teardown := mock.NewProfileFromTmpDir(t, "app_init_test")
		defer teardown()

		cmd := &command{inputs: inputs{
			Name:            "test-app",
			DeploymentModel: deploymentModelLocal,
			Location:        locationSydney,
		}}

		assert.Nil(t, cmd.Handler(profile, nil))

		config, err := ioutil.ReadFile(filepath.Join(profile.WorkingDirectory, realm.FileAppConfig))
		assert.Nil(t, err)
		assert.Equal(t, `{
    "config_version": 20200603,
    "name": "test-app",
    "location": "AU",
    "deployment_model": "LOCAL",
    "security": {},
    "custom_user_data_config": {
        "enabled": false
    },
    "sync": {
        "development_mode_enabled": false
    }
}
`, string(config))
	})

	t.Run("Should initialze a templated app when from type is specified to app", func(t *testing.T) {
		profile, teardown := mock.NewProfileFromTmpDir(t, "app_init_test")
		defer teardown()

		app := realm.App{
			ID:          primitive.NewObjectID().Hex(),
			GroupID:     primitive.NewObjectID().Hex(),
			ClientAppID: "test-app-abcde",
			Name:        "test-app",
		}

		var zipCloser *zip.ReadCloser
		defer func() { zipCloser.Close() }()

		client := mock.RealmClient{}
		client.FindAppsFn = func(filter realm.AppFilter) ([]realm.App, error) {
			return []realm.App{app}, nil
		}
		client.ExportFn = func(groupID, appID string, req realm.ExportRequest) (string, *zip.Reader, error) {
			zipPkg, err := zip.OpenReader("testdata/project.zip")

			zipCloser = zipPkg
			return "", &zipPkg.Reader, err
		}

		cmd := &command{
			inputs: inputs{
				FromType: fromApp,
			},
			realmClient: client,
		}

		assert.Nil(t, cmd.Handler(profile, nil))

		t.Run("Should have the expected contents in the app config file", func(t *testing.T) {
			config, err := ioutil.ReadFile(filepath.Join(profile.WorkingDirectory, realm.FileAppConfig))
			assert.Nil(t, err)
			assert.Equal(t, `{
    "config_version": 20200603,
    "name": "from-app",
    "location": "IE",
    "deployment_model": "GLOBAL",
    "security": {},
    "custom_user_data_config": {
        "enabled": false
    },
    "sync": {
        "development_mode_enabled": false
    }
}
`, string(config))
		})

		t.Run("Should have the expected contents in the api key auth provider config file", func(t *testing.T) {
			config, err := ioutil.ReadFile(filepath.Join(profile.WorkingDirectory, realm.FileAuthProvider("api-key")))
			assert.Nil(t, err)
			assert.Equal(t, `{
    "name": "api-key",
    "type": "api-key",
    "disabled": false
}
`, string(config))
		})
	})

	t.Run("Should error out when from type is set to template", func(t *testing.T) {
		profile := mock.NewProfile(t)

		cmd := &command{inputs: inputs{
			FromType: fromTemplate,
		}}

		assert.Equal(t,
			errors.New("initializing from templates is not yet supported"),
			cmd.Handler(profile, nil),
		)
	})
}

func TestAppInitFeedback(t *testing.T) {
	t.Run("Feedback should print a message that app initialization was successful", func(t *testing.T) {
		out, ui := mock.NewUI()

		cmd := &command{}

		err := cmd.Feedback(nil, ui)
		assert.Nil(t, err)

		assert.Equal(t, "01:23:45 UTC INFO  Successfully initialized app\n", out.String())
	})
}
