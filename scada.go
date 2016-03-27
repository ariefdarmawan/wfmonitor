package wfmonitor

import (
    "time"
    "github.com/eaciit/toolkit"
    "github.com/eaciit/orm"
)

type Scada struct {
    orm.ModelBase `bson:"-" json:"-"`
    ID string `bson:"_id"`
	Timestamp                                                         time.Time
	Turbine                                                           string
	Speed, Direction, Nacel, Temp, FailureTime, ConnectTime, FullTime float32
	Power                                                             float64
}

func (s *Scada) TableName() string{
    return "scadas"
}

func (s *Scada) RecordID() interface{}{
    return s.ID
}

func (s *Scada) PreSave()error{
    s.ID = toolkit.Sprintf("%s-%s", s.Turbine, toolkit.Date2String(s.Timestamp, "YYYYMMddHHmmss"))
    return nil
}
