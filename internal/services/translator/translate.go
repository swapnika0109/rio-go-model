package translator

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type Translator struct {
	client *translate.Client
	logger *log.Logger
}

func NewTranslator() *Translator {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create translate client: %v", err)
	}
	return &Translator{
		client: client,
		logger: log.New(log.Writer(), "Translator: ", log.LstdFlags),
	}
}

func (t *Translator) Translate(text string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := t.client.Translate(ctx, []string{text}, language.English, nil)
	if err != nil {
		t.logger.Printf("client.Translate: %v", err)
		return "", err
	}
	t.logger.Printf("Translated text: %s", resp[0].Text)
	return resp[0].Text, nil
}
