package process

import (
	"bytes"
	"flag"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/imMohika/gohangyourself/cmd/sub/config/internal"
	"github.com/imMohika/gohangyourself/log"
	"github.com/tomwright/dasel"
	"os"
	"path/filepath"
	"strings"
)

type SubCommand struct {
}

func (s SubCommand) Handle(args []string) {
	flags := flag.NewFlagSet("compile", flag.ExitOnError)

	var pathGlob string
	flags.StringVar(&pathGlob, "path", "", "./plugins/**/*.tmpl.yml")
	flags.StringVar(&pathGlob, "p", "", "./plugins/**/*.tmpl.yml")

	var dataFile string
	flags.StringVar(&dataFile, "data", "", "./data.yml")
	flags.StringVar(&dataFile, "d", "", "./data.yml")

	err := flags.Parse(args)
	log.Error(err, "error parsing flags")

	basepath, pattern := doublestar.SplitPattern(pathGlob)
	files, err := doublestar.Glob(os.DirFS(basepath), pattern)
	log.Error(err, "Error globbing files", "pattern", pathGlob)

	data, err := dasel.NewFromFile(dataFile, "yaml")
	log.Error(err, "Error reading data file", "file", dataFile)

	for _, file := range files {
		content, err := os.ReadFile(file)
		log.Error(err, "Error reading file", "file", file)

		var buffer bytes.Buffer
		err = internal.ProcessFile(file, content, data, &buffer)
		log.Error(err, "Error processing file", "file", file)

		baseName := filepath.Base(file)
		outputFileName := strings.Replace(baseName, ".tmpl", "", 1)
		outputPath := filepath.Join(filepath.Dir(file), outputFileName)

		err = os.WriteFile(outputPath, buffer.Bytes(), 0600)
		log.Error(err, "Error writing file", "file", file, "path", outputPath)

		log.Info("Processed file", "input", file, "output", outputPath)
	}

	log.Info("Finished processing files", "input", pathGlob, "len", len(files))
}
