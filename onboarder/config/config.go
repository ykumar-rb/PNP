package config

type ConfigEnvironment struct {
	EnvironmentName string `json:"EnvironmentName"`
	Mac []string `json:"Mac"`
	InstructionFileName string `json:"InstructionFileName"`
	AutoUpdate bool	`json:"AutoUpdate"`
}

type ClientEnv struct {
	ClientConfigFile string
	AutoUpdate bool
}
