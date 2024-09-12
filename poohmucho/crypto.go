package poohmucho

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"time"
	"unsafe"
)

type Nonce [16]byte

type Cipher struct {
	aes128Key [16]byte
	block     cipher.Block
}

type Signer struct {
	hmac hash.Hash
}

func unsafeNonEscape(a []byte) []byte {
	var p1 unsafe.Pointer = unsafe.Pointer(unsafe.SliceData(a))
	var p2 unsafe.Pointer = unsafe.Pointer(uintptr(p1) ^ 0)
	return unsafe.Slice((*byte)(p2), len(a))
}

func incrmentCtr(ctr *[16]byte) {
	for i := 0; i < 16; i++ {
		ctr[i]++
		if ctr[i] != 0 {
			break
		}
	}
}

// z = x ^ y, 16 bytes
func xor16(x unsafe.Pointer, y unsafe.Pointer, z unsafe.Pointer) {
	type uint128 [2]uint64

	var ux = (*uint128)(x)
	var uy = (*uint128)(y)
	var uz = (*uint128)(z)

	uz[0] = ux[0] ^ uy[0]
	uz[1] = ux[1] ^ uy[1]
}

func NewNonce() Nonce {
	var nonce Nonce
	var seconds uint64
	rand.Read(nonce[:8])
	seconds = uint64(time.Now().Unix())
	binary.BigEndian.PutUint64(nonce[8:16], seconds)
	return nonce
}

func hmacHash(key string, message string) [32]byte {
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

func NewKeyId(key string) [16]byte {
	var sum = hmacHash(key, key)
	var keyId = ([16]byte)(sum[0:16])
	return keyId
}

func NewCipher(key string) Cipher {
	const HMAC_MESSAGE = "CIPHER-v6FRDZ2NgQmTdIigf97hK5FwlIBvYuxBgPYgsLVj"

	var sum [32]byte = hmacHash(key, HMAC_MESSAGE)
	return Cipher{
		aes128Key: [16]byte(sum[0:16]),
	}
}

func (c *Cipher) Encrypt(data []byte, nonce Nonce, b *bytes.Buffer) error {
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

func (c *Cipher) Decrypt(data []byte, nonce Nonce, b *bytes.Buffer) error {
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

func (c *Cipher) transform(data []byte, nonce Nonce, b *bytes.Buffer, block cipher.Block) error {
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
		incrmentCtr(&ctr)
		block.Encrypt(xor_ne, ctr_ne) // use non escaping slice
		xor16(unsafe.Pointer(&data[index]), unsafe.Pointer(&xor), unsafe.Pointer(&buf))
		b.Write(buf_ne)
		index += 16
	}

	if index < len(data) {
		incrmentCtr(&ctr)
		block.Encrypt(xor_ne, ctr_ne) // use non escaping slice
		for i := index; i < len(data); i++ {
			buf[i-index] = data[i] ^ xor[i-index]
		}
		b.Write(buf[:len(data)-index])
	}

	return nil
}

func NewSigner(key string) Signer {
	const HMAC_MESSAGE = "SINGER-qSQD7uBp2fqWDgABRdERNWQr7lBViICJxPiwigYh" // len = 40

	var sum [32]byte = hmacHash(key, HMAC_MESSAGE)
	return Signer{hmac: hmac.New(sha256.New, sum[:])}
}

func (s *Signer) Signature(data []byte) [20]byte {
	var sum [20]byte
	s.hmac.Reset()
	s.hmac.Write(data)
	s.hmac.Sum(sum[:0])
	return sum
}
