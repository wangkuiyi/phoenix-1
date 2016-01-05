package srvs

// Bootstrap notifies workers and aggregators to load and build
// vocabulary and vsharder.  Aggregators will load existing model if
// there is any, or initialize a new model.
func (m *Master) Bootstrap(iter int) {
	parallel.For(0, len(m.workers), 1, func(i int) {
		var dumb int
		if e := m.workers[i].Call("Worker.Bootstrap", iter, &dumb); e != nil {
			log.Panicf("Worker(%s).Bootstrap(%v) failed: %v", m.workers[i].Addr, iter, e)
		}
	})

	parallel.For(0, len(m.
}

func (w *Worker) Bootstrap(iter int, _ *int) error {

}
