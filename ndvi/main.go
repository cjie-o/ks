package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Data struct {
	Time []int
	Lon  []float64
	Lat  []float64
	Ndvi []interface{}
	data []string
}

func main() {
	var data *Data
	tmp, err1 := ioutil.ReadFile("./ndvi.json")
	if err1 != nil {
		panic(" 读取文件异常")
	}

	data = &Data{}
	json.Unmarshal(tmp, data)

	txt := ""
	data.data = make([]string, len(data.Ndvi))
	for i := 0; i < len(data.Ndvi)/len(data.Time); i++ {

		low := i * len(data.Lat) * len(data.Lon)
		high := low + len(data.Lat)*len(data.Lon)

		for j := low; j < high; j++ {
			if data.Ndvi[j] == nil {
				continue
			}
			data.data[j] = strconv.FormatFloat(data.Ndvi[j].(float64), 'f', -1, 32)
		}

		txt += fmt.Sprintln(strings.Join(data.data[low:high], ","))
	}
	ioutil.WriteFile("./t.txt", []byte(txt), 0755)
}
