package file

const (
	// TODO: should we limit ourselves to working directory only?
	WorkingDir    = "."
	ConfigRootDir = WorkingDir + "/configs"
	SecretsDir    = ConfigRootDir + "/shared/secrets"
	JwtSecretPath = ConfigRootDir + "/shared/secrets/jwt.hex"
	PidDir        = "/tmp" // until a script for /tmp/lukso or other dir is provided
	ClientsDir    = WorkingDir + "/clients"
)
