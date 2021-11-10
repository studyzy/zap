module github.com/studyzy/zap/zaptest/devin/v2

go 1.15

require (
	chainmaker.org/chainmaker/common/v2 v2.1.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.19.1
	github.com/studyzy/zap/zaptest/v2 v2.0.0
)

replace (
	studyzy/zap => ../../
	github.com/studyzy/zap/zaptest/v2 => ../

)
