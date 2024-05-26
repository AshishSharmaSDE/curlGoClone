package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// RequestConfig holds the configuration for the HTTP request.
type RequestConfig struct {
	Method  string
	URL     string
	Data    string
	Headers map[string]string
	Output  string
}

// ReadFile reads the content of a file and returns it as a string.
func ReadFile(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return string(content), nil
}

// ParseArguments parses command line arguments into a RequestConfig struct.
func ParseArguments() RequestConfig {
	// Define flags for command line arguments
	method := flag.String("X", "GET", "HTTP method (GET, POST, PUT, DELETE)")
	data := flag.String("d", "", "HTTP request body data")
	headers := flag.String("H", "", "HTTP headers in the format 'Key: Value'")
	output := flag.String("o", "", "Output file to save response")

	// Parse the flags
	flag.Parse()

	// URL is a positional argument
	url := flag.Arg(0)

	// Initialize RequestConfig
	config := RequestConfig{
		Method:  *method,
		URL:     url,
		Data:    *data,
		Headers: make(map[string]string),
		Output:  *output,
	}

	// Parse headers
	if *headers != "" {
		headerPairs := strings.Split(*headers, ",")
		for _, header := range headerPairs {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				config.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	return config
}

// MakeRequest makes an HTTP request based on the given RequestConfig.
func MakeRequest(config RequestConfig) error {
	// Prepare the request body
	var requestBody *bytes.Reader
	if config.Data != "" {
		requestBody = bytes.NewReader([]byte(config.Data))
	} else {
		requestBody = bytes.NewReader([]byte{})
	}

	// Create the request
	req, err := http.NewRequest(config.Method, config.URL, requestBody)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	// Output the response
	if config.Output != "" {
		err = os.WriteFile(config.Output, body, 0644)
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
		fmt.Printf("Response data saved to %s\n", config.Output)
	} else {
		fmt.Println("Response data:", string(body))
	}

	return nil
}

// main is the entry point of the application.
func main() {
	// Parse the command line arguments
	config := ParseArguments()

	// Validate URL
	if config.URL == "" {
		fmt.Println("Error: URL is required")
		os.Exit(1)
	}

	// If data is from a file, read the file content
	if strings.HasPrefix(config.Data, "@") {
		filepath := strings.TrimPrefix(config.Data, "@")
		data, err := ReadFile(filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		config.Data = data
	}

	// Make the HTTP request
	err := MakeRequest(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
