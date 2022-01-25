package egolog

import (
	"bytes"
	"fmt"
	info "github.com/egovorukhin/egoappinfo"
	"github.com/egovorukhin/egoconf"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Logger struct {
	//fileName string
	buf    bytes.Buffer
	config Config
	*log.Logger
}

type Config struct {
	Console  bool   `yaml:"console" json:"console" xml:"Console"`
	DirPath  string `yaml:"file_path" json:"dir_path" xml:"DirPath"`
	FileName string `yaml:"file_name" json:"file_name" xml:"FileName"`
	//FileSize int    `yaml:"file_size" json:"file_size" xml:"FileSize"`
	Info     Flags     `yaml:"info" json:"info" xml:"Info"`
	Error    Flags     `yaml:"error" json:"error" xml:"Error"`
	Debug    Flags     `yaml:"debug" json:"debug" xml:"Debug"`
	Api      *Api      `yaml:"api,omitempty" json:"api,omitempty" xml:"Api,omitempty"`
	Rotation *Rotation `yaml:"rotation,omitempty" json:"rotation,omitempty" xml:"Rotation,omitempty"`
}

type Rotation struct {
	Size   int    `yaml:"size" json:"size" xml:"Size"`
	Format string `yaml:"format" json:"format" xml:"Format"`
	Path   string `yaml:"path,omitempty" json:"path,omitempty" xml:"Path,omitempty"`
	//Count int `yaml:"count" json:"count" xml:"count"`
}

type Flags string

const (
	INFO  = "Info"
	ERROR = "Error"
	DEBUG = "Debug"
)

var logger *Logger

func InitLogger(configPath string, app *info.Application) error {

	if app == nil {
		// Инициализируем структуру Application
		app = info.New()
	}

	// Если путь не абсолютный, то подставляем путь относительно приложения
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(app.Executable.Dir, configPath)
	}

	cfg := Config{}
	err := egoconf.Load(configPath, &cfg)
	if err != nil {
		return err
	}

	// Устанавливаем размер файла
	/*if cfg.FileSize == 0 {
		cfg.FileSize = 10
	}*/

	// Устанавливаем путь по умолчанию
	if cfg.DirPath == "" {
		cfg.DirPath = app.Executable.Dir
	} else {
		if !filepath.IsAbs(cfg.DirPath) {
			cfg.DirPath = filepath.Join(app.Executable.Dir, cfg.DirPath)
		}
	}

	// Если не указано имя основного файла, то даем имя приложения
	if cfg.FileName == "" {
		cfg.FileName = strings.ToLower(app.Name) //strings.ReplaceAll(filepath.Base(app.Executable.File), filepath.Ext(app.Executable.File), "")
	}

	if cfg.Api != nil {
		cfg.Api.App = app
	}

	logger = &Logger{
		config: cfg,
		Logger: new(log.Logger),
	}

	// Устанавливаем Writer
	if cfg.Console {
		// вывод в консоль
		logger.SetOutput(os.Stdout)
	} else {
		// Буфер для записи в файл
		logger.SetOutput(&logger.buf)
	}

	return nil
}

// Выводим данные
func (l *Logger) print(prefix, filename string /*callDepth int,*/, sending bool, message interface{}, v ...interface{}) {

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
	err := l.Output(3, m)
	if err != nil {
		log.Println(err)
	}

	// Отправка в стороннюю систему
	if l.config.Api != nil && sending {
		go func() {
			resp, err := l.config.Api.send(prefix, m)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Response: %s\n", resp)
		}()
	}

	// Сохранение в файл
	if !l.config.Console {
		err := logger.save(filename)
		if err != nil {
			log.Println(err)
		}
	}
}

func (l *Logger) Flags(prefix string) {

	l.SetPrefix(prefix + ": ")
	s := "3"
	switch prefix {
	case INFO:
		s = string(l.config.Info)
		break
	case ERROR:
		s = string(l.config.Error)
		break
	case DEBUG:
		s = string(l.config.Debug)
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
func Info(sending bool, message interface{}, v ...interface{}) {
	logger.print(INFO, "", sending, message, v...)
}

// InfoFn Используем шаблоны в конфиг файле для каждого из префиксов
func InfoFn(filename string, sending bool, message interface{}, v ...interface{}) {
	logger.print(INFO, filename, sending, message, v...)
}

// Error Префикс
func Error(sending bool, message interface{}, v ...interface{}) {
	logger.print(ERROR, "", sending, message, v...)
}

// ErrorFn Используем шаблоны в конфиг файле для каждого из префиксов
func ErrorFn(filename string, sending bool, message interface{}, v ...interface{}) {
	logger.print(ERROR, filename, sending, message, v...)
}

// Debug Префикс
func Debug(sending bool, message interface{}, v ...interface{}) {
	logger.print(DEBUG, "", sending, message, v...)
}

// DebugFn Используем шаблоны в конфиг файле для каждого из префиксов
func DebugFn(filename string, sending bool, message interface{}, v ...interface{}) {
	logger.print(DEBUG, filename, sending, message, v...)
}
