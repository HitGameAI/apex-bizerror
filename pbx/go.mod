module github.com/kamalyes/go-bizerror/pbx

go 1.25.1

require (
	github.com/kamalyes/go-bizerror v0.2.0
	github.com/kamalyes/go-pbmo v0.3.0
	github.com/stretchr/testify v1.12.1
	google.golang.org/grpc v1.83.0
)

require (
	github.com/kamalyes/go-i18n v0.3.1 // indirect
	github.com/kamalyes/go-toolbox v0.18.1 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

// 独立成模块的理由同 gwx：让主模块的 go.mod 不出现 go-pbmo 与 protobuf
// replace 只在本仓库内生效，供 tag 发布前的本地与 CI 编译
// [EN] Same rationale as gwx; the replace applies only inside this repository
replace github.com/kamalyes/go-bizerror => ../
