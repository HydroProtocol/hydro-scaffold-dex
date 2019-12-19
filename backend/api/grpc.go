

type grpcServer struct {
	e *backend
}
type grpcClient struct {
	conn     protocol.Exchange_ConnectionServer
	loggedIn bool
	user     string
}
