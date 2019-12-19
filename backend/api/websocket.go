
Websocket API
Websocket Endpoint: /socket

There are 8 channels on the matching engine websocket API:

orders
ohlcv
orderbook
trades
price_board
markets
notification
func NewWSServer(allowedOrigins []string, srv *Server) *http.Server {
	return &http.Server{Handler: srv.WebsocketHandler(allowedOrigins)}
}
func (s *Server) WebsocketHandler(allowedOrigins []string) http.Handler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  wsReadBuffer,
		WriteBufferSize: wsWriteBuffer,
		WriteBufferPool: wsBufferPool,
		CheckOrigin:     wsHandshakeValidator(allowedOrigins),
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Debug("WebSocket upgrade failed", "err", err)
			return
		}
		codec := newWebsocketCodec(conn)
		s.ServeCodec(codec, 0)
	})
}
type webSocketResponseWriter struct {
	writtenHeaders  bool
	wsConn          *websocket.Conn
	headers         http.Header
	flushedHeaders  http.Header
	timeOutInterval time.Duration
	timer           *timer.Timer
}
[func] -ping

type webSocketWrappedReader struct {
	wsConn          *websocket.Conn
	respWriter      *webSocketResponseWriter
	remainingBuffer []byte
	remainingError  error
	cancel          context.CancelFunc
}

[func]-headers
