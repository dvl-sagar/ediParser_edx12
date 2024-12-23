package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func ediToJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	file, _, err := r.FormFile("ediFile")
	if err != nil {
		res := Response{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		r, _ := json.Marshal(res)
		w.Write(r)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		res := Response{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		r, _ := json.Marshal(res)
		w.Write(r)
		return
	}

	result := EdiToJsonService(data)

	res := Response{
		Message: "Success",
		Data:    result,
	}
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(res)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
