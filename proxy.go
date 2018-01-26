package main

import (
	"fmt"
	"net"
	"os"

	"github.com/golang/glog"
	"github.com/hyperhq/runv/agent"
	"github.com/hyperhq/runv/agent/proxy"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var proxyCommand = cli.Command{
	Name:     "proxy",
	Usage:    "[internal command] proxy hyperstart API into vm",
	HideHelp: true,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "vmid",
			Usage: "the vm name",
		},
		cli.StringFlag{
			Name:  "hyperstart-ctl-sock",
			Usage: "the vm's ctl sock address to be connected",
		},
		cli.StringFlag{
			Name:  "hyperstart-stream-sock",
			Usage: "the vm's stream sock address to be connected",
		},
		cli.StringFlag{
			Name:  "proxy-hyperstart",
			Usage: "gprc sock address to be created for proxying hyperstart API",
		},
	},
	Before: func(context *cli.Context) error {
		return cmdPrepare(context, false, false)
	},
	Action: func(context *cli.Context) (err error) {
		if context.String("vmid") == "" || context.String("hyperstart-ctl-sock") == "" ||
			context.String("hyperstart-stream-sock") == "" || context.String("proxy-hyperstart") == "" {
			return err
		}

		glog.Infof("agent.NewJsonBasedHyperstart")
		h, _ := agent.NewJsonBasedHyperstart(context.String("vmid"), context.String("hyperstart-ctl-sock"), context.String("hyperstart-stream-sock"), 1, true, false)

		var s *grpc.Server
		grpcSock := context.String("proxy-hyperstart")
		glog.Infof("proxy.NewServer")
		s, err = proxy.NewServer(grpcSock, h)
		if err != nil {
			glog.Errorf("proxy.NewServer() failed with err: %#v", err)
			return err
		}
		if _, err := os.Stat(grpcSock); !os.IsNotExist(err) {
			return fmt.Errorf("%s existed, someone may be in service", grpcSock)
		}
		glog.Infof("net.Listen() to grpcsock: %s", grpcSock)
		l, err := net.Listen("unix", grpcSock)
		if err != nil {
			glog.Errorf("net.Listen() failed with err: %#v", err)
			return err
		}

		glog.Infof("proxy: grpc api on %s", grpcSock)
		if err = s.Serve(l); err != nil {
			glog.Errorf("proxy serve grpc with error: %v", err)
		}

		return err
	},
}
