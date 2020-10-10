package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"net"
	"sync"
	"time"
)

const radix = 62

var digitalAry62 = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func digTo62(_val int64, _digs byte, _sb *bytes.Buffer) {
	hi := int64(1) << (_digs * 4)
	i := hi | (_val & (hi - 1))

	negative := i < 0
	if !negative {
		i = -i
	}

	skip := true
	for i <= -radix {
		if skip {
			skip = false
		} else {
			offset := -(i % radix)
			_sb.WriteByte(digitalAry62[int(offset)])
		}
		i = i / radix
	}
	_sb.WriteByte(digitalAry62[int(-i)])

	if negative {
		_sb.WriteByte('-')
	}
}

func suidToShortS(data []byte) string {
	// [16]byte
	buf := make([]byte, 22)
	sb := bytes.NewBuffer(buf)
	sb.Reset()

	var msb int64
	for i := 0; i < 8; i++ {
		msb = msb<<8 | int64(data[i])
	}

	var lsb int64
	for i := 8; i < 16; i++ {
		lsb = lsb<<8 | int64(data[i])
	}

	digTo62(msb>>12, 8, sb)
	digTo62(msb>>16, 4, sb)
	digTo62(msb, 4, sb)
	digTo62(lsb>>48, 4, sb)
	digTo62(lsb, 12, sb)

	return sb.String()
}

// UUID layout variants.
const (
	VariantNCS = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

// Difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and Unix epoch (January 1, 1970).
const epochStart = 122192928000000000

// Used in string method conversion
const dash byte = '-'

// UUID v1/v2 storage.
var (
	storageMutex  sync.Mutex
	storageOnce   sync.Once
	clockSequence uint16
	lastTime      uint64
	hardwareAddr  [6]byte
)

func initClockSequence() {
	buf := make([]byte, 2)
	safeRandom(buf)
	clockSequence = binary.BigEndian.Uint16(buf)
}

func initHardwareAddr() {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) >= 6 {
				copy(hardwareAddr[:], iface.HardwareAddr)
				return
			}
		}
	}

	// Initialize hardwareAddr randomly in case
	// of real network interfaces absence
	safeRandom(hardwareAddr[:])

	// Set multicast bit as recommended in RFC 4122
	hardwareAddr[0] |= 0x01
}

func initStorage() {
	initClockSequence()
	initHardwareAddr()
}

func safeRandom(dest []byte) {
	if _, err := rand.Read(dest); err != nil {
		panic(err)
	}
}

// Returns difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and current time.
// This is default epoch calculation function.
func unixTimeFunc() uint64 {
	return epochStart + uint64(time.Now().UnixNano()/100)
}

// UUID representation compliant with specification
// described in RFC 4122.
type SUID struct {
	value []byte // [16]byte
}

// The nil UUID is special form of UUID that is specified to have all
// 128 bits set to zero.
var SUIDNil = &SUID{make([]byte, 16)}

// Equal returns true if u1 and u2 equals, otherwise returns false.
func Equal(u1 *SUID, u2 *SUID) bool {
	return bytes.Equal(u1.value, u2.value)
}

// Version returns algorithm version used to generate UUID.
func (u *SUID) Version() uint {
	return uint(u.value[6] >> 4)
}

// Variant returns UUID layout variant.
func (u *SUID) Variant() uint {
	switch {
	case (u.value[8] & 0x80) == 0x00:
		return VariantNCS
	case (u.value[8]&0xc0)|0x80 == 0x80:
		return VariantRFC4122
	case (u.value[8]&0xe0)|0xc0 == 0xc0:
		return VariantMicrosoft
	}
	return VariantFuture
}

// Bytes returns bytes slice representation of UUID.
func (u *SUID) Bytes() []byte {
	return u.value
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u *SUID) StringFull() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u.value[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], u.value[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], u.value[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], u.value[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], u.value[10:])

	return string(buf)
}

func (this *SUID) String() string {
	return suidToShortS(this.value)
}

// SetVersion sets version bits.
func (this *SUID) SetVersion(v byte) {
	this.value[6] = (this.value[6] & 0x0f) | (v << 4)
}

// SetVariant sets variant bits as described in RFC 4122.
func (this *SUID) SetVariant() {
	this.value[8] = (this.value[8] & 0xbf) | 0x80
}

// Returns UUID v1/v2 storage state.
// Returns epoch timestamp, clock sequence, and hardware address.
func getStorage() (uint64, uint16, []byte) {
	storageOnce.Do(initStorage)

	storageMutex.Lock()
	defer storageMutex.Unlock()

	timeNow := unixTimeFunc()
	// Clock changed backwards since last UUID generation.
	// Should increase clock sequence.
	if timeNow <= lastTime {
		clockSequence++
	}
	lastTime = timeNow

	return timeNow, clockSequence, hardwareAddr[:]
}

// NewV1 returns UUID based on current timestamp and MAC address.
func NewSUID() *SUID {
	value := make([]byte, 16)
	this := SUID{value}

	t, q, h := getStorage()

	binary.BigEndian.PutUint32(value[0:], uint32(t))
	binary.BigEndian.PutUint16(value[4:], uint16(t>>32))
	binary.BigEndian.PutUint16(value[6:], uint16(t>>48))
	binary.BigEndian.PutUint16(value[8:], q)

	copy(this.value[10:], h)

	this.SetVersion(1)
	this.SetVariant()

	return &this
}

func WrapSUID(_buf []byte) *SUID {
	return &SUID{value: _buf}
}
