package poohmucho

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"unsafe"
)

import (
	"github.com/Pooh-Mucho/go-aigc"
)

type Client struct {
	endpoint string
	apiKey   string
	keyId    [16]byte
	cipher   Cipher
	signer   Signer
	client   *http.Client
}

func (c *Client) InvokeModel(ctx context.Context, input *bytes.Buffer, output *bytes.Buffer) error {
	var err error
	var nonce Nonce
	var inputSignature [20]byte
	var encryptedInput *bytes.Buffer
	var encryptedOutput *bytes.Buffer
	var outputSignatureHeader string
	var request *http.Request
	var response *http.Response

	nonce = NewNonce()

	encryptedInput = aigc.AllocBuffer()
	defer aigc.FreeBuffer(encryptedInput)

	err = c.cipher.Encrypt(input.Bytes(), nonce, encryptedInput)
	if err != nil {
		return fmt.Errorf("[Client.InvokeModel] encrypt %w", err)
	}

	inputSignature = c.signer.Signature(encryptedInput.Bytes())
	request, err = http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, encryptedInput)
	if err != nil {
		return fmt.Errorf("[Client.InvokeModel] new request %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("X-Auth-KeyId", hex.EncodeToString(c.keyId[:]))
	request.Header.Set("X-Auth-Nonce", hex.EncodeToString(nonce[:]))
	request.Header.Set("X-Auth-Signature", hex.EncodeToString(inputSignature[:]))

	response, err = c.client.Do(request)
	if err != nil {
		return fmt.Errorf("[Client.InvokeModel] do request %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		io.Copy(output, response.Body)
		return fmt.Errorf("[Client.InvokeModel] http status %s", response.Status)
	}

	encryptedOutput = aigc.AllocBuffer()
	defer aigc.FreeBuffer(encryptedOutput)

	_, err = io.Copy(encryptedOutput, response.Body)
	if err != nil {
		return fmt.Errorf("[Client.InvokeModel] read response %w", err)
	}
	outputSignatureHeader = response.Header.Get("X-Auth-Signature")
	if outputSignatureHeader != "" {
		var decodeLen int
		var outputSignature [20]byte
		var expected [20]byte
		var buff = unsafe.Slice(unsafe.StringData(outputSignatureHeader), len(outputSignatureHeader))

		decodeLen, err = hex.Decode(buff, outputSignature[:])
		if err != nil {
			return fmt.Errorf("[Client.InvokeModel] decode signature %w", err)
		}
		if decodeLen != 20 {
			return fmt.Errorf("[Client.InvokeModel] invalid signature length %d", decodeLen)
		}
		expected = c.signer.Signature(encryptedOutput.Bytes())
		if outputSignature != expected {
			return fmt.Errorf("[Client.InvokeModel] invalid signature")
		}
	}

	err = c.cipher.Decrypt(encryptedOutput.Bytes(), nonce, output)
	if err != nil {
		return fmt.Errorf("[Client.InvokeModel] decrypt %w", err)
	}

	return nil
}

func NewClient(endpoint string, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		endpoint: endpoint,
		apiKey:   apiKey,
		keyId:    NewKeyId(apiKey),
		cipher:   NewCipher(apiKey),
		signer:   NewSigner(apiKey),
		client:   httpClient,
	}
}
