package chat

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pooh-Mucho/go-aigc"
	"github.com/Pooh-Mucho/go-aigc/poohmucho"
	"hash"
	"net/http"
	"time"
	"unsafe"
)

/* Model name format: "vendor:model[:date]"
 * For example:
 *   openai:gpt-4o:2024-08-16
 *   openai:gpt-4-turbo:2024-04-09
 *   openai:gpt-3.5-turbo:2024-01-25
 *   azure:gpt-4-turbo:2024-01-25
 *   anthropic:claude-3.5-sonnet:2024-06-20
 *   bedrock:claude-3-haiku:2024-03-07
 *   bedrock:llama-3-1-70b
 *   aliyun:qwenmax:2024-04-28
 */

type poohmuchoNonce [16]byte

type poohmuchoCipher struct {
	aes128Key [16]byte
	block     cipher.Block
}

type poohmuchoSigner struct {
	hmac hash.Hash
}

type poohmuchoModelRequest struct {
	Model             string      `json:"model"`
	Messages          []Message   `json:"messages"`
	Tools             []Tool      `json:"tools,omitempty"`
	MaxTokens         *int32      `json:"max_tokens,omitempty"`
	Temperature       *float64    `json:"temperature,omitempty"`
	TopP              *float64    `json:"top_p,omitempty"`
	ToolChoice        *ToolChoice `json:"tool_choice,omitempty"`
	ParallelToolCalls *bool       `json:"parallel_tool_calls,omitempty"`

	maxTokensValue         int32
	temperatureValue       float64
	topPValue              float64
	toolChoiceValue        ToolChoice
	parallelToolCallsValue bool
}

type poohmuchoModelResponse struct {
	Id           string    `json:"id"`
	Messages     []Message `json:"messages"`
	FinishReason string    `json:"finish_reason"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	}
	ContentFilterResult string `json:"content_filter_result"`
	ErrorMessage        string `json:"error_message"`
}

type poohmuchoModel struct {
	ModelId     string
	Endpoint    string
	ApiKey      string
	Proxy       string
	Retries     int
	RequestLog  func([]byte)
	ResponseLog func([]byte)

	httpClient  aigc.HttpClient
	modelClient *poohmucho.Client
}

func unsafeNonEscape(a []byte) []byte {
	var p1 unsafe.Pointer = unsafe.Pointer(unsafe.SliceData(a))
	var p2 unsafe.Pointer = unsafe.Pointer(uintptr(p1) ^ 0)
	return unsafe.Slice((*byte)(p2), len(a))
}

func poohmuchoIncrmentCtr(ctr *[16]byte) {
	for i := 0; i < 16; i++ {
		ctr[i]++
		if ctr[i] != 0 {
			break
		}
	}
}

// z = x ^ y, 16 bytes
func poohmuchoXor16(x unsafe.Pointer, y unsafe.Pointer, z unsafe.Pointer) {
	type uint128 [2]uint64

	var ux = (*uint128)(x)
	var uy = (*uint128)(y)
	var uz = (*uint128)(z)

	uz[0] = ux[0] ^ uy[0]
	uz[1] = ux[1] ^ uy[1]
}

func poohmuchoNewNonce() poohmuchoNonce {
	var nonce poohmuchoNonce
	var seconds uint64
	rand.Read(nonce[:8])
	seconds = uint64(time.Now().Unix())
	binary.BigEndian.PutUint64(nonce[8:16], seconds)
	return nonce
}

func poohmuchoHMAC(key string, message string) [32]byte {
	// HMAC(key, message) = Hash((key ^ outer_pad) || Hash((key ^ inner_pad) || message))

	const INNER_PAD = 0x36
	const OUTER_PAD = 0x5c

	var buf = []byte{119: 0}
	var sum [32]byte

	buf = append(buf, key...)
	for i := 0; i < len(buf); i++ {
		buf[i] ^= OUTER_PAD
	}
	buf = append(buf, key...)
	for i := len(key); i < len(buf); i++ {
		buf[i] ^= INNER_PAD
	}
	// Hash inner
	sum = sha256.Sum256(buf[len(key):])
	buf = append(buf[:len(key)], sum[:]...)
	buf = append(buf, message...)

	sum = sha256.Sum256(buf)
	return sum
}

func poohmuchoKeyId(key string) [16]byte {
	var sum = poohmuchoHMAC(key, key)
	var keyId = ([16]byte)(sum[0:16])
	return keyId
}

func poohmuchoNewCipher(key string) poohmuchoCipher {
	const HMAC_MESSAGE = "CIPHER-v6FRDZ2NgQmTdIigf97hK5FwlIBvYuxBgPYgsLVj"

	var sum [32]byte = poohmuchoHMAC(key, HMAC_MESSAGE)
	return poohmuchoCipher{
		aes128Key: [16]byte(sum[0:16]),
	}
}

func (c *poohmuchoCipher) Encrypt(data []byte, nonce poohmuchoNonce, b *bytes.Buffer) error {
	var err error
	var block cipher.Block = c.block

	if block == nil {
		block, err = aes.NewCipher(c.aes128Key[:])
		if err != nil {
			return err
		}
		c.block = block
	}

	return c.transform(data, nonce, b, block)
}

func (c *poohmuchoCipher) Decrypt(data []byte, nonce poohmuchoNonce, b *bytes.Buffer) error {
	var err error
	var block cipher.Block = c.block

	if block == nil {
		block, err = aes.NewCipher(c.aes128Key[:])
		if err != nil {
			return err
		}
		c.block = block
	}

	return c.transform(data, nonce, b, block)
}

func (c *poohmuchoCipher) transform(data []byte, nonce poohmuchoNonce, b *bytes.Buffer, block cipher.Block) error {
	var (
		ctr    [16]byte = nonce
		xor    [16]byte
		buf    [16]byte
		ctr_ne []byte = unsafeNonEscape(ctr[:]) // make non escaping ctr slice
		xor_ne []byte = unsafeNonEscape(xor[:]) // make non escaping xor slice
		buf_ne []byte = unsafeNonEscape(buf[:]) // make non escaping xor slice
		end16         = (len(data) / 16) * 16
		index         = 0
	)

	for index < end16 {
		poohmuchoIncrmentCtr(&ctr)
		block.Encrypt(xor_ne, ctr_ne) // use non escaping slice
		poohmuchoXor16(unsafe.Pointer(&data[index]), unsafe.Pointer(&xor), unsafe.Pointer(&buf))
		b.Write(buf_ne)
		index += 16
	}

	if index < len(data) {
		poohmuchoIncrmentCtr(&ctr)
		block.Encrypt(xor_ne, ctr_ne) // use non escaping slice
		for i := index; i < len(data); i++ {
			buf[i-index] = data[i] ^ xor[i-index]
		}
		b.Write(buf[:len(data)-index])
	}

	return nil
}

func poohmuchoNewSigner(key string) poohmuchoSigner {
	const HMAC_MESSAGE = "SINGER-qSQD7uBp2fqWDgABRdERNWQr7lBViICJxPiwigYh" // len = 40

	var sum [32]byte = poohmuchoHMAC(key, HMAC_MESSAGE)
	return poohmuchoSigner{hmac: hmac.New(sha256.New, sum[:])}
}

func (s *poohmuchoSigner) Signature(data []byte) [20]byte {
	var sum [20]byte
	s.hmac.Reset()
	s.hmac.Write(data)
	s.hmac.Sum(sum[:0])
	return sum
}

func (r *poohmuchoModelRequest) load(request *ModelRequest) error {
	r.Messages = request.Messages
	r.Tools = request.Tools
	if request.MaxTokens.Valid {
		r.maxTokensValue = request.MaxTokens.Value
		r.MaxTokens = &r.maxTokensValue
	} else {
		r.maxTokensValue = 0
		r.MaxTokens = nil
	}
	if request.Temperature.Valid {
		r.temperatureValue = request.Temperature.Value
		r.Temperature = &r.temperatureValue
	} else {
		r.temperatureValue = 0
		r.Temperature = nil
	}
	if request.TopP.Valid {
		r.topPValue = request.TopP.Value
		r.TopP = &r.topPValue
	} else {
		r.topPValue = 0
		r.TopP = nil
	}
	if request.ToolChoice != nil {
		r.toolChoiceValue = *request.ToolChoice
		r.ToolChoice = &r.toolChoiceValue
	} else {
		r.ToolChoice = nil
	}
	if request.ParallelToolCalls.Valid {
		r.parallelToolCallsValue = request.ParallelToolCalls.Value
		r.ParallelToolCalls = &r.parallelToolCallsValue
	} else {
		r.parallelToolCallsValue = false
		r.ParallelToolCalls = nil
	}
	return nil
}

func (r *poohmuchoModelResponse) dump(response *ModelResponse) error {
	response.Id = r.Id
	response.Messages = r.Messages
	response.Usage.InputTokens = r.Usage.InputTokens
	response.Usage.OutputTokens = r.Usage.OutputTokens
	response.FinishReason = FinishReason(r.FinishReason)
	response.ContentFilterResult = r.ContentFilterResult
	return nil
}

func (m *poohmuchoModel) requestToJson(request *ModelRequest, jsonBuffer *bytes.Buffer) error {
	var err error
	var poohmuchoRequest poohmuchoModelRequest
	var encoder *json.Encoder

	err = poohmuchoRequest.load(request)
	if err != nil {
		return fmt.Errorf("[poohmuchoModel.requestToJson] %w", err)
	}

	poohmuchoRequest.Model = m.ModelId

	encoder = json.NewEncoder(jsonBuffer)
	encoder.SetEscapeHTML(false)

	err = encoder.Encode(poohmuchoRequest)
	if err != nil {
		return fmt.Errorf("[poohmuchoModel.requestToJson] %w", err)
	}
	return nil
}

func (m *poohmuchoModel) jsonToResponse(jsonBuffer *bytes.Buffer) (*ModelResponse, error) {
	var err error
	var poohmuchoResponse poohmuchoModelResponse

	err = json.Unmarshal(jsonBuffer.Bytes(), &poohmuchoResponse)
	if err != nil {
		return nil, fmt.Errorf("[poohmuchoModel.jsonToResponse] %w", err)
	}

	if poohmuchoResponse.ErrorMessage != "" {
		return nil, fmt.Errorf("[poohmuchoModel.jsonToResponse] server error %s", poohmuchoResponse.ErrorMessage)
	}

	var response = ModelResponse{}
	err = poohmuchoResponse.dump(&response)
	if err != nil {
		return nil, fmt.Errorf("[poohmuchoModel.jsonToResponse] %w", err)
	}
	return &response, nil
}

func (m *poohmuchoModel) GetModelId() string {
	return m.ModelId
}

func (m *poohmuchoModel) Complete(ctx context.Context, request *ModelRequest) (*ModelResponse, error) {
	var (
		err          error
		requestJson  *bytes.Buffer
		responseJson *bytes.Buffer
		response     *ModelResponse
	)

	requestJson = aigc.AllocBuffer()
	defer aigc.FreeBuffer(requestJson)

	err = m.requestToJson(request, requestJson)
	if err != nil {
		return nil, fmt.Errorf("[poohmuchoModel.Complete] %w", err)
	}

	if m.RequestLog != nil {
		m.RequestLog(requestJson.Bytes())
	}

	responseJson = aigc.AllocBuffer()
	defer aigc.FreeBuffer(responseJson)

	err = m.modelClient.InvokeModel(ctx, requestJson, responseJson)

	if err != nil {
		if responseJson.Len() > 0 {
			return nil, fmt.Errorf("[poohmuchoModel.Complete] %w %s", err, responseJson.String())
		}
		return nil, fmt.Errorf("[poohmuchoModel.Complete] %w", err)
	}

	if m.ResponseLog != nil {
		m.ResponseLog(responseJson.Bytes())
	}

	response, err = m.jsonToResponse(responseJson)

	if err != nil {
		return nil, fmt.Errorf("[poohmuchoModel.Complete] %w", err)
	}

	return response, nil
}

func newPoohMuchoChatModel(modelId string, opts *aigc.ModelOptions) (*poohmuchoModel, error) {
	if opts.ApiKey == "" {
		return nil, errors.New("PoohMucho api key is required")
	}
	if opts.Endpoint == "" {
		return nil, errors.New("PoohMucho endpoint is required")
	}
	if opts.ApiVersion != "" {
		return nil, errors.New("PoohMucho api version is not supported")
	}

	var err error
	var httpClient *http.Client
	var model = &poohmuchoModel{
		ModelId:  modelId,
		Endpoint: opts.Endpoint,
		ApiKey:   opts.ApiKey,
		Proxy:    opts.Proxy,
		Retries:  opts.Retries,
	}
	model.httpClient = aigc.HttpClient{Proxy: opts.Proxy, Retries: opts.Retries}

	httpClient, err = model.httpClient.Client()
	if err != nil {
		return nil, fmt.Errorf("[newPoohMuchoChatModel] %w", err)
	}

	model.modelClient = poohmucho.NewClient(opts.Endpoint, opts.ApiKey, httpClient)

	return model, nil
}
