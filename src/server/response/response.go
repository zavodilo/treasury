package response

import (
	"encoding/json"
	"net/http"
)

type UpdateResponse struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
	Code   int    `json:"code"`
}

type StateResponse struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
}

type EmptySearchResponse struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
}

func JsonResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
