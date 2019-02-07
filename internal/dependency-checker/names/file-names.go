package names

type ImportsExports [2]*Names

func NewImportsExports(imports *Names, exports *Names) *ImportsExports {
	return &ImportsExports{imports, exports}
}

type FileNames map[string]*Names
