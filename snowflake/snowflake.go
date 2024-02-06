// Package snowflake provides a very simple Twitter snowflake generator and parser.
package snowflake

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
	// You may customize this to set a different epoch for your application.
	// Epoch int64 = 1288834974657
	Epoch int64 = 1704038400000 // 2024-01-01 00:00:00

	ShortId bool

	// NodeBits holds the number of bits to use for Node
	// Remember, you have a total 22 bits to share between Node/Step
	NodeBits uint8 = 10

	// StepBits holds the number of bits to use for Step
	// Remember, you have a total 22 bits to share between Node/Step
	StepBits uint8 = 12

	// DEPRECATED: the below four variables will be removed in a future release.
	mu        sync.Mutex
	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask        = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift       = NodeBits + StepBits
	nodeShift       = StepBits
)

const encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"

var decodeBase32Map [256]byte

const encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var decodeBase58Map [256]byte

// A JSONSyntaxError is returned from UnmarshalJSON if an invalid ID is provided.
type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid snowflake ID %q", string(j.original))
}

// ErrInvalidBase58 is returned by ParseBase58 when given an invalid []byte
var ErrInvalidBase58 = errors.New("invalid base58")

// ErrInvalidBase32 is returned by ParseBase32 when given an invalid []byte
var ErrInvalidBase32 = errors.New("invalid base32")

// Create maps for decoding Base58/Base32.
// This speeds up the process tremendously.
func init() {

	for i := 0; i < len(decodeBase58Map); i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}

	for i := 0; i < len(decodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

// A Node struct holds the basic information needed for a snowflake generator
// node
type Node struct {
	mu        sync.Mutex
	epoch     time.Time
	time      int64
	shortTime int64
	node      int64
	step      int64
	shortStep int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

// An ID is a custom type used for a snowflake ID.  This is used so we can
// attach methods onto the ID.
type ID struct {
	id      int64
	shortId bool
}

// NewNode returns a new snowflake node that can be used to generate snowflake
// IDs
func NewNode(node uint16) (*Node, error) {
	if node > 1023 {
		return nil, errors.New("Node set between 0 and 1023")
	}
	if NodeBits+StepBits > 22 {
		return nil, errors.New("Remember, you have a total 22 bits to share between Node/Step")
	}
	// re-calc in case custom NodeBits or StepBits were set
	// DEPRECATED: the below block will be removed in a future release.
	mu.Lock()
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits
	mu.Unlock()

	n := Node{}
	n.node = int64(node)
	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits
	n.nodeShift = StepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

func (n *Node) getTime() int64 {
	if ShortId {
		return int64(time.Since(n.epoch).Seconds())
	} else {
		return time.Since(n.epoch).Milliseconds()
	}
}

func (n *Node) getPreTime() int64 {
	if ShortId {
		return n.shortTime
	} else {
		return n.time
	}
}

func (n *Node) setPreTime(time int64) {
	if ShortId {
		n.shortTime = time
	} else {
		n.time = time
	}
}

func (n *Node) getStep() int64 {
	if ShortId {
		return n.shortStep
	} else {
		return n.step
	}
}
func (n *Node) setStep(step int64) {
	if ShortId {
		n.shortStep = step
	} else {
		n.step = step
	}
}

func (n *Node) ShortId() ID {
	ShortId = true
	return n.generate()
}

func (n *Node) NextId() ID {
	ShortId = false
	return n.generate()
}

// Generate creates and returns a unique snowflake ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node ID
func (n *Node) generate() ID {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := n.getTime()
	preTime := n.getPreTime()
	step := n.getStep()

	if now == preTime {
		step = (step + 1) & n.stepMask

		if step == 0 {
			for now <= preTime {
				now = n.getTime()
			}
		}
	} else {
		step = 0
	}

	n.setPreTime(now)
	n.setStep(step)

	id := ID{
		id:      (now)<<n.timeShift | (n.node << n.nodeShift) | (step),
		shortId: ShortId,
	}
	return id
}

func (f ID) WithShortId() ID {
	return ID{id: f.id, shortId: true}
}

// Int64 returns an int64 of the snowflake ID
func (f ID) Int64() int64 {
	return int64(f.id)
}

// ParseInt64 converts an int64 into a snowflake ID
func ParseInt64(id int64) ID {
	return ID{id: id}
}

// String returns a string of the snowflake ID
func (f ID) String() string {
	return strconv.FormatInt(int64(f.id), 10)
}

// ParseString converts a string into a snowflake ID
func ParseString(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return ID{id: i}, err

}

// Base2 returns a string base2 of the snowflake ID
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f.id), 2)
}

// ParseBase2 converts a Base2 string into a snowflake ID
func ParseBase2(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 2, 64)
	return ID{id: i}, err
}

// Base32 uses the z-base-32 character set but encodes and decodes similar
// to base58, allowing it to create an even smaller result string.
// NOTE: There are many different base32 implementations so becareful when
// doing any interoperation.
func (f ID) Base32() string {

	if f.id < 32 {
		return string(encodeBase32Map[f.id])
	}

	b := make([]byte, 0, 12)
	for f.id >= 32 {
		b = append(b, encodeBase32Map[f.id%32])
		f.id /= 32
	}
	b = append(b, encodeBase32Map[f.id])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase32 parses a base32 []byte into a snowflake ID
// NOTE: There are many different base32 implementations so becareful when
// doing any interoperation.
func ParseBase32(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase32Map[b[i]] == 0xFF {
			return ID{id: -1}, ErrInvalidBase32
		}
		id = id*32 + int64(decodeBase32Map[b[i]])
	}

	return ID{id: id}, nil
}

// Base36 returns a base36 string of the snowflake ID
func (f ID) Base36() string {
	return strconv.FormatInt(int64(f.id), 36)
}

// ParseBase36 converts a Base36 string into a snowflake ID
func ParseBase36(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 36, 64)
	return ID{id: i}, err
}

// Base58 returns a base58 string of the snowflake ID
func (f ID) Base58() string {

	if f.id < 58 {
		return string(encodeBase58Map[f.id])
	}

	b := make([]byte, 0, 11)
	for f.id >= 58 {
		b = append(b, encodeBase58Map[f.id%58])
		f.id /= 58
	}
	b = append(b, encodeBase58Map[f.id])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase58 parses a base58 []byte into a snowflake ID
func ParseBase58(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase58Map[b[i]] == 0xFF {
			return ID{id: -1}, ErrInvalidBase58
		}
		id = id*58 + int64(decodeBase58Map[b[i]])
	}

	return ID{id: id}, nil
}

// Base64 returns a base64 string of the snowflake ID
func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

// ParseBase64 converts a base64 string into a snowflake ID
func ParseBase64(id string) (ID, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return ID{id: -1}, err
	}
	return ParseBytes(b)

}

// Bytes returns a byte slice of the snowflake ID
func (f ID) Bytes() []byte {
	return []byte(f.String())
}

// ParseBytes converts a byte slice into a snowflake ID
func ParseBytes(id []byte) (ID, error) {
	i, err := strconv.ParseInt(string(id), 10, 64)
	return ID{id: i}, err
}

// IntBytes returns an array of bytes of the snowflake ID, encoded as a
// big endian integer.
func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f.id))
	return b
}

// ParseIntBytes converts an array of bytes encoded as big endian integer as
// a snowflake ID
func ParseIntBytes(id [8]byte) ID {
	return ID{id: int64(binary.BigEndian.Uint64(id[:]))}
}

// Time returns an int64 unix timestamp in milliseconds of the snowflake ID time
// DEPRECATED: the below function will be removed in a future release.
func (f ID) Timestamp() int64 {
	if f.shortId {
		return (int64(f.id) >> timeShift) + int64(Epoch/1000)
	} else {
		return (int64(f.id) >> timeShift) + Epoch
	}
}

func (f ID) Time() time.Time {
	if f.shortId {
		return time.Unix(f.Timestamp(), 0)
	} else {
		return time.UnixMilli(f.Timestamp())
	}
}

func (f ID) TimeFormat(layout string) string {
	return f.Time().Format(layout)
}

// Node returns an int64 of the snowflake ID node number
// DEPRECATED: the below function will be removed in a future release.
func (f ID) Node() int64 {
	return int64(f.id) & nodeMask >> nodeShift
}

// Step returns an int64 of the snowflake step (or sequence) number
// DEPRECATED: the below function will be removed in a future release.
func (f ID) Step() int64 {
	return int64(f.id) & stepMask
}

// MarshalJSON returns a json byte array string of the snowflake ID.
func (f ID) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	buff = append(buff, '"')
	buff = strconv.AppendInt(buff, int64(f.id), 10)
	buff = append(buff, '"')
	return buff, nil
}

// UnmarshalJSON converts a json byte array of a snowflake ID into an ID type.
func (f *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		return JSONSyntaxError{b}
	}

	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}

	*f = ID{id: i}
	return nil
}
