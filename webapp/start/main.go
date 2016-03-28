package main

import (
    "github.com/eaciit/knot/knot.v1"
    "eaciit/wfmonitor/webapp"
)

func main() {
    app := webapp.App()
    knot.StartApp(app, "localhost:9100")
}