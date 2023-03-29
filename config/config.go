package config

const configPath = "./cli-config.yml"

type Config struct {
	ExecutionClient string
	ConsensusClient string
}

func CreateConfigFile(selectedExecution, selectedConsensus string) (err error) {
	return
}
