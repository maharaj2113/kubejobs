#Kubernetes Translation Job Service
This Go application exposes a RESTful API that allows users to send a text message for translation. When a POST request is made to the /translate endpoint, the application creates a Kubernetes Job using the Kubernetes API. The Job then uses the Google Translate API to translate the text and outputs the result.

#Prerequisites
Go 1.20+: Install Go from golang.org.
Kubernetes Cluster: Ensure that you have access to a Kubernetes cluster and the kubectl command is configured.
Docker: Ensure Docker is installed and running.
Kubernetes Configuration: Ensure your Kubernetes configuration file is set up correctly and accessible via the KUBECONFIG environment variable.
RapidAPI Account: Sign up for a free RapidAPI account and subscribe to the Google Translate API to obtain an API key.

#Install Dependencies:
Run the following command to install the necessary Go modules: go mod tidy

#Set Environment Variables:
Set your RapidAPI key and Kubernetes configuration file path:
export RAPIDAPI_KEY="your-rapidapi-key"
export KUBECONFIG="/path/to/your/kubeconfig"

#Start the Go application:
go run main.go

#API Usage
POST /translate
This endpoint accepts a JSON payload to translate text from one language to another.

curl --request POST \
  --url http://localhost:8080/translate \
  --header "Content-Type: application/json" \
  --data '{
    "q": "The Great Pyramid of Giza is the oldest and largest of the three pyramids in the Giza pyramid complex.",
    "source": "en",
    "target": "es",
    "format": "text"
  }'

Request Body:

q (string): The text to translate.
source (string): The source language code (e.g., "en" for English).
target (string): The target language code (e.g., "es" for Spanish).
format (string): The format of the input text (e.g., "text").
Response:

The server responds with an HTTP status code 202 Accepted to indicate that the Job was created successfully.

