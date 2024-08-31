package internal

import (
	"bytes"
	"github.com/imMohika/gohangyourself/log"
	"github.com/tomwright/dasel"
	"os"
	"text/template"
)

func ProcessFile(name string, content []byte, data *dasel.Node, buffer *bytes.Buffer) error {
	tmpl, err := template.New(name).Funcs(template.FuncMap{
		"env": os.Getenv,
		"data": func(query string) string {
			got, err := data.Query(query)
			log.Error(err, "Error querying data", "query", query)
			return got.String()
		},
	}).Parse(string(content))

	if err != nil {
		return err
	}

	err = tmpl.Execute(buffer, nil)
	if err != nil {
		return err
	}

	return nil
}
