package boostrap

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func FindConfig() (Config, error) {
	dir, err := filepath.Abs("./conf.json")
	if err != nil {
		log.Fatal(err)
	}

	bufSize := 1024

	config := make([]byte, bufSize)
	file, err := os.OpenFile(dir, os.O_RDWR|os.O_CREATE, 0666)
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

	var data struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}
	err = json.Unmarshal([]byte(actualConfig), &data)
	Check(err)

	return data, err
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
