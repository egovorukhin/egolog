package egolog

type InfoLog map[string]interface{}

const (
	InfoLogPath    = "path"
	InfoLogName    = "name"
	InfoLogPrefix  = "prefix"
	InfoLogMessage = "message"
)

// Get Получение значения по ключу из таблицы
func (info InfoLog) Get(key string) interface{} {
	if value, ok := info[key]; ok {
		return value
	}
	return nil
}

// Path Вернуть путь до файла
func (info InfoLog) Path() string {
	value := info.Get(InfoLogPath)
	if value != nil {
		return value.(string)
	}
	return ""
}

// Name Вернуть имя файла
func (info InfoLog) Name() string {
	value := info.Get(InfoLogName)
	if value != nil {
		return value.(string)
	}
	return ""
}

// Prefix Вернуть префикс лога
func (info InfoLog) Prefix() string {
	value := info.Get(InfoLogPrefix)
	if value != nil {
		return value.(string)
	}
	return ""
}

// Message Вернуть сообщение лога
func (info InfoLog) Message() string {
	value := info.Get(InfoLogMessage)
	if value != nil {
		return value.(string)
	}
	return ""
}
