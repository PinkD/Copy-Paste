package cpst

import (
	"crypto/sha1"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

var encodeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
var charLen int
var encodeCharArray []string

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

func initCharArray() {
	for _, char := range encodeChars {
		encodeCharArray = append(encodeCharArray, string(char))
	}
	charLen = len(encodeCharArray)
}

type codeGenerator struct {
	r  *redisClient
	db *dB
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
	code := g.db.GetCount()
	if code != 0 {
		g.r.setCount(code + 1)
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
	code, err = g.r.genCode()
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
	encodeChars = chars
	initCharArray()
	fmt.Println(charLen)
}
