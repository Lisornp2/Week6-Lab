package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"moviehub/mw"
	"moviehub/store"
)

func (h *Handler) GetVote(w http.ResponseWriter, r *http.Request) {
	memberID, ok := mw.MemberID(r.Context())
	if !ok {
		writeError(w, 401, "UNAUTHORIZED", "กรุณาเข้าสู่ระบบ")
		return
	}

	movieID, ok := pathID(r, "id")
	if !ok {
		writeError(w, 400, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}

	score, err := h.Store.GetVote(r.Context(), memberID, movieID)
	if errors.Is(err, store.ErrNotFound) {
		writeJSON(w, 200, map[string]any{"score": nil})
		return
	}
	if err != nil {
		internal(w, err)
		return
	}

	writeJSON(w, 200, map[string]any{"score": score})
}

func (h *Handler) PutVote(w http.ResponseWriter, r *http.Request) {
	memberID, ok := mw.MemberID(r.Context())
	if !ok {
		writeError(w, 401, "UNAUTHORIZED", "กรุณาเข้าสู่ระบบ")
		return
	}

	movieID, ok := pathID(r, "id")
	if !ok {
		writeError(w, 400, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}

	var input struct {
		Score int `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "BAD_JSON", "ข้อมูลไม่ถูกต้อง")
		return
	}
	if input.Score < 1 || input.Score > 10 {
		writeError(w, 400, "BAD_SCORE", "คะแนนต้องอยู่ระหว่าง 1 ถึง 10")
		return
	}

	err := h.Store.SaveVote(r.Context(), memberID, movieID, input.Score)
	if err != nil {
		internal(w, err)
		return
	}

	writeJSON(w, 200, map[string]any{
		"success": true,
		"score":   input.Score,
	})
}

func (h *Handler) DeleteVote(w http.ResponseWriter, r *http.Request) {
	memberID, ok := mw.MemberID(r.Context())
	if !ok {
		writeError(w, 401, "UNAUTHORIZED", "กรุณาเข้าสู่ระบบ")
		return
	}

	movieID, ok := pathID(r, "id")
	if !ok {
		writeError(w, 400, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}

	if err := h.Store.RemoveVote(r.Context(), memberID, movieID); err != nil {
		internal(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
