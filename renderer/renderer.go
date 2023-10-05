package renderer

type Renderer interface {
	RenderFolder(folder string, outputFile *string) error
	RenderPage(pageFile string) (string, error)
}
