// +build k8srequired

package bootstrap

import (
	"context"
	"testing"
)

// TestMetrics TODO
//
func TestMetrics(t *testing.T) {
	ctx := context.Background()

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", "TODO: test metrics")
	}
}
