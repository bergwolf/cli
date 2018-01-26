package main

import (
	runvcli "github.com/hyperhq/runv/cli"
	_ "github.com/hyperhq/runv/cli/nsenter"
	"github.com/urfave/cli"
)

var nsListenCommand = cli.Command{
	Name:     "network-nslisten",
	Usage:    "[internal command] collection net namespace's network configuration",
	HideHelp: true,
	Before: func(context *cli.Context) error {
		return cmdPrepare(context, false, false)
	},
	Action: func(context *cli.Context) error {
		runvcli.DoListen()
		return nil
	},
}
