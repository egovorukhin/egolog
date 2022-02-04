package main

import (
	"errors"
	"github.com/egovorukhin/egolog"
	"log"
	"time"
)

const app = "app"

func main() {

	// Инициализируем Logger
	err := egolog.InitLogger("log.yml", nil)
	if err != nil {
		log.Fatal(err)
	}

	egolog.Info("Старт приложения")
	egolog.Error(errors.New("Какая то ошибка"))

	egolog.InfoSend("Старт приложения 1")
	egolog.ErrorSend("Какая то ошибка 1")

	egolog.ErrorFn(app, true, "Ошибка: %v %s", "Какая то ошибка", "Еще что то")
	egolog.DebugFn(app, false, "Какая то ошибка")

	time.Sleep(time.Second * 10)

	egolog.Info(true, "Остановка приложения")

}
