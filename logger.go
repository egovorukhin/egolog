package egolog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Logger struct {
	buf      bytes.Buffer
	callback Handler
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

type Handler func(prefix, message string)

type Flags string

const (
	INFO  = "Info"
	ERROR = "Error"
	DEBUG = "Debug"
)

var logger *Logger

func InitLogger(cfg Config, callback ...Handler) error {

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

// Выводим данные
func (l *Logger) print(prefix, filename string, isHandler bool, message interface{}, v ...interface{}) {

	// Если сообщение нет, то выходим
	if message == nil {
		return
	}

	if l == nil {
		log.Println("Необходимо инициализировать структуру Logger.  Функция InitLogger()")
		return
	}

	// Устанавливаем флаг
	logger.Flags(prefix)

	m := ""
	// Если в массиве v есть элементы, то message используем как формат
	if v != nil {
		if reflect.ValueOf(message).Kind() == reflect.String {
			// CallDepth - глубина стека, количество кадров стека для вызывающего файла.
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

	// Вывод в консоль
	log.Printf("%s: %s", prefix, m)
	// Сохранение в файл
	err = logger.save(filename)
	if err != nil {
		log.Println(err)
	}

	// Выполнение обработчика
	if isHandler && l.callback != nil {
		go l.callback(prefix, m)
	}
}

func (l *Logger) Flags(prefix string) {

	l.SetPrefix(prefix + ": ")
	s := "3"
	switch prefix {
	case INFO:
		s = string(l.Info)
		break
	case ERROR:
		s = string(l.Error)
		break
	case DEBUG:
		s = string(l.Debug)
		break
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
