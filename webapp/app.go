package webapp

import (
    _ "github.com/eaciit/dbox/dbc/mongo"
    "github.com/eaciit/dbox"
    "github.com/eaciit/orm"
    "github.com/eaciit/knot/knot.v1"
    "os"
)

func App() *knot.App{
    app := knot.NewApp("wfm")
    wd, _ := os.Getwd()
    wd += "/../"
    app.ViewsPath = wd+"views/"
    app.LayoutTemplate="_layout.html"
    app.Static("static",wd+"assets")
    app.Register(&Dashboard{})
    app.Register(&SDL{})
    app.Register(&Analytic{})
    app.Register(&Forecast{})
    app.DefaultOutputType = knot.OutputHtml
    return app
}

func conn() dbox.IConnection{
    conn, _ := dbox.NewConnection("mongo",&dbox.ConnectionInfo{"localhost:27123","ecwfmdemo","","",nil})
    conn.Connect()
    return conn
}

func DB() *orm.DataContext{
    return orm.New(conn())
}