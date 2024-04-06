package support

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type Message struct {
	Message string
}

func printMessage(message Message) {
	if message.Message != "" {
		fmt.Print(message.Message)
	}
}

func GetLine(message Message) string {
	var line string

	printMessage(message)

	_, err := fmt.Scanln(&line)
	if err != nil {
		panic(err)
	}

	return line
}

func GetInt(message Message) int {
	printMessage(message)

	var n string
	_, err := fmt.Scanln(&n)
	if err != nil {
		panic(err)
	}

	i, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}
	return i
}

func GetPassword(message Message) []byte {
	printMessage(message)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	return password
}

func GetToken(email string) string {
	resp, err := http.Post("https://doppler-beta.up.railway.app/v1/monitor/get-token", "application/json", strings.NewReader(""))

	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	return resp.Status
}
