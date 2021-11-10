module github.com/studyzy/zap/zaptest/v2

go 1.15

require (
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.10.0
	studyzy/zap v0.0.0-00010101000000-000000000000
)

require (
	github.com/benbjohnson/clock v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace studyzy/zap => ../
