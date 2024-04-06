package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/lampguard/osmonitor.git/boostrap"
	"github.com/lampguard/osmonitor.git/parsers"
	"github.com/lampguard/osmonitor.git/support"
	flag "github.com/spf13/pflag"
)

func collateStats() [][]byte {
	out, err := exec.Command("bash", "-c", "iostat | tail -1").Output()
	out2, err2 := exec.Command("bash", "-c", "df | tail -2  | head -1").Output()
	out3, err3 := exec.Command("uptime").Output()
	boostrap.Check(err)
	boostrap.Check(err2)
	boostrap.Check(err3)

	iostat := parsers.ParseOut(string(out))
	df := parsers.ParseOut(string(out2))
	uptime := parsers.ParseOut(string(out3))

	up, err := parsers.ParseUptime(uptime)
	boostrap.Check(err)

	DF, err := parsers.ParseDf(df)
	boostrap.Check(err)

	io, err := parsers.ParseIoStat(iostat)
	boostrap.Check(err)

	var stat [][]byte = [][]byte{io, up, DF}
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

	var email string = support.GetLine(support.Message{
		Message: "E-mail address: ",
	})

	fmt.Println(email)
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
