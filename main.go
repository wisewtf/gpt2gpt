package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type OpenAIRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	// Set your API key as a variable
	apiKey := "SET_YOUR_API_KEY"

	// Check if a query was provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("No query was provided.")
		os.Exit(1)
	}

	// Set the query as a variable
	query := os.Args[1]

	// Create a new OpenAI request
	req := &OpenAIRequest{
		Model:     "text-davinci-002",
		Prompt:    query,
		MaxTokens: 4000,
	}

	// Convert the request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Send the request and store the response in a variable
	client := &http.Client{}
	request, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmarshal the response into an OpenAI response struct
	var openAIResp OpenAIResponse
	err = json.Unmarshal(respBody, &openAIResp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the text
	fmt.Println(openAIResp.Choices[0].Text)
}
