package parsers

import (
	"encoding/json"
	"strconv"
	"strings"
)

// iostat
//
//								 disk0       cpu    load average
//	    KB/t  tps  MB/s  us sy id   1m   5m   15m
type IoStat struct {
	MBPS float64 `json:"mbps"`
	CPU  struct {
		User   float64 `json:"user"`
		System float64 `json:"system"`
		Idle   float64 `json:"idle"`
	} `json:"cpu"`
	LoadAvg struct {
		OneM     float64 `json:"1m"`
		FiveM    float64 `json:"5m"`
		FifteenM float64 `json:"15m"`
	} `json:"load_average"`
}

func ParseIoStat(data []string) ([]byte, error) {
	var iostat IoStat

	fifteenM := strings.Trim(data[8], "\n")

	iostat.MBPS, _ = strconv.ParseFloat(data[2], 64)
	iostat.CPU.User, _ = strconv.ParseFloat(data[3], 64)
	iostat.CPU.System, _ = strconv.ParseFloat(data[4], 64)
	iostat.CPU.Idle, _ = strconv.ParseFloat(data[5], 64)
	iostat.LoadAvg.FifteenM, _ = strconv.ParseFloat(fifteenM, 64)
	iostat.LoadAvg.OneM, _ = strconv.ParseFloat(data[6], 64)
	iostat.LoadAvg.FiveM, _ = strconv.ParseFloat(data[7], 64)

	return json.Marshal(iostat)
}

// df
// Filesystem     512-blocks      Used Available Capacity iused      ifree %iused  Mounted on
// /dev/disk1s1s1  976490576  19999072 338663312     6%  403755 1693316560    0%   /
type Df struct {
	Used      int    `json:"used"`
	Available int    `json:"available"`
	Capacity  string `json:"capacity"`
}

func ParseDf(data []string) ([]byte, error) {
	var df Df

	used, _ := strconv.Atoi(data[2])
	available, _ := strconv.Atoi(data[3])
	df.Available = available
	df.Used = used
	df.Capacity = data[4]
	return json.Marshal(df)
}

// uptime
// 16:20  up 1 day, 19:44, 1 user, load averages: 3.38 3.10 3.14
type Uptime struct {
	Time   string `json:"time"`
	Uptime string `json:"uptime"`
	UpFor  string `json:"upfor"`
	Users  int    `json:"users"`
}

func ParseUptime(data []string) ([]byte, error) {
	var uptime Uptime

	uptime.Time = data[0]
	uptime.Uptime = data[2] + " " + data[3]
	uptime.UpFor = data[4]

	users, _ := strconv.Atoi(data[5])
	uptime.Users = users

	return json.Marshal(uptime)
}

func ParseOut(stats string) []string {
	var r []string
	temp := strings.Split(stats, " ")
	for _, o := range temp {
		if o != "" && o != "," {
			r = append(r, o)
		}
	}

	return r
}
