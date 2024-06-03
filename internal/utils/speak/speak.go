package speak

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/ln64-git/voxctl/internal/types"
	"github.com/ln64-git/voxctl/internal/utils/convo"
)

func ProcessSpeakText(state *types.AppState) {
	for result := range state.SpeakTextChan {
		var textResult types.TextResponse
		err := json.Unmarshal([]byte(result), &textResult)
		if err != nil {
			log.Printf("Failed to parse JSON: %v", err)
			continue
		}
		text := strings.TrimSpace(textResult.Text)
		if text != "" {
			state.SpeakText += text + " "
			// Handle conversation mode
			// if state.ConversationMode && len(strings.Fields(state.SpeakText)) >= 3 {
			if state.ConversationMode {
				convo.HandleConversation(state)
			}
		}
	}
}