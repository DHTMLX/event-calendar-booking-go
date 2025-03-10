package data

import (
	"errors"

	"gorm.io/gorm"
)

type doctorsEventDAO struct {
	db *gorm.DB
}

func newDoctorsEventDAO(db *gorm.DB) *doctorsEventDAO {
	return &doctorsEventDAO{db}
}

func (d *doctorsEventDAO) GetOne(id int) (DoctorEvent, error) {
	data := DoctorEvent{}
	err := d.db.Find(&data, id).Error
	return data, err
}

func (d *doctorsEventDAO) GetAll() ([]DoctorEvent, error) {
	sch := make([]DoctorEvent, 0)
	err := d.db.Find(&sch).Error
	return sch, err
}

func (d *doctorsEventDAO) Add(doctorID int, text string, from, to int, date, start, end int64, recurring bool, rrule string, recID int, original string, status string) (int, error) {
	if date == 0 {
		return 0, errors.New("date argument not defined")
	}

	event := DoctorEvent{
		DoctorID:         doctorID,
		From:             from,
		To:               to,
		Date:             date,
		Start:            start,
		End:              end,
		Recurring:        recurring,
		Rrule:            rrule,
		RecurringEventID: recID,
		OriginalStart:    original,
		Status:           status,
	}

	err := d.db.Save(&event).Error

	return event.ID, err
}

func (d *doctorsEventDAO) Update(id, doctorID int, text string, from, to int, date, start, end int64, recurring bool, rrule string, recID int, original string, status string) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	event, err := d.GetOne(id)
	if err != nil {
		return err
	}

	newEvent := DoctorEvent{
		DoctorID:         doctorID,
		From:             from,
		To:               to,
		Date:             date,
		Start:            start,
		End:              end,
		Recurring:        recurring,
		Rrule:            rrule,
		RecurringEventID: recID,
		OriginalStart:    original,
		Status:           status,
	}

	if event.DoctorID == newEvent.DoctorID {
		// add ID to update existing
		newEvent.ID = event.ID
		err = tx.Save(&newEvent).Error
	} else {
		// delete event at all for this old doctor and create event for doctorID
		err = tx.Delete(&DoctorEvent{}, id).Error
		if err != nil {
			return err
		}

		err = tx.Create(&newEvent).Error
	}

	return err
}

func (d *doctorsEventDAO) Delete(id int) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Delete(&DoctorEvent{}, "recurring_event_id = ?", id).Error
	if err != nil {
		return err
	}

	err = tx.Delete(&DoctorEvent{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}
