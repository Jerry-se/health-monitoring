package types

type ErrorCode uint32

const (
	ErrCodeParam    ErrorCode = iota + 1 // 参数错误
	ErrCodeParse                         // 解析错误
	ErrCodeSign                          // 签名错误
	ErrCodeDatabase                      // 数据库错误
)
