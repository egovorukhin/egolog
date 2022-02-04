package egolog

import (
	"encoding/json"
	info "github.com/egovorukhin/egoappinfo"
	"github.com/egovorukhin/egorest"
	"net/http"
	"strings"
	"time"
)

type Api struct {
	Url         string            `yaml:"url" json:"url" xml:"Url"`
	Method      string            `yaml:"method" json:"method" xml:"Method"`
	Timeout     int               `yaml:"timeout" json:"timeout" xml:"Timeout"`
	Template    string            `yaml:"template" json:"template" xml:"Template"`
	App         *info.Application `yaml:"-" json:"-" xml:"-"`
	Body        string            `yaml:"body,omitempty" json:"body,omitempty" xml:"Body,omitempty"`
	ContentType string            `yaml:"content_type" json:"content_type" xml:"ContentType"`
	BasicAuth   *BasicAuth        `yaml:"basic_auth,omitempty" json:"basic_auth,omitempty" xml:"BasicAuth,omitempty"`
	Proxy       string            `yaml:"proxy,omitempty" json:"proxy,omitempty" xml:"Proxy,omitempty"`
	Headers     []egorest.Header  `yaml:"headers" json:"headers" xml:"Headers"`
}

type BasicAuth struct {
	User string `yaml:"user" json:"user" xml:"User"`
	Pass string `yaml:"pass" json:"pass" xml:"Pass"`
}

// Отправляем данные по ошибки в любую систему используя api
func (a *Api) send(prefix, message string) (resp *http.Response, err error) {

	// Инициализация клиента
	client, err := egorest.NewClientByUri(a.Url)
	if err != nil {
		return
	}

	// Установка таймаута
	client.SetTimeout(a.Timeout)
	if a.BasicAuth != nil {
		client.SetBasicAuth(a.BasicAuth.User, a.BasicAuth.Pass)
	}
	// Установка прокси сервера
	if a.Proxy != "" {
		client.SetProxy(a.Proxy)
	}

	// Инициализация запроса
	req := egorest.NewRequest(a.Method, "")
	// Сериализация тела запроса
	if a.Body != "" {
		// Готовим шаблон
		message = strings.Trim(message, "\n")
		template := message
		if a.Template != "" {
			template = strings.ReplaceAll(a.Template, "%prefix", prefix)
			template = strings.ReplaceAll(template, "%name", a.App.Name)
			template = strings.ReplaceAll(template, "%version", a.App.Version.String())
			template = strings.ReplaceAll(template, "%host", a.App.Hostname)
			template = strings.ReplaceAll(template, "%system", a.App.System)
			template = strings.ReplaceAll(template, "%time", time.Now().Format(time.RFC3339))
			template = strings.ReplaceAll(template, "%message", message)
		}

		// Преобразуем json в map
		var body map[string]interface{}
		b := strings.ReplaceAll(a.Body, "%template", template)
		err = json.Unmarshal([]byte(b), &body)
		if err != nil {
			return
		}

		switch egorest.ContentType(a.ContentType) {
		case egorest.XML:
			req = req.Xml(body)
			break
		default:
			req = req.Json(body)
			break
		}
	}
	// Установка дополнительных заголовков
	req.SetHeader(a.Headers...)

	// Отправка запроса
	resp, err = client.Send(req)
	if err != nil {
		return
	}

	return resp, nil
}
