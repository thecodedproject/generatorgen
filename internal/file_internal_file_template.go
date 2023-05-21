package internal

import (
	"path/filepath"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileInternalFileTemplate(
	d generatorDef,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {

		imports := tmpl.UnnamedImports(
			"path/filepath",
			"github.com/thecodedproject/gopkg",
			"github.com/thecodedproject/gopkg/tmpl",
		)

		return []gopkg.FileContents{
			{
				Filepath: filepath.Join(d.OutputPath, "internal", "file_template.go"),
				PackageName: d.InternalImport.Alias,
				PackageImportPath: d.InternalImport.Import,
				Imports: imports,
				Functions: []gopkg.DeclFunc{
					{
						Name: "fileTemplate",
						Args: []gopkg.DeclVar{
							{
								Name: "d",
								Type: gopkg.TypeNamed{
									Name: "pkgDef",
								},
							},
						},
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeFunc{
								ReturnArgs: tmpl.UnnamedReturnArgs(
									gopkg.TypeArray{
										ValueType: gopkg.TypeNamed{
											Name: "FileContents",
											Import: "github.com/thecodedproject/gopkg",
										},
									},
									gopkg.TypeError{},
								),
							},
						),
						BodyTmpl: `
	return func() ([]gopkg.FileContents, error) {

		// REMOVE THIS RETURN TO ACTUALLY GENERATE FILES
		return nil, nil

		imports := tmpl.UnnamedImports(
			// ... Add imports ...
		)

		return []gopkg.FileContents{
			{
				Filepath: filepath.Join(d.OutputPath, "some", "path.go"),
				PackageName: d.Import.Alias,
				PackageImportPath: d.Import.Import,
				Imports: imports,
				Functions: []gopkg.DeclFunc{
					{
						Name: "HelloWorld",
					},
				},
			},
		}, nil
	}
`,
					},
				},
			},
		}, nil
	}
}

