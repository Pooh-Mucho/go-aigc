package embedding

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

import (
	"github.com/Pooh-Mucho/go-aigc"
)

const (
	azureOpenAIDefaultApiVersion = "2024-06-01"

	openaiDefaultEndpoint = "https://api.openai.com/v1/embeddings"
)

// OpenAI Embedding models:
// +-------------------------------------------+------------+---------------+-----------+
// | Model                  | Max Input Tokens | Dimensions | Performance(on MTEB EVAL) |
// +------------------------+------------------+------------+---------------------------+
// | text-embedding-ada-002 |             8192 |       1536 |                     61.0% |
// | text-embedding-3-small |             8192 |       1536 |                     62.3% |
// | text-embedding-3-large |             8192 |       3072 |                     64.6% |
// +------------------------+------------------+------------+---------------------------+

type openaiModelRequest struct {
	// OpenAI should set Model. Azure API should not.
	Model string `json:"model,omitempty"`
	// Input text to embed, encoded as a string or array of string. The input
	// must not exceed the max input tokens for the model (8192 tokens for
	// text-embedding-ada-002)
	Input []string `json:"input,omitempty"`
	// The format to return the embeddings in. Can be either float or base64.
	// Default is float.
	EncodingFormat string `json:"encoding_format,omitempty"`
	// The number of dimensions the resulting output embeddings should have.
	// Only supported in text-embedding-3 and later models.
	Dimensions int `json:"dimensions,omitempty"`
	// A unique identifier representing your end-user, which can help OpenAI to
	// monitor and detect abuse
	User string `json:"user,omitempty"`
}

// The embedding object
type openaiEmbedding struct {
	// The object type, which is always "embedding".
	Object string `json:"object"`
	// The index of the embedding in the list of embeddings.
	Index int `json:"index"`
	// The embedding vector, which is a list of floats. The length of vector
	// depends on the model:
	Embedding []float32 `json:"embedding"`
}

type openaiModelResponse struct {
	// Always "list"
	Object string `json:"object,omitempty"`
	// A list of embedding objects.
	Data []openaiEmbedding `json:"data,omitempty"`
	// The model used to generate the embeddings.
	Model string `json:"model,omitempty"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens,omitempty"`
		TotalTokens  int `json:"total_tokens,omitempty"`
	}
}

type openaiEmbeddingModel struct {
	ModelId     string
	Endpoint    string
	ApiKey      string
	Proxy       string
	Retries     int
	RequestLog  func([]byte)
	ResponseLog func([]byte)

	client aigc.HttpClient
}

type azureOpenAIEmbeddingModel struct {
	ModelId     string
	Endpoint    string
	ApiKey      string
	ApiVersion  string
	Proxy       string
	Retries     int
	RequestLog  func([]byte)
	ResponseLog func([]byte)

	client aigc.HttpClient
}

func (r *openaiModelRequest) load(request *ModelRequest) error {
	if len(request.Document) == 0 {
		return fmt.Errorf("[openaiModelRequest:load] document is empty")
	}
	r.Input = []string{request.Document}
	return nil
}

func (r *openaiModelResponse) dump(response *ModelResponse) error {
	if len(r.Data) == 0 {
		return fmt.Errorf("[openaiModelResponse:dump] data is empty")
	}
	response.Embedding = r.Data[0].Embedding
	response.Tokens = r.Usage.TotalTokens
	return nil
}

func (m *openaiEmbeddingModel) getModelUrl() string {
	if m.Endpoint == "" {
		return openaiDefaultEndpoint
	}
	return m.Endpoint
}

func (m *openaiEmbeddingModel) requestToJson(request *ModelRequest, jsonBuffer *bytes.Buffer) error {
	var err error
	var openaiRequest openaiModelRequest

	err = openaiRequest.load(request)
	if err != nil {
		return fmt.Errorf("[openaiEmbeddingModel.requestToJson] %w", err)
	}
	openaiRequest.Model = m.ModelId

	err = aigc.EncodeJson(jsonBuffer, openaiRequest)
	if err != nil {
		return fmt.Errorf("[openaiEmbeddingModel.requestToJson] %w", err)
	}
	return nil
}

func (m *openaiEmbeddingModel) jsonToResponse(jsonBuffer *bytes.Buffer) (*ModelResponse, error) {
	var err error
	var openaiResponse openaiModelResponse

	err = aigc.DecodeJson(jsonBuffer, &openaiResponse)
	if err != nil {
		return nil, fmt.Errorf("[openaiEmbeddingModel.jsonToResponse] %w", err)
	}

	var response = ModelResponse{}
	err = openaiResponse.dump(&response)
	if err != nil {
		return nil, fmt.Errorf("[openaiEmbeddingModel.jsonToResponse] %w", err)
	}
	return &response, nil
}

func (m *openaiEmbeddingModel) GetModelId() string {
	return m.ModelId
}

func (m *openaiEmbeddingModel) GetDistanceType() VectorDistanceType {
	return CosineDistance
}

func (m *openaiEmbeddingModel) Distance(vector1, vector2 []float32) (float32, error) {
	// Because OpenAI embeddings are normalized, we can use dot product instead
	// of cosine similarity.
	return VectorDotProduct(vector1, vector2)
	// return VectorCosineSimilarity(vector1, vector2)
}

func (m *openaiEmbeddingModel) Embedding(ctx context.Context, request *ModelRequest) (*ModelResponse, error) {
	var err error
	var modelUrl string
	var requestJson *bytes.Buffer
	var responseJson *bytes.Buffer
	var response *ModelResponse
	var httpRequest *http.Request
	var httpResponse *http.Response

	modelUrl = m.getModelUrl()
	requestJson = aigc.AllocBuffer()

	err = m.requestToJson(request, requestJson)
	if err != nil {
		return nil, fmt.Errorf("[openaiEmbeddingModel.Embedding] %w", err)
	}
	if m.RequestLog != nil {
		m.RequestLog(requestJson.Bytes())
	}

	httpRequest, err = http.NewRequestWithContext(ctx, http.MethodPost, modelUrl, requestJson)
	if err != nil {
		return nil, fmt.Errorf("[openaiEmbeddingModel.Embedding] create http request %w", err)
	}

	httpRequest.Header.Set("Authorization", "Bearer "+m.ApiKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err = m.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("[openaiEmbeddingModel.Embedding] do http request %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[openaiEmbeddingModel.Embedding] http error %s %s",
			httpResponse.Status, aigc.HttpResponseText(httpResponse))
	}

	responseJson = aigc.AllocBuffer()
	defer aigc.FreeBuffer(responseJson)
	_, err = io.Copy(responseJson, httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("[openaiEmbeddingModel.Embedding] read http response %w", err)
	}

	if m.ResponseLog != nil {
		m.ResponseLog(responseJson.Bytes())
	}

	response, err = m.jsonToResponse(responseJson)

	if err != nil {
		return nil, fmt.Errorf("[openaiEmbeddingModel.Embedding] %w", err)
	}

	return response, nil
}

func (m *azureOpenAIEmbeddingModel) getModelUrl() string {
	var builder strings.Builder

	builder.WriteString(m.Endpoint)
	if !strings.HasSuffix(m.Endpoint, "/") {
		builder.WriteByte('/')
	}
	builder.WriteString("openai/deployments/")
	builder.WriteString(m.ModelId)
	builder.WriteString("/embeddings?api-version=")
	if m.ApiVersion != "" {
		builder.WriteString(m.ApiVersion)
	} else {
		builder.WriteString(azureOpenAIDefaultApiVersion)
	}
	return builder.String()
}

func (m *azureOpenAIEmbeddingModel) requestToJson(request *ModelRequest, jsonBuffer *bytes.Buffer) error {
	var err error
	var openaiRequest openaiModelRequest

	err = openaiRequest.load(request)
	if err != nil {
		return fmt.Errorf("[azureOpenAIEmbeddingModel.requestToJson] %w", err)
	}

	err = aigc.EncodeJson(jsonBuffer, openaiRequest)
	if err != nil {
		return fmt.Errorf("[azureOpenAIEmbeddingModel.requestToJson] %w", err)
	}
	return nil
}

func (m *azureOpenAIEmbeddingModel) jsonToResponse(jsonBuffer *bytes.Buffer) (*ModelResponse, error) {
	var err error
	var openaiResponse openaiModelResponse

	err = aigc.DecodeJson(jsonBuffer, &openaiResponse)
	if err != nil {
		return nil, fmt.Errorf("[azureOpenAIEmbeddingModel.jsonToResponse] %w", err)
	}

	var response = ModelResponse{}
	err = openaiResponse.dump(&response)
	if err != nil {
		return nil, fmt.Errorf("[azureOpenAIEmbeddingModel.jsonToResponse] %w", err)
	}
	return &response, nil
}

func (m *azureOpenAIEmbeddingModel) GetModelId() string {
	return m.ModelId
}

func (m *azureOpenAIEmbeddingModel) GetDistanceType() VectorDistanceType {
	return CosineDistance
}

func (m *azureOpenAIEmbeddingModel) Distance(vector1, vector2 []float32) (float32, error) {
	// Because OpenAI embeddings are normalized, we can use dot product instead
	// of cosine similarity.
	return VectorDotProduct(vector1, vector2)
	// return VectorCosineSimilarity(vector1, vector2)
}

func (m *azureOpenAIEmbeddingModel) Embedding(ctx context.Context, request *ModelRequest) (*ModelResponse, error) {
	var err error
	var modelUrl string
	var requestJson *bytes.Buffer
	var responseJson *bytes.Buffer
	var response *ModelResponse
	var httpRequest *http.Request
	var httpResponse *http.Response

	modelUrl = m.getModelUrl()
	requestJson = aigc.AllocBuffer()

	err = m.requestToJson(request, requestJson)
	if err != nil {
		return nil, fmt.Errorf("[azureOpenAIEmbeddingModel.Embedding] %w", err)
	}
	if m.RequestLog != nil {
		m.RequestLog(requestJson.Bytes())
	}

	httpRequest, err = http.NewRequestWithContext(ctx, http.MethodPost, modelUrl, requestJson)
	if err != nil {
		return nil, fmt.Errorf("[azureOpenAIEmbeddingModel.Embedding] create http request %w", err)
	}
	httpRequest.Header.Set("api-key", m.ApiKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err = m.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("[azureOpenAIEmbeddingModel.Embedding] do http request %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[azureOpenAIEmbeddingModel.Embedding] http error %s %s",
			httpResponse.Status, aigc.HttpResponseText(httpResponse))
	}

	responseJson = aigc.AllocBuffer()
	defer aigc.FreeBuffer(responseJson)
	_, err = io.Copy(responseJson, httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("[azureOpenAIEmbeddingModel.Embedding] read http response %w", err)
	}

	if m.ResponseLog != nil {
		m.ResponseLog(responseJson.Bytes())
	}

	response, err = m.jsonToResponse(responseJson)

	if err != nil {
		return nil, fmt.Errorf("[azureOpenAIEmbeddingModel.Embedding] %w", err)
	}

	return response, nil
}

func newOpenAIEmbeddingModel(modelId string, opts *aigc.ModelOptions) (*openaiEmbeddingModel, error) {
	var model *openaiEmbeddingModel

	if opts.ApiKey == "" {
		return nil, errors.New("openai api key is required")
	}
	if opts.ApiVersion != "" {
		return nil, errors.New("openai api version is not supported")
	}

	model = &openaiEmbeddingModel{
		ModelId:     modelId,
		Endpoint:    opts.Endpoint,
		ApiKey:      opts.ApiKey,
		Proxy:       opts.Proxy,
		Retries:     opts.Retries,
		RequestLog:  opts.RequestLog,
		ResponseLog: opts.ResponseLog,
	}

	model.client = aigc.HttpClient{
		Proxy:   opts.Proxy,
		Retries: opts.Retries,
	}

	return model, nil
}

func newAzureOpenAIEmbeddingModel(modelId string, opts *aigc.ModelOptions) (*azureOpenAIEmbeddingModel, error) {
	var model *azureOpenAIEmbeddingModel

	if opts.Endpoint == "" {
		return nil, errors.New("azure endpoint is required")
	}

	if opts.ApiKey == "" {
		return nil, errors.New("azure api key is required")
	}

	model = &azureOpenAIEmbeddingModel{
		ModelId:     modelId,
		Endpoint:    opts.Endpoint,
		ApiKey:      opts.ApiKey,
		ApiVersion:  opts.ApiVersion,
		Proxy:       opts.Proxy,
		Retries:     opts.Retries,
		RequestLog:  opts.RequestLog,
		ResponseLog: opts.ResponseLog,
	}

	model.client = aigc.HttpClient{
		Proxy:   opts.Proxy,
		Retries: opts.Retries,
	}

	return model, nil
}
