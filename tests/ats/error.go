//go:build functional
// +build functional

package ats

import "github.com/giantswarm/microerror"

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}
