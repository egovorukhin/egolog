package egolog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Сохранение логов в файл
func (l *Logger) save(filename string) error {

	// Возвращаем путь к директории с логами
	if _, err := os.Stat(l.DirPath); os.IsNotExist(err) {
		err = os.MkdirAll(l.DirPath, 0777)
		if err != nil {
			return err
		}
	}

	if filename == "" {
		filename = l.FileName
	}

	// Формируем полный путь к файлу логов
	fullPath := filepath.Join(l.DirPath, filename+".log")

	// Проверяем путь на корректность
	info, err := os.Stat(fullPath)
	if !os.IsNotExist(err) {

		if info == nil {
			return errors.New("Неверный формат пути!")
		}

		// Проверяем размер файла и удаляем если превышает установленный размер
		if l.Rotation != nil && info.Size() > int64(l.Rotation.Size)*1024 {
			path := l.DirPath
			if l.Rotation.Path != "" {
				path = l.Rotation.Path
			}
			format := strings.ReplaceAll(l.Rotation.Format, "%name", filename)
			format = strings.ReplaceAll(format, "%time", time.Now().Format("2006-01-02T15:04:05"))
			err = os.Rename(fullPath, filepath.Join(path, format+".log"))
			//err = os.Remove(fullPath)
			if err != nil {
				return err
			}
		}
	}

	l.write(fullPath)

	return nil
}

// Запись в файл
func (l *Logger) write(path string) {

	fmt.Println(path)
	// Открываем файл и раздаем права
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	fmt.Println(file.Name())

	// Пишем в файл данные
	_, err = file.Write(l.buf.Bytes())
	// Очистка буфера, чтобы не писать повторяющиеся данные
	defer l.buf.Reset()
	if err != nil {
		log.Println(err)
	}
}
