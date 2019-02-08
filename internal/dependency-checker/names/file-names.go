package names

type ImportsExports [2]*Names

func NewImportsExports(imports *Names, exports *Names) *ImportsExports {
	return &ImportsExports{imports, exports}
}

//noinspection GoNameStartsWithPackageName
type NamesByFile map[string]*Names

type FileNames struct {
	Path  Path
	Names *Names
}

func NewFileNames(path Path, names *Names) *FileNames {
	return &FileNames{Path: path, Names: names}
}

type FileAnalysis struct {
	Path             string
	Imports, Exports *Names
}

func NewFileAnalysis(path string, imports *Names, exports *Names) *FileAnalysis {
	return &FileAnalysis{Path: path, Imports: imports, Exports: exports}
}

type Class string
type ClassesByFile map[Path][]Class

//func (C *ClassesByFile) Flip() FilesByClass {
//	F := make(FilesByClass)
//
//	C.fo
//}

type Path string

type FilesByClass map[Class][]Path

//func (F *FilesByClass) Flip() ClassesByFile {
//
//}
