package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Data struct {
	Time []int
	Lon  []float64
	Lat  []float64
	PlantMtemp, PlantMpre, PlantM, PlantS,
	Ndvi []interface{}
}

func main() {
	var data *Data
	tmp, err1 := ioutil.ReadFile("./t.json")
	if err1 != nil {
		panic(" 读取文件异常")
	}

	data = &Data{}
	json.Unmarshal(tmp, data)

	dt := reflect.TypeOf(data).Elem()
	log.Println(dt.Kind())
	for i := 0; i < dt.NumField(); i++ {
		fname := dt.Field(i).Name
		if fname == "Time" || fname == "Lon" || fname == "Lat" {
			continue
		}
		log.Println(fname)
		d1 := reflect.ValueOf(data).Elem().Field(i).Interface().([]interface{})
		if len(d1) == 0 {
			continue
		}
		txt := ""
		rest := make([]string, len(d1))
		for i := 0; i < len(data.Time); i++ {

			low := i * len(data.Lat) * len(data.Lon)
			high := low + len(data.Lat)*len(data.Lon)

			for j := low; j < high; j++ {
				if d1[j] == nil {
					continue
				}
				rest[j] = strconv.FormatFloat(d1[j].(float64), 'f', -1, 32)
			}
			txt += fmt.Sprintln(i, ",", strings.Join(rest[low:high], ","))
		}
		if txt != "" {
			ioutil.WriteFile("./"+fname+".txt", []byte(txt), 0755)
		}
	}
}
