all: build run

install:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0

install_work_v:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

ver:
	oapi-codegen -version

genService:
	oapi-codegen --config=configs/oapi/configShopService.yaml api/shopAPI.yml

init:
	go mod init "AvitoWinter"

tidy:
	go mod tidy

up:
	docker compose -f docker-compose.yaml up --build

down:
	docker compose down -v

nano:
	nano ~/.docker/config.json

.PHONY: all, install, ver, genUserService, nano