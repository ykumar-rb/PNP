package server

import (
	"io/ioutil"
	"github.com/BurntSushi/toml"
)

type Package struct {
	Name string `toml:"name"`
	Version string `toml:"version"`
	CheckInstalledCmd []string `toml:"checkInstalledCmd"`
	CheckInstalledVersion []string `toml:"checkInstalledVersion"`
	UnInstallInstructions []string `toml:"unInstallInstructions"`
	InstallFromFile string `toml:"installFromFile"`
	InstallInstructions []string `toml:"installInstructions"`
	UpdateRepo []string `toml:"updateRepo"`
}

type PackageInfo struct {
	Packages []Package `toml:"packages"`
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