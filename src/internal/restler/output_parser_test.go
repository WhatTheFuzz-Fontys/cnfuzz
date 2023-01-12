package restler

import (
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"testing"
)

func TestParseRestlerOutput(t *testing.T) {
	l := logger.CreateDebugLogger()
	// relative path from src/cmd/restlerwrapper/main.go
	result := ParseRestlerOutput(l, "../../../hack/example_restler_output")
	// TODO add some useful asserts
	assert.NotNil(t, result)
}
