package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var Conf map[string]interface{}

func init() {
	confData, err := ioutil.ReadFile("conf/conf.json")
	if err != nil {
		log.Fatal("error reading configuration file")
	}
	err = json.Unmarshal(confData, &Conf)
	if err != nil {
		log.Fatal("configuration file format error", err)
	}
	log.Println("%#v", Conf)
}
