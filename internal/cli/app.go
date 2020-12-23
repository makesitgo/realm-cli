package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/10gen/realm-cli/internal/cloud/realm"
	"github.com/AlecAivazis/survey/v2/core"
)

const (
	ExportedJSONPrefix = ""
	ExportedJSONIndent = "    "
)

var (
	maxDirectoryContainSearchDepth = 8
)

// AppData is the partial form of an exported MongoDB Realm application configuration data
type AppData struct {
	ID   string `json:"client_app_id,omitempty"`
	Name string `json:"name"`
}

// AppConfig is the exported MongoDB Realm application configuration data
type AppConfig struct {
	ConfigVersion realm.AppConfigVersion `json:"config_version"`
	AppData
	Location             AppLocation             `json:"location"`
	DeploymentModel      AppDeploymentModel      `json:"deployment_model"`
	Security             AppSecurityConfig       `json:"security"`
	CustomUserDataConfig AppCustomUserDataConfig `json:"custom_user_data_config"`
	Sync                 AppSyncConfig           `json:"sync"`
}

type AppSecurityConfig struct{}

type AppCustomUserDataConfig struct {
	Enabled bool `json:"enabled"`
}

type AppSyncConfig struct {
	DevelopmentModeEnabled bool `json:"development_mode_enabled"`
}

// ResolveAppData resolves the MongoDB Realm application based on the current working directory
// Empty data is successfully returned if this is called outside of a project directory
func ResolveAppData(wd string) (AppData, error) {
	appDir, appDirOK, appDirErr := ResolveAppDirectory(wd)
	if appDirErr != nil {
		return AppData{}, appDirErr
	}
	if !appDirOK {
		return AppData{}, nil
	}

	path := filepath.Join(appDir, realm.FileAppConfig)

	data, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return AppData{}, readErr
	}

	if len(data) == 0 {
		return AppData{}, fmt.Errorf("failed to read app data at %s", path)
	}

	var appData AppData
	if err := json.Unmarshal(data, &appData); err != nil {
		return AppData{}, err
	}
	return appData, nil
}

// ResolveAppDirectory searches upwards from the current working directory
// for the root directory of a MongoDB Realm application project
func ResolveAppDirectory(wd string) (string, bool, error) {
	wd, wdErr := filepath.Abs(wd)
	if wdErr != nil {
		return "", false, wdErr
	}

	for i := 0; i < maxDirectoryContainSearchDepth; i++ {
		path := filepath.Join(wd, realm.FileAppConfig)
		if _, err := os.Stat(path); err == nil {
			return filepath.Dir(path), true, nil
		}

		if wd == "/" {
			break
		}
		wd = filepath.Clean(filepath.Join(wd, ".."))
	}

	return "", false, nil
}

type AppDeploymentModel string

// String returns the deployment model display
func (dm AppDeploymentModel) String() string { return string(dm) }

// Type returns the AppDeploymentModel type
func (dm AppDeploymentModel) Type() string { return "string" }

// Set validates and sets the deployment model value
func (dm *AppDeploymentModel) Set(val string) error {
	newDeploymentModel := AppDeploymentModel(val)

	if !isValidDeploymentModel(newDeploymentModel) {
		return errInvalidDeploymentModel
	}

	*dm = newDeploymentModel
	return nil
}

// WriteAnswer validates and sets the deployment model value
func (dm *AppDeploymentModel) WriteAnswer(name string, value interface{}) error {
	var newDeploymentModel AppDeploymentModel

	switch v := value.(type) {
	case core.OptionAnswer:
		newDeploymentModel = AppDeploymentModel(v.Value)
	}

	if !isValidDeploymentModel(newDeploymentModel) {
		return errInvalidDeploymentModel
	}
	*dm = newDeploymentModel
	return nil
}

const (
	AppDeploymentModelNil    AppDeploymentModel = ""
	AppDeploymentModelGlobal AppDeploymentModel = "GLOBAL"
	AppDeploymentModelLocal  AppDeploymentModel = "LOCAL"
)

var (
	errInvalidDeploymentModel = func() error {
		allDeploymentModels := []string{AppDeploymentModelGlobal.String(), AppDeploymentModelLocal.String()}
		return fmt.Errorf("unsupported value, use one of [%s] instead", strings.Join(allDeploymentModels, ", "))
	}()
)

func isValidDeploymentModel(dm AppDeploymentModel) bool {
	switch dm {
	case
		AppDeploymentModelNil, // allow AppDeploymentModel to be optional
		AppDeploymentModelGlobal,
		AppDeploymentModelLocal:
		return true
	}
	return false
}

type AppLocation string

// String returns the AppLocation display
func (l AppLocation) String() string { return string(l) }

// Type returns the AppLocation type
func (l AppLocation) Type() string { return "string" }

// Set validates and sets the AppLocation value
func (l *AppLocation) Set(val string) error {
	newLocation := AppLocation(val)

	if !isValidLocation(newLocation) {
		return errInvalidLocation
	}

	*l = newLocation
	return nil
}

// WriteAnswer validates and sets the AppLocation value
func (l *AppLocation) WriteAnswer(name string, value interface{}) error {
	var newLocation AppLocation

	switch v := value.(type) {
	case core.OptionAnswer:
		newLocation = AppLocation(v.Value)
	}

	if !isValidLocation(newLocation) {
		return errInvalidLocation
	}
	*l = newLocation
	return nil
}

const (
	AppLocationNil       AppLocation = ""
	AppLocationVirginia  AppLocation = "US-VA"
	AppLocationOregon    AppLocation = "US-OR"
	AppLocationFrankfurt AppLocation = "DE-FF"
	AppLocationIreland   AppLocation = "IE"
	AppLocationSydney    AppLocation = "AU"
	AppLocationMumbai    AppLocation = "IN-MB"
	AppLocationSingapore AppLocation = "SG"
)

var (
	errInvalidLocation = func() error {
		allLocations := []string{
			AppLocationVirginia.String(),
			AppLocationOregon.String(),
			AppLocationFrankfurt.String(),
			AppLocationIreland.String(),
			AppLocationSydney.String(),
			AppLocationMumbai.String(),
			AppLocationSingapore.String(),
		}
		return fmt.Errorf("unsupported value, use one of [%s] instead", strings.Join(allLocations, ", "))
	}()
)

func isValidLocation(l AppLocation) bool {
	switch l {
	case
		AppLocationNil, // allow AppLocation to be optional
		AppLocationVirginia,
		AppLocationOregon,
		AppLocationFrankfurt,
		AppLocationIreland,
		AppLocationSydney,
		AppLocationMumbai,
		AppLocationSingapore:
		return true
	}
	return false
}
