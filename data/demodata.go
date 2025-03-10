package data

import (
	"fmt"
	"math/rand"
	"strings"

	"gorm.io/gorm"
)

func dataDown(tx *gorm.DB) {
	must(tx.Exec("DELETE FROM `doctors`").Error)
	must(tx.Exec("DELETE FROM `reviews`").Error)
	must(tx.Exec("DELETE FROM `doctor_events`").Error)
	must(tx.Exec("DELETE FROM `occupied_slots`").Error)
}

var (
	firstNames = []string{
		"Emma",
		"Olivia",
		"James",
		"Mia",
		"Amelia",
		"Alexander",
		"Harper",
		"William",
		"Abigail",
		"Lily",
	}

	lastNames = []string{
		"Johnson",
		"Smith",
		"Brown",
		"Wilson",
		"Jackson",
		"King",
		"Scott",
		"Green",
		"Adams",
		"Baker",
	}
)

const (
	nameFormat  = "%s %s"
	emailFormat = "%s.%s@event-calendar.booking"

	endUnix     = int64(253402300799000) // max UnixMilli
	allDayMilli = int64(86400000)        // 24 * 60 * 60 * 1000
)

func dataUp(tx *gorm.DB) {
	today := DateNow() // date only
	todayMilli := today.UnixMilli()
	tWeekDay := int(today.Weekday())

	recurringEvent := func(from, to int, rrule string) DoctorEvent {
		if to < from {
			to += 24 * 60 // one day
		}

		return DoctorEvent{
			From:      from,
			To:        to,
			Date:      todayMilli,
			Start:     todayMilli,
			End:       endUnix,
			Recurring: true,
			Rrule:     "INTERVAL=1;FREQ=WEEKLY;BYDAY=" + strings.ToUpper(rrule),
		}
	}

	routineEvent := func(from, to int, date int64, days int) []DoctorEvent {
		routines := make([]DoctorEvent, days)
		for i := range routines {
			routines[i] = DoctorEvent{
				From: from,
				To:   to,
				Date: date + int64(i)*allDayMilli,
			}
		}

		return routines
	}

	nextWeekDay := func(day int, weeks ...int) int64 {
		var week int
		if len(weeks) > 0 {
			week = weeks[0]
		}

		next := (7 + day - tWeekDay) % 7
		return today.AddDate(0, 0, next+7*week).UnixMilli()
	}

	newSlot := func(date int64, time int) OccupiedSlot {
		first := firstNames[rand.Intn(len(firstNames))]
		last := lastNames[rand.Intn(len(lastNames))]

		return OccupiedSlot{
			Date:        date + int64(time)*60*1000,
			ClientName:  fmt.Sprintf(nameFormat, first, last),
			ClientEmail: fmt.Sprintf(emailFormat, strings.ToLower(first), strings.ToLower(last)),
		}
	}

	doctors := []Doctor{
		{
			Name:     "Dr. Conrad Hubbard",
			Category: "Psychiatrist",
			Subtitle: "2 years of experience",
			Details:  "Desert Springs Hospital (Schroeders Avenue 90, Fannett, Ethiopia)",
			SlotSize: 20,
			Price:    "$45",
			ImageURL: "https://snippet.dhtmlx.com/codebase/data/booking/01/img/11.jpg",
			Gap:      20,
			Review: Review{
				Count: 1245,
				Stars: 4,
			},
			DoctorEvents: append(
				[]DoctorEvent{
					// every week day 9:00-17:00 (except sun, sat - holidays)
					recurringEvent(9*60, 17*60, "MO,TU,WE,TH,FR"),
				},
				// next tue, wed, thu 2:00-6:00
				routineEvent(2*60, 6*60, nextWeekDay(2), 3)...,
			),
			OccupiedSlots: []OccupiedSlot{
				newSlot(nextWeekDay(1), 9*60+40),    // next mon 9:40
				newSlot(nextWeekDay(2), 11*60),      // next tue 11:00
				newSlot(nextWeekDay(2), 15*60),      // next tue 15:00
				newSlot(nextWeekDay(3, 1), 11*60),   // after next wed 11:00
				newSlot(nextWeekDay(4), 3*60+20),    // next thu 3:20
				newSlot(nextWeekDay(4), 16*60+20),   // next thu 16:20
				newSlot(nextWeekDay(4, 1), 5*60+20), // after next thu 5:20
				newSlot(nextWeekDay(5), 13*60+20),   // next fri 13:20
			},
			Active: true,
		},
		{
			Name:     "Dr. Debra Weeks",
			Category: "Allergist",
			Subtitle: "7 years of experience",
			Details:  "Silverstone Medical Center (Vanderbilt Avenue 13, Chestnut, New Zealand)",
			SlotSize: 45,
			Price:    "$120",
			ImageURL: "https://snippet.dhtmlx.com/codebase/data/booking/01/img/03.jpg",
			Gap:      5,
			Review: Review{
				Count: 6545,
				Stars: 4,
			},
			DoctorEvents: append(
				[]DoctorEvent{
					// mon, wed 7:00-15:00
					recurringEvent(7*60, 15*60, "MO,WE"),
				},
				// next wed 18:00-22:00
				routineEvent(18*60, 22*60, nextWeekDay(3), 1)...,
			),
			OccupiedSlots: []OccupiedSlot{
				newSlot(nextWeekDay(1), 7*60+50),  // next mon 7:50
				newSlot(nextWeekDay(2), 13*60+40), // next tue 13:40
				newSlot(nextWeekDay(3), 11*60+10), // next wed 11:10
				newSlot(nextWeekDay(3), 19*60+40), // next wed 19:40
				newSlot(nextWeekDay(4), 14*60+30), // next thu 14:30
			},
			Active: true,
		},
		{
			Name:     "Dr. Barnett Mueller",
			Category: "Ophthalmologist",
			Subtitle: "6 years of experience",
			Details:  "Navy Street 1, Kiskimere, United States",
			SlotSize: 25,
			Price:    "$35",
			ImageURL: "https://snippet.dhtmlx.com/codebase/data/booking/01/img/02.jpg",
			Gap:      0,
			Review: Review{
				Count: 184,
				Stars: 3,
			},
			DoctorEvents: append(
				[]DoctorEvent{
					// mon, wed, fri 9:00-17:00
					recurringEvent(9*60, 17*60, "MO,WE,FR"),
				},
				// next tue, wed, thu 04:00-08:00
				routineEvent(4*60, 8*60, nextWeekDay(2), 3)...,
			),
			OccupiedSlots: []OccupiedSlot{
				newSlot(nextWeekDay(1), 13*60+10),    // next mon 13:10
				newSlot(nextWeekDay(1, 1), 12*60+45), // after next mon 12:45
				newSlot(nextWeekDay(3), 9*60+25),     // next wed 9:25
				newSlot(nextWeekDay(5), 11*60+55),    // next fri 11:55
				newSlot(nextWeekDay(5, 1), 11*60+30), // after next fri 11:30
				newSlot(nextWeekDay(6), 16*60+10),    // next sat 16:10
			},
			Active: true,
		},
		{
			Name:     "Dr. Myrtle Wise",
			Category: "Ophthalmologist",
			Subtitle: "4 years of experience",
			Details:  "Prescott Place 5, Freeburn, Bulgaria",
			SlotSize: 25,
			Price:    "$40",
			ImageURL: "https://snippet.dhtmlx.com/codebase/data/booking/01/img/01.jpg",
			Gap:      5,
			Review: Review{
				Count: 829,
				Stars: 5,
			},
			DoctorEvents: append(
				[]DoctorEvent{
					// tue, thu 12:00-20:00
					recurringEvent(12*60, 20*60, "TU,TH"),
				},
				// next fri, sat 4:00-8:00
				routineEvent(4*60, 8*60, nextWeekDay(5), 2)...,
			),
			OccupiedSlots: []OccupiedSlot{
				newSlot(nextWeekDay(2), 7*60),     // next tue 7:00
				newSlot(nextWeekDay(2), 10*60),    // next tue 10:00
				newSlot(nextWeekDay(4), 9*60+30),  // next thu 9:30
				newSlot(nextWeekDay(5), 7*60+30),  // next fri 7:30
				newSlot(nextWeekDay(6), 11*60+30), // next sat 11:30
				newSlot(nextWeekDay(6), 5*60),     // next sat 5:00
				newSlot(nextWeekDay(0), 12*60),    // next sun 12:00
			},
			Active: true,
		},
		{
			Name:     "Dr. Browning Peck",
			Category: "Dentist",
			Subtitle: "11 years of experience",
			SlotSize: 60,
			Details:  "Seacoast Terrace 174, Belvoir, Mauritania",
			Price:    "$175",
			ImageURL: "https://snippet.dhtmlx.com/codebase/data/booking/01/img/12.jpg",
			Gap:      10,
			Review: Review{
				Count: 391,
				Stars: 5,
			},
			DoctorEvents: []DoctorEvent{
				//  thu, fri, sat, sun  9:00-17:00
				recurringEvent(9*60, 17*60, "TH,FR,SA,SU"),
			},
			OccupiedSlots: []OccupiedSlot{
				newSlot(nextWeekDay(4), 11*60+20), // next thu 11:20
				newSlot(nextWeekDay(5), 14*60+50), // next fri 14:50
				newSlot(nextWeekDay(6), 9*60),     // next sat 9:00
				newSlot(nextWeekDay(6), 13*60+20), // next sat 13:20
				newSlot(nextWeekDay(0), 14*60+50), // next sun 14:50
			},
			Active: true,
		},
	}

	err := tx.Create(doctors).Error
	if err != nil {
		panic(err)
	}
}
