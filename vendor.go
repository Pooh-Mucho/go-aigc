package aigc

type VendorId string

var Vendors = struct {
	OpenAI      VendorId
	Anthropic   VendorId
	ZhipuAI     VendorId
	MidJourney  VendorId
	StabilityAI VendorId
	Runway      VendorId
	Amazon      VendorId
	Microsoft   VendorId
	Google      VendorId
	Alibaba     VendorId
	Baidu       VendorId
	ByteDance   VendorId
	HuggingFace VendorId
	Ollama      VendorId
	PoohMucho   VendorId
}{
	OpenAI:      "OpenAI",
	Anthropic:   "Anthropic",
	ZhipuAI:     "zhipuai",
	MidJourney:  "MidJourney",
	StabilityAI: "StabilityAI",
	Runway:      "Runway",
	Amazon:      "AWS",
	Microsoft:   "Azure",
	Google:      "GCP",
	Alibaba:     "AliYun",
	Baidu:       "Baidu",
	ByteDance:   "ByteDance",
	HuggingFace: "HuggingFace",
	Ollama:      "Ollama",
	PoohMucho:   "PoohMucho",
}
