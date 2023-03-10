package helper

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Json encode escaping HTML and setting the
// Content-Type as application/json.
func Write(w http.ResponseWriter, statusCode int, data interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
