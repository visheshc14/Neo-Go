package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"io"
	"os"
	"strconv"
	"time"
)

/// Generate a function that does nothing but return the given piece of static content
func generateStaticServerFn(content []byte) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Return 404 for all non-GET methods
		if string(ctx.Method()) != "GET" {
			ctx.NotFound()
			return
		}

		// Return the content
		ctx.SetBodyStream(bytes.NewReader(content), len(content))
	}
}

const (
	DEFAULT_PORT                       = 5000
	DEFAULT_STDIN_READ_TIMEOUT_SECONDS = 60
)

func main() {
	// Retrieve settings from ENV if present
	host, ok := os.LookupEnv("HOST")

	var port int
	portRaw, ok := os.LookupEnv("PORT")
	if !ok {
		port = DEFAULT_PORT
	} else {
		parsedPort, err := strconv.Atoi(portRaw)
		if err != nil {
			log.Fatal("Failed to parse port")
		}
		port = parsedPort
	}

	var stdinReadTimeoutSeconds int
	stdinReadTimeoutSecondsRaw, ok := os.LookupEnv("STDIN_READ_TIMEOUT_SECONDS")
	if !ok {
		stdinReadTimeoutSeconds = DEFAULT_STDIN_READ_TIMEOUT_SECONDS
	} else {
		parsedStdinReadTimeoutSeconds, err := strconv.Atoi(stdinReadTimeoutSecondsRaw)
		if err != nil {
			log.Fatal("Failed to parse STDIN read timeout seconds")
		}
		stdinReadTimeoutSeconds = parsedStdinReadTimeoutSeconds
	}

	filePath, ok := os.LookupEnv("FILE")

	// Retrieve settings from flags
	hostPtr := flag.String("host", "", "Host")
	portPtr := flag.Int("port", -1, "Port")
	stdinReadTimeoutSecondsPtr := flag.Int(
		"stdin-read-timeout-seconds",
		-1,
		"Amount of seconds to wait for input on STDIN to serve",
	)
	filePathPtr := flag.String("file", "", "File to read")
	flag.Parse()

	// Override ENV with flags if necessary
	if *hostPtr != "" {
		host = *hostPtr
	}

	if *portPtr != -1 {
		port = *portPtr
	}

	if *filePathPtr != "" {
		filePath = *filePathPtr
	}

	if *stdinReadTimeoutSecondsPtr != -1 {
		stdinReadTimeoutSeconds = *stdinReadTimeoutSecondsPtr
	}

	// TODO: Validate settings (ex. port has to be non-negative)

	// Build address from host and port
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Info(fmt.Sprintf("Server configured to run @ [%s]", addr))

	// Create content to be filled in later
	var content []byte

	// Check if filepath was provided
	if filePath != "" {
		// Load file into content

		// Check if the file exists
		if _, err := os.Stat(filePath); err != nil {
			log.Fatal(fmt.Sprintf("Failed fo find file [%s]", filePath))
		}

		// Read from file
		contentBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed fo read file @ [%s]", filePath))
		}

		log.Info(fmt.Sprintf("Reading file from path [%s]", filePath))
		content = contentBytes
	} else {
		// Attempt to read content from STDIN

		log.Info(fmt.Sprintf("No file path provided, waiting for input on STDIN (max %d seconds)...", stdinReadTimeoutSeconds))
		// No file path is present, attempt to read from STDIN

		stdinContent, err := readStdinWithTimeout(stdinReadTimeoutSeconds)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to read from STDIN after waiting %d seconds", stdinReadTimeoutSeconds))
		}

		content = stdinContent
	}

	// Ensure we have *some* content at this point
	if content == nil {
		log.Fatal("No file contents -- please ensure you've specified a file or fed in data via STDIN")
	}

	handler := generateStaticServerFn(content)

	// Start the server
	fmt.Println("Starting the server...")
	log.Fatal(fasthttp.ListenAndServe(addr, handler))
}

func readStdinWithTimeout(timeoutSeconds int) ([]byte, error) {
	// Create channel for timeout
	ch := make(chan int)

	var result []byte

	// Spawn goroutine to attempt to read STDIN
	go func() {
		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to read from STDIN after waiting %d seconds", timeoutSeconds))
		}
		result = stdinBytes
		ch <- 1
	}()

	// Wait for ReadAll or timeout
	select {
	// Read STDIN
	case <-ch:
		log.Info(fmt.Sprintf("Successfully read input from STDIN"))
		return result, nil

		// Timeout
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		return nil, errors.New("Failed to read from STDIN")
	}

}
