package main

import (
	"fmt"
	"github.com/cloudflare/conf"
	"log"
	"os"
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
	c, err := conf.ReadConfigFile("config.conf")
	e(err, true)
	MONGO_URL := c.GetString("MONGO_URL", "127.0.0.1:27017")
	MONGO_DB := c.GetString("MONGO_DB", "test")
	fmt.Println(MONGO_URL, MONGO_DB)
}
