package token

type File struct {
	name string
}

func NewFile(name string) *File {
	return &File{
		name: name,
	}
}

func (f *File) Name() string {
	return f.name
}
