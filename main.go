package main

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/JonyBepary/chainkun/internal"
)

func main() {
	// CHECK ENV VARIABLE OR DEFAULT CONFIG PATH AVAILABLE
	// Check if the CONFIG_FILE environment variable is set
	var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()

	if configPath, ok := os.LookupEnv("CHAINKUN_ENV"); ok {
		internal.InitChaincode(configPath)
	} else {
		internal.InitChaincode(path.Join(os.Getenv("HOME"), ".chainkun"))
	}
	// }()

	wg.Wait()

	argsWithoutProg := os.Args[1:]
	cfgPath := internal.GetCfgPath(argsWithoutProg)

	cfg, err := internal.InitConfig(cfgPath)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Println(cfg)

}
