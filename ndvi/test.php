<?php
$a = json_decode(
    file_get_contents('./ndvi.json')
);
foreach ($a as $key => $value) {
    if ($key === "time" || $key === "lat" || $key === "lon") {
        continue;
    }


    for ( $i = 0; i < count($value)/count($a["time"]); $i++ ) { 


		$low = $i * count($a["lat"])*count($a["lon"]);
		$high = $low + count($a["lat"])*count($a["lon"]);

		for j := low; j < high; j++ {
			if data.Ndvi[j] == nil {
				continue
			}
			data.data[j] = strconv.FormatFloat(data.Ndvi[j].(float64), 'f', -1, 32)
		}
		
		txt += fmt.Sprintln(strings.Join(data.data[low:high], ","))
	}
}
