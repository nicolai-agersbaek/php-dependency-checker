package names

type ImportsExports [2]*Names

func NewImportsExports(imports *Names, exports *Names) *ImportsExports {
	return &ImportsExports{imports, exports}
}

type FileNames map[string]*Names

type FileAnalysis struct {
	Path             string
	Imports, Exports *Names
}

func NewFileAnalysis(path string, imports *Names, exports *Names) *FileAnalysis {
	return &FileAnalysis{Path: path, Imports: imports, Exports: exports}
}
