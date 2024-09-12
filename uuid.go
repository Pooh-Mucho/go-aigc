package aigc

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"net"
	"sync"
	"sync/atomic"
	_ "unsafe"
)

type UUID [16]byte

var (
	uuidSeq    uint64 = 1
	uuidHwSig  []byte
	uuidHwInit sync.Once

	EmptyUUID = UUID{}
)

//go:linkname runtime_nanotime runtime.nanotime
func runtime_nanotime() int64

func NewUUID() UUID {
	/* initialize hardware addresses */
	uuidHwInit.Do(
		func() {
			var ifaces, err = net.Interfaces()
			if err == nil {
				for _, iface := range ifaces {
					uuidHwSig = append(uuidHwSig, iface.HardwareAddr...)
				}
			}
		})

	var buffer = []byte{79: 0}
	buffer = buffer[:0]
	/* hardware signature */
	buffer = append(buffer, uuidHwSig...)
	/* sequence */
	var seq = atomic.AddUint64(&uuidSeq, 1)
	var seqBuf [8]byte
	binary.BigEndian.PutUint64(seqBuf[:], seq)
	buffer = append(buffer, seqBuf[:]...)
	/* timestamp */
	var nano = runtime_nanotime()
	var nanoBuf [8]byte
	binary.BigEndian.PutUint64(nanoBuf[:], uint64(nano))
	buffer = append(buffer, nanoBuf[:]...)
	/* random */
	var randBuf [16]byte
	rand.Reader.Read(randBuf[:])
	buffer = append(buffer, randBuf[:]...)

	var hash = sha1.Sum(buffer)

	var result UUID
	copy(result[:], hash[:])
	return result
}

// String returns a canonical RFC-4122 string representation of the UUID:
//
//	xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u UUID) String() string {
	return u.StringD()
}

// StringN returns a non-canonical string representation of the UUID:
//
//	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u UUID) StringN() string {
	buf := make([]byte, 32)

	hex.Encode(buf[0:8], u[0:4])
	hex.Encode(buf[8:12], u[4:6])
	hex.Encode(buf[12:16], u[6:8])
	hex.Encode(buf[16:20], u[8:10])
	hex.Encode(buf[20:], u[10:])

	return string(buf)
}

// StringD returns a non-canonical string representation of the UUID:
//
//	xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u UUID) StringD() string {
	var buf = []byte{35: 0}

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}
