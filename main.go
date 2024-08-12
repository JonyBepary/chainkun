package main

import (
	"fmt"
	"os"
	"path"

	"github.com/JonyBepary/chainkun/internal"
	"github.com/davecgh/go-spew/spew"
)

// check .chainkun folder exists and .chainkun/config.yaml exists
// if not create .chainkun folder and create config.yaml
func CheckChainkunFolder(roo string) bool {
	// check if chainkun folder exists
	_, err := os.Stat(path.Join(roo, "config.yaml"))
	if err == nil {
		fmt.Println("CHAINKUN FOLDER and Config file EXISTS")
		return true
	}

	// check config.yaml exists
	return false

}

func main() {
	// CHECK ENV VARIABLE OR DEFAULT CONFIG PATH AVAILABLE
	// Check if the CONFIG_FILE environment variable is set
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	if ok := CheckChainkunFolder(path.Join(os.Getenv("HOME"), ".chainkun")); !ok {
		fmt.Println("CHAINKUN FOLDER DOES NOT EXISTS")
		if configPath, ok := os.LookupEnv("CHAINKUN_ENV"); ok {
			internal.InitChaincode(configPath)
		} else {
			internal.InitChaincode(path.Join(os.Getenv("HOME"), ".chainkun"))
		}
	}
	var cfg internal.ConfigStruct
	var err error
	// if configFile pass then read the config file else read the default config file
	if len(os.Args) > 1 {
		cfg, err = internal.ReadCfgFile(os.Args[1])
		if err != nil {
			fmt.Println("Error reading config file:", err)
		}
	} else {
		cfg, err = internal.ReadCfgFile(path.Join(os.Getenv("HOME"), ".chainkun", "config.yaml"))
		if err != nil {
			fmt.Println("Error reading config file:", err)
		}
	}
	spew.Dump(cfg)
	internal.FabricDownload(&cfg)

}
