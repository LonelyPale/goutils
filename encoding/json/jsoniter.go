// +build jsoniter

package json

import "github.com/json-iterator/go"

// 导出 json-iterator/go 方法到 goutils/json
// goutils 和 gin 使用 jsoniter 代替默认的 json，需要在 build 时，加入 -tags jsoniter。
// go build -tags jsoniter -o main.exe
// 详见：gin/internal/json/jsoniter.go

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Marshal is exported by goutils/json package.
	Marshal = json.Marshal

	// Unmarshal is exported by goutils/json package.
	Unmarshal = json.Unmarshal

	// MarshalIndent is exported by goutils/json package.
	MarshalIndent = json.MarshalIndent

	// NewDecoder is exported by goutils/json package.
	NewDecoder = json.NewDecoder

	// NewEncoder is exported by goutils/json package.
	NewEncoder = json.NewEncoder
)

type RawMessage = jsoniter.RawMessage
