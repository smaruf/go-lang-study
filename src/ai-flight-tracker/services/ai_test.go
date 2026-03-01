package services_test

import (
	"strings"
	"testing"

	. "github.com/smaruf/go-lang-study/src/ai-flight-tracker/services"
)

func newTestAI() *AIService {
	opensky := NewOpenSkyClient("https://opensky-network.org/api", "", "", 10)
	pricing := NewPricingService(300)
	return NewAIService("http://localhost:11434/api/generate", "llama2", 5, opensky, pricing)
}

func TestAskRoutesQuestion(t *testing.T) {
	ai := newTestAI()
	answer := ai.Ask("What routes are available from Dhaka?", "")
	// Even without Ollama, we expect a non-empty answer (fallback to context)
	if answer == "" {
		t.Error("expected non-empty answer")
	}
	if !strings.Contains(answer, "DAC") && !strings.Contains(answer, "Dhaka") && !strings.Contains(answer, "AI unavailable") {
		t.Errorf("expected answer to mention DAC or Dhaka, got: %s", answer)
	}
}

func TestAskGeneralQuestion(t *testing.T) {
	ai := newTestAI()
	answer := ai.Ask("What flights are available?", "")
	if answer == "" {
		t.Error("expected non-empty answer")
	}
}
