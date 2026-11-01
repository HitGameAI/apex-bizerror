module github.com/HitGameAI/apex-bizerror/pbx

go 1.25.1

require (
	github.com/HitGameAI/apex-bizerror v0.0.0
	github.com/HitGameAI/apex-toolbox v0.19.4
	github.com/stretchr/testify v1.12.1
	google.golang.org/grpc v1.83.0
)

require (
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

// 独立成模块的理由同 gwx：让主模块的 go.mod 不出现 protobuf / apex-toolbox/pkg/pbmo
//
// [EN] Same rationale as gwx: keep protobuf and apex-toolbox/pkg/pbmo out of the root module.

replace github.com/HitGameAI/apex-bizerror => ../
