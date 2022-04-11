package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"strings"
)

type Data struct {
	Time []int
	Lon  []float64
	Lat  []float64
	Pre, Tmp, Ndvi,
	PlantS,
	PlantMtemp, PlantMpre, PlantM []interface{}
	Data [][][]map[string]float64
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

	// d, _ := json.Marshal(data)
	// log.Println(string(d))
	// a := com(data)

	PlantMtemp, _ := json.Marshal(data.PlantMtemp)
	PlantMpre, _ := json.Marshal(data.PlantMpre)
	PlantM, _ := json.Marshal(data.PlantM)
	PlantS, _ := json.Marshal(data.PlantS)

	time, _ := json.Marshal(data.Time[0 : len(data.Time)/12])
	lon, _ := json.Marshal(data.Lon)
	lat, _ := json.Marshal(data.Lat)
	// fmt.Println()
	// d, err := json.Marshal(data)
	// log.Println(string(d), err)

	plant, _ := ioutil.ReadFile("data.txt")
	plant1 := strings.ReplaceAll(string(plant), "plantMtempdata", string(PlantMtemp[1:len(PlantMtemp)-1]))
	plant1 = strings.ReplaceAll(plant1, "plantMpredata", string(PlantMpre[1:len(PlantMpre)-1]))
	plant1 = strings.ReplaceAll(plant1, "plantMdata", string(PlantM[1:len(PlantM)-1]))
	plant1 = strings.ReplaceAll(plant1, "plantSdata", string(PlantS[1:len(PlantS)-1]))

	plant1 = strings.ReplaceAll(plant1, "timedata", string(time[1:len(time)-1]))
	plant1 = strings.ReplaceAll(plant1, "latdata", string(lat[1:len(lat)-1]))
	plant1 = strings.ReplaceAll(plant1, "londata", string(lon[1:len(lon)-1]))
	plant1 = strings.ReplaceAll(plant1, "null", "_")
	ioutil.WriteFile("./t.txt", []byte(plant1), 0755)
	// log.Println(plant1)
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
	a1 := make([]interface{}, len(D.Lon)*len(D.Lat)*len(D.Time)/10/12)
	a2 := make([]interface{}, len(D.Lon)*len(D.Lat)*len(D.Time)/10/12)
	a3 := make([]interface{}, len(D.Lon)*len(D.Lat)*len(D.Time)/10/12)
	b := make([]interface{}, len(D.Lon)*len(D.Lat)*len(D.Time)/10/12)

	D.PlantMtemp = a1
	D.PlantMpre = a2
	D.PlantM = a3
	D.PlantS = b

	for i := 0; i < len(D.Lon); i++ {
		for j := 0; j < len(D.Lat); j++ {
			for k := 0; k < len(D.Time)/12/10; k++ {
				is_continu := false
				var T, P float64
				T = 0
				P = 0
				for l := 0; l < 10*12; l++ {
					location := l*len(D.Lat)*len(D.Lon) + k*len(D.Lat)*len(D.Lon)*10*12 + j*len(D.Lon) + i
					if D.Tmp[location] == nil || D.Pre[location] == nil {
						is_continu = true
						log.Println("lon : ", i, "lat:", j, "time:", l)
						break
					}
					T += D.Tmp[location].(float64) / 10
					P += D.Pre[location].(float64) / 10 * 12
				}

				if is_continu {
					continue
				}

				tmp := 3000 / (1 + math.Exp(1.315-0.119*T))
				pre := 3000 * (1 - math.Exp(-0.000664*P))
				location := k*len(D.Lat)*len(D.Lon) + j*len(D.Lon) + i
				a1[location] = tmp
				a2[location] = pre

				// s
				s := func() float64 {
					l := 300 + 25*T + 0.05*math.Pow(T, 3)
					v := 1.05 * P / math.Sqrt(1+math.Pow(1.05*P/l, 2))
					s := 3000 * (1 - math.Exp(-0.0009695*(v-20)))
					return s
				}
				b[location] = s()

				if tmp > pre {
					a3[location] = pre
					continue
				}
				a3[location] = tmp
			}
		}
	}
}

// 1.TSP( t) =3 000/ [ 1 ×+exp( 1.315 -0.119t)]( 4)
// 2. TSP( p ) =3 000[ 1 -exp( -0.000 664p )] ( 5)
