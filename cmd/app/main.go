package main

import "AvitoWinter/internal/app"

func main() {

	err := app.RunApp()
	if err != nil {
		return
	}

}
