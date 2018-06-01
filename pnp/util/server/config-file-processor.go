package server

import (
	"io/ioutil"
	"github.com/BurntSushi/toml"
)

type Package struct {
	Name string `toml:"name"`
	Version string `toml:"version"`
	CheckInstalledCmd []string `toml:"checkInstalledCmd"`
	InstallFromFile string `toml:"installFromFile"`
	InstallInstructions []string `toml:"installInstructions"`
	UpdateRepo []string `toml:"updateRepo"`
}

type PackageInfo struct {
	Packages []Package `toml:"packages"`
}

func GetConfigFromToml(file string, configStruct interface{}) (err error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	
	if err = toml.Unmarshal(b, configStruct); err != nil {
		return err
	}

	return nil
}