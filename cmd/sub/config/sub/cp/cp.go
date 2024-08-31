package cp

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
	flags := flag.NewFlagSet("cp", flag.ExitOnError)

	var inputGlob string
	flags.StringVar(&inputGlob, "input", "", "./in/**/*.yml")
	flags.StringVar(&inputGlob, "i", "", "./in/**/*.yml")

	var outputDir string
	flags.StringVar(&outputDir, "out", "", "./out")
	flags.StringVar(&outputDir, "o", "", "./out")

	var dataFile string
	flags.StringVar(&dataFile, "data", "", "./data.yml")
	flags.StringVar(&dataFile, "d", "", "./data.yml")

	err := flags.Parse(args)
	log.Fatal(err, "error parsing flags")

	basepath, pattern := doublestar.SplitPattern(inputGlob)
	files, err := doublestar.Glob(os.DirFS(basepath), pattern)
	log.Fatal(err, "Error globbing files", "pattern", inputGlob)

	err = os.MkdirAll(outputDir, 0755)
	log.Fatal(err, "Error creating output directory", "dir", outputDir)

	data, err := dasel.NewFromFile(dataFile, "yaml")
	log.Fatal(err, "Error reading data file", "file", dataFile)

	for _, file := range files {
		fullFilePath := filepath.Join(basepath, file)

		content, err := os.ReadFile(fullFilePath)
		log.Error(err, "Error reading file", "file", fullFilePath)
		if err != nil {
			continue
		}

		var buffer bytes.Buffer
		err = internal.ProcessFile(filepath.Base(file), content, data, &buffer)
		log.Error(err, "Error processing file", "file", file)
		if err != nil {
			continue
		}

		relativePath, err := filepath.Rel(basepath, fullFilePath)
		log.Error(err, "Error getting relative path", "file", fullFilePath)
		if err != nil {
			continue
		}

		outputFileName := strings.Replace(relativePath, ".tmpl", "", 1)
		outputPath := filepath.Join(outputDir, outputFileName)

		err = os.MkdirAll(filepath.Dir(outputPath), 0755)
		log.Error(err, "Error creating output directories", "dir", filepath.Dir(outputPath))
		if err != nil {
			continue
		}

		err = os.WriteFile(outputPath, buffer.Bytes(), 0600)
		log.Error(err, "Error writing file", "file", fullFilePath, "path", outputPath)
		if err != nil {
			continue
		}

		log.Info("Processed file", "input", fullFilePath, "output", outputPath)
	}

	log.Info("Finished processing files", "input", inputGlob, "output", outputDir, "len", len(files))
}
