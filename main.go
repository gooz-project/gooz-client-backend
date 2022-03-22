package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
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
	err := json.NewDecoder(r.Body).Decode(&newPayload)
	if err != nil {
		panic(err.Error())
	}
	if newPayload.Timer == 0 {
		newPayload.Timer = 10
	}
	replaceCommand(newPayload.Cmd, "writer.py")
	log.Println("Replace Command has been completed")
	log.Println("Command is running on the board...")
	output := runCommand("writer.py", newPayload.Timer)
	log.Println("Command Run has been completed")
	response := outPayload{
		Workspace: newPayload.Workspace,
		ComPort:   newPayload.ComPort,
		Cmd:       newPayload.Cmd,
		Output:    output,
	}
	clearCommand(newPayload.Cmd, "writer.py")
	log.Println("Command Script has been changed to default")
	json.NewEncoder(w).Encode(response)
}
func main() {
	log.Println("Program has been stated")
	log.Println("Backend is running on 6667 Port")
	http.HandleFunc("/exec", handler)
	log.Fatal(http.ListenAndServe(":6667", nil))
}

func replaceCommand(cmd string, filename string) {
	data, readerr := os.ReadFile(filename)
	if readerr != nil {
		panic(readerr.Error())
	}
	formattedData := string(data)
	formattedData = strings.Replace(formattedData, "changewillbecmd", cmd, 1)
	writeErr := os.WriteFile(filename, []byte(formattedData), 0644)
	if writeErr != nil {
		panic(writeErr.Error())
	}
}

func clearCommand(cmd string, filename string) {
	data, readerr := os.ReadFile(filename)
	if readerr != nil {
		panic(readerr.Error())
	}
	formattedData := string(data)
	formattedData = strings.Replace(formattedData, cmd, "changewillbecmd", 1)
	writeErr := os.WriteFile(filename, []byte(formattedData), 0644)
	if writeErr != nil {
		panic(writeErr.Error())
	}
}

func runCommand(filename string, timer int) (outPutBoard string) {
	cmd := exec.Command("ampy", "-p", newPayload.ComPort, "run", filename)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		panic(err.Error())
	}
	if err = cmd.Start(); err != nil {
		panic(err.Error())
	}
	tmp := make([]byte, 4096)
	time.Sleep(time.Duration(timer) * time.Second)
	for {
		_, err := stdout.Read(tmp)
		if err != nil {
			break
		}
	}
	var newTmp []byte
	for _, v := range tmp {
		if v != 0 {
			newTmp = append(newTmp, v)
		}
	}
	return string(newTmp)
}
