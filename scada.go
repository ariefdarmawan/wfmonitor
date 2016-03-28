package wfmonitor

import (
	"time"

	"github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
)

type Scada struct {
	orm.ModelBase                                                     `bson:"-" json:"-"`
	ID                                                                string `bson:"_id"`
	Timestamp                                                         time.Time
	TransDate                                                         time.Time
	Turbine                                                           string
	Speed, Direction, Nacel, Temp, FailureTime, ConnectTime, FullTime float32
	Power                                                             float64
    Created, LastUpdated                                              time.Time
}

func (s *Scada) TableName() string {
	return "scadas"
}

func (s *Scada) RecordID() interface{} {
	return s.ID
}

func (s *Scada) PreSave() error {
	s.ID = toolkit.Sprintf("%s-%s", s.Turbine, toolkit.Date2String(s.Timestamp, "YYYYMMddHHmmss"))
	return nil
}

type Summary struct {
	orm.ModelBase                                                     `bson:"-" json:"-"`
	ID                                                                string `bson:"_id"`
	TransDate                                                         time.Time
	Turbine                                                           string
	Speed, Direction, Nacel, Temp, FailureTime, ConnectTime, FullTime float32
	Power                                                             float64
	ChildCount                                                        int
	//Created, LastUpdated                                              time.Time
}

func (s *Summary) TableName() string {
	return "summaries"
}

func (s *Summary) RecordID() interface{} {
	return s.ID
}

func (s *Summary) PreSave() error {
	s.ID = toolkit.Sprintf("%s-%s", s.Turbine, toolkit.Date2String(s.TransDate, "YYYYMMdd"))
	return nil
}
