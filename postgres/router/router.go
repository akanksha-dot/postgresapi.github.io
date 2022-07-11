package router

import (
	"postgres/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/country", middleware.GetCountry).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newcountry", middleware.CreateCountry).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/updatecountry/{id}", middleware.UpdateCountry).Methods("Put", "OPTIONS")
	router.HandleFunc("/api/deletecountry/{id}", middleware.DeleteCountry).Methods("DELETE", "OPTIONS")

	return router

}
