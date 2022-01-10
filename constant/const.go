package constant

const (
	ServiceName          = `bar`
	Host                 = `localhost`
	GrpcPort             = `9081`
	HttpPort             = `8081`
	CircuitBreakerTimout = 1000 * 30 // change the second operand, this means 30 second timeout
	ZipkinHostPort       = `:0`
)
