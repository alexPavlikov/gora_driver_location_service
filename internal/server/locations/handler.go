package locations

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexPavlikov/gora_driver_location_service/internal/models"
)

// Handler получающий координаты водителя
func DriverPostCord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		slog.Info("start DriverPostCord")

		r.ParseForm()

		var driver models.Driver
		var err error

		driver.ID, err = strconv.Atoi(r.FormValue("driver_id"))
		if err != nil {
			slog.Error("not correct input driver_id" + err.Error())
			http.Error(w, "Invalid argument", http.StatusBadRequest)
		}

		driver.Cord.Longitude = r.FormValue("longitude")
		driver.Cord.Latitude = r.FormValue("latitude")

		// TODO: call func write to kafka

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
