.PHONY: buildWebhook setWebhook run

buildWebhook:
	go build -o ./cmd/setTelegramWebhookUrl ./cmd/setTelegramWebhookUrl.go

setWebhook: buildWebhook
	./cmd/setTelegramWebhookUrl $(url)
	rm -f ./cmd/setTelegramWebhookUrl

run:
	go run ./cmd/app/main.go

