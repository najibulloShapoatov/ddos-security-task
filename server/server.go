package main

import (
	"bufio"
	"ddos-security/hashcash"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	difficultyFlag = flag.Int("difficulty", getEnvAsInt("DIFFICULTY", 20), "difficulty for hashcash algorithm")
	timeoutFlag    = flag.Duration("timeout", getEnvAsDuration("TIMEOUT", 15*time.Second), "timeout for client operations")
	portFlag       = flag.String("port", getEnv("PORT", "7890"), "port to listen on")
)

func main() {
	loadQuotes("/data")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", *portFlag))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("TCP server listening on port 7890...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Set timeout for client operations
	conn.SetDeadline(time.Now().Add(*timeoutFlag))

	hc := hashcash.New("ddos-security", *difficultyFlag)

	_, err := fmt.Fprintln(conn, hc.String())
	if err != nil {
		fmt.Println("Error sending challenge to client:", err)
		return
	}

	// Read response
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response from client:", err)
		return
	}
	response = strings.TrimSpace(response)

	hc.Solution = response
	// Verify response
	if err := hc.Verify(); err == nil {
		// Send quote if PoW is valid
		quote := Quotes[rand.Intn(len(Quotes))]
		_, err := fmt.Fprintln(conn, quote)
		if err != nil {
			fmt.Println("Error sending quote to client:", err)
		}
	} else {
		fmt.Fprintln(conn, "Invalid PoW")
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(name string, fallback int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
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
