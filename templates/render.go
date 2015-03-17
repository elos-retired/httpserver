package templates

import (
	"net/http"

	"github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

func Show(name TemplateName) httprouter.Handle {
	return Template(name, nil)
}

func Template(name TemplateName, data interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		renderTemplate(transfer.NewHTTPConnection(w, r, nil), name, data)
	}
}

func Render(c *transfer.HTTPConnection, name TemplateName, data interface{}) error {
	return renderTemplate(c, name, data)
}

func renderTemplate(c *transfer.HTTPConnection, name TemplateName, data interface{}) error {
	t, ok := templates[name]

	if !ok {
		return NewNotFoundError(name)
	}

	if err := t.Execute(c.ResponseWriter(), data); err != nil {
		return NewRenderError(err)
	}

	return nil
}
