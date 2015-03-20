package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"log"
	"os"
)

const (
	MONGO_URL = "127.0.0.1:27017"
	MONGO_DB  = "devmgmt"
)

func e(err error, fatal bool) (ret bool) {
	if err != nil {
		log.Println(err.Error())
		if fatal {
			os.Exit(1)
		}
		return true
	} else {
		return false
	}
}

func main() {
	ms, err := mgo.Dial(MONGO_URL)
	e(err, true)
	defer ms.Close()
	ms.SetMode(mgo.Monotonic, true)
	mdb := ms.DB(MONGO_DB)

	m := martini.Classic()
	m.Map(mdb)

	m.Get("/", func(mdb *mgo.Database) string {
		names, err := mdb.CollectionNames()
		e(err, false)

		body, err := json.Marshal(names)
		e(err, false)
		return string(body)
	})

	m.Get("/:col", func(mdb *mgo.Database, params martini.Params) string {
		mcol := mdb.C(params["col"])
		items := []map[string]interface{}{}
		iter := mcol.Find(nil).Iter()
		for {
			item := map[string]interface{}{}
			if iter.Next(&item) {
				items = append(items, item)
			} else {
				break
			}
		}
		body, err := json.Marshal(items)
		e(err, false)
		return string(body)
	})

	m.Run()
}
