package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"moviehub/mw"
	"moviehub/store"
)

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	movieID, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}

	memberID, ok := mw.MemberID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "กรุณาเข้าสู่ระบบ")
		return
	}

	var input struct {
		Text string `json:"text"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "ข้อมูลรีวิวไม่ถูกต้อง")
		return
	}

	input.Text = strings.TrimSpace(input.Text)
	if len([]rune(input.Text)) < 10 {
		writeError(w, http.StatusBadRequest, "REVIEW_TOO_SHORT", "รีวิวต้องมีอย่างน้อย 10 ตัวอักษร")
		return
	}

	review, err := h.Store.SaveReview(r.Context(), memberID, movieID, input.Text)
	if err != nil {
		if err == store.ErrNotFound {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "ไม่พบข้อมูลที่ต้องการ")
			return
		}
		internal(w, err)
		return
	}

	writeJSON(w, http.StatusOK, review)
}
