package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	const testKey = "TEST_ENV_VAR"
	const testValue = "test value"

	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	if got := getEnv(testKey, "default"); got != testValue {
		t.Errorf("getEnv() = %v, want %v", got, testValue)
	}

	if got := getEnv("NON_EXISTENT_VAR", "default"); got != "default" {
		t.Errorf("getEnv() = %v, want 'default'", got)
	}
}

func TestGetEnvAsDuration(t *testing.T) {
	const testKey = "TEST_DURATION_VAR"
	testValue := "2m"
	expectedDuration, _ := time.ParseDuration(testValue)

	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	if got := getEnvAsDuration(testKey, 1*time.Minute); got != expectedDuration {
		t.Errorf("getEnvAsDuration() = %v, want %v", got, expectedDuration)
	}

	if got := getEnvAsDuration("NON_EXISTENT_VAR", 1*time.Minute); got != 1*time.Minute {
		t.Errorf("getEnvAsDuration() = %v, want 1m", got)
	}
}

// mockTCPServer starts a mock TCP server for testing.
// It sends a predetermined challenge and expects a response.
func mockTCPServer(t *testing.T, port string, challenge string, expectedResponse string) {
	l, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		t.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer l.Close()

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("Failed to accept connection: %v", err)
		}
		defer conn.Close()

		// Send challenge
		fmt.Fprintln(conn, challenge)

		// Read response
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		if response != expectedResponse+"\n" {
			t.Errorf("Expected response %q, got %q", expectedResponse, response)
		}
	}()
}

// TestConnect2TCPServer tests the connect2TCPServer method.
func TestConnect2TCPServer(t *testing.T) {
	// Setting up a mock TCP server
	testPort := "7892" // Use a different port than your actual server
	testChallenge := "test-challenge"
	expectedResponse := "test-response"
	mockTCPServer(t, testPort, testChallenge, expectedResponse)

	// Give the server a moment to start
	time.Sleep(time.Second)

	// Call the method to test
	logs := connectToServer("127.0.0.1", testPort)

	t.Logf("%v", logs)
}
