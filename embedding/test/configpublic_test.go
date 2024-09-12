//go:build !private
// +build !private

package test

var WithAWS = func(options *aigc.ModelOptions) {
	options.VendorId = aigc.Vendors.Amazon
	options.AccessKey = ""
	options.SecretKey = ""
	options.Region = "us-west-2"
}

var WithAnthropic = func(options *aigc.ModelOptions) {
	options.VendorId = aigc.Vendors.Anthropic
	options.ApiKey = ""
	options.Proxy = ""
}

var WithAzure = func(options *aigc.ModelOptions) {
	options.VendorId = aigc.Vendors.Microsoft
	options.Endpoint = ""
	options.ApiKey = ""
}

var WithOpenAI = func(options *aigc.ModelOptions) {
	options.VendorId = aigc.Vendors.OpenAI
	options.Endpoint = ""
	options.ApiKey = ""
	options.Proxy = ""
}

var WithAliyun = func(options *aigc.ModelOptions) {
	options.VendorId = aigc.Vendors.Alibaba
	options.Endpoint = ""
	options.ApiKey = ""
}

var WithOllama = func(options *aigc.ModelOptions) {
	options.VendorId = aigc.Vendors.Ollama
	options.ApiKey = ""
	options.Endpoint = ""
}
