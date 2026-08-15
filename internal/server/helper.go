package server

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	// Извлекаем соответствующий набор шаблонов из кэша в зависимости от названия страницы
	// (например, 'home.page.tmpl'). Если в кэше нет записи запрашиваемого шаблона, то
	// вызывается вспомогательный метод serverError(), который мы создали ранее.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("Шаблон %s не существует!", name))
		return
	}

	// Рендерим файлы шаблона, передавая динамические данные из переменной `td`.
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}
