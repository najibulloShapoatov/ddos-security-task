package main

import (
	"bufio"
	"context"
	"ddos-security/hashcash"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	timeoutFlag = flag.Duration("timeout", getEnvAsDuration("TIMEOUT", 15*time.Second), "timeout for client operations")
	portFlag    = flag.String("port", getEnv("PORT", "7891"), "port to listen on")
)

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp string
	Message   string
}

// ConnectionRequest contains the details for the TCP connection.
type ConnectionRequest struct {
	Host string
	Port string
}

func main() {
	// Handle root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
	})

	// Handle the connection request
	http.HandleFunc("/connect", handleConnectionRequest)

	// Start the web server
	log.Printf("Starting web server on http://0.0.0.0:%s", *portFlag)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *portFlag), nil))
}

func handleConnectionRequest(w http.ResponseWriter, r *http.Request) {
	// Parse the connection request
	var connReq ConnectionRequest
	err := json.NewDecoder(r.Body).Decode(&connReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform the TCP connection and PoW, return logs
	logs := connectToServer(connReq.Host, connReq.Port)
	json.NewEncoder(w).Encode(logs)
}

func connectToServer(host string, port string) (logs []LogEntry) {

	logFunc := func(format string, args ...any) {
		logs = append(logs, LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Message:   fmt.Sprintf(format, args...),
		})
	}

	serverAddress := fmt.Sprintf("%s:%s", host, port)

	logFunc("Open connection to %s by TCP with timeout %.2f seconds", serverAddress, timeoutFlag.Seconds())

	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeoutFlag)
	defer cancel()

	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", serverAddress)
	if err != nil {
		logFunc("Error connecting to server: %s", err.Error())
		return
	}
	defer func() {
		logFunc("Closing connection to %s", serverAddress)
		if err := conn.Close(); err != nil {
			logFunc("Error closing connection: %s", err.Error())
		}
	}()

	// Read the challenge from the server
	challenge, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		logFunc("Error reading challenge from server: %s", err.Error())
		return
	}
	challenge = strings.TrimSpace(challenge)
	logFunc("Get Challenge: %s", challenge)

	parse, err := hashcash.Parse(challenge)
	if err != nil {
		logFunc("Error solving challenge: %s", err.Error())
		return
	}

	if err := parse.Solve(ctx); err != nil {
		logFunc("Error solving challenge: %s", err.Error())
		return
	}

	logFunc("Solve challenge: %s", parse.Solution)

	if _, err := fmt.Fprintln(conn, parse.Solution); err != nil {
		logFunc("Error sending solution to server: %s", err.Error())
		return
	}

	// Read the quote from the server
	quote, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		logFunc("Error reading quote from server: %s", err.Error())
		return
	}
	logFunc("Quote from server: %s", quote)
	return
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsDuration(name string, fallback time.Duration) time.Duration {
	valueStr := getEnv(name, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return fallback
}
