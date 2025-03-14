package api

import (
	"event-calendar-booking/service"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

var Debug = true

type API struct {
	sAll   *service.ServiceAll
	format *render.Render
}

func NewAPI(service *service.ServiceAll) *API {
	format := render.New()
	return &API{service, format}
}

func (api *API) InitRoutes(r chi.Router) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		api.response(w, "Server launched successfully!", nil)
	})

	r.Get("/units", func(w http.ResponseWriter, r *http.Request) {
		units, err := api.sAll.Units.GetAll()
		api.response(w, units, err)
	})

	r.Get("/calendars", func(w http.ResponseWriter, r *http.Request) {
		doctors, err := api.sAll.Doctors.GetDoctorsList()
		api.response(w, doctors, err)
	})

	r.Put("/calendars/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := numberParam(r, "id")
		upd := service.DoctorUpdate{}
		err := parseForm(w, r, &upd)
		if err != nil {
			api.errResponse(w, err.Error())
			return
		}
		err = api.sAll.Doctors.Update(id, upd)

		api.response(w, &response{}, err)
	})

	r.Get("/events", func(w http.ResponseWriter, r *http.Request) {
		data, err := api.sAll.Events.GetAll()
		api.response(w, data, err)
	})

	r.Post("/events", func(w http.ResponseWriter, r *http.Request) {
		event := service.Event{}
		err := parseForm(w, r, &event)
		if err != nil {
			api.errResponse(w, err.Error())
			return
		}
		id, err := api.sAll.Events.Add(event)

		api.response(w, &response{id}, err)
	})

	r.Put("/events/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := numberParam(r, "id")
		event := service.Event{}
		err := parseForm(w, r, &event)
		if err != nil {
			api.errResponse(w, err.Error())
			return
		}
		err = api.sAll.Events.Update(id, event)

		api.response(w, &response{}, err)
	})

	r.Delete("/events/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := numberParam(r, "id")
		err := api.sAll.Events.Delete(id)
		api.response(w, &response{}, err)
	})

	r.Get("/reservations", func(w http.ResponseWriter, r *http.Request) {
		reservations, err := api.sAll.Reservations.GetAll()
		api.response(w, reservations, err)
	})

	r.Post("/reservations", func(w http.ResponseWriter, r *http.Request) {
		reservation := service.Reservation{}
		err := parseForm(w, r, &reservation)
		if err != nil {
			api.errResponse(w, err.Error())
			return
		}
		id, err := api.sAll.Reservations.Add(reservation)

		api.response(w, &response{id}, err)
	})
}

func (api *API) response(w http.ResponseWriter, data any, err error) {
	if err != nil {
		api.errResponse(w, err.Error())
	} else {
		api.format.JSON(w, 200, data)
	}
}

func (api *API) errResponse(w http.ResponseWriter, msg string) {
	if Debug {
		fmt.Println(msg)
	}
	api.format.Text(w, 500, msg)
}
