package internal

import (
	"path/filepath"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileInternalGenerate(
	d generatorDef,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {

		imports := tmpl.UnnamedImports(
			"flag",
			"path",
			"github.com/thecodedproject/gopkg/tmpl",
		)

		return []gopkg.FileContents{
			{
				Filepath: filepath.Join(d.OutputPath, "internal", "generate.go"),
				PackageName: d.InternalImport.Alias,
				PackageImportPath: d.InternalImport.Import,
				Imports: imports,
				Vars: []gopkg.DeclVar{
					{
						Name: "outputPath",
						Type: gopkg.TypeUnnamedLiteral{},
						LiteralValue: `flag.String("outdir,o", ".", "output directory for generated files")`,
					},
				},
				Types: []gopkg.DeclType{
					{
						Name: "pkgDef",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{
									Name: "OutputPath",
									Type: gopkg.TypeString{},
								},
								{
									Name: "Import",
									Type: gopkg.TypeNamed{
										Name: "ImportAndAlias",
										Import: "github.com/thecodedproject/gopkg",
									},
								},
							},
						},
					},
				},
				Functions: internalGenerateFuncs,
			},
		}, nil
	}
}

var internalGenerateFuncs = []gopkg.DeclFunc{
	{
		Name: "Generate",
		ReturnArgs: tmpl.UnnamedReturnArgs(
			gopkg.TypeError{},
		),
		BodyTmpl: `
	flag.Parse()

	d, err := createPkgDef()
	if err != nil {
		return err
	}

	var files []gopkg.FileContents
	files, err = tmpl.AppendFileContents(
		files,
		// Add file creators here...
		fileTemplate(d),
	)
	if err != nil {
		return err
	}

	return gopkg.LintAndGenerate(files)
`,
	},
	{
		Name: "createPkgDef",
		ReturnArgs: tmpl.UnnamedReturnArgs(
			gopkg.TypeNamed{
				Name: "pkgDef",
			},
			gopkg.TypeError{},
		),
		BodyTmpl: `
	importPath, err := gopkg.PackageImportPath(*outputPath)
	if err != nil {
		return pkgDef{}, err
	}

	pkgName := path.Base(importPath)

	return pkgDef{
		OutputPath: *outputPath,
		Import: gopkg.ImportAndAlias{
			Import: importPath,
			Alias: pkgName,
		},
	}, nil
`,
	},
}
