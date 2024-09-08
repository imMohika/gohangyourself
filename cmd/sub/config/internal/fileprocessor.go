package internal

import (
	"bytes"
	"github.com/imMohika/gohangyourself/log"
	"github.com/tomwright/dasel"
	"os"
	"strings"
	"text/template"
)

func ProcessFile(name string, content []byte, data *dasel.Node, buffer *bytes.Buffer) error {
	tmpl, err := template.New(name).Funcs(template.FuncMap{
		"env": os.Getenv,
		"data": func(query string) string {
			got, err := data.Query(query)
			log.Error(err, "Error querying data", "query", query)
			log.Info("kind", got.Value.Type())
			gotType, err := got.Query(".[@]")
			if err != nil {
				log.Warn("Error getting node type, returning node as string", "node", got, "err", err)
				return got.String()
			}

			if gotType.String() == "array" {
				vals, err := got.QueryMultiple(".[*]")
				if err != nil {
					log.Warn("Node type is array but couldn't get it's values, returning node as string", "node", got, "err", err)
					return got.String()
				}
				res := "["
				for _, val := range vals {
					res += val.String() + ", "
				}
				res = strings.TrimSuffix(res, ", ") + "]"
				return res
			}

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
