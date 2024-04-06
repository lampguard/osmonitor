package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/lampguard/osmonitor.git/boostrap"
	flag "github.com/spf13/pflag"
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

func parseIoStat(data []string) ([]byte, error) {
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

func parseDf(data []string) ([]byte, error) {
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

func parseOut(stats string) []string {
	var r []string
	temp := strings.Split(stats, " ")
	for _, o := range temp {
		if o != "" && o != "," {
			r = append(r, o)
		}
	}

	return r
}

func parseUptime(data []string) ([]byte, error) {
	var uptime Uptime

	uptime.Time = data[0]
	uptime.Uptime = data[2] + " " + data[3]
	uptime.UpFor = data[4]

	users, _ := strconv.Atoi(data[5])
	uptime.Users = users

	return json.Marshal(uptime)
}

func collateStats() [][]byte {
	out, err := exec.Command("bash", "-c", "iostat | tail -1").Output()
	out2, err2 := exec.Command("bash", "-c", "df | tail -2  | head -1").Output()
	out3, err3 := exec.Command("uptime").Output()
	boostrap.Check(err)
	boostrap.Check(err2)
	boostrap.Check(err3)

	iostat := parseOut(string(out))
	df := parseOut(string(out2))
	uptime := parseOut(string(out3))

	up, err := parseUptime(uptime)
	boostrap.Check(err)

	DF, err := parseDf(df)
	boostrap.Check(err)

	io, err := parseIoStat(iostat)
	boostrap.Check(err)

	var stat [][]byte = [][]byte{io, up, DF}
	// fmt.Println(iostat, len(iostat))
	// fmt.Println(df, len(df))
	// fmt.Println(uptime, len(uptime))

	// stat += fmt.Sprintf(`%d: %v`, len(iostat[0]), df[0])

	return stat
}

func reportStats() {
	type Stats struct {
		Iostat string `json:"iostat"`
		Uptime string `json:"uptime"`
		DF     string `json:"df"`
	}
	stat := collateStats()
	var stats Stats = Stats{
		Iostat: string(stat[0]),
		Uptime: string(stat[1]),
		DF:     string(stat[2]),
	}

	jStats, _ := json.Marshal(stats)

	response, err := http.Post("https://doppler-beta.up.railway.app/v1/logs", "application/json", bytes.NewReader(jStats))

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var p []byte

	fmt.Println("Status:", response.Status)
	// fmt.Println(response.Body.Read())
	for i, err := response.Body.Read(p); err == nil; {
		if i == 0 {
			break
		}
		fmt.Println(i)
		fmt.Println(string(p))
	}

	// for {
	// 	time.Sleep(time.Second * 5)
	// }
}

func StartService() {
	fmt.Println("No command provided starting service")
	config, err := boostrap.FindConfig()
	boostrap.Check(err)

	if config.Token == "" {
		log.Fatal("User not signed in")
	}

	reportStats()
}

func Login() {
	fmt.Println("Logging in")
}

func main() {
	login := flag.Bool("login", false, "")
	flag.Parse()

	if *login {
		Login()
		return
	}

	if flag.NArg() == 0 {
		StartService()
		return
	}
}
