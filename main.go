package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"strings"
)

type Data struct {
	Time     []int
	Lon      []float64
	Lat      []float64
	Pre, Tmp []interface{}
}

func main() {
	var data *Data
	tmp, err1 := ioutil.ReadFile("./tmp.json")
	pre, err2 := ioutil.ReadFile("./pre.json")
	if err1 != nil || err2 != nil {
		panic(" 读取文件异常")
	}

	data = &Data{}
	json.Unmarshal(tmp, data)
	json.Unmarshal(pre, data)
	a := com(data)

	a1, _ := json.Marshal(a)
	time, _ := json.Marshal(data.Time[0 : len(data.Time)/12])
	lon, _ := json.Marshal(data.Lon)
	lat, _ := json.Marshal(data.Lat)
	// fmt.Println()
	// d, err := json.Marshal(data)
	// log.Println(string(d), err)

	ndvi, _ := ioutil.ReadFile("ndvi.txt")
	ndvi1 := strings.Replace(string(ndvi), "ndvidata", string(a1[1:len(a1)-1]), 1)
	ndvi1 = strings.ReplaceAll(ndvi1, "timedata", string(time[1:len(time)-1]))
	ndvi1 = strings.ReplaceAll(ndvi1, "latdata", string(lat[1:len(lat)-1]))
	ndvi1 = strings.ReplaceAll(ndvi1, "londata", string(lon[1:len(lon)-1]))
	ioutil.WriteFile("./t.txt", []byte(ndvi1), 0755)
}

func com(data *Data) []float64 {
	len := len(data.Tmp) / 12
	a := make([]float64, len)
	// var a [len]float64
	for i := 0; i < len; i++ {
		is_continu := false
		var T, P float64
		T = 0
		P = 0
		for j := 12 * i; j < 12*(i+1); j++ {
			if data.Tmp[j] == nil || data.Pre[j] == nil {
				is_continu = true
				break
			}
			T += data.Tmp[j].(float64) / 12
			P += data.Pre[j].(float64) / 12
		}
		if is_continu {
			continue
		}

		tmp := 3000 / (1 + math.Exp(1.315-0.119*T))
		pre := 3000 * (1 - math.Exp(-0.000664*P))
		if tmp > pre {
			a[i] = pre
			continue
		}
		a[i] = tmp
	}
	return a
}

// 1.TSP( t) =3 000/ [ 1 ×+exp( 1.315 -0.119t)]( 4)
// 2. TSP( p ) =3 000[ 1 -exp( -0.000 664p )] ( 5)
