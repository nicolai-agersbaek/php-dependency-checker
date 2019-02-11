package names

type ImportsExports [2]*Names

func NewImportsExports(imports *Names, exports *Names) *ImportsExports {
	return &ImportsExports{imports, exports}
}

//noinspection GoNameStartsWithPackageName
type NamesByFile map[string]*Names

func ConvertToNames(F NamesByFile) *Names {
	// FIXME: Missing tests!
	N := NewNames()

	for _, nn := range F {
		N = N.Merge(nn)
	}

	N.Clean()

	return N
}

//noinspection GoNameStartsWithPackageName
type NamesByFileData struct {
	m NamesByFile
}

func NewNamesByFileData() *NamesByFileData {
	return &NamesByFileData{m: make(NamesByFile)}
}

func (n *NamesByFileData) Put(file string, names *Names) {
	n.m[file] = names
}

func (n *NamesByFileData) Data() NamesByFile {
	return n.m
}

func (n *NamesByFileData) Merge(m *NamesByFileData) {
	if m != nil {
		for k, v := range m.Data() {
			n.Put(k, v)
		}
	}
}

type FileAnalysis struct {
	Path             string
	Imports, Exports *Names
	Error            error
}

func NewFileAnalysisExp(path string, imports *Names, exports *Names, error error) *FileAnalysis {
	return &FileAnalysis{Path: path, Imports: imports, Exports: exports, Error: error}
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
