package algo

type Sampler struct {
	Model, Diff *Model
}

func (s *Sampler) SampleDocument(d *Document, vshard int) error {
	// TODO(y): Implement standard Gibbs sampling algorithm here, so we can check the general framework.
	return nil
}
