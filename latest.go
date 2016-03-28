package wfmonitor

import (
	"errors"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
)

type Latest struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            string `bson:"_id"`
	Reference1    string
	Reference2    string
	Timestamp     time.Time
}

func (s *Latest) TableName() string {
	return "latest"
}

func (s *Latest) RecordID() interface{} {
	return s.ID
}

func BuildLatestRead(){
    latest := new(Latest)
    latest.ID = "ingestion"
    latest.Timestamp = time.Now()
    DB().Save(latest)
}

func BuildLatest(dates []time.Time) {
	var maxdate time.Time
	for i, d := range dates {
		if i == 0 {
			maxdate = d
		} else if maxdate.Before(d) {
			maxdate = d
		}
	}

	latest := new(Latest)
	c, e := DB().Find(latest, toolkit.M{}.Set("where", dbox.Eq("_id", "scada")).
		Set("order", []string{"-timestamp"}).
		Set("limit", 1))
	if e != nil {
		return
	}
	if c.Count() == 0 {
		latest.ID = "scada"
		latest.Timestamp = maxdate
		DB().Save(latest)
	} else {
		e = c.Fetch(latest, 1, false)
		if e != nil {
			return
		}

		if latest.Timestamp.Before(maxdate) {
			latest.Timestamp = maxdate
			DB().Save(latest)
		}
	}
}

func BuildSummary(dates []time.Time) error {
	for _, date := range dates {
		var scadas []Scada
		c, e := DB().Find(new(Scada), toolkit.M{}.Set("where", dbox.Eq("transdate", date)))
        if e!=nil {
            continue
        }
		if c.Count()==0{
            continue
        }
        c.Fetch(&scadas, 0, false)

        toolkit.Println("Fetching ", len(scadas))
		summaries := map[string]*Summary{}
		for _, scada := range scadas {
			var summary *Summary
			sid := toolkit.Sprintf("%s-%s", scada.Turbine, toolkit.Date2String(scada.TransDate.UTC(), "YYYYMMdd"))
			if _, exist := summaries[sid]; !exist {
				summary = new(Summary)
				summary.ID = sid
				summary.Turbine = scada.Turbine
				summary.TransDate = scada.TransDate.UTC()
                summaries[sid]=summary
			} else {
				summary = summaries[sid]
			}
			oldcount := float32(summary.ChildCount)
			newcount := float32(summary.ChildCount + 1)
			summary.ChildCount++
			summary.Power += scada.Power
			summary.ConnectTime += scada.ConnectTime
			summary.FailureTime += scada.FailureTime
			summary.FullTime += scada.FullTime
			summary.Speed = (summary.Speed*oldcount + scada.Speed) / newcount
			summary.Temp = (summary.Temp*oldcount + scada.Temp) / newcount
			summary.Nacel = scada.Nacel
			summary.Direction = scada.Direction
            //summary.Created = time.Now()
		}

		for _, summary := range summaries {
			e := DB().Save(summary)
			if e != nil {
				return errors.New("Unable to build summary for " + summary.ID)
			}
		}
	}
    
    return nil
}
