package wfmonitor

import (
    "github.com/eaciit/orm"
    "github.com/eaciit/dbox"
)

var db *orm.DataContext

func SetDb(conn dbox.IConnection){
    db = orm.New(conn)
}

func DB() *orm.DataContext{
    return db
}