package internal

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// ConfigStruct represents the overall configuration structure
type ConfigStruct struct {
	HyperledgerFabric struct {
		Fabric     string   `yaml:"fabric"`
		CA         string   `yaml:"ca"`
		Components []string `yaml:"components"`
	} `yaml:"hyperledger-fabric"`
	CCTool struct {
		Repo   string `yaml:"repo"`
		Branch string `yaml:"branch"`
	} `yaml:"cc-tool"`
	Chaincode struct {
		Name     string `yaml:"name"`
		Label    string `yaml:"label"`
		Version  string `yaml:"version"`
		Sequence int    `yaml:"sequence"`
	} `yaml:"chaincode"`
	Network struct {
		CCAs    bool     `yaml:"ccas"`
		OrgQnty int      `yaml:"org_qnty"`
		CCAAs   bool     `yaml:"ccaas"`
		CCasTls string   `yaml:"CcasTls"`
		Orgs    []string `yaml:"orgs"`
	} `yaml:"network"`
	System struct {
		Path string `yaml:"path"`
		Go   string `yaml:"go"`
	} `yaml:"system"`
}

func (CS *ConfigStruct) CreateConfigFile(filePath string) error {
	CS = &ConfigStruct{
		HyperledgerFabric: struct {
			Fabric     string   `yaml:"fabric"`
			CA         string   `yaml:"ca"`
			Components []string `yaml:"components"`
		}{
			Fabric:     "2.5.9",
			CA:         "1.5.12",
			Components: []string{"docker", "binary"},
		},
		CCTool: struct {
			Repo   string `yaml:"repo"`
			Branch string `yaml:"branch"`
		}{
			Repo:   "https://github.com/hyperledger-labs/cc-tools-demo",
			Branch: "main",
		},
		Chaincode: struct {
			Name     string `yaml:"name"`
			Label    string `yaml:"label"`
			Version  string `yaml:"version"`
			Sequence int    `yaml:"sequence"`
		}{
			Name:     "chainkun",
			Label:    "1.0",
			Version:  "0.1",
			Sequence: 1,
		},
		Network: struct {
			CCAs    bool     `yaml:"ccas"`
			OrgQnty int      `yaml:"org_qnty"`
			CCAAs   bool     `yaml:"ccaas"`
			CCasTls string   `yaml:"CcasTls"`
			Orgs    []string `yaml:"orgs"`
		}{
			CCAs:    false,
			OrgQnty: 3,
			CCAAs:   false,
			CCasTls: "1.3",
			Orgs:    []string{"org1MSP", "org2MSP", "org3MSP"},
		},
		System: struct {
			Path string `yaml:"path"`
			Go   string `yaml:"go"`
		}{
			Path: "/home/jony/.chainkun",
			Go:   "1.22.5",
		},
	}
	yamlData, err := yaml.Marshal(CS)
	if err != nil {
		return fmt.Errorf("error marshalling config to YAML: %w", err)
	}

	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func ReadCfgFile(cfgPath string) (ConfigStruct, error) {

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return ConfigStruct{}, fmt.Errorf("config file does not exist: %w", err)
	}
	yamlFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return ConfigStruct{}, fmt.Errorf("error reading config file: %w", err)
	}
	var cfg ConfigStruct
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return ConfigStruct{}, fmt.Errorf("error unmarshalling config file: %w", err)
	}
	return cfg, nil
}
