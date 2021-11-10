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
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerOption configures the test logger built by NewLogger.
type LoggerOption interface {
	ApplyLoggerOption(*LoggerOptions)
}

type LoggerOptions struct {
	Level      zapcore.LevelEnabler
	ZapOptions []zap.Option
}

type LoggerOptionFunc func(*LoggerOptions)

func (f LoggerOptionFunc) ApplyLoggerOption(opts *LoggerOptions) {
	f(opts)
}

// Level controls which messages are logged by a test Logger built by
// NewLogger.
func Level(enab zapcore.LevelEnabler) LoggerOption {
	return LoggerOptionFunc(func(opts *LoggerOptions) {
		opts.Level = enab
	})
}

// WrapOptions adds zap.Option's to a test Logger built by NewLogger.
func WrapOptions(zapOpts ...zap.Option) LoggerOption {
	return LoggerOptionFunc(func(opts *LoggerOptions) {
		opts.ZapOptions = zapOpts
	})
}

// NewLogger builds a new Logger that logs all messages to the given
// testing.TB.
//
//   logger := zaptest.NewLogger(t)
//
// Use this with a *testing.T or *testing.B to get logs which get printed only
// if a test fails or if you ran go test -v.
//
// The returned logger defaults to logging debug level messages and above.
// This may be changed by passing a zaptest.Level during construction.
//
//   logger := zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel))
//
// You may also pass zap.Option's to customize test logger.
//
//   logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.AddCaller()))
func NewLogger(t TestingT, opts ...LoggerOption) *zap.Logger {
	cfg := LoggerOptions{
		Level: zapcore.DebugLevel,
	}
	for _, o := range opts {
		o.ApplyLoggerOption(&cfg)
	}

	writer := NewTestingWriter(t)
	zapOptions := []zap.Option{
		// Send zap errors to the same writer and mark the test as failed if
		// that happens.
		zap.ErrorOutput(writer.WithMarkFailed(true)),
	}
	zapOptions = append(zapOptions, cfg.ZapOptions...)

	return zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			writer,
			cfg.Level,
		),
		zapOptions...,
	)
}

// TestingWriter is a WriteSyncer that writes to the given testing.TB.
type TestingWriter struct {
	t TestingT

	// If true, the test will be marked as failed if this TestingWriter is
	// ever used.
	markFailed bool
}

func NewTestingWriter(t TestingT) TestingWriter {
	return TestingWriter{t: t}
}

// WithMarkFailed returns a copy of this TestingWriter with markFailed set to
// the provided value.
func (w TestingWriter) WithMarkFailed(v bool) TestingWriter {
	w.markFailed = v
	return w
}

func (w TestingWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	// Strip trailing newline because t.Log always adds one.
	p = bytes.TrimRight(p, "\n")

	// Note: t.Log is safe for concurrent use.
	w.t.Logf("%s", p)
	if w.markFailed {
		w.t.Fail()
	}

	return n, nil
}

func (w TestingWriter) Sync() error {
	return nil
}

func PrintLog(t *testing.T){
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	log := NewLogger(ts, WrapOptions(zap.AddCaller(), zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.DebugLevel),
		zap.Fields(zap.String("k1", "v1"))))
	mylog1:=&mylog{log: log}
	mylog1.Info("received work order")
	log.Debug("starting work")
	log.Warn("work may fail")
	log.Error("work failed", zap.Error(errors.New("great sadness")))

}

func PrintLog2(log  *zap.Logger){
	log.Debug("ssssssss work")
	log.Warn("eeeeeeee fail")
}

type mylog struct {
	log  *zap.Logger
}

func (l *mylog)Info(str string)  {
	l.log.Info(str)
}

// testLogSpy is a testing.TB that captures logged messages.
type testLogSpy struct {
	testing.TB

	failed   bool
	Messages []string
}

func newTestLogSpy(t testing.TB) *testLogSpy {
	return &testLogSpy{TB: t}
}

func (t *testLogSpy) Fail() {
	t.failed = true
}

func (t *testLogSpy) Failed() bool {
	return t.failed
}

func (t *testLogSpy) FailNow() {
	t.Fail()
	t.TB.FailNow()
}

func (t *testLogSpy) Logf(format string, args ...interface{}) {
	// Log messages are in the format,
	//
	//   2017-10-27T13:03:01.000-0700	DEBUG	your message here	{data here}
	//
	// We strip the first part of these messages because we can't really test
	// for the timestamp from these tests.
	m := fmt.Sprintf(format, args...)
	m = m[strings.IndexByte(m, '\t')+1:]
	t.Messages = append(t.Messages, m)
	t.TB.Log(m)
}

func (t *testLogSpy) AssertMessages(msgs ...string) {
	assert.Equal(t.TB, msgs, t.Messages, "logged messages did not match")
}

func (t *testLogSpy) AssertPassed() {
	t.assertFailed(false, "expected test to pass")
}

func (t *testLogSpy) AssertFailed() {
	t.assertFailed(true, "expected test to fail")
}

func (t *testLogSpy) assertFailed(v bool, msg string) {
	assert.Equal(t.TB, v, t.failed, msg)
}
