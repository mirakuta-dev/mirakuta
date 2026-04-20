package preset

type Preset struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Extends     []string `yaml:"extends"`
	Languages   []string `yaml:"languages"`
	Editors     []Editor `yaml:"editors"`
	Terminals   []string `yaml:"terminals"`
	Git         Git      `yaml:"git"`
	Tools       []Tool   `yaml:"tools"`
}

type Editor struct {
	ID         string   `yaml:"id"`
	Extensions []string `yaml:"extensions"`
}

type Git struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type Tool struct {
	ID  string `yaml:"id"`
	Via string `yaml:"via"`
	Pkg string `yaml:"pkg"`
}
