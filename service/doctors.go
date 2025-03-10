package service

import (
	"event-calendar-booking/data"
	"fmt"
)

type doctorsService struct {
	dao *data.DAO
}

type DoctorUpdate struct {
	Active bool `json:"active"`
}

func (s *doctorsService) GetDoctorsList() ([]data.Doctor, error) {
	doctors, err := s.dao.Doctors.GetAll(false)
	return doctors, err
}

func (s *doctorsService) Update(id int, upd DoctorUpdate) error {
	doctor, err := s.dao.Doctors.GetOne(id)
	if err != nil {
		return err
	}

	if doctor.ID == 0 {
		return fmt.Errorf("doctor with id %d not found", doctor.ID)
	}

	doctor.Active = upd.Active

	err = s.dao.Doctors.Update(doctor)
	return err
}
