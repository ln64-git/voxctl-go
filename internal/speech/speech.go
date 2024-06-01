package speech

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/ln64-git/voxctl/external/azure"
	"github.com/ln64-git/voxctl/internal/audio"
)

// AzureSpeechRequest represents a request to synthesize speech.
type AzureSpeechRequest struct {
	Text      string `json:"text"`
	Gender    string `json:"gender"`
	VoiceName string `json:"voiceName"`
}

// SanitizeText sanitizes and escapes text for JSON compatibility.
func SanitizeText(input string) string {
	// Replace newlines, carriage returns, and tabs with a space
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\r", " ")
	input = strings.ReplaceAll(input, "\t", " ")

	// Replace multiple spaces with a single space
	input = strings.Join(strings.Fields(input), " ")

	// Escape double quotes
	input = strings.ReplaceAll(input, `"`, `\"`)

	return input
}

// SpeechRequestToJSON converts a SpeechRequest to a JSON string.
func (r AzureSpeechRequest) SpeechRequestToJSON() string {
	sanitizedText := SanitizeText(r.Text)
	return fmt.Sprintf(`{"text":"%s","gender":"%s","voiceName":"%s"}`, sanitizedText, r.Gender, r.VoiceName)
}

// ProcessSpeech processes the speech request by synthesizing and playing the speech.
func ProcessSpeech(req AzureSpeechRequest, azureSubscriptionKey, azureRegion string, audioPlayer *audio.AudioPlayer) error {
	segments := SegmentedText(req.Text)
	for _, segment := range segments {
		audioData, err := azure.SynthesizeSpeech(azureSubscriptionKey, azureRegion, segment, req.Gender, req.VoiceName)
		if err != nil {
			log.Errorf("%s", err)
			return err
		}
		audioPlayer.Play(audioData)
		log.Infof("Speech processed: %s", segment) // Example log message
	}
	return nil
}

func SegmentTextFromChannel(tokenChan <-chan string, sentenceChan chan<- string) {
	defer close(sentenceChan)
	var builder strings.Builder

	for token := range tokenChan {
		builder.WriteString(token)
		if strings.ContainsAny(token, ",.!?") {
			sentence := builder.String()
			sentenceChan <- sentence
			builder.Reset()
		}
	}
}

// SegmentedText splits text into segments based on punctuation.
func SegmentedText(text string) []string {
	var sentences []string
	var currentSentence string
	for i, char := range text {
		if char == ',' || char == '.' || char == '!' || char == '?' {
			if currentSentence != "" {
				sentences = append(sentences, currentSentence)
				currentSentence = ""
			}
		} else {
			currentSentence += string(char)
			if i == len(text)-1 && currentSentence != "" {
				sentences = append(sentences, currentSentence)
			}
		}
	}
	return sentences
}
