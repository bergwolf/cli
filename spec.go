package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	runvcli "github.com/hyperhq/runv/cli"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

var specCommand = cli.Command{
	Name:  "spec",
	Usage: "create a new specification file",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Usage: "path to the root of the bundle directory",
		},
	},
	Before: func(context *cli.Context) error {
		return cmdPrepare(context, false, false)
	},
	Action: func(context *cli.Context) {
		spec := specs.Spec{
			Version: specs.Version,
			Root: &specs.Root{
				Path:     "rootfs",
				Readonly: true,
			},
			Process: &specs.Process{
				Terminal: true,
				User:     specs.User{},
				Args: []string{
					"sh",
				},
				Env: []string{
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
					"TERM=xterm",
				},
				Cwd: "/",
			},
			Hostname: "shell",
			Linux: &specs.Linux{
				Resources: &specs.LinuxResources{},
			},
		}

		checkNoFile := func(name string) error {
			_, err := os.Stat(name)
			if err == nil {
				return fmt.Errorf("File %s exists. Remove it first", name)
			}
			if !os.IsNotExist(err) {
				return err
			}
			return nil
		}

		bundle := context.String("bundle")
		if bundle != "" {
			if err := os.Chdir(bundle); err != nil {
				fmt.Printf("Failed to chdir to bundle dir:%s\nerror:%v\n", bundle, err)
				return
			}
		}
		if err := checkNoFile(runvcli.SpecConfig); err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		data, err := json.MarshalIndent(&spec, "", "\t")
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		if err := ioutil.WriteFile(runvcli.SpecConfig, data, 0666); err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
	},
}
