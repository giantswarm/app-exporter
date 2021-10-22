//go:build functional || smoke
// +build functional smoke

package ats

import "github.com/giantswarm/microerror"

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}
