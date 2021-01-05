// +build !jsoniter

package json

import "encoding/json"

// 导出 encoding/json 方法到 goutils/json

var (
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

type RawMessage = json.RawMessage
