package mw

import (
    "net/http"
    "strings"
)

// CORS อนุญาตเฉพาะ origin ที่อยู่ในรายชื่อ
func CORS(allowed []string, next http.Handler) http.Handler {
    ok := map[string]bool{}
    for _, o := range allowed {
        ok[strings.TrimSpace(o)] = true
    }
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if ok[origin] {
            h := w.Header()
            h.Set("Access-Control-Allow-Origin", origin)
            h.Add("Vary", "Origin")
            h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
            h.Set("Access-Control-Max-Age", "600")
        }
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // อนุญาต Origin
        w.Header().Set("Access-Control-Allow-Origin", "https://moviehub-650710873.vercel.app")
        
        // อนุญาต Methods ทุกตัวที่ frontend ใช้
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        
        // *** สำคัญมาก: ต้องเพิ่ม Authorization เข้าไปด้วย ***
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // ถ้าเบราว์เซอร์ยิง Preflight request (OPTIONS) มา ให้ตอบ 200 OK กลับไปทันที
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}