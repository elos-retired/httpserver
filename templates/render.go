package templates

import (
	"net/http"

	"github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

/*
	Show is for rendering templates that require
	no specific data
*/
func Show(name Name) httprouter.Handle {
	return Template(name, nil)
}

/*
	Template is a httprouter.Handle curried function to inject
	the template name and data

	You can only really use this if the data is constant.
*/
func Template(name Name, data interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := transfer.NewHTTPConnection(w, r, nil)
		CatchError(c, renderTemplate(c, name, data))
	}
}

/*
	Render will render the named template witht he provided data

	note: data can absolutely be nil, but if that is the case, consider
	using Show
*/
func Render(c *transfer.HTTPConnection, name Name, data interface{}) error {
	return renderTemplate(c, name, data)
}

/*
	renderTemplate is the internally used implementation of rendering a named
	template with the supplied data
*/
func renderTemplate(c *transfer.HTTPConnection, name Name, data interface{}) error {
	err := parseHTMLTemplates(templateSets)
	if err != nil {
		return err
	}

	t, ok := templates[name]

	if !ok {
		return NewNotFoundError(name)
	}

	if err := t.Execute(c.ResponseWriter(), data); err != nil {
		return NewRenderError(err)
	}

	return nil
}
