module github.com/kamalyes/go-bizerror/gwx

go 1.25.1

require (
	github.com/kamalyes/go-bizerror v0.2.0
	github.com/kamalyes/go-i18n v0.3.0
	github.com/kamalyes/go-rpc-gateway v1.6.9
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/grpc v1.83.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// 本子包独立成模块，目的是让 go-bizerror 主模块的 go.mod 不出现 go-rpc-gateway：
// 同一模块内无法按包区分依赖，只有拆模块才能让「只要 gRPC 错误转换」的服务
// 连 go.sum 都看不到网关那一整套
//
// 下面的 replace 只在本仓库内生效（Go 会忽略依赖方 go.mod 里的 replace），
// 作用是让 tag 发布之前，本地与 CI 就能对着主模块的当前代码编译测试
//
// [EN] A separate module so the parent's go.mod never mentions the gateway;
// dependencies cannot be scoped per package inside one module. The replace below
// applies only inside this repository, since Go ignores replace directives coming
// from dependency modules
replace github.com/kamalyes/go-bizerror => ../
