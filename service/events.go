package service

import (
	"event-calendar-booking/common"
	"event-calendar-booking/data"
	"fmt"
	"time"
)

type eventsService struct {
	dao *data.DAO
}

type Event struct {
	DoctorID int    `json:"type"`
	Text     string `json:"text"`

	StartDate *common.JDate `json:"start_date"`
	EndDate   *common.JDate `json:"end_date"`

	Recurring bool          `json:"recurring"`
	Rrule     string        `json:"RRULE"`
	STDate    *common.JDate `json:"STDATE"`
	DTEnd     *common.JDate `json:"DTEND"`

	RecurringEventID int    `json:"recurringEventId"`
	OriginalStart    string `json:"originalStartTime"`
	Status           string `json:"status"`
}

type EventStr struct {
	ID       int    `json:"id"`
	DoctorID int    `json:"type"`
	Text     string `json:"text"`

	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`

	Recurring bool   `json:"recurring,omitempty"`
	Rrule     string `json:"RRULE,omitempty"`
	STDate    string `json:"STDEND,omitempty"`
	DTEnd     string `json:"DTEND,omitempty"`

	RecurringEventID int    `json:"recurringEventId,omitempty"`
	OriginalStart    string `json:"originalStartTime,omitempty"`
	Status           string `json:"status,omitempty"`
}

const (
	strFormat = "2006-01-02 15:04:05"
	endDate   = "9999-02-01 00:00:00"
	endUnix   = int64(253402300799000)
)

// returns records for the event-calendar Doctors View
func (s *eventsService) GetAll() ([]EventStr, error) {
	events, err := s.dao.DoctorsEvent.GetAll()
	out := make([]EventStr, 0)

	for _, event := range events {
		fh := event.From / 60
		fm := event.From % 60
		th := event.To / 60
		tm := event.To % 60

		y, m, d := time.UnixMilli(event.Date).UTC().Date()

		var stdate, dtend string
		if event.Recurring {
			stdate = time.UnixMilli(event.Start).UTC().Format(strFormat)
			dtend = endDate
		}

		r := EventStr{
			ID:               event.ID,
			DoctorID:         event.DoctorID,
			StartDate:        time.Date(y, m, d, fh, fm, 0, 0, time.UTC).Format(strFormat),
			EndDate:          time.Date(y, m, d, th, tm, 0, 0, time.UTC).Format(strFormat),
			Recurring:        event.Recurring,
			Rrule:            event.Rrule,
			STDate:           stdate,
			DTEnd:            dtend,
			RecurringEventID: event.RecurringEventID,
			OriginalStart:    event.OriginalStart,
			Status:           event.Status,
		}

		out = append(out, r)
	}

	return out, err
}

// adds doctor's event
func (s *eventsService) Add(data Event) (int, error) {
	if err := data.validate(); err != nil {
		return 0, err
	}

	date := data.StartDate.Truncate(24 * time.Hour).UnixMilli()

	from := data.StartDate.Hour()*60 + data.StartDate.Minute()
	to := from + data.duration()

	var start, end int64
	if data.Recurring {
		start = data.STDate.UnixMilli()
		end = endUnix // data.DTEnd.UnixMilli()
	}

	id, err := s.dao.DoctorsEvent.Add(
		data.DoctorID,
		data.Text,
		from,
		to,
		date,
		start,
		end,
		data.Recurring,
		data.Rrule,
		data.RecurringEventID,
		data.OriginalStart,
		data.Status,
	)
	return id, err
}

// updates doctor's event
func (s *eventsService) Update(eventID int, data Event) error {
	event, err := s.dao.DoctorsEvent.GetOne(eventID)
	if err != nil {
		return err
	}

	if event.ID == 0 {
		return fmt.Errorf("event with id %d not found", eventID)
	}

	if err := data.validate(); err != nil {
		return err
	}

	date := data.StartDate.Truncate(24 * time.Hour).UnixMilli()

	from := data.StartDate.Hour()*60 + data.StartDate.Minute()
	to := from + data.duration()

	var start, end int64
	if data.Recurring {
		start = data.STDate.UnixMilli()
		end = endUnix // data.DTEnd.UnixMilli()
	}

	err = s.dao.DoctorsEvent.Update(
		eventID,
		data.DoctorID,
		data.Text,
		from,
		to,
		date,
		start,
		end,
		data.Recurring,
		data.Rrule,
		data.RecurringEventID,
		data.OriginalStart,
		data.Status,
	)
	return err
}

// delets doctor's event
func (s *eventsService) Delete(id int) error {
	return s.dao.DoctorsEvent.Delete(id)
}

func (e Event) validate() error {
	if e.StartDate.UnixMilli() < data.Now().UnixMilli() {
		return fmt.Errorf("cannot set work time in the past")
	}
	if e.StartDate.UnixMilli() >= e.EndDate.UnixMilli() {
		return fmt.Errorf("invalid time interval")
	}
	return nil
}

// in minutes
func (e Event) duration() int {
	diff := e.EndDate.Sub(e.StartDate.Time)
	return int(diff.Minutes())
}
