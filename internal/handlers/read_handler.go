package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/ln64-git/voxctl/internal/types"
	"github.com/ln64-git/voxctl/internal/utils/read"
)

func HandleReadText(w http.ResponseWriter, r *http.Request, state *types.AppState) {
	// Process the Azure speech request
	var speechReq read.AzureSpeechRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&speechReq)
	if err != nil {
		log.Errorf("Failed to process speech request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read the text using the processed request
	err = read.ReadText(speechReq, state.AzureSubscriptionKey, state.AzureRegion, state.AudioPlayer)
	if err != nil {
		log.Errorf("Failed to process speech: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
