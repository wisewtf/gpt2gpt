# GPT2GPT

A Golang implementation of a simple ChatGPT API request written by ChatGPT (and fixed by a human) to be used from the CLI.

## The Story

A coworker of mine told me about ChatGPT today. I erroneously thought it was some version of DALL-E, or something similar. After trying it out I was left fascinated by the sheer amount of possibilities it gave me. I started writing some random queries and at some point, I thought... what if I make it write a BASH script so I could send queries from my terminal?

And that's what I did. My original query was: "Write a bash script to send queries to the ChatGPT API, specify the queries with an argument and parse the output using jq"

This is what I received:

```bash
#!/bin/bash

# Set your API key as a variable
api_key="API_KEY"

# Check if a query was provided as an argument
if [ $# -eq 0 ]; then
  echo "No query was provided."
  exit 1
fi

# Set the query as a variable
query=$1

# Send the request and store the response in a variable
response=$(curl -X POST "https://api.openai.com/v1/chatgpt/generations" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $api_key" \
  -d "{
    \"model\": \"chatgpt\",
    \"prompt\": \"$query\",
    \"max_tokens\":4000
  }")

# Extract the text key from the response using jq
text=$(echo "$response" | jq -r '.data[0].text')

# Print the text
echo "$text"
```

Right off the bat, this did not work. Mainly because:

1. The endpoint is wrong.
2. There is no `data` object in the response.
3. The `chatgpt` model does not exist.

So after some human finessing, I came up with this working script:

```bash
#!/bin/bash

# Set your API key as a variable
api_key="API_KEY"

# Check if a query was provided as an argument
if [ $# -eq 0 ]; then
  echo "No query was provided."
  exit 1
fi

# Set the query as a variable
query=$1

# Send the request and store the response in a variable
response=$(curl -X POST "https://api.openai.com/v1/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $api_key" \
  -d "{
    \"model\": \"text-davinci-002\",
    \"prompt\": \"$query\",
    \"max_tokens\":4000
  }")

# Extract the text key from the response using jq
text=$(echo "$response" | jq -r '.choices[0].text')

# Print the text
echo "$text"
```

And then I thought... "How boring... what if I just make it convert the script to Golang so I can build a binary file?". So I asked ChatGPT to do just that.

This is what I got:

```golang
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
    Model    string `json:"model"`
    Prompt   string `json:"prompt"`
    MaxTokens int    `json:"max_tokens"`
}

type OpenAIResponse struct {
    Choices []struct {
        Text string `json:"text"`
    } `json:"choices"`
}

func main() {
    // Set your API key as a variable
    apiKey := "API_KEY"

    // Check if a query was provided as an argument
    if len(os.Args) < 2 {
        fmt.Println("No query was provided.")
        os.Exit(1)
    }

    // Set the query as a variable
    query := os.Args[1]

    // Create a new OpenAI request
    req := &OpenAIRequest{
        Model:    "text-davinci-002",
        Prompt:   query,
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
    req, err = http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(reqBody))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
    resp, err := client.Do(req)
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
```

There was only 1 issue with this AI-generated code: the `req` was being used improperly. We guess that the AI knew where all 'requests' should be and just assigned the `req` variable to everything related to those. But at some point, this breaks the compiler because when assigning the `http.NewRequest` method to `req` it could not do it as `req` was already assigned as a pointer to a struct: `OpenAIRequest`.

So we just changed `req` to something else:

```golang
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
```

And that was it. The code worked and we could ask it simple questions.

## Conclusions

As a simple test, this is truly fascinating and it makes me wonder what will happen in 10 years if this technology matures further. I played with it a bit more and this AI was able to tell me how to create complex Active Directory GPOs and complex Ansible Playbooks which have worked dead-on without any further fixing.

It still requires some human intervention, which is completely understandable in my opinion. But going further, I hope it will always require human intervention. IMHO, this technology is useful only for people who know what they are doing. In terms of IT, at least. I have the same line of thought when it comes to things like Docker. I always say you should not run WordPress (or anything else) in Docker unless you know how it works and how to install it without it (Docker). This AI is the same. I asked it to write me a complete Ansible Playbook configuring a Linux VM with secured SSH, SSH keys in a specific user, install a LEMP stack and configure a WordPress virtual host using `php-fpm` configured to be generally optimized for WordPress and allow the upload of files over 1GB (Not a sane value).

It did everything perfectly, it even configured `max_body_size` in the virtual host correctly even when I erroneously only asked it to configure `php-fpm` for that, and not the virtual host.

## Can I use this code?

You can. You are limited to 96 characters for your query. I have not bothered optimizing anything, I didn't even allow flags or os variables for the API key, meaning I have not built it to be released inside this repo. But you can do whatever you want. (And you can't 'talk' with it. It won't remember your previous query like the web version. But it does use the same model.)

It's not my code anyway. Or is it? :)

## Disclaimer

I did the absolute bare minimum whatsoever to make this work. I barely read the API documentation and I tried fixing everything using ChatGPT itself. The only thing I had to think about was how it messed up the Golang code. That's it.
