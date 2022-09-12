package egolog

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

const app = "app"

func Test(t *testing.T) {

	cfg := Config{
		DirPath: "logs",
		Info:    "3",
		Error:   "3 | 16",
		Debug:   "1 | 4 | 8",
		Rotation: &Rotation{
			Size:   10240,
			Format: "%name_%time",
		},
	}

	callback := func(filename, prefix, m string, message interface{}, v ...interface{}) {
		fmt.Printf("name: %s, prefix: %s, message: %s", filename, prefix, m)
		m = fmt.Sprintf(message.(string), v...)
		fmt.Printf("name: %s, prefix: %s, message: %s", filename, prefix, m)
	}

	// Инициализируем Logger
	err := InitLogger(cfg, callback)
	if err != nil {
		log.Fatal(err)
	}

	Info("Старт приложения")
	Error(errors.New("Какая то ошибка"))

	Infocb("Старт приложения 1")
	Errorcb("Какая то ошибка 1")

	Errorfn(app, "Ошибка: %v %s", "Какая то ошибка", "Еще что то")
	Debugfn(app, "Какая то ошибка")

	Errorfncb(app, "Ошибка: %v %s", "Какая то ошибка", "Еще что то")
	Debugfncb(app, "Какая то ошибка")

	Info("Остановка приложения")
}
