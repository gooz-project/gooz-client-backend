package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type payload struct {
	Workspace string `json:"workspace"`
	ComPort   string `json:"comPort"`
	Cmd       string `json:"cmd"`
	Timer     int    `json:"timer"`
}

type outPayload struct {
	Workspace string `json:"workspace"`
	ComPort   string `json:"comPort"`
	Cmd       string `json:"cmd"`
	Output    string `json:"output"`
}

var newPayload payload

func handler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	json.NewDecoder(r.Body).Decode(&newPayload)
	if newPayload.Cmd != "" {
		log.Println("Command is running on the board...")
		output := runCommand()
		log.Println("Command Run has been completed")
		response := outPayload{
			Workspace: newPayload.Workspace,
			ComPort:   newPayload.ComPort,
			Cmd:       newPayload.Cmd,
			Output:    output,
		}
		fmt.Println(response.Output)
		log.Println("Command Script has been changed to default")
		json.NewEncoder(w).Encode(response)
	} else {
		response := outPayload{
			Workspace: newPayload.Workspace,
			ComPort:   newPayload.ComPort,
			Cmd:       newPayload.Cmd,
			Output:    "You has been sent empty command",
		}
		json.NewEncoder(w).Encode(response)
	}

}

func main() {
	command := exec.Command("npm", "start")
	command.Dir = "../gooz-client"
	err := command.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client has been started")
	log.Println("Program has been started")
	log.Println("Backend is running on 5000 Port")
	http.HandleFunc("/exec", handler)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func runCommand() (outPutBoard string) {
	c := &serial.Config{Name: newPayload.ComPort, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	data := []byte(newPayload.Cmd + "\r")
	n, err := s.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 1000)
	buf := make([]byte, 2048)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	senderData := string(buf[:n])
	sender := strings.Split(senderData, "\n")
	sendData := ""
	for k, v := range sender {
		if k != 0 && k != len(sender)-1 {
			sendData += v
			sendData += "\n"
		}
	}
	println(sendData)
	s.Close()
	return sendData
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
