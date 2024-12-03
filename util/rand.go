package util

import (
	"math/rand"
	"time"
)

// 数字 + 大小写字母
const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numBytes = "0123456789"
const lowerBytes = "abcdefghijklmnopqrstuvwxyz"
const upperBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	alphabetIdxBits = 6
	alphabetIdxMask = 1<<alphabetIdxBits - 1
	alphabetIdxMax  = 63 / alphabetIdxBits
	numIdxBits      = 4
	numIdxMask      = 1<<numIdxBits - 1
	numIdxMax       = 63 / numIdxBits
	lowerIdxBits    = 5
	lowerIdxMask    = 1<<lowerIdxBits - 1
	lowerIdxMax     = 63 / lowerIdxBits
	upperIdxBits    = 5
	upperIdxMask    = 1<<upperIdxBits - 1
	upperIdxMax     = 63 / upperIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandGroup(p ...uint32) int {
	if p == nil {
		panic("args not found")
	}

	r := make([]uint32, len(p))
	for i := 0; i < len(p); i++ {
		if i == 0 {
			r[0] = p[0]
		} else {
			r[i] = r[i-1] + p[i]
		}
	}

	rl := r[len(r)-1]
	if rl == 0 {
		return 0
	}

	rn := uint32(rand.Int63n(int64(rl)))
	for i := 0; i < len(r); i++ {
		if rn < r[i] {
			return i
		}
	}

	panic("bug")
}

func RandInterval(b1, b2 int32) int32 {
	if b1 == b2 {
		return b1
	}

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	return int32(rand.Int63n(max-min+1) + min)
}

func RandIntervalN(b1, b2 int32, n uint32) []int32 {
	if b1 == b2 {
		return []int32{b1}
	}

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	l := max - min + 1
	if int64(n) > l {
		n = uint32(l)
	}

	r := make([]int32, n)
	m := make(map[int32]int32)
	for i := uint32(0); i < n; i++ {
		v := int32(rand.Int63n(l) + min)

		if mv, ok := m[v]; ok {
			r[i] = mv
		} else {
			r[i] = v
		}

		lv := int32(l - 1 + min)
		if v != lv {
			if mv, ok := m[lv]; ok {
				m[v] = mv
			} else {
				m[v] = lv
			}
		}

		l--
	}

	return r
}

// RandString 随机字符串（数字 + 大小写字母）
func RandString(n ...int) string {
	num := 64
	if len(n) > 0 {
		num = n[0]
	}
	b := make([]byte, num)
	for i, cache, remain := num-1, src.Int63(), alphabetIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), alphabetIdxMax
		}
		if idx := int(cache & alphabetIdxMask); idx < len(alphabet) {
			b[i] = alphabet[idx]
			i--
		}
		cache >>= alphabetIdxBits
		remain--
	}
	return Bytes2String(b)
}

// RandStrWithNum 随机生成数字字符串
func RandStrWithNum(n ...int) string {
	num := 64
	if len(n) > 0 {
		num = n[0]
	}
	b := make([]byte, num)
	for i, cache, remain := num-1, src.Int63(), numIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), numIdxMax
		}
		if idx := int(cache & numIdxMask); idx < len(numBytes) {
			b[i] = numBytes[idx]
			i--
		}
		cache >>= numIdxBits
		remain--
	}
	return Bytes2String(b)
}

// RandStrWithLower 随机生成小写字母字符串
func RandStrWithLower(n ...int) string {
	num := 64
	if len(n) > 0 {
		num = n[0]
	}
	b := make([]byte, num)
	for i, cache, remain := num-1, src.Int63(), lowerIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), lowerIdxMax
		}
		if idx := int(cache & lowerIdxMask); idx < len(lowerBytes) {
			b[i] = lowerBytes[idx]
			i--
		}
		cache >>= lowerIdxBits
		remain--
	}
	return Bytes2String(b)
}

// RandStrWithUpper 随机生成大写字母字符串
func RandStrWithUpper(n ...int) string {
	num := 64
	if len(n) > 0 {
		num = n[0]
	}
	b := make([]byte, num)
	for i, cache, remain := num-1, src.Int63(), upperIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), upperIdxMax
		}
		if idx := int(cache & upperIdxMask); idx < len(upperBytes) {
			b[i] = upperBytes[idx]
			i--
		}
		cache >>= upperIdxBits
		remain--
	}
	return Bytes2String(b)
}
