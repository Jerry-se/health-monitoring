package types

type WsHeader struct {
	Version   uint32 `json:"version"`   // 协议版本，暂时用 0
	Timestamp int64  `json:"timestamp"` // 时间戳
	Id        uint64 `json:"id"`        // 消息 ID
	Type      uint32 `json:"type"`      // 消息类型 WsMessageType
	PubKey    []byte `json:"pub_key"`   // 公钥，验证消息安全完整，暂时不需要
	Sign      []byte `json:"sign"`      // 签名，验证消息安全完整，暂时不需要
}

type WsRequest struct {
	WsHeader
	Body []byte `json:"body"` // 消息体，由 WsOnlineRequest 或者 WsMachineInfoRequest 编码后的字节流
}

type WsResponse struct {
	WsHeader
	Code    uint32 `json:"code"`    // 返回的错误码
	Message string `json:"message"` // 返回的错误描述
	Body    []byte `json:"body"`    // 消息体，暂时用不到
}

type WsMessageType uint32

const (
	WsMtOnline WsMessageType = iota + 1
	WsMtMachineInfo
)

type WsOnlineRequest struct {
	NodeId string `json:"node_id"`
}

type ModelInfo struct {
	Model string `json:"model" bson:"model"`
}

type WsMachineInfoRequest struct {
	Project        string      `json:"project" bson:"project"`
	Models         []ModelInfo `json:"models" bson:"models"`
	GPUName        string      `json:"gpu_name" bson:"gpu_name"`
	UtilizationGPU int         `json:"utilization_gpu" bson:"utilization_gpu"` // GPU 使用率，乘以 100 取整
	MemoryTotal    int64       `json:"memory_total" bson:"memory_total"`       // 显存总大小，单位 MB 或者 MiB
	MemoryUsed     int64       `json:"memory_used" bson:"memory_used"`         // 已用显存，单位 MB 或者 MiB
}
