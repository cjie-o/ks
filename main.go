package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
)

type Data struct {
	Time                  []int
	Lon                   []float64
	Lat                   []float64
	Pre, Tmp, Ndvi, Plant []interface{}
	Data                  [][][]map[string]float64
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

	//
	data.com()

	d, _ := json.Marshal(data)
	log.Println(string(d))
	// a := com(data)

	// a1, _ := json.Marshal(a)
	// time, _ := json.Marshal(data.Time[0 : len(data.Time)/12])
	// lon, _ := json.Marshal(data.Lon)
	// lat, _ := json.Marshal(data.Lat)
	// // fmt.Println()
	// // d, err := json.Marshal(data)
	// // log.Println(string(d), err)

	// ndvi, _ := ioutil.ReadFile("ndvi.txt")
	// ndvi1 := strings.Replace(string(ndvi), "ndvidata", string(a1[1:len(a1)-1]), 1)
	// ndvi1 = strings.ReplaceAll(ndvi1, "timedata", string(time[1:len(time)-1]))
	// ndvi1 = strings.ReplaceAll(ndvi1, "latdata", string(lat[1:len(lat)-1]))
	// ndvi1 = strings.ReplaceAll(ndvi1, "londata", string(lon[1:len(lon)-1]))
	// ioutil.WriteFile("./t.txt", []byte(ndvi1), 0755)
	// // log.Println(ndvi1)
}
func (D *Data) marsh() {
	D.Data = make([][][]map[string]float64, len(D.Lon))
	for i := 0; i < len(D.Lon); i++ {
		D.Data[i] = make([][]map[string]float64, len(D.Lon))
		for j := 0; j < len(D.Lat); j++ {
			D.Data[i][j] = make([]map[string]float64, len(D.Time))
			for k := 0; k < len(D.Time); k++ {
				location := k*len(D.Lat)*len(D.Lon) + j*len(D.Lon) + i
				if D.Tmp[location] == nil || D.Pre[location] == nil {
					// log.Println("i : ", i, "j:", j)
					continue
				}
				D.Data[i][j][k] = map[string]float64{}
				D.Data[i][j][k]["tmp"] = D.Tmp[location].(float64)
				D.Data[i][j][k]["pre"] = D.Pre[location].(float64)
			}
		}
	}

}

func (D *Data) com() {
	a := make([]interface{}, len(D.Pre)/12)
	D.Plant = a
	for i := 0; i < len(D.Lon); i++ {
		for j := 0; j < len(D.Lat); j++ {
			for k := 0; k < len(D.Time)/12; k++ {
				is_continu := false
				var T, P float64
				T = 0
				P = 0

				for l := 0; l < 12; l++ {
					location := (12*k+l)*len(D.Lat)*len(D.Lon) + j*len(D.Lon) + i
					if D.Tmp[location] == nil || D.Pre[location] == nil {
						is_continu = true
						log.Println("lon : ", i, "lat:", j, "time:", 12*k+l)
						break
					}
					T += D.Tmp[location].(float64) / 12
					P += D.Pre[location].(float64) / 12
				}
				if is_continu {
					continue
				}

				tmp := 3000 / (1 + math.Exp(1.315-0.119*T))
				pre := 3000 * (1 - math.Exp(-0.000664*P))
				if tmp > pre {
					a[k] = pre
					continue
				}
				a[k] = tmp
			}
		}
	}
}

// 1.TSP( t) =3 000/ [ 1 ×+exp( 1.315 -0.119t)]( 4)
// 2. TSP( p ) =3 000[ 1 -exp( -0.000 664p )] ( 5)
