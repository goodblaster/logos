package logos

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogos_FromContext(t *testing.T) {
	buf := &bytes.Buffer{}
	customLog := NewLogger(LevelDebug, JSONFormatter(), buf)
	SetDefaultLogger(NewLogger(LevelInfo, JSONFormatter(), &bytes.Buffer{}))

	// Store logger in context
	ctx := WithLogger(context.Background(), customLog)

	// Retrieve logger from context
	logger := FromContext(ctx)
	logger.Log(LevelDebug, "test message")

	m := Map(buf)
	assert.Equal(t, "test message", m["msg"])
	assert.Equal(t, "debug", m["level"])
}

func TestLogos_FromContext_NoLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	defaultLog := NewLogger(LevelInfo, JSONFormatter(), buf)
	SetDefaultLogger(defaultLog)

	// Context without logger should return DefaultLogger
	ctx := context.Background()
	logger := FromContext(ctx)
	logger.Log(LevelInfo, "default message")

	m := Map(buf)
	assert.Equal(t, "default message", m["msg"])
	assert.Equal(t, "info", m["level"])
}

func TestLogos_FromContext_NilContext(t *testing.T) {
	buf := &bytes.Buffer{}
	defaultLog := NewLogger(LevelInfo, JSONFormatter(), buf)
	SetDefaultLogger(defaultLog)

	// Nil context should return DefaultLogger
	var nilCtx context.Context //nolint:staticcheck // intentionally testing nil context behavior
	logger := FromContext(nilCtx)
	logger.Log(LevelInfo, "nil context message")

	m := Map(buf)
	assert.Equal(t, "nil context message", m["msg"])
}

func TestLogos_WithLogger(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	log1 := NewLogger(LevelDebug, JSONFormatter(), buf1)
	log2 := NewLogger(LevelInfo, JSONFormatter(), buf2)

	// Store first logger in context
	ctx := WithLogger(context.Background(), log1)
	retrieved1 := FromContext(ctx)
	retrieved1.Log(LevelDebug, "from log1")

	m1 := Map(buf1)
	assert.Equal(t, "from log1", m1["msg"])

	// Store second logger in context (replaces first)
	ctx = WithLogger(ctx, log2)
	retrieved2 := FromContext(ctx)
	retrieved2.Log(LevelInfo, "from log2")

	m2 := Map(buf2)
	assert.Equal(t, "from log2", m2["msg"])

	// First buffer should not have new message
	assert.Empty(t, buf1.String())
}
