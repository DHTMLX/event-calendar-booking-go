# Event Calendar - Booking Demo Backend

## How to use

```
go build
./event-calendar-booking
```

# API

### GET /units

Returns all neccessary information to build booking dataset. Using `slots + usedslots` approach

#### Response example

```js
{
  "id": 1,
  "title": "Dr. Conrad Hubbard",
  "category": "Psychiatrist",
  "subtitle": "2 years of experience",
  "details": "Desert Springs Hospital (Schroeders Avenue 90, Fannett, Ethiopia)",
  "preview": "",
  "price": 120,
  "slots": [
    {
      "from": "09:00",
      "to": "14:00",
      "size": 45,
      "gap": 5,
      "days": [0, 1, 2, 3, 4, 5, 6] // reccuring events
    },
    {
      "from": "15:30",
      "to": "20:00",
      "size": 45,
      "gap": 5,
      "dates": [1695254400000] // Thu Sep 21 2023
    },
    ...
  ],
  "usedSlots": [
    1695367800000, // Fri Sep 22 2023 10:30:00 AM
    ...
  ]
}
```

### GET /calendars

Returns a list of doctors (without images)

#### Response example

```js
[
  {
    "id": 1,
    "name": "Dr. Conrad Hubbard",
    "subtitle": "2 years of experience",
    "details": "Desert Springs Hospital (Schroeders Avenue 90, Fannett, Ethiopia)",
    "category": "Psychiatrist",
    "price": 45,
    "gap": 20,
    "slot_size": 20,
    "active": true
  },
  ...
]
```

### GET /events

Returns a list of doctor's schedule (excluding expired dates).
You can show this data on Doctors view in Booking-Event-Calendar Demo

#### Response example

```js
[
  {
    "id": 1,
    "type": 1,
    "start_date": "2024-10-28T09:00:00Z",
    "end_date": "2024-10-28T17:00:00Z",
    "recurring": true, // recurring event
    "RRULE": "FREQ=WEEKLY",
    "STDATE": "2024-10-28T00:00:00Z",
    "DTEND": "9999-02-01T00:00:00.000Z"
  },
  {
    "id": 2,
    "type": 1,
    "start_date": "2024-10-29T18:00:00Z",
    "end_date": "2024-10-29T22:00:00Z"
  },
  {
    "id": 3,
    "type": 1,
    "start_date": "2024-10-30T18:00:00Z",
    "end_date": "2024-10-30T22:00:00Z"
  },
  // extension
  {
    "id": 4,
    "type": 1,
    "start_date": "2024-10-31T09:00:00Z",
    "end_date": "2024-10-31T17:00:00Z",
    "recurringEventId": 1,
    "originalStartTime": "2024-10-30T09:00:00Z"
  },
  // removed extension
  {
    "id": 5,
    "type": 1,
    "start_date": "2024-11-2T09:00:00Z",
    "end_date": "2024-11-2T17:00:00Z",
    "recurringEventId": 1,
    "originalStartTime": "2024-11-01T09:00:00Z",
    "status": "cancelled"
  }
  ...
]
```

### POST /events

Creates a new doctor's schedule with **concrete date** (Doctors view)

#### Body

```js
{
  "type": 1,
  "start_date": "2024-10-31T10:30:00Z"
  "end_date": "2024-10-31T14:30:00Z",
}
```

Creates a new doctor's schedule with **recurring days** (Doctors view)

#### Body

```js
{
  "type": 1,
  "start_date": "2024-10-28T10:30:00Z",
  "end_date": "2024-10-28T14:30:00Z",
  "recurring": true,
  "RRULE": "FREQ=WEEKLY",
  "STDATE": "2024-10-28T00:00:00Z",
  "DTEND": "9999-02-01T00:00:00.000Z"
}
```

### Response example

Returns an ID of created schedule (Doctors view)

```js
{
  "id": 10
}
```

### PUT /events/{id}

Updates doctor's schedule

#### Body

```js
{
  "type": 1,
  "start_date": "2024-10-31T12:20:00Z",
  "end_date": "2024-10-31T16:55:00Z"
}
```

Updates **recurring** doctor's schedule

#### Body

```js
{
  "type": 1,
  "start_date": "2024-10-31T10:30:00Z",
  "end_date": "2024-10-31T14:30:00Z",
  "recurring": true,
  "RRULE": "FREQ=WEEKLY",
  "STDATE": "2024-10-31T00:00:00Z",
  "DTEND": "9999-02-01T00:00:00.000Z"
}
```

### Response example

Returns an ID of updated schedule (Doctors view)

```js
{
  "id": 10
}
```

#### URL Params:

- id [required] - ID of the schedule to be updated

### DELETE /events/{id}

Deletes doctor's schedule (Doctors view)


#### URL Params:

- id [required] - ID of the schedule to be deleted

### GET /reservations

Returns all occupied slots (Clients view)

#### Response example

```js
[
    {
        "id": 1,
        "type": 2,
        "date": 1730289600000,
        "client_name": "Alan",
        "client_email": "alan@gmail.com",
        "client_details": ""
    },
    {
        "id": 2,
        "type": 3,
        "date": 1730356200000,
        "client_name": "Viron",
        "client_email": "viron@hr.com",
        "client_details": ""
    }
    ...
]
```

### POST /reservations

Creates reservation (Booking view)

#### Body

```js
{
  "doctor": 2,
  "date": 1730289600000,
  "form": {
    "name": "Alan",
    "email": "alan@gmail.com",
    "details": ""
  }
}
```

# Features

### Booking schedules

If the schedule encompasses midnight and there is enough time for a time slot after it, then the schedule is divided into two parts

### Used slots

Booking processes only matches exact used slots for the doctor. If the booked slot does not match any of the slots, the two closest relevant slots will be booked instead

# Config

```yaml
db:
  path: db.sqlite    # path to the database
  resetonstart: true # reset data on server restart
server:
  url: "http://localhost:3000"
  port: ":3000"
  cors:
    - "*"
  resetFrequence: 120 # every 2 hours restart data (value in minutes)
```
