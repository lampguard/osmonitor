package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/lampguard/osmonitor/boostrap"
	"github.com/lampguard/osmonitor/parsers"
	"github.com/lampguard/osmonitor/support"
	flag "github.com/spf13/pflag"
)

// var base_url = "https://doppler-beta.up.railway.app/v1"

var base_url = "http://localhost:3000/v1"

var client *http.Client = &http.Client{
	Timeout: time.Second * 20,
}

var config *boostrap.Config

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

	for {
		stat := collateStats()
		var stats Stats = Stats{
			Iostat: string(stat[0]),
			Uptime: string(stat[1]),
			DF:     string(stat[2]),
		}

		jStats, _ := json.Marshal(stats)

		request, err := http.NewRequest("POST", fmt.Sprintf("%v/cli/logs", base_url), bytes.NewReader(jStats))
		if err != nil {
			log.Fatal(err)
		}

		request.Header.Add("Content-type", "application/json")
		request.Header.Add("Auth-Token", "")
		response, err := client.Do(request)

		if err != nil {
			log.Fatal(err)
		}

		_, err = io.ReadAll(request.Body)

		if err != nil {
			log.Println(err)
		}

		response.Body.Close()
		time.Sleep(time.Second * 5)
	}
}

func StartService() {
	fmt.Println("No command provided starting service")
	cfg, err := boostrap.FindConfig()
	boostrap.Check(err)

	config = cfg

	if config.Token == "" {
		log.Fatal("User not signed in")
	}

	reportStats()
}

func Login() {
	config, err := boostrap.FindConfig()
	if err != nil {
		config = &boostrap.Config{
			Token: "",
			Email: "",
		}
	}

	var email string = support.GetLine(support.Message{
		Message: "E-mail address: ",
	})
	password := support.GetPassword(support.Message{
		Message: "Password: ",
	})

	fmt.Println("logging in...")

	resp, err := http.Post(fmt.Sprintf("%s/cli/login", base_url), "application/json", strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, string(password))))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var p []byte
	p, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	type LoginData struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}
	loginData := &LoginData{}

	err = json.Unmarshal(p, loginData)
	if err != nil {
		panic(err)
	}

	if loginData.Status != 0 && loginData.Status != 200 {
		fmt.Println(loginData.Message)
		return
	}

	config.Token = loginData.Token
	config.Email = email
	configstr, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	os.WriteFile(boostrap.ConfigDir, configstr, 0666)
	fmt.Println("Login suceeded")
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
