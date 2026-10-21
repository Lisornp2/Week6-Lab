package mw

import (
	"net/http"
	"strings"
)

// CORS อนุญาต origin ที่อยู่ในรายชื่อ หรือ * สำหรับทุก origin
func CORS(allowed []string, next http.Handler) http.Handler {
	ok := map[string]bool{}
	allowAll := false

	for _, o := range allowed {
		trimmed := strings.TrimSpace(o)
		if trimmed == "*" {
			allowAll = true
		}
		ok[trimmed] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// ตรวจสอบว่าอนุญาตทั้งหมด หรือ Origin ตรงกับใน Whitelist
		if allowAll || ok[origin] {
			h := w.Header()
			if allowAll {
				h.Set("Access-Control-Allow-Origin", "*")
			} else {
				h.Set("Access-Control-Allow-Origin", origin)
				h.Add("Vary", "Origin")
			}
			h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			h.Set("Access-Control-Max-Age", "600")
		}

		// 处理 Preflight Request (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}