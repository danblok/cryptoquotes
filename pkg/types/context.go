package types

type (
	// RequestID is a type for the request id.
	RequestID string

	// RemoteAddr is a type for the context value of the remote addr key.
	RemoteAddr string

	// TransportType is the type of a transport layer that was used.
	TransportType string
)

const (
	// RequestIDKey is the key for RequestID.
	RequestIDKey RequestID = "request_id"

	// RemoteAddrKey is the key for RemoteAddr.
	RemoteAddrKey RemoteAddr = "remote_addr"

	// TransportTypeKey is the key for Transport type.
	TransportTypeKey TransportType = "transport_type"

	// HTTPTransport is the HTTP transport type.
	HTTPTransport = "HTTP"

	// GRPCTransport is the GRPC transport type.
	GRPCTransport = "GRPC"
)
