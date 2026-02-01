package handler

import (
	"encoding/json"
	"github.com/sunwyx/currency-converter/internal/service"
	_ "log"
	"net/http"
	"strconv"
)

type Handler struct {
	conv *service.Converter
}

func NewHandler(conv *service.Converter) *Handler {
	return &Handler{
		conv: conv,
	}
}

func (h *Handler) Convert(w http.ResponseWriter, r *http.Request){
	base := r.URL.Query().Get("from")
	target := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	result, err := h.conv.Convert(r.Context(), base, target, amount)
	if err != nil {
		h.conv.Log.Error("convert failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"from": base,
		"to": target,
		"amount": amount,
		"result": result,
		})
}
