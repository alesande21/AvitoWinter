package main

import (
	"AvitoWinter/internal/app"
	log2 "github.com/sirupsen/logrus"
)

func main() {

	err := app.RunApp()
	if err != nil {
		log2.Errorf("app.RunApp%v", err)
		return
	}

}
