package handler

import (
	"encoding/json"
	"github.com/sunwyx/currency-converter/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	conv *service.Converter
	log *slog.Logger
}

func NewHandler(conv *service.Converter, log *slog.Logger) *Handler {
	return &Handler{
		conv: conv,
		log: log,
	}
}

func (h *Handler) Convert(w http.ResponseWriter, r *http.Request){
	base := r.URL.Query().Get("from")
	if base == "" {
		http.Error(w, "Base currency is required", http.StatusBadRequest)
		return
	}
	toParams := r.URL.Query().Get("to")
	if toParams == "" {
		http.Error(w, "Target currencies are required", http.StatusBadRequest)
		return
	}
	amountStr := r.URL.Query().Get("amount")

	target := strings.Split(toParams, ",")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}
	for i := range target {
		target[i] = strings.ToUpper(strings.TrimSpace(target[i]))
	}
	result, err := h.conv.Convert(r.Context(), base, target, amount)
	if err != nil {
		h.log.Error("convert failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"from": base,
		"amount": amount,
		"to": result,
		})
	if err != nil {
		h.log.Error("Encode failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
