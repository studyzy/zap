// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package zaptest

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"studyzy/zap/ztest"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
)

func TestTestLogger(t *testing.T) {
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	log := NewLogger(ts)

	log.Info("received work order")
	log.Debug("starting work")
	log.Warn("work may fail")
	log.Error("work failed", zap.Error(errors.New("great sadness")))

	assert.Panics(t, func() {
		log.Panic("failed to do work")
	}, "log.Panic should panic")

	ts.AssertMessages(
		"INFO	received work order",
		"DEBUG	starting work",
		"WARN	work may fail",
		`ERROR	work failed	{"error": "great sadness"}`,
		"PANIC	failed to do work",
	)
}

func TestTestLoggerSupportsLevels(t *testing.T) {
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	log := NewLogger(ts, Level(zap.WarnLevel))

	log.Info("received work order")
	log.Debug("starting work")
	log.Warn("work may fail")
	log.Error("work failed", zap.Error(errors.New("great sadness")))

	assert.Panics(t, func() {
		log.Panic("failed to do work")
	}, "log.Panic should panic")

	ts.AssertMessages(
		"WARN	work may fail",
		`ERROR	work failed	{"error": "great sadness"}`,
		"PANIC	failed to do work",
	)
}

func PrintLog(t *testing.T){
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	log := NewLogger(ts, WrapOptions(zap.AddCaller(), zap.Fields(zap.String("k1", "v1"))))

	log.Info("received work order")
	log.Debug("starting work")
	log.Warn("work may fail")
	log.Error("work failed", zap.Error(errors.New("great sadness")))

}

func TestTestLoggerSupportsWrappedZapOptions(t *testing.T) {
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	log := NewLogger(ts, WrapOptions(zap.AddCaller(), zap.Fields(zap.String("k1", "v1"))))

	log.Info("received work order")
	log.Debug("starting work")
	log.Warn("work may fail")
	log.Error("work failed", zap.Error(errors.New("great sadness")))

	assert.Panics(t, func() {
		log.Panic("failed to do work")
	}, "log.Panic should panic")

	ts.AssertMessages(
		`INFO	zaptest/logger_test.go:89	received work order	{"k1": "v1"}`,
		`DEBUG	zaptest/logger_test.go:90	starting work	{"k1": "v1"}`,
		`WARN	zaptest/logger_test.go:91	work may fail	{"k1": "v1"}`,
		`ERROR	zaptest/logger_test.go:92	work failed	{"k1": "v1", "error": "great sadness"}`,
		`PANIC	zaptest/logger_test.go:95	failed to do work	{"k1": "v1"}`,
	)
}

func TestTestingWriter(t *testing.T) {
	ts := newTestLogSpy(t)
	w := NewTestingWriter(ts)

	n, err := io.WriteString(w, "hello\n\n")
	assert.NoError(t, err, "WriteString must not fail")
	assert.Equal(t, 7, n)
}

func TestTestLoggerErrorOutput(t *testing.T) {
	// This test verifies that the test logger logs internal messages to the
	// testing.T and marks the test as failed.

	ts := newTestLogSpy(t)
	defer ts.AssertFailed()

	log := NewLogger(ts)

	// Replace with a core that fails.
	log = log.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(zapcore.AddSync(ztest.FailWriter{})),
			zapcore.DebugLevel,
		)
	}))

	log.Info("foo") // this fails

	if assert.Len(t, ts.Messages, 1, "expected a log message") {
		assert.Regexp(t, `write error: failed`, ts.Messages[0])
	}
}
