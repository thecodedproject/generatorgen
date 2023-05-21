package internal

import (
	"path/filepath"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileMain(
	d generatorDef,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {

		imports := tmpl.UnnamedImports(
			"log",
		)
		imports = append(imports, d.InternalImport)

		return []gopkg.FileContents{
			{
				Filepath: filepath.Join(d.OutputPath, "main.go"),
				PackageName: d.Import.Alias,
				PackageImportPath: d.Import.Import,
				Imports: imports,
				Functions: []gopkg.DeclFunc{
					{
						Name: "main",
						BodyTmpl: `
	err := internal.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}
`,
					},
				},
			},
		}, nil
	}
}

