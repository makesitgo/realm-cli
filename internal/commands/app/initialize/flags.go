package initialize

import (
	"github.com/10gen/realm-cli/internal/cli"
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
	flagDeploymentModelDefault = cli.AppDeploymentModelGlobal

	flagLocation        = "location"
	flagLocationShort   = "l"
	flagLocationUsage   = `select the Realm app's location, available options: ["US-VA", "local"]`
	flagLocationDefault = cli.AppLocationVirginia
)
