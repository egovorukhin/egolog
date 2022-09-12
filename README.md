# EgoLog
Логирование важная часть приложения. Данный модуль позваоляет сохранять логи в именованых файлах и позволяет реализовать ```callback``` функцию, которую можно использовать для отправки например в elastic. Логер можно настроить указав размер для перезаписи файла, а также ротацию переполневшегося файла.
### Установка
```
 go get github.com/egovorukhin/egolog
```
### Описание

* **Конфигурация**
```golang
  type Config struct {
    DirPath   string // директория с логами, по умолчанию для Windows корень приложения, для Linux /var/log/{app_name}
    FileName  string // имя основного файла, по умолчанию - имя приложения
    CallDepth int // глубина вложености вызываемых фукций
    Info      Flags // флаг уровня логирования
    Error     Flags // флаг уровня логирования
    Debug     Flags // флаг уровня логирования
    Rotation  *Rotation //ротация файлов
  }

  type Rotation struct {
    Size   int // критический размер файла
    Format string // формат записи при переносе файла
    Path   string // путь куда будет перемещён файла
  }
```

* **Инициализация**
```golang
  import "github.com/egovorukhin/egolog"
  ...
  cfg := egolog.Config{
		DirPath:  config.DirPath,
		FileName: "app",
		Info:     egolog.Flags(config.Info),
		Error:    egolog.Flags(config.Error),
		Debug:    egolog.Flags(config.Debug),
		Rotation: config.Rotation,
	}
	return egolog.InitLogger(cfg, callback)
  ...
  func callback(info egolog.InfoLog) {
	  ...
  }
  
```
* **Уровень логирования**

Для уровня логирования можно использовать переменные в конфигурации, так как под капотом используется стандартный модуль golang [log](https://pkg.go.dev/log#pkg-constants), то и указывать нужно соответсвенно.
```yaml
  info: 1 | 4
  error: 1 | 4 | 8
  debug: 1 | 4 | 8
```
Либо указываем в самом коде, строкой.
```golang
egolog.Config{
		DirPath:  config.DirPath,
		FileName: "app",
		Info:     "1 | 4",
		Error:    "1 | 4 | 8",
		Debug:    "1 | 4 | 8",
		Rotation: config.Rotation,
	}
```
* **Использование**
```golang
  // Пишем лог в один файл
  egolog.Info("Старт приложения")
  egolog.Error("Ошибка")
  egolog.Debug("Ошибка развернутая")
  
  // Пишем лог в один файл и вызываем обработчик callback
  egolog.Infocb("Старт приложения")
  egolog.Errorcb("Ошибка")
  egolog.Debugcb("Ошибка развернутая")
  
  // Пишем лог в именованый файл, пример: app.log
  egolog.Infofn("app", "Старт приложения")
  egolog.Errorfn("app", "Ошибка")
  egolog.Debugfn("app", "Ошибка развернутая")
  
  // Пишем лог в именованый файл, пример: app.log и вызываем обработчик callback
  egolog.Infofncb("app", "Старт приложения")
  egolog.Errorfncb("app", "Ошибка")
  egolog.Debugfncb("app", "Ошибка развернутая")
```
* **Обработчик Callback**

Для того чтобы сработал обработчик его нужно реализовать, переопределить. Функция передает структуру ```InfoLog```. ```InfoLog - map[string]interface{}```, которая имеет несколько значений и имеет данные о записи лога. Возможно что со временем будут добавляться новые константы.
```golang
  func callback(info egolog.InfoLog) {
	  ...
  }
```
```
  const (
    InfoLogPath    = "path"
    InfoLogName    = "name"
    InfoLogPrefix  = "prefix"
    InfoLogMessage = "message"
  )
```
### Пример
```golang
const app = "app"

func main() {

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

  callback := func(infoLog InfoLog) {
    fmt.Printf("path: %s, name: %s, prefix: %s, message: %s", infoLog["path"], infoLog["name"], infoLog.Prefix(), infoLog[InfoLogMessage])
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
```
