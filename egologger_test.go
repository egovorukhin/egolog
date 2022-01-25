package egolog

import (
	"errors"
	"log"
	"testing"
)

type Main struct{}

func (m Main) GetError(text string) error {
	return errors.New(text)
}

func Test(t *testing.T) {

	// Инициализируем Logger
	err := InitLogger("log.yml")
	if err != nil {
		log.Fatal(err)
	}

	Info("Старт приложения")

	m := Main{}
	Error(m.GetError("Какая то ошибка"))
	Error("Ошибка: %s", m.GetError("Какая то ошибка"))
	Debug(m.GetError("Какая то ошибка"))

	Info("Остановка приложения")

}
