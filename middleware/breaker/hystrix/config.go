package hystrix

import "github.com/afex/hystrix-go/hystrix"

type CommandConfig = hystrix.CommandConfig

type Settings = hystrix.Settings

func Configure(cmds map[string]CommandConfig) {
	hystrix.Configure(cmds)
}

func ConfigureCommand(name string, config CommandConfig) {
	hystrix.ConfigureCommand(name, config)
}

func GetCircuitSettings() map[string]*Settings {
	return hystrix.GetCircuitSettings()
}
