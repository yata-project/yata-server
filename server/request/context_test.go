package request

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	// Test the logger returned when no logger is set on the context.
	assert.Equal(t, logrus.Fields{"requestID": "UNKNOWN"}, Logger(context.Background()).Data)

	// Test when a logger is set on the context.
	lggr := logrus.WithField("foo", "bar")
	assert.Equal(t, lggr, Logger(WithLogger(context.Background(), lggr)))
}
