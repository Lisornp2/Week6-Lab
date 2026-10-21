package handler

import (
	"net/http"

	"moviehub/mw"
	"moviehub/store"
)

func (h *Handler) ListWishlist(w http.ResponseWriter, r *http.Request) {
	memberID, ok := mw.MemberID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "กรุณาเข้าสู่ระบบ")
		return
	}

	items, err := h.Store.ListWishlist(r.Context(), memberID)
	if err != nil {
		internal(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) AddWishlist(w http.ResponseWriter, r *http.Request) {
	memberID, ok := mw.MemberID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "กรุณาเข้าสู่ระบบ")
		return
	}

	movieID, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}

	err := h.Store.AddWishlist(r.Context(), memberID, movieID)
	if err != nil {
		if err == store.ErrNotFound {
			writeError(w, http.StatusNotFound, "MOVIE_NOT_FOUND", "ไม่พบหนังเรื่องนี้")
			return
		}
		internal(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *Handler) RemoveWishlist(w http.ResponseWriter, r *http.Request) {
	memberID, ok := mw.MemberID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "กรุณาเข้าสู่ระบบ")
		return
	}

	movieID, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}

	if err := h.Store.RemoveWishlist(r.Context(), memberID, movieID); err != nil {
		internal(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
