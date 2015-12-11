package util

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

type UUId int64

const base58 = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

func NewUUId(count int64) (u *UUId) {
	u = new(UUId)
	*u = UUId(count)
	return
}

func (u *UUId) Int64() int64 {
	return int64(*u)
}

func (u UUId) Base58() string {
	rest := u.Int64()
	result := ""

	for rest != 0 {
		result = result + string(base58[rest%58])
		rest /= 58
	}

	return result
}

func (u *UUId) String() string {
	return fmt.Sprintf("uuid=%d", u.Int64())
}

func GenIdFromBase58(input string) (*UUId, error) {
	sum := int64(0)
	base := len(input)
	for _, b := range []byte(input) {
		sum *= int64(base)
		index := bytes.IndexByte([]byte(base58), b)
		if index == -1 {
			return nil, fmt.Errorf("Invalid string.")
		}
		sum += int64(index)
	}
	return NewUUId(sum), nil
}

func IsValidBase58(input string) bool {
	for _, b := range []byte(input) {
		if bytes.IndexByte([]byte(base58), b) == -1 {
			return false
		}
	}

	return true
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

//func genUUID8() (u []byte, err error) {
//const size = 10
//if size <= 8 {
//err = fmt.Errorf("size %d is too small for uuid. Code would break.\n")
//return
//}
//// Taken from github.com/nu7hatch/gouuid
//const ReservedRFC4122 byte = 0x40
//u = make([]byte, size)

//// Set all bits to randomly (or pseudo-randomly) chosen values.
//_, err = rand.Read(u[:])
//if err != nil {
//return
//}

//// Set the two most significant bits (bits 6 and 7) of the
//// clock_seq_hi_and_reserved to zero and one, respectively.
//u[8] = (u[8] | ReservedRFC4122) & 0x7F

//// Set the four most significant bits (bits 12 through 15) of the
//// time_hi_and_version field to the 4-bit version number.
//u[6] = (u[6] & 0xF) | (4 << 4)

//return
//}
