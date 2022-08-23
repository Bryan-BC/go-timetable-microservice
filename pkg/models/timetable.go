package models

import (
	"encoding/json"

	"database/sql/driver"

	"github.com/Bryan-BC/go-timetable-microservice/pkg/pb"
)

type CourseArr []*pb.TimetableSchedule

type Timetable struct {
	Id      int64     `json:"id" gorm:"primary_key"`
	Courses CourseArr `json:"courses"`
}

func (ca *CourseArr) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), ca)
}

func (ca CourseArr) Value() (driver.Value, error) {
	return json.Marshal(ca)
}
