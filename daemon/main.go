package main

import (
	"bufio"
	"eaciit/wfmonitor"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
	"github.com/eaciit/uklam"
)

var (
	log    *toolkit.LogEngine
	cdone  chan bool
	status string
	path   = "/users/ariefdarmawan/Dropbox/pvt/Temp/bhesada"
)

func main() {
	toolkit.Println("WF Monitor Deamon v0.5 (c) EACIIT")
	toolkit.Println("Run http://localhost:8888/stop to stop the daemon")
	toolkit.Println("")
	log, _ = toolkit.NewLog(true, false, "", "", "")
	defer func() {
		log.Info("Closing daemon")
	}()

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

	winbox := prepareInbox()
	winbox.Start()
	defer winbox.Stop()

	wrun := prepareRunning()
	wrun.Start()
	defer wrun.Stop()

	cdone = make(chan bool)
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		cdone <- true
		w.Write([]byte("Daemon will be stopped"))
	})

	go func() {
		e = http.ListenAndServe(":8888", nil)
		if e != nil {
			log.Error("Can not start daemon REST server. " + e.Error())
			cdone <- true
		}
	}()

	for {
		select {
		case <-cdone:
			status = "Closing"
			return

		default:
			//-- do nothing
		}
	}
}

func prepareInbox() *uklam.FSWalker {
	w := uklam.NewFS(filepath.Join(path, "inbox"))
	w.EachFn = func(dw uklam.IDataWalker, in toolkit.M, info os.FileInfo, r *toolkit.Result) {
		//go func(info os.FileInfo) {
		sourcename := filepath.Join(path, "inbox", info.Name())
		dstname := filepath.Join(path, "running", info.Name())
		log.Info(toolkit.Sprintf("Processing " + sourcename))
		e := uklam.FSCopy(sourcename, dstname, true)
		if e != nil {
			log.Error("Processing " + sourcename + " NOK " + e.Error())
		} else {
			log.Info("Processing " + sourcename + " OK ")
		}
		//time.Sleep(100 * time.Millisecond)
		//}(info)
	}
	return w
}

func prepareRunning() *uklam.FSWalker {
	w2 := uklam.NewFS(filepath.Join(path, "running"))
	w2.EachFn = func(dw uklam.IDataWalker, in toolkit.M, info os.FileInfo, r *toolkit.Result) {
		//go func(info os.FileInfo) {
		//time.Sleep(10 * time.Millisecond)
		sourcename := filepath.Join(path, "running", info.Name())
		dstnameOK := filepath.Join(path, "success", info.Name())
		dstnameNOK := filepath.Join(path, "fail", info.Name())
		log.Info(toolkit.Sprintf("Reading %s", sourcename))
		e := streamDo(sourcename)
		//time.Sleep(100*time.Millisecond)
		if e == nil {
			uklam.FSCopy(sourcename, dstnameOK, true)
			log.Info(toolkit.Sprintf("%s OK", sourcename))
		} else {
			uklam.FSCopy(sourcename, dstnameNOK, true)
			log.Error(toolkit.Sprintf("%s NOK: %s", sourcename, e.Error()))
		}
		//}(info)
	}
	return w2
}

func streamDo(src string) error {
	f, _ := os.Open(src)
	defer f.Close()

	b := bufio.NewScanner(f)
	b.Split(bufio.ScanLines)
	dates := map[time.Time]int{}

	i := 0
	for b.Scan() {
		if i > 0 {
			str := strings.Split(b.Text(), ",")
			scada := new(wfmonitor.Scada)
			scada.Timestamp = toolkit.String2Date(str[0], "YYYYMMddHHmmss")
			scada.TransDate = toolkit.String2Date(toolkit.Date2String(scada.Timestamp, "YYYYMMdd"), "YYYYMMdd")
			scada.Turbine = str[1][len(str[1])-6:]
			scada.Speed = toolkit.ToFloat32(str[2], 2, toolkit.RoundingAuto)
			scada.Direction = toolkit.ToFloat32(str[3], 2, toolkit.RoundingAuto)
			if scada.Direction < 0 {
				scada.Direction = 360 + scada.Direction
			} else if scada.Direction >= 360 {
				scada.Direction = scada.Direction - 360
			}
			scada.Nacel = toolkit.ToFloat32(str[4], 2, toolkit.RoundingAuto)
			scada.FailureTime = toolkit.ToFloat32(str[6], 2, toolkit.RoundingAuto) / float32(60)
			scada.ConnectTime = toolkit.ToFloat32(str[7], 2, toolkit.RoundingAuto) / float32(60)
			scada.FullTime = scada.FailureTime + scada.ConnectTime
			scada.Power = toolkit.ToFloat64(str[5], 2, toolkit.RoundingAuto) * float64(scada.FullTime) / float64(60) / float64(1000)
			scada.Temp = toolkit.ToFloat32(str[8], 2, toolkit.RoundingAuto)
            scada.Created = time.Now()
			esave := wfmonitor.DB().Save(scada)
			if esave == nil {
				if _, exist := dates[scada.TransDate]; !exist {
					dates[scada.TransDate] = 1
				}
			}
			//toolkit.Println(toolkit.JsonString(scada))
		}
		i++
	}

	//go func() {
		if len(dates) > 0 {
			dateslice := func() []time.Time {
				var ds []time.Time
				for d, _ := range dates {
					ds = append(ds, d)
				}
				return ds
			}()

			log.Info("Building logs")
			wfmonitor.BuildLatest(dateslice)

			log.Info("Building Summary Model")
			wfmonitor.BuildSummary(dateslice)
		}
	//}()
	return nil
}
