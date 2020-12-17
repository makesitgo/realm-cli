package initialize

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2/core"
)

const (
	flagProject      = "project"
	flagProjectUsage = "the MongoDB cloud project id"

	flagFrom      = "from"
	flagFromShort = "s"
	flagFromUsage = "choose an application or template to initialize the Realm app with"

	flagName      = "name"
	flagNameShort = "n"
	flagNameUsage = "set the name of the Realm app to be initialized"

	flagDeploymentModel        = "deployment-model"
	flagDeploymentModelShort   = "d"
	flagDeploymentModelUsage   = `select the Realm app's deployment model, available options: ["global", "local"]`
	flagDeploymentModelDefault = deploymentModelGlobal

	flagLocation        = "location"
	flagLocationShort   = "l"
	flagLocationUsage   = `select the Realm app's location, available options: ["US-VA", "local"]`
	flagLocationDefault = locationVirginia
)

type deploymentModel string

// String returns the deployment model display
func (dm deploymentModel) String() string { return strings.ToUpper(string(dm)) }

// Type returns the deploymentModel type
func (dm deploymentModel) Type() string { return "string" }

// Set validates and sets the deployment model value
func (dm *deploymentModel) Set(val string) error {
	newDeploymentModel := deploymentModel(val)

	if !isValidDeploymentModel(newDeploymentModel) {
		return errInvalidDeploymentModel
	}

	*dm = newDeploymentModel
	return nil
}

// WriteAnswer validates and sets the deployment model value
func (dm *deploymentModel) WriteAnswer(name string, value interface{}) error {
	var newDeploymentModel deploymentModel

	switch v := value.(type) {
	case core.OptionAnswer:
		newDeploymentModel = deploymentModel(v.Value)
	}

	if !isValidDeploymentModel(newDeploymentModel) {
		return errInvalidDeploymentModel
	}
	*dm = newDeploymentModel
	return nil
}

const (
	deploymentModelNil    deploymentModel = ""
	deploymentModelGlobal deploymentModel = "global"
	deploymentModelLocal  deploymentModel = "local"
)

var (
	errInvalidDeploymentModel = func() error {
		allDeploymentModels := []string{deploymentModelGlobal.String(), deploymentModelLocal.String()}
		return fmt.Errorf("unsupported value, use one of [%s] instead", strings.Join(allDeploymentModels, ", "))
	}()
)

func isValidDeploymentModel(dm deploymentModel) bool {
	switch dm {
	case
		deploymentModelNil, // allow deploymentModel to be optional
		deploymentModelGlobal,
		deploymentModelLocal:
		return true
	}
	return false
}

type location string

// String returns the location display
func (l location) String() string { return string(l) }

// Type returns the location type
func (l location) Type() string { return "string" }

// Set validates and sets the location value
func (l *location) Set(val string) error {
	newLocation := location(val)

	if !isValidLocation(newLocation) {
		return errInvalidLocation
	}

	*l = newLocation
	return nil
}

// WriteAnswer validates and sets the location value
func (l *location) WriteAnswer(name string, value interface{}) error {
	var newLocation location

	switch v := value.(type) {
	case core.OptionAnswer:
		newLocation = location(v.Value)
	}

	if !isValidLocation(newLocation) {
		return errInvalidLocation
	}
	*l = newLocation
	return nil
}

const (
	locationNil       location = ""
	locationVirginia  location = "US-VA"
	locationOregon    location = "US-OR"
	locationFrankfurt location = "DE-FF"
	locationIreland   location = "IE"
	locationSydney    location = "AU"
	locationMumbai    location = "IN-MB"
	locationSingapore location = "SG"
)

var (
	errInvalidLocation = func() error {
		allLocations := []string{
			locationVirginia.String(),
			locationOregon.String(),
			locationFrankfurt.String(),
			locationIreland.String(),
			locationSydney.String(),
			locationMumbai.String(),
			locationSingapore.String(),
		}
		return fmt.Errorf("unsupported value, use one of [%s] instead", strings.Join(allLocations, ", "))
	}()
)

func isValidLocation(l location) bool {
	switch l {
	case
		locationNil, // allow location to be optional
		locationVirginia,
		locationOregon,
		locationFrankfurt,
		locationIreland,
		locationSydney,
		locationMumbai,
		locationSingapore:
		return true
	}
	return false
}
