package main

import "order_service/internal/app"

func main() {
	err := app.Init()
	if err != nil {
		return
	}

	app.Run()
}
