module github.com/HitGameAI/apex-bizerror/gwx

go 1.25.1

require (
	github.com/HitGameAI/apex-bizerror v0.0.0
	github.com/HitGameAI/apex-rpc-gateway v1.9.8
	github.com/HitGameAI/apex-toolbox v0.19.4
	github.com/stretchr/testify v1.12.1
)

require (
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260831171406-18b4a7587f8a // indirect
	google.golang.org/grpc v1.84.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

// 本子包独立成模块，目的是让 apex-bizerror 主模块的 go.mod 不出现 apex-rpc-gateway：
// 同一模块内无法按包区分依赖，只有拆模块才能让「只要 gRPC 错误转换」的服务
// 连 go.sum 都看不到网关那一整套
//
// [EN] This subpackage stays a separate module so the apex-bizerror root module's go.mod
// never pulls in apex-rpc-gateway; Go cannot scope dependencies per package inside one module.

replace github.com/HitGameAI/apex-bizerror => ../
