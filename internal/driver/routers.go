package driver

import "net/http"

func HandlerRequest() {
	http.HandleFunc("/driver_post_cord", DriverPostCord)
}
