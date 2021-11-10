package devin

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/studyzy/zap/zaptest/v2"

	rotatelogs "chainmaker.org/chainmaker/common/v2/log/file-rotatelogs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
func NewLogger(t zaptest.TestingT, opts ...zaptest.LoggerOption) *zap.Logger {
	cfg := zaptest. LoggerOptions{
		Level: zapcore.DebugLevel,
	}
	for _, o := range opts {
		o.ApplyLoggerOption(&cfg)
	}

	writer := zaptest. NewTestingWriter(t)
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

func TestTestLoggerSupportsWrappedZapOptions(t *testing.T) {
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	//log := NewLogger(ts, zaptest.WrapOptions(zap.AddCaller()))
	log:=newLogger(zap.NewAtomicLevel())
	log.Info("received work order")
	log.Debug("starting work")
	log.Warn("work may fail")
	log.Error("work failed", zap.Error(errors.New("great sadness")))

	zaptest.PrintLog2(log)
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
func  TestAgain(t *testing.T)  {
	zaptest.PrintLog(t)
}


func newLogger( level zap.AtomicLevel) *zap.Logger {
	var (
		hook io.Writer
		//ok   bool
		//err  error
	)
	hook, _ = getHook("./log", 10,1, 1000)


	var syncer zapcore.WriteSyncer

		syncer = zapcore.AddSync(hook)


	var encoderConfig zapcore.EncoderConfig

		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "line",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    CustomLevelEncoder,
			EncodeTime:     CustomTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}


	var encoder zapcore.Encoder

		encoder = zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(
		encoder,
		syncer,
		level,
	)

	//chainId := "chain1"

	var name string
			name = "Module1"


	logger := zap.New(core).Named(name)
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)


		logger = logger.WithOptions(zap.AddCaller())


		logger = logger.WithOptions(zap.AddStacktrace(zapcore.InfoLevel))

	logger = logger.WithOptions(zap.AddCallerSkip(1))
	return logger
}


func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
//
//// nolint: deadcode, unused
//func showColor(color color, msg string) string {
//	return fmt.Sprintf("\033[%dm%s\033[0m", int(color), msg)
//}
//
//func showColorBold(color color, msg string) string {
//	return fmt.Sprintf("\033[%d;1m%s\033[0m", int(color), msg)
//}
//
//func getColorChainId(chainId string) string {
//	c := crc32.ChecksumIEEE([]byte(chainId))
//	color := colorList[int(c)%len(colorList)]
//	return showColorBold(color, chainId)
//}

func getHook(filename string, maxAge, rotationTime int, rotationSize int64) (io.Writer, error) {

	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H",
		rotatelogs.WithRotationTime(time.Hour*time.Duration(rotationTime)),
		//filename+".%Y%m%d%H%M",
		//rotatelogs.WithRotationSize(rotationSize*ROTATION_SIZE_MB),
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxAge)),
	)

	if err != nil {
		return nil, err
	}

	return hook, nil
}