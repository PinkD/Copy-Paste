package cpst

import (
	"crypto/sha1"
	"fmt"
	_ "github.com/lib/pq"
	"math"
	"strings"
	"sync/atomic"
)

var encodeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
var charLen int
var encodeCharArray []string
var codeLen = 6

func NumberToChar(number uint64) (code string) {
	for bit := 0; bit < codeLen; bit++ {
		unit := number % uint64(charLen)
		code = encodeCharArray[unit] + code
		number /= uint64(charLen)
	}
	return
}

func CharToNumber(code string) (number uint64) {
	for _, c := range code {
		number *= uint64(charLen)
		index := indexInArray(string(c))
		if index == -1 {
			return 0
		} else {
			number += uint64(index)
		}
	}
	return
}

func indexInArray(char string) int {
	for i, ch := range encodeCharArray {
		if ch == char {
			return i
		}
	}
	return -1
}

func zeroCodeCount(code string) int {
	return strings.Count(code, encodeCharArray[0])
}

var initFlag = false

func initCharArray() {
	for _, char := range encodeChars {
		encodeCharArray = append(encodeCharArray, string(char))
	}
	charLen = len(encodeCharArray)
	initFlag = true
}

type codeGenerator struct {
	r        *redisClient
	db       *dB
	count    uint64
	maxCount uint64
}

func (g *codeGenerator) genCode() (uint64, error) {
	count := atomic.AddUint64(&g.count, 1) - 1
	if g.maxCount <= count {
		return g.maxCount, fmt.Errorf("Sorry, server can only store %d records\n", g.maxCount)
	}
	return count, nil
}

func (g *codeGenerator) sha1(data string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(data)))
}

func newCodeGenerator(redisAddr, dbAddr string) *codeGenerator {
	initCharArray()
	g := &codeGenerator{
		r:  newRedis(redisAddr),
		db: newDB(dbAddr),
	}
	g.maxCount = uint64(math.Pow(float64(charLen), float64(codeLen)))
	count, err := g.db.GetCount()
	g.count = count
	if err == nil {
		g.count++
		fmt.Printf("Record count: %d/%d, %.2f%% used\n", g.count, g.maxCount, float64(g.count)/float64(g.maxCount)*100)
	}
	return g
}

func (g *codeGenerator) save(sha, content string) (code uint64, err error) {
	code, err = g.r.ContainsContent(sha, content)
	if err == nil && code != 0 {
		return //in redis
	}
	code, err = g.db.ContainsContent(sha, content)
	if err == nil && code != 0 {
		data := &contentData{
			Code:    code,
			Sha:     sha,
			Content: content,
		}
		_ = g.r.SaveContent(data)
		return //in db
	}
	//create new
	code, err = g.genCode()
	if err != nil {
		return
	}
	data := &contentData{
		Code:    code,
		Sha:     sha,
		Content: content,
	}
	_ = g.r.SaveContent(data)
	err = g.db.SaveContent(data)
	return
}

func (g *codeGenerator) getContent(code uint64) (content string, err error) {
	content, err = g.r.GetContent(code)
	if err == nil && len(content) != 0 {
		return //in redis
	}
	content, err = g.db.GetContent(code)
	if err == nil && len(content) != 0 {
		_ = g.r.SaveContent(&contentData{
			Code:    code,
			Sha:     g.sha1(content),
			Content: content,
		})
	}
	return
}

func SetEncodeChars(chars string) {
	if initFlag {
		panic("Please call SetEncodeChars before initialization")
	}
	encodeChars = chars
}

//MaxInt64 is 9223372036854775807, bigger than 62^10
//So when len(encodeChars) is 62, max code len is 10
func SetCodeLen(length int) {
	if initFlag {
		panic("Please call SetCodeLen before initialization")
	}
	codeLen = length
}
