package service

import "event-calendar-booking/data"

type ServiceAll struct {
	Doctors      *doctorsService
	Events       *eventsService
	Reservations *reservationsService
	Units        *unitsService
}

func NewService(dao *data.DAO) *ServiceAll {
	return &ServiceAll{
		Doctors:      &doctorsService{dao},
		Reservations: &reservationsService{dao},
		Events:       &eventsService{dao},
		Units:        &unitsService{dao},
	}
}
