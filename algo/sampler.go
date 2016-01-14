package algo

type Sampler struct {
	Model, Diff *Model
}

func (s *Sampler) SampleDocument(d *Document, vshard int) error {
	return nil
}
