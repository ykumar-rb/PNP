package common

import (
	"io/ioutil"
	"encoding/json"
	"github.com/BurntSushi/toml"
)

type Package struct {
	Name                  string        `json:"name"`
	Version               string        `json:"version"`
	CheckInstalledCmd     []string      `json:"checkInstalledCmd"`
	CheckInstalledVersion []string      `json:"checkInstalledVersion"`
	UnInstallInstructions []string      `json:"unInstallInstructions"`
	InstallFromFile       string        `json:"installFromFile"`
	InstallInstructions   []string      `json:"installInstructions"`
	UpdateRepo            []string      `json:"updateRepo"`
}

type PackageInfo struct {
	Packages []Package `json:"packages"`
}

type DeploySDP struct {
	ExportEnvParams []string `toml:"exportEnvParams"`
	CheckSDPUser []string `toml:"checkSDPUser"`
	CheckSDPMasterStatus []string `toml:"checkSDPMasterStatus"`
	CheckSDPSatelliteStatus []string `toml:"checkSDPSatelliteStatus"`
	IsSDPArtifactPresent []string `toml:"isSDPArtifactPresent"`
	IsSDPArtifactLatest []string `toml:"isSDPArtifactLatest"`
	DeleteOutdatedArtifact []string `toml:"deleteOutdatedArtifact"`
	DownloadLatestSDPArtifact []string `toml:"downloadLatestSDPArtifact"`
	ExtractSDPArtifact []string `toml:"extractSDPArtifact"`
	InstallSDPMaster []string `toml:"installSDPMaster"`
	InstallSDPSatellite []string `toml:"installSDPSatellite"`
	CleanInstallSDPMaster []string `toml:"cleanInstallSDPMaster"`
	CleanInstallSDPSatellite []string `toml:"cleanInstallSDPSatellite"`
	PlatformCleanUp []string `toml:"platformCleanUp"`
}

type PlatformDeploy struct {
	DeployInfo DeploySDP `toml:"deploySDP"`
}

func GetConfigFromJson(file string, configStruct interface{}) (err error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(b, configStruct); err != nil {
		return err
	}

	return nil
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