.PHONY: buildWebhook setWebhook run test translate

buildWebhook:
	@go build -o ./cmd/setTelegramWebhookUrl ./cmd/setTelegramWebhookUrl.go

setWebhook: buildWebhook
	@./cmd/setTelegramWebhookUrl $(url)
	@rm -f ./cmd/setTelegramWebhookUrl

run:
	@go run ./cmd/app/main.go

test:
	@go test -v ./...

translate:
	@goi18n merge i18n/active.*.toml
	@echo "Translate all the messages in the translate.*.toml files"
	@echo "Run goi18n merge active.*.toml translate.*.toml"
	@echo "Replace files in i18n dir with the new files"

