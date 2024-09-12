package aigc

type ModelId string

type ModelOptions struct {
	VendorId    VendorId
	Endpoint    string
	Region      string
	ApiKey      string
	AccessKey   string
	SecretKey   string
	ApiVersion  string
	Proxy       string
	Retries     int
	RequestLog  func([]byte)
	ResponseLog func([]byte)
}

/*
type ModelCredentials struct {
	ApiKey    string
	AccessKey string
	SecretKey string
}

type ModelClient interface {
	InvokeModel(options *ModelOptions, modelId ModelId, request []byte) ([]byte, error)
}

type httpModelClient struct {
}
*/

type ModelOptionFunc func(*ModelOptions)

func WithVendor(vendorId VendorId) func(options *ModelOptions) {
	return func(o *ModelOptions) {
		o.VendorId = vendorId
	}
}

func WithEndpoint(endpoint string) func(options *ModelOptions) {
	return func(o *ModelOptions) {
		o.Endpoint = endpoint
	}
}

func WithRegion(region string) func(options *ModelOptions) {
	return func(o *ModelOptions) {
		o.Region = region
	}
}

func WithApiKey(apiKey string) func(*ModelOptions) {
	return func(o *ModelOptions) {
		o.ApiKey = apiKey
	}
}

func WithAccessKeySecretKey(accessKey string, secretKey string) func(*ModelOptions) {
	return func(o *ModelOptions) {
		o.AccessKey = accessKey
		o.SecretKey = secretKey
	}
}

func WithApiVersion(apiVersion string) func(*ModelOptions) {
	return func(o *ModelOptions) {
		o.ApiVersion = apiVersion
	}
}

func WithProxy(proxy string) func(*ModelOptions) {
	return func(o *ModelOptions) {
		o.Proxy = proxy
	}
}

func WithRetries(retries int) func(*ModelOptions) {
	return func(o *ModelOptions) {
		o.Retries = retries
	}
}

func WithRequestLog(requestLog func([]byte)) func(*ModelOptions) {
	return func(o *ModelOptions) {
		o.RequestLog = requestLog
	}
}

func WithResponseLog(responseLog func([]byte)) func(*ModelOptions) {
	return func(o *ModelOptions) {
		o.ResponseLog = responseLog
	}
}
