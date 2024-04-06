package boostrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

var ConfigDir string

func FindConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dir, err := filepath.Abs(fmt.Sprintf("%v/.osmonitor/conf.json", home))
	if err != nil {
		log.Fatal(err)
	}

	ConfigDir = dir
	bufSize := 1024

	config := make([]byte, bufSize)
	file, err := os.OpenFile(dir, os.O_RDWR, 0666)
	if err != nil {
		// exec.Command("bash", "-c", fmt.Sprintf("mkdir -p %s/.osmonitor", home)).Output()

		// 2nd pass
		// file, err = os.OpenFile(dir, os.O_RDWR|os.O_CREATE, 0666)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		return &Config{}, errors.New("please run with 'init' or 'login' to login")
	}
	Check(err)
	defer file.Close()

	configStr := []byte{}

	for n, err := file.Read(config); err == nil; {
		configStr = append(configStr, config...)
		if n == 0 || n < bufSize {
			break
		}
	}
	Check(err)
	actualConfig := string(strings.Split(string(configStr), string('\x00'))[0])

	var data Config
	err = json.Unmarshal([]byte(actualConfig), &data)
	Check(err)

	return &data, err
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
