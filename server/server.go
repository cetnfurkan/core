package server

// Server is an interface for server implementations
// that can be started and stopped gracefully using
// the Start and Stop methods.
type Server interface {

	// Start starts the server and does not block the
	// calling goroutine. It panics if the server fails
	// to start.
	Start()

	// Stop stops the server gracefully and does not
	// block the calling goroutine. It returns an error
	// if the server fails to stop.
	//
	// The server should stop all its operations,
	// release all resources before returning from this method and
	// should not accept any new requests after this method is called.
	Stop() error
}
