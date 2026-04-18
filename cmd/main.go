package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pbm-org/pbm/internal/builder"
	"github.com/pbm-org/pbm/internal/config"
)

func main() {
	init := flag.Bool("init", false, "init file pbm.yaml")
	build := flag.Bool("build", false, "build protobuf file")
	update := flag.Bool("update", false, "update remote dep proto file")
	clean := flag.Bool("clean", false, "clean remote cache dep proto file")
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}
	if *init {
		err := config.InitConfig()
		if err != nil {
			fmt.Println("init failed", err)
			os.Exit(1)
		}
	} else {

		pbmCfg, err := config.PbmConfigFromFile("pbm.yaml")
		if err != nil {
			fmt.Println("read pbm.yaml failed", err)
			os.Exit(1)
		}
		if *build {
			err = builder.CheckPbCfg(pbmCfg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cmds, err := builder.PbBuildCmd(pbmCfg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = builder.PbmCmd(cmds)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if *update {
			err = builder.CleabPbmDep(pbmCfg)
			if err != nil {
				fmt.Println("clean pbm dep failed", err)
				os.Exit(1)
			}
			err = builder.CheckPbCfg(pbmCfg)
			if err != nil {
				fmt.Println("check pbm failed", err)
				os.Exit(1)
			}
		}
		if *clean {
			err = builder.CleabPbmDep(pbmCfg)
			if err != nil {
				fmt.Println("clean pbm dep failed", err)
				os.Exit(1)
			}
		}
	}

}
