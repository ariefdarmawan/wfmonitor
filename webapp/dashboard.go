package webapp

import (
	"eaciit/wfmonitor"
	"time"

	"github.com/eaciit/crowd.dev"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type Dashboard struct {
}

func (d *Dashboard) Index(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputTemplate
	return struct{}{}
}

type timedata struct {
	TimeStamp                             time.Time
	Power, Uptime, Downtime, Speed, Count float32
}

func (d *Dashboard) GetStreamData(ctx *knot.WebContext)interface{}{
    ctx.Config.OutputType = knot.OutputJson
	result := toolkit.NewResult()
    
    db := DB()
    defer db.Close()
    
    scadas := []wfmonitor.Scada{}
    c, _ := db.Find(new(wfmonitor.Scada), toolkit.M{}.
        Set("order",[]string{"-created"}))
    if c.Count() > 0{
        c.Fetch(&scadas,0,false)
    }
    result.Data = scadas
    return result
}

func (d *Dashboard) GetDaily(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	result := toolkit.NewResult()

	model := &struct {
		LastUpdate                     time.Time
		Power, Speed, Uptime, Downtime float32
		TimeDatas                      []timedata
	}{}
	result.Data = model
	latest := new(wfmonitor.Latest)
	db := DB()
	defer db.Close()

	c, e := db.Find(latest, toolkit.M{}.
		Set("where", dbox.Eq("_id", "scada")))
	if e != nil {
		return result.SetErrorTxt("Can not get latest data. " + e.Error())
	}

	defer c.Close()
	if c.Count() == 0 {
		return result.SetErrorTxt("No latest date information received")
	}

	c.Fetch(latest, 1, false)
	model.LastUpdate = latest.Timestamp.UTC()

	var scadas []wfmonitor.Scada
	c, e = db.Find(new(wfmonitor.Scada), toolkit.M{}.
		Set("where", dbox.Eq("transdate", model.LastUpdate.UTC())))
	c.Fetch(&scadas, 0, false)

	timedatas := prepTimeData(model.LastUpdate)
	for _, scada := range scadas {
		ts := scada.Timestamp.UTC()
		timedatas[ts].Power += float32(scada.Power)
		timedatas[ts].Downtime += scada.FailureTime
		timedatas[ts].Uptime += scada.ConnectTime
		timedatas[ts].Speed = (timedatas[ts].Speed*timedatas[ts].Count + scada.Speed) / (timedatas[ts].Count + 1)
		timedatas[ts].Count += 1

		model.Power += float32(scada.Power)
		model.Uptime += scada.ConnectTime
		model.Downtime += scada.FailureTime
		model.Speed = (model.Speed*144 + timedatas[ts].Speed/144)
	}

	model.TimeDatas = func() []timedata {
		ret := []timedata{}
		for _, v := range timedatas {
			ret = append(ret, *v)
		}
		crowd.From(&ret).Sort(crowd.SortAscending, func(x interface{}) interface{} {
			return x.(timedata).TimeStamp
		})
		return ret
	}()

	return result
}

func prepTimeData(dt time.Time) map[time.Time]*timedata {
	ret := map[time.Time]*timedata{}
	di := dt
	for {
		td := new(timedata)
		td.TimeStamp = di
		ret[di] = td
		di = di.Add(10 * time.Minute)
		_, _, ddi := di.Date()
		_, _, ddt := dt.Date()
		if ddi != ddt {
			return ret
		}
	}
	return ret
}
