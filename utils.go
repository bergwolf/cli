package main

import (
	"errors"
	"os"
	"strings"
	"syscall"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

func firstExistingFile(candidates []string) string {
	for _, file := range candidates {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}
	return ""
}

func getDefaultBundlePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return cwd
}

// find it from spec.Process.Env, runv's env(todo) and context.GlobalString
func chooseKernelFromConfigs(context *cli.Context, spec *specs.Spec) string {
	for k, env := range spec.Process.Env {
		slices := strings.Split(env, "=")
		if len(slices) == 2 && slices[0] == "hypervisor.kernel" {
			// remove kernel env because this is only allow to be used by runv
			spec.Process.Env = append(spec.Process.Env[:k], spec.Process.Env[k+1:]...)
			return slices[1]
		}
	}
	return context.GlobalString("kernel")
}

func chooseInitrdFromConfigs(context *cli.Context, spec *specs.Spec) string {
	for k, env := range spec.Process.Env {
		slices := strings.Split(env, "=")
		if len(slices) == 2 && slices[0] == "hypervisor.initrd" {
			// remove initrd env because this is only allow to be used by runv
			spec.Process.Env = append(spec.Process.Env[:k], spec.Process.Env[k+1:]...)
			return slices[1]
		}
	}
	return context.GlobalString("initrd")
}

func chooseBiosFromConfigs(context *cli.Context, spec *specs.Spec) string {
	for k, env := range spec.Process.Env {
		slices := strings.Split(env, "=")
		if len(slices) == 2 && slices[0] == "hypervisor.bios" {
			// remove bios env because this is only allow to be used by runv
			spec.Process.Env = append(spec.Process.Env[:k], spec.Process.Env[k+1:]...)
			return slices[1]
		}
	}
	return context.GlobalString("bios")
}

func chooseCbfsFromConfigs(context *cli.Context, spec *specs.Spec) string {
	for k, env := range spec.Process.Env {
		slices := strings.Split(env, "=")
		if len(slices) == 2 && slices[0] == "hypervisor.cbfs" {
			// remove cbfs env because this is only allow to be used by runv
			spec.Process.Env = append(spec.Process.Env[:k], spec.Process.Env[k+1:]...)
			return slices[1]
		}
	}
	return context.GlobalString("cbfs")
}

func osProcessWait(process *os.Process) (int, error) {
	state, err := process.Wait()
	if err != nil {
		return -1, err
	}
	if state.Success() {
		return 0, nil
	}

	ret := -1
	if status, ok := state.Sys().(syscall.WaitStatus); ok {
		ret = status.ExitStatus()
		if ret != 0 {
			err = errors.New("")
		}
	}
	return ret, err
}
