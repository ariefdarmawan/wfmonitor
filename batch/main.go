package main

import (
    "eaciit/wfmonitor"
   "github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
    "os"
)

var log *toolkit.LogEngine

func main() {
    log, _ = toolkit.NewLog(true,false,"","","")
    defer log.Close()
    
    conn, _ := dbox.NewConnection("mongo", &dbox.ConnectionInfo{"localhost:27123", "ecwfmdemo", "", "", nil})
	e := conn.Connect()
	if e != nil {
		log.AddLog("Unable to connect to database: "+e.Error(), "ERROR")
		os.Exit(100)
	} else {
		log.AddLog("Connected to database", "INFO")
	}

	wfmonitor.SetDb(conn)
	defer wfmonitor.DB().Close()

    scadas := []wfmonitor.Scada{}
    c, _ := wfmonitor.DB().Find(new(wfmonitor.Scada),nil)
    e = c.Fetch(&scadas,0,false)
    if e!=nil{
        log.Error("Fetch Scadas " + e.Error())
        os.Exit(100)
    }
    
    wfmonitor.DB().Connection.NewQuery().From("summaries").Delete().Exec(toolkit.M{}.Set("multiexec",true))
    summaries := map[string]*wfmonitor.Summary{}
    for _, scada := range scadas{
        var summary *wfmonitor.Summary
        sid := toolkit.Sprintf("%s-%s", scada.Turbine, toolkit.Date2String(scada.Timestamp, "YYYYMMdd"))
        if _, exist := summaries[sid]; !exist{
            summary = new(wfmonitor.Summary)
            summary.ID = sid
            summary.Turbine = scada.Turbine
            summary.TransDate = scada.TransDate
            summaries[sid]=summary
        } else {
            summary = summaries[sid]
        }
        oldcount := float32(summary.ChildCount)
        newcount := float32(summary.ChildCount+1)
        summary.ChildCount++
        summary.Power += scada.Power
        summary.ConnectTime += scada.ConnectTime
        summary.FailureTime += scada.FailureTime
        summary.FullTime += scada.FullTime
        summary.Speed = (summary.Speed * oldcount + scada.Speed) / newcount
	    summary.Temp = (summary.Temp * oldcount + scada.Temp) / newcount
        summary.Nacel = scada.Nacel
        summary.Direction = scada.Direction
	    e = wfmonitor.DB().Save(summary)
        if e!=nil {
            log.Error("Unable to save " + toolkit.JsonString(summary) + " " + e.Error())
        } else {
            log.Info("Saving " + summary.ID)
        }
    } 
}