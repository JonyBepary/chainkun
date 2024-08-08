package internal

import (
	"fmt"
	"os"
	"path"
)

func checkFileExists(filename string) (string, bool) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return "", false
	}
	filePath := path.Join(pwd, filename)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, true
	} else if os.IsNotExist(err) {
		return "", false
	} else {
		fmt.Println("Error checking file existence:", err)
		return "", false
	}
}

func GetCfgPath(argsWithoutProg []string) string {
	if len(argsWithoutProg) > 0 {
		if cfgPath, ok := checkFileExists(argsWithoutProg[0]); ok {
			return cfgPath
		}
	} else {
		if cfgPath, ok := checkFileExists("config.yaml"); ok {
			return cfgPath
		} else {
			fmt.Println("ACCESSING ROOT CONFIG FILE ...")
			if cfg, ok := os.LookupEnv("CHAINKUN_ENV"); ok {
				return cfg
			} else {
				return path.Join(os.Getenv("HOME"), ".chainkun", "config.yaml")
			}
		}
	}
	return ""
}

func InitChaincode(root string) {
	fmt.Println("INIT CHAINCODE")
	// check if chainkun folder exists
	_, err := os.Stat(root)
	if err == nil {
		fmt.Println("CHAINKUN FOLDER EXISTS")
	} else {
		// create folder and create chaincode
		err := os.MkdirAll(root, 0755)
		if err != nil {
			panic(err)
		}
	}
	// check if config exists
	_, err = os.Stat(path.Join(root, "config.yaml"))
	if err == nil {
		fmt.Println("CONFIG FILE EXISTS")
	} else {
		var CS ConfigStruct
		CS.CreateConfigFile(path.Join(root, "config.yaml"))
		err := downloadBinaries(CS)
		if err != nil {
			panic(err)
		}

	}
	//
}

func downloadBinaries(config ConfigStruct) error {
	// CHECK ENV VARIABLE OR DEFAULT CONFIG PATH AVAILABLE
	// Check if the CONFIG_FILE environment variable is set
	var configPath string
	var ok bool
	if configPath, ok = os.LookupEnv("CHAINKUN_ENV"); !ok {

		configPath = path.Join(os.Getenv("HOME"), ".chainkun")
	}

	_, err := os.Stat(path.Join(configPath, config.HyperledgerFabric.Fabric))
	if err == nil {
		fmt.Println("FABRIC BINARIES ALREADY DOWNLOADED")
		return nil
	}

	// create folder and create mkdir
	err = os.MkdirAll(path.Join(configPath, "fabric_bin", config.HyperledgerFabric.Fabric), 0755)
	if err != nil {
		panic(err)
	}

	return nil
}
