.PHONY: buildWebhook setWebhook run run/live test translate audit

MAIN_PACKAGE_PATH := ./cmd/app
BINARY_NAME := wordy

build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

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
	@echo "Replace files in i18n dir with the new active.*.toml files"

run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
			--build.cmd "make build" --build.bin "/tmp/bin/${BINARY_NAME}" --build.delay "100" \
			--build.exclude_dir "" \
			--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
			--misc.clean_on_exit "true"

audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

