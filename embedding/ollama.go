package embedding

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Pooh-Mucho/go-aigc"
	"io"
	"net/http"
	"strings"
)

var ollamaModelMapping = map[string]string{
	string(Models.BaaiBgeM3):         "bge-m3:567m",
	string(Models.NomicEmbedText):    "nomic-embed-text:v1.5",
	string(Models.NomicEmbedTextV15): "nomic-embed-text:v1.5",
	string(Models.MxbaiEmbedLarge):   "mxbai-embed-large:335m",
	string(Models.MxbaiEmbedLargeV1): "mxbai-embed-large:335m",
}

type ollamaModelRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type ollamaModelResponse struct {
	Model           string      `json:"model"`
	Embeddings      [][]float32 `json:"embeddings"`
	TotalDuration   int         `json:"total_duration,omitempty"`
	LoadDuration    int         `json:"load_duration,omitempty"`
	PromptEvalCount int         `json:"prompt_eval_count,omitempty"`
}

type ollamaEmbeddingModel struct {
	ModelId     string
	Endpoint    string
	ApiKey      string
	Proxy       string
	Retries     int
	RequestLog  func([]byte)
	ResponseLog func([]byte)

	client aigc.HttpClient
}

func (r *ollamaModelRequest) load(request *ModelRequest) error {
	if len(request.Document) == 0 {
		return fmt.Errorf("[ollamaModelRequest:load] document is empty")
	}
	r.Input = []string{request.Document}
	return nil
}

func (r *ollamaModelResponse) dump(response *ModelResponse) error {
	if len(r.Embeddings) == 0 {
		return fmt.Errorf("[ollamaModelResponse:dump] embeddings is empty")
	}
	response.Embedding = r.Embeddings[0]
	response.Tokens = r.PromptEvalCount
	return nil
}

// http://host:port/api/embed
func (m *ollamaEmbeddingModel) getModelUrl() string {
	var url = m.Endpoint
	if strings.HasSuffix(url, "/api/embed") {
		return url
	}
	if strings.HasSuffix(url, "/") {
		return url + "api/embed"
	} else {
		return url + "/api/embed"
	}
}

func (m *ollamaEmbeddingModel) GetModelDeploymentId() string {
	var modelId, ok = ollamaModelMapping[m.ModelId]
	if ok {
		return modelId
	}
	return m.ModelId
}

func (m *ollamaEmbeddingModel) requestToJson(request *ModelRequest, jsonBuffer *bytes.Buffer) error {
	var err error
	var ollamaRequest ollamaModelRequest

	err = ollamaRequest.load(request)
	if err != nil {
		return fmt.Errorf("[ollamaEmbeddingModel.requestToJson] %w", err)
	}
	ollamaRequest.Model = m.GetModelDeploymentId()

	err = aigc.EncodeJson(jsonBuffer, ollamaRequest)
	if err != nil {
		return fmt.Errorf("[ollamaEmbeddingModel.requestToJson] %w", err)
	}
	return nil
}

func (m *ollamaEmbeddingModel) jsonToResponse(jsonBuffer *bytes.Buffer) (*ModelResponse, error) {
	var err error
	var ollamaResponse ollamaModelResponse

	err = aigc.DecodeJson(jsonBuffer, &ollamaResponse)
	if err != nil {
		return nil, fmt.Errorf("[ollamaEmbeddingModel.jsonToResponse] %w", err)
	}

	var response = ModelResponse{}
	err = ollamaResponse.dump(&response)
	if err != nil {
		return nil, fmt.Errorf("[ollamaEmbeddingModel.jsonToResponse] %w", err)
	}
	return &response, nil
}

func (m *ollamaEmbeddingModel) GetModelId() string {
	return m.ModelId
}

func (m *ollamaEmbeddingModel) GetDistanceType() VectorDistanceType {
	return CosineDistance
}

func (m *ollamaEmbeddingModel) Distance(vector1, vector2 []float32) (float32, error) {
	// Because most embeddings are normalized, we can use dot product instead
	// of cosine similarity.
	return VectorDotProduct(vector1, vector2)
	// return VectorCosineSimilarity(vector1, vector2)
}

func (m *ollamaEmbeddingModel) Embedding(ctx context.Context, request *ModelRequest) (*ModelResponse, error) {
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
		return nil, fmt.Errorf("[ollamaEmbeddingModel.Embedding] %w", err)
	}
	if m.RequestLog != nil {
		m.RequestLog(requestJson.Bytes())
	}

	httpRequest, err = http.NewRequestWithContext(ctx, http.MethodPost, modelUrl, requestJson)
	if err != nil {
		return nil, fmt.Errorf("[ollamaEmbeddingModel.Embedding] create http request %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	if m.ApiKey != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+m.ApiKey)
	}

	httpResponse, err = m.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("[ollamaEmbeddingModel.Embedding] do http request %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[ollamaEmbeddingModel.Embedding] http error %s %s",
			httpResponse.Status, aigc.HttpResponseText(httpResponse))
	}

	responseJson = aigc.AllocBuffer()
	defer aigc.FreeBuffer(responseJson)
	_, err = io.Copy(responseJson, httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("[ollamaEmbeddingModel.Embedding] read http response %w", err)
	}

	if m.ResponseLog != nil {
		m.ResponseLog(responseJson.Bytes())
	}

	response, err = m.jsonToResponse(responseJson)

	if err != nil {
		return nil, fmt.Errorf("[ollamaEmbeddingModel.Embedding] %w", err)
	}

	return response, nil
}

func newOllamaEmbeddingModel(modelId string, opts *aigc.ModelOptions) (*ollamaEmbeddingModel, error) {
	var model *ollamaEmbeddingModel

	if opts.Endpoint == "" {
		return nil, errors.New("ollama endpoint is required")
	}
	if opts.ApiVersion != "" {
		return nil, errors.New("ollama api version is not supported")
	}

	model = &ollamaEmbeddingModel{
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
