package internal

import (
	"flag"
	"path"
	"path/filepath"
	"os"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

var (
	outputPath = flag.String("outdir,o", ".", "directory to output generator")
)

func Generate() error {

	flag.Parse()

	d, err := createGeneratorDef()
	if err != nil {
		return err
	}

	var files []gopkg.FileContents
	files, err = tmpl.AppendFileContents(
		files,
		fileInternalFileTemplate(d),
		fileInternalGenerate(d),
		fileMain(d),
		fileMainTest(d),
	)
	if err != nil {
		return err
	}

	return gopkg.LintAndGenerate(files)
}

type generatorDef struct {
	OutputPath string
	Import gopkg.ImportAndAlias
	InternalImport gopkg.ImportAndAlias
	GeneratorMainFileExists bool
}

func createGeneratorDef() (generatorDef, error) {

	var mainExists bool
	if _, err := os.Stat(filepath.Join(*outputPath, "main.go")); err == nil {
		mainExists = true
	}

	importPath, err := gopkg.PackageImportPath(*outputPath)
	if err != nil {
		return generatorDef{}, err
	}

	generatorName := path.Base(importPath)

	return generatorDef{
		OutputPath: *outputPath,
		Import: gopkg.ImportAndAlias{
			Import: importPath,
			Alias: generatorName,
		},
		InternalImport: gopkg.ImportAndAlias{
			Import: path.Join(importPath, "internal"),
			Alias: "internal",
		},
		GeneratorMainFileExists: mainExists,
	}, nil
}
