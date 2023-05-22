package internal

import (
	"path/filepath"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileMainTest(
	d generatorDef,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {


		imports := tmpl.UnnamedImports(
			"os",
			"os/exec",
			"path/filepath",
			"regexp",
			"sort",
			"io",
			"github.com/stretchr/testify/require",
		)
		imports = append(imports, gopkg.ImportAndAlias{
			Import: "github.com/sebdah/goldie/v2",
			Alias: "goldie",
		})

		return []gopkg.FileContents{
			{
				Filepath: filepath.Join(d.OutputPath, "main_test.go"),
				PackageName: "main_test",
				PackageImportPath: d.Import.Import,

				Imports: imports,

				Functions: mainTestFuncs,
			},
		}, nil
	}
}

var mainTestFuncs = []gopkg.DeclFunc{
	{
		Name: "TestGenerate",
		Args: []gopkg.DeclVar{
			{
				Name: "t",
				Type: gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: "T",
						Import: "testing",
					},
				},
			},
		},
		BodyTmpl: `
	testDirs := listDirectoriesInPath(t, "./examples")

	for _, testDir := range testDirs {
		t.Run(testDir, func(t *testing.T) {

			generatedFiles := runGenerateAndGetGeneratedFileBuffers(t, testDir)

			runGoTestAndCheckOutput(t, testDir)

			checkGeneratedFilePaths(t, generatedFiles)

			checkGeneratedFileBuffers(t, testDir, generatedFiles)

			removeGeneratedFiles(t, generatedFiles)
		})
	}
`,
	},
	{
		Name: "listDirectoriesInPath",
		Args: []gopkg.DeclVar{
			{
				Name: "t",
				Type: gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: "T",
						Import: "testing",
					},
				},
			},
			{
				Name: "path",
				Type: gopkg.TypeString{},
			},
		},
		ReturnArgs: tmpl.UnnamedReturnArgs(
			gopkg.TypeArray{
				ValueType: gopkg.TypeString{},
			},
		),
		BodyTmpl: `
	entries, err := os.ReadDir(path)
	require.NoError(t, err)

	dirs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join(path, e.Name()))
		}
	}

	return dirs
`,
	},
	{
		Name: "runGenerateAndGetGeneratedFileBuffers",
		Args: []gopkg.DeclVar{
			{
				Name: "t",
				Type: gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: "T",
						Import: "testing",
					},
				},
			},
			{
				Name: "testDir",
				Type: gopkg.TypeString{},
			},
		},
		ReturnArgs: tmpl.UnnamedReturnArgs(
			gopkg.TypeMap{
				KeyType: gopkg.TypeString{},
				ValueType: gopkg.TypeArray{
					ValueType: gopkg.TypeByte{},
				},
			},
		),
		BodyTmpl: `
	initalFiles, err := listFilesRecursively(testDir)
	require.NoError(t, err)

	cmd := exec.Command("go", "generate", "./" + testDir)

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)

	output, err := io.ReadAll(stdout)
	require.NoError(t, err)

	errOutput, err := io.ReadAll(stderr)
	require.NoError(t, err)

	err = cmd.Wait()
	require.NoError(t, err, string(output) + string(errOutput))

	postGenFiles, err := listFilesRecursively(testDir)
	require.NoError(t, err)

	generatedFiles := make(map[string][]byte, 0)
	for f := range postGenFiles {
		if !initalFiles[f] {
			fileBuffer, err := os.ReadFile(f)
			require.NoError(t, err)

			generatedFiles[f] = fileBuffer
		}
	}

	return generatedFiles
`,
	},
	{
		Name: "runGoTestAndCheckOutput",
		Args: []gopkg.DeclVar{
			{
				Name: "t",
				Type: gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: "T",
						Import: "testing",
					},
				},
			},
			{
				Name: "testDir",
				Type: gopkg.TypeString{},
			},
		},
		BodyTmpl: `
	cmd := exec.Command("go", "test", "-v", "-count=1", "./" + testDir + "/...")

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)

	testOutput, err := io.ReadAll(stdout)
	require.NoError(t, err)

	errOutput, err := io.ReadAll(stderr)
	require.NoError(t, err)

	_ = cmd.Wait()
	testOutput = []byte(string(testOutput) + string(errOutput))

	timeRegex1, err := regexp.Compile(` + "`" + `\([0-9]\.[0-9][0-9]s\)` + "`" + `)
	require.NoError(t, err)
	timeRegex2, err := regexp.Compile(` + "`" + `[0-9]\.[0-9][0-9][0-9]s` + "`" + `)
	require.NoError(t, err)

	testOutput = timeRegex1.ReplaceAll(testOutput, []byte("(X.XXs)"))
	testOutput = timeRegex2.ReplaceAll(testOutput, []byte("X.XXXs"))

	t.Run("go_test", func(t *testing.T) {
		g := goldie.New(t)
		g.Assert(t, t.Name(), []byte(testOutput))
	})
`,
	},
	{
		Name: "checkGeneratedFilePaths",
		Args: []gopkg.DeclVar{
			{
				Name: "t",
				Type: gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: "T",
						Import: "testing",
					},
				},
			},
			{
				Name: "generatedFiles",
				Type: gopkg.TypeMap{
					KeyType: gopkg.TypeString{},
					ValueType: gopkg.TypeArray{
						ValueType: gopkg.TypeByte{},
					},
				},
			},
		},
		BodyTmpl: `
	var generatedFilesPaths []string
	for path := range generatedFiles {
		generatedFilesPaths = append(generatedFilesPaths, path)
	}

	sort.Slice(generatedFilesPaths, func(i, j int) bool {
		return generatedFilesPaths[i] < generatedFilesPaths[j]
	})

	var generatedFilesBuffer string
	for _, f := range generatedFilesPaths {
		generatedFilesBuffer += f + "\n"
	}

	t.Run("generated_file_paths", func(t *testing.T) {
			g := goldie.New(t)
			g.Assert(t, t.Name(), []byte(generatedFilesBuffer))
	})
`,
	},
	{
		Name: "checkGeneratedFileBuffers",
		Args: []gopkg.DeclVar{
			{
				Name: "t",
				Type: gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: "T",
						Import: "testing",
					},
				},
			},
			{
				Name: "testDir",
				Type: gopkg.TypeString{},
			},
			{
				Name: "generatedFiles",
				Type: gopkg.TypeMap{
					KeyType: gopkg.TypeString{},
					ValueType: gopkg.TypeArray{
						ValueType: gopkg.TypeByte{},
					},
				},
			},
		},
		BodyTmpl: `
	for filePath, fileBuffer := range generatedFiles {

		testName, err := filepath.Rel(testDir, filePath)
		require.NoError(t, err)

		t.Run(testName, func(t *testing.T) {
			g := goldie.New(t)
			g.Assert(t, t.Name(), fileBuffer)
		})
	}
`,
	},
	{
		Name: "removeGeneratedFiles",
		Args: []gopkg.DeclVar{
			{
				Name: "t",
				Type: gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: "T",
						Import: "testing",
					},
				},
			},
			{
				Name: "generatedFiles",
				Type: gopkg.TypeMap{
					KeyType: gopkg.TypeString{},
					ValueType: gopkg.TypeArray{
						ValueType: gopkg.TypeByte{},
					},
				},
			},
		},
		BodyTmpl: `
	for path := range generatedFiles {
		err := os.Remove(path)
		require.NoError(t, err)
	}
`,
	},
	{
		Name: "listFilesRecursively",
		Args: []gopkg.DeclVar{
			{
				Name: "path",
				Type: gopkg.TypeString{},
			},
		},
		ReturnArgs: tmpl.UnnamedReturnArgs(
			gopkg.TypeMap{
				KeyType: gopkg.TypeString{},
				ValueType: gopkg.TypeBool{},
			},
			gopkg.TypeError{},
		),
		BodyTmpl: `
	files := make(map[string]bool)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			files[path] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
`,
	},

}
