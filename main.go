package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"log"
	"os"
	"net/http"
	"strconv"
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

	m.Get("/:col", func(mdb *mgo.Database, params martini.Params, req *http.Request) string {
		req.ParseForm()
		mcol := mdb.C(params["col"])
		limit := 20
		offset := 0
		if t := req.FormValue("limit");t != ""{
			limit, err = strconv.Atoi(t)
			e(err, false)
		}
		if t := req.FormValue("offset"); t != ""{
			offset, err = strconv.Atoi(t)
			e(err, false)
		}

		total_count, err := mcol.Find(nil).Count()
		e(err, false)

		items := []map[string]interface{}{}

		iter := mcol.Find(nil).Skip(offset).Limit(limit).Iter()
		for {
			item := map[string]interface{}{}
			if iter.Next(&item) {
				items = append(items, item)
			} else {
				break
			}
		}

		dict := map[string]interface{}{
			"meta": map[string]int{"total_count": total_count, "offset": offset, "limit": limit,},
			"data": items,
			"msg": "success",
			"code": 0,
		}

		body, err := json.Marshal(dict)
		e(err, false)
		return string(body)
	})

	m.Run()
}
