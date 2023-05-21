package internal

import (
	"path/filepath"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileSome(d pkgDef) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {

		imports := tmpl.UnnamedImports(
			"my_import",
		)

		return []gopkg.FileContents{
			{
				Filepath: filepath.Join(d.OutputPath, "some.go"),
				PackageName: d.Import.Alias,
				PackageImportPath: d.Import.Import,
				Imports: imports,
			},
		}, nil
	}
}

