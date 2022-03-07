package common

import (
	"encoding/json"
	"net/http"
)

func CreateResponse(rw http.ResponseWriter, code int, message string, data interface{}) (int, error) {
	rw.WriteHeader(code)

	response, _ := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	})

	return rw.Write(response)
}
