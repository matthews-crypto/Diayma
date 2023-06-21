package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func handleChatbotRequest(w http.ResponseWriter, r *http.Request) {
	projectID := "YOUR_PROJECT_ID"
	lang := "en" // Language code for the conversation

	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Dialogflow session client: %v", err)
	}

	sessionID := "unique-session-id"
	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	query := r.FormValue("message")

	request := dialogflowpb.DetectIntentRequest{
		Session: sessionPath,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         query,
					LanguageCode: lang,
				},
			},
		},
	}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		log.Fatalf("Failed to detect intent: %v", err)
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()

	fmt.Fprintf(w, "Response: %s", fulfillmentText)
}

func main() {
	http.HandleFunc("/chatbot", handleChatbotRequest)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
