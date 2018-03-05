package filter

type FilterSet struct {
	Word string
	idx  int64
}

type FilterConfig struct {
	replaceWord string
}

type Filter interface {
	Conf(config FilterConfig)
	Build(words []string) error
	Search(text string) ([]FilterSet, error)
	Ban(text string) bool
	Replace() error
}
