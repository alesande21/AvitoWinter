all: build run

install:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0

install_work_v:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

ver:
	oapi-codegen -version

genUserService:
	oapi-codegen --config=configs/oapi/configUserService.yaml api/shopAPI.yml

init:
	go mod init "AvitoWinter"

tidy:
	go mod tidy

.PHONY: all, install, ver, genUserService