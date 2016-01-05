package srvs

// Bootstrap notifies workers and aggregators to load and build
// vocabulary and vsharder.  Aggregators will load existing model if
// there is any, or initialize a new model.
func (wf *Master) Bootstrap(iter int) {
}
