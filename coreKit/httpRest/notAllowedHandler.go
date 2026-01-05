package httpRest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func NotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		resp := map[string]interface{}{
			"code": http.StatusMethodNotAllowed,
			"msg":  fmt.Sprintf("路由[%s]，当前不支持 %s 方法", r.URL.Path, r.Method),
		}
		bs, _ := json.Marshal(resp)
		if _, err := w.Write(bs); err != nil {
			panic(fmt.Sprintf("NotAllowedHandler 写响应失败: %v", err))
		}
	})
}
