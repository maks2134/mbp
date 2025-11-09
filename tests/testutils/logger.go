package testutils

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill"
)

func SetupTestLogger(t *testing.T) watermill.LoggerAdapter {
	return watermill.NewStdLogger(false, false)
}
