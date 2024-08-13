package cp

import (
	"bytes"
	"flag"
	"github.com/imMohika/gohangyourself/log"
	"github.com/tomwright/dasel"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type SubCommand struct {
}

func (s SubCommand) Handle(args []string) {
	flags := flag.NewFlagSet("cp", flag.ExitOnError)

	var inputGlob string
	flags.StringVar(&inputGlob, "input", "", "./config/*.yml")
	flags.StringVar(&inputGlob, "i", "", "./config/*.yml")

	var outputDir string
	flags.StringVar(&outputDir, "out", "", "./data")
	flags.StringVar(&outputDir, "o", "", "./data")

	var dataFile string
	flags.StringVar(&dataFile, "data", "", "./data.yml")
	flags.StringVar(&dataFile, "d", "", "./data.yml")

	err := flags.Parse(args)
	log.Error(err, "error parsing flags")

	files, err := filepath.Glob(inputGlob)
	log.Error(err, "Error globbing files", "pattern", inputGlob)

	err = os.MkdirAll(outputDir, 0755)
	log.Error(err, "Error creating output directory", "dir", outputDir)

	data, err := dasel.NewFromFile(dataFile, "yaml")
	log.Error(err, "Error reading data file", "file", dataFile)

	for _, file := range files {
		content, err := os.ReadFile(file)
		log.Error(err, "Error reading file", "file", file)

		tmpl, err := template.New(filepath.Base(file)).Funcs(template.FuncMap{
			"env": os.Getenv,
			"data": func(query string) string {
				got, err := data.Query(query)
				log.Error(err, "Error querying data", "query", query)
				return got.String()
			},
		}).Parse(string(content))
		log.Error(err, "Error parsing template", "file", file)

		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, nil)
		log.Error(err, "Error executing template", "file", file)

		baseName := filepath.Base(file)
		outputFileName := strings.Replace(baseName, ".tmpl", "", 1)
		outputPath := filepath.Join(outputDir, outputFileName)

		err = os.WriteFile(outputPath, buffer.Bytes(), 0600)
		log.Error(err, "Error writing file", "file", file, "path", outputPath)

		log.Info("Processed file", "input", file, "output", outputPath)
	}

	log.Info("Finished processing files", "input", inputGlob, "output", outputDir, "len", len(files))
}
