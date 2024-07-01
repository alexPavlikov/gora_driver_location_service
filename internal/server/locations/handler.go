package locations

import "net/http"

// Handler получающий координаты водителя
func DriverPostCord(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
