package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/pbm-org/pbm/internal/builder"
	"github.com/pbm-org/pbm/internal/config"
)

func main() {
	init := flag.Bool("init", false, "init file pbm.yaml")
	build := flag.Bool("build", false, "build protobuf file")
	update := flag.Bool("update", false, "update remote dep proto file")
	clean := flag.Bool("clean", false, "clean remote cache dep proto file")
	v := flag.Bool("v", false, "show debug output")
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}
	if *v {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	if *init {
		err := config.InitConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		pbmCfg, err := config.PbmConfigFromFile("pbm.yaml")
		if err != nil {
			fmt.Println("read pbm.yaml failed", err)
			return
		}
		if *build {
			err = builder.CheckPbCfg(pbmCfg)
			if err != nil {
				fmt.Println(err)
				return
			}
			cmds, err := builder.PbBuildCmd(pbmCfg)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = builder.RunPbmCmd(cmds)
			if err != nil {
				fmt.Println(err)
				return
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
