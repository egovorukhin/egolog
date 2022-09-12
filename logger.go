package egolog

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Logger struct {
	buf bytes.Buffer
	//callback Handler
	callback Callback
	FullPath string
	Config
	*log.Logger
}

type Config struct {
	DirPath   string
	FileName  string
	CallDepth int
	Info      Flags
	Error     Flags
	Debug     Flags
	Rotation  *Rotation
}

type Rotation struct {
	Size   int
	Format string
	Path   string
}

type InfoLog map[string]interface{}

type Callback func(InfoLog)

type Flags string

const (
	INFO  = "Info"
	ERROR = "Error"
	DEBUG = "Debug"
)

var logger *Logger

func InitLogger(cfg Config, callback ...Callback) error {

	path, err := os.Executable()
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(path)
	fileName := strings.ReplaceAll(filepath.Base(path), " ", "_")

	// Устанавливаем путь по умолчанию
	if cfg.DirPath == "" {
		switch strings.ToLower(runtime.GOOS) {
		case "windows":
			cfg.DirPath = filepath.Join(dirPath, "logs")
			break
		case "linux":
			cfg.DirPath = filepath.Join("/var/log/", fileName)
		}
	} else {
		if !filepath.IsAbs(cfg.DirPath) {
			cfg.DirPath = filepath.Join(dirPath, cfg.DirPath)
		}
	}

	// Если не указано имя основного файла, то даем имя приложения
	if cfg.FileName == "" {
		cfg.FileName = strings.ToLower(fileName)
	}

	if cfg.CallDepth == 0 {
		cfg.CallDepth = 3
	}

	logger = &Logger{
		Config: cfg,
		Logger: new(log.Logger),
	}

	if callback != nil {
		logger.callback = callback[0]
	}

	// Устанавливаем Writer
	logger.SetOutput(&logger.buf)

	return nil
}

// SetCallback Устанавливаем обработчик для обратной функции, начиная с версии v0.2.1
func SetCallback(callback Callback) {
	logger.callback = callback
}

func (l *Logger) createPathAndRotation(filename string) error {

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
	l.FullPath = filepath.Join(l.DirPath, filename+".log")

	// Проверяем путь на корректность
	info, err := os.Stat(l.FullPath)
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
			err = os.Rename(l.FullPath, filepath.Join(path, format+".log"))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *Logger) Write(data []byte) (int, error) {

	// Открываем файл и раздаем права
	file, err := os.OpenFile(l.FullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Пишем в файл данные
	return file.Write(data)
}

// Выводим данные
func (l *Logger) print(prefix, filename string, isHandler bool, message interface{}, v ...interface{}) {

	// Если сообщение нет, то выходим
	if message == nil {
		return
	}

	if l == nil {
		log.Println("Необходимо инициализировать структуру Logger. Функция InitLogger()")
		return
	}

	// Устанавливаем флаг
	logger.Flags(prefix)

	m := ""
	// Если в массиве v есть элементы, то message используем как формат
	if v != nil {
		if reflect.ValueOf(message).Kind() == reflect.String {
			m = fmt.Sprintf(message.(string), v...)
		}
	} else {
		m = fmt.Sprintln(message)
	}

	// CallDepth - глубина стека, количество кадров стека для вызывающего файла.
	err := l.Output(l.CallDepth, m)
	if err != nil {
		log.Println(err)
	}
	// Устанавливаем путь для файла сохранения
	err = logger.createPathAndRotation(filename)
	if err != nil {
		log.Println(err)
	}

	// Вывод в консоль
	log.Printf("%s: %s", prefix, m)

	// Выполнение обработчика
	if isHandler && l.callback != nil {
		infoLog := InfoLog{
			"name":    filename,
			"path":    l.FullPath,
			"prefix":  prefix,
			"message": m,
		}
		go l.callback(infoLog)
	}
}

func (l *Logger) Flags(prefix string) {

	l.SetPrefix(prefix + ": ")
	s := "3"
	switch prefix {
	case INFO:
		s = string(l.Info)
	case ERROR:
		s = string(l.Error)
	case DEBUG:
		s = string(l.Debug)
	}

	flags := strings.Split(s, "|")

	f := 0
	for _, flag := range flags {
		i, err := strconv.Atoi(strings.Trim(flag, " "))
		if err != nil {
			continue
		}
		f = f | i
	}
	l.SetFlags(f)
}

// Info Используем шаблоны в конфиг файле для каждого из префиксов
func Info(message interface{}, v ...interface{}) {
	logger.print(INFO, "", false, message, v...)
}

// Error Префикс
func Error(message interface{}, v ...interface{}) {
	logger.print(ERROR, "", false, message, v...)
}

// Debug Префикс
func Debug(message interface{}, v ...interface{}) {
	logger.print(DEBUG, "", false, message, v...)
}

// Infofn Сохранение в файл с указанием имени файла
func Infofn(filename string, message interface{}, v ...interface{}) {
	logger.print(INFO, filename, false, message, v...)
}

// Errorfn Используем шаблоны в конфиг файле для каждого из префиксов
func Errorfn(filename string, message interface{}, v ...interface{}) {
	logger.print(ERROR, filename, false, message, v...)
}

// Debugfn Используем шаблоны в конфиг файле для каждого из префиксов
func Debugfn(filename string, message interface{}, v ...interface{}) {
	logger.print(DEBUG, filename, false, message, v...)
}

// Infocb вызов обработчика
func Infocb(message interface{}, v ...interface{}) {
	logger.print(INFO, "", true, message, v...)
}

// Errorcb вызов обработчика
func Errorcb(message interface{}, v ...interface{}) {
	logger.print(ERROR, "", true, message, v...)
}

// Debugcb вызов обработчика
func Debugcb(message interface{}, v ...interface{}) {
	logger.print(DEBUG, "", true, message, v...)
}

// Infofncb вызов обработчика
func Infofncb(filename string, message interface{}, v ...interface{}) {
	logger.print(INFO, filename, true, message, v...)
}

// Errorfncb вызов обработчика
func Errorfncb(filename string, message interface{}, v ...interface{}) {
	logger.print(ERROR, filename, true, message, v...)
}

// Debugfncb вызов обработчика
func Debugfncb(filename string, message interface{}, v ...interface{}) {
	logger.print(DEBUG, filename, true, message, v...)
}
