package httpRest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func UnauthorizedHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"code": http.StatusUnauthorized,
		"msg":  fmt.Sprintf("授权认证失败：%v", err),
	}
	bs, _ := json.Marshal(resp)
	if _, err := w.Write(bs); err != nil {
		panic(fmt.Sprintf("UnauthorizedHandler 写响应失败: %v", err))
	}
}
