package egolog

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Сохранение логов в файл
func (l *Logger) save(filename string) error {

	// Возвращаем путь к директории с логами
	if _, err := os.Stat(l.config.DirPath); os.IsNotExist(err) {
		err = os.MkdirAll(l.config.DirPath, 0777)
		if err != nil {
			return err
		}
	}

	if filename == "" {
		filename = l.config.FileName
	}

	// Формируем полный путь к файлу логов
	fullFileName := filepath.Join(l.config.DirPath, filename+".log")

	// Проверяем путь на корректность
	info, err := os.Stat(fullFileName)
	if !os.IsNotExist(err) {

		if info == nil {
			return errors.New("Неверный формат пути!")
		}

		// Проверяем размер файла и удаляем если превышает установленный размер
		if l.config.Rotation != nil /*&& l.config.Rotation.Count > 0*/ && info.Size() > int64(l.config.Rotation.Size)*1024*1024 {
			path := l.config.DirPath
			if l.config.Rotation.Path != "" {
				path = l.config.Rotation.Path
			}
			format := strings.ReplaceAll(l.config.Rotation.Format, "%name", filename)
			format = strings.ReplaceAll(format, "%time", time.Now().Format("2006-01-02T15:04:05"))
			err = os.Rename(fullFileName, filepath.Join(path, format+".log"))
			//err = os.Remove(fullFileName)
			if err != nil {
				return err
			}
		}
	}

	// Используем mutex за нормальную конкуренцию за память
	go l.write(fullFileName)

	return nil
}

// Запись в файл
func (l *Logger) write(path string) {

	// Открываем файл и раздаем права
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}

	// Пишем в файл данные
	_, err = file.Write(l.buf.Bytes())
	// Очистка буфера, чтобы не писать повторяющиеся данные
	defer l.buf.Reset()
	if err != nil {
		log.Println(err)
	}
}
