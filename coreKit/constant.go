package coreKit

// 这里的常量值可以根据实际项目需求进行调整
const (
	HeaderUserID    = "X-User-ID" // 从客户端请求头中需要带过来的用户ID字段
	CtxKeyJwtUserId = "jwtUserId" // ctx 中存储用户ID的 key
)
