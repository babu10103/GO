package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

// func ParseBody(r *http.Request, x interface{}) {
// 	if body, err := ioutil.ReadAll(r.Body); err != nil {
// 		if err := json.Unmarshal([]byte(body), x); err != nil {
// 			return
// 		}
// 	}
// }

// ParseBody reads the request body and unmarshals it into the provided value.
//
// r: the HTTP request.
// v: the value to unmarshal the body into.
// error: any error that occurred during the operation.
func ParseBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(body), v)
}
