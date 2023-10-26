//go:build functional
// +build functional

package ats

import (
	"fmt"
	"os"
)

const (
	EnvVarKubeConfigPath = "KUBECONFIG"
)

var (
	kubeConfigPath string
)

func init() {
	kubeConfigPath = os.Getenv(EnvVarKubeConfigPath)
	if kubeConfigPath == "" {
		panic(fmt.Sprintf("env var '%s' must not be empty", EnvVarKubeConfigPath))
	}
}

func KubeConfigPath() string {
	return kubeConfigPath
}
