package internal

import (
	//"errors"
	"flag"
	"path"

	//"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	//"github.com/thecodedproject/gopkg/tmpl"
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
	files, err = appendFileContents(
		files,
		fileMainTest(d),
	)
	if err != nil {
		return err
	}

	return gopkg.LintAndGenerate(files)
}

func appendFileContents(
	files []gopkg.FileContents,
	fileFuncs ...func() ([]gopkg.FileContents, error),
) ([]gopkg.FileContents, error) {

	for _, fileFunc := range fileFuncs {
		newFiles, err := fileFunc()
		if err != nil {
			return nil, err
		}

		files = append(files, newFiles...)
	}
	return files, nil
}

type generatorDef struct {
	OutputPath string
	Import gopkg.ImportAndAlias
}

func createGeneratorDef() (generatorDef, error) {

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
	}, nil
}
