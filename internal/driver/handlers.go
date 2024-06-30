package driver

import "net/http"

func DriverPostCord(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
