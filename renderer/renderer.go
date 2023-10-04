package renderer

type Renderer interface {
	RenderFolder(folder string) error
	RenderPage(pageFile string) (string, error)
}
