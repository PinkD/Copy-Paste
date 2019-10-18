package cpst

type contentData struct {
	Code    uint64
	Sha     string
	Content string
}

type dataProcess interface {
	ContainsContent(sha, content string) (uint64, error)
	SaveContent(data *contentData) error
	GetContent(code uint64) (string, error)
}
