package picture

type Processor interface {
	Process(image *Index) error
}
