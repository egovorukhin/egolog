package main

import (
	"errors"
	"github.com/egovorukhin/egolog"
	"log"
)

const app = "app"

func main() {

	// Инициализируем Logger
	err := egolog.InitLogger("log.yml", nil)
	if err != nil {
		log.Fatal(err)
	}

	egolog.Info(true, "Старт приложения")
	egolog.Error(true, errors.New("Какая то ошибка"))

	egolog.ErrorFn(app, true, "Ошибка: %v %s", "Какая то ошибка", "Еще что то")
	egolog.DebugFn(app, false, "Какая то ошибка")

	egolog.Info(true, "Остановка приложения")

}
