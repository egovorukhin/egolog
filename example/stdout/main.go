package main

import (
	"errors"
	"github.com/egovorukhin/egolog"
	"log"
)

func main() {

	// Инициализируем Logger
	err := egolog.InitLogger("log.yml")
	if err != nil {
		log.Fatal(err)
	}

	egolog.Info("Старт приложения")

	egolog.Error(errors.New("Какая то ошибка"))
	egolog.Error("Ошибка: %v %s", "Какая то ошибка", "Еще что то")
	egolog.Debug("Какая то ошибка")

	egolog.Info("Остановка приложения")

}
