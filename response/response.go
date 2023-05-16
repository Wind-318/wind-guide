package response

const (
	// MaxNumber is 1061109567.
	MaxNumber = 0x3f3f3f3f
)

const (
	RequestID      = "request_id"
	ServiceID      = "service_id"
	ServiceAddr    = "service_addr"
	ServicePort    = "service_port"
	ServiceName    = "service_name"
	ServiceVersion = "service_version"
	HealthCheckUrl = "health_check_url"
)

// Response message.
const (
	SuccessMessage             = "success"
	BadRequestMessage          = "bad request"
	ServiceIDEixstsMessage     = "service id exists"
	NotFoundMessage            = "not found"
	InternalServerErrorMessage = "internal server error"
	MarshalFailedMessage       = "marshal failed"
	UnMarshalFailedMessage     = "unmarshal failed"
)

// Response code.
const (
	SuccessCode             = "2000"
	BadRequestCode          = "4000"
	ServiceIDEixstsCode     = "4001"
	NotFoundCode            = "4004"
	InternalServerErrorCode = "5000"
	MarshalFailedCode       = "5001"
	UnMarshalFailedCode     = "5002"
)
