package egolog

import (
	"errors"
	"fmt"
	"testing"
)

const app = "app"

func Test(t *testing.T) {

	cfg := Config{
		Escaped: true,
		DirPath: "logs",
		Info:    "3",
		Error:   "3 | 16",
		Debug:   "1 | 4 | 8",
		Rotation: &Rotation{
			Size:   10240,
			Format: "%name_%time",
		},
	}

	callback := func(infoLog InfoLog) {
		fmt.Printf("path: %s, name: %s, prefix: %s, message: %s", infoLog["path"], infoLog["name"], infoLog["prefix"], infoLog["message"])
	}

	// Инициализируем Logger
	err := InitLogger(cfg, callback)
	if err != nil {
		t.Fatal(err)
	}

	s := `-Привет мир!
-Как дела?

	-Хорошо`

	Info("Старт приложения")
	Error(errors.New(s))

	Infocb("Старт приложения 1")
	Errorcb("Какая то ошибка 1")

	Errorfn(app, "Ошибка: %v %s", "Какая то ошибка", "Еще что то")
	Debugfn(app, "Какая то ошибка")

	Errorfncb(app, "Ошибка: %v %s", "Какая то ошибка", "Еще что то")
	Debugfncb(app, "Какая то ошибка")

	Info("Остановка приложения")
}
