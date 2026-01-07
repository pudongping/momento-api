package jwtToken

import (
	"github.com/golang-jwt/jwt/v4"
)

// GenJwtToken 生成 JWT Token
// args:
// secretKey 签名秘钥
// customKey 自定义存储用户 ID 的 key
// iat 签名生成时间
// seconds 签名过期时间，单位秒
// userId 当前用户 ID
//
// return:
// token 生成的 JWT Token 字符串
// err 生成过程中可能出现的错误
func GenJwtToken(secretKey, customKey string, iat, seconds, userId int64) (token string, err error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds                  // 签名过期时间
	claims["iat"] = iat                            // 生成签名的时间
	claims[customKey] = userId                     // 当前用户 id
	jwtInstance := jwt.New(jwt.SigningMethodHS256) // 使用 HS256 算法生成的 token
	jwtInstance.Claims = claims
	return jwtInstance.SignedString([]byte(secretKey)) // 生成签名字符串
}
