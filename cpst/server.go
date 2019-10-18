package cpst

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/bytes"
	"html/template"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

const maxContentLen = 256 * bytes.KB

var codeLen = 6

type server struct {
	g *codeGenerator
	e *echo.Echo
}

func (s *server) index(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("Usage:\n    cat filename | curl -F \"content=<-\" %s/new\n", c.Request().Host))
}

func (s *server) getContent(c echo.Context) error {
	code := c.Param("code")
	l := len(code)
	if len(code) != codeLen {
		l = utf8.RuneCountInString(code)
		if utf8.RuneCountInString(code) != codeLen {
			return echo.ErrNotFound
		} //else code is utf8
	}
	highlight := c.Param("highlight")
	if len(highlight) == 0 {
		highlight = "text"
	}
	codeNumber := CharToNumber(code)
	if codeNumber == 0 && zeroCodeCount(code) != l {
		return echo.ErrNotFound
	}
	content, err := s.g.getContent(codeNumber)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "fail to get content")
	}
	if len(content) == 0 { //not found
		c.Logger().Errorf("got empty content for code %s", code)
		return echo.ErrNotFound
	}
	return c.Render(http.StatusOK, "content", map[string]string{
		"highlight": highlight,
		"content":   content,
	})
}

func (s *server) newContent(c echo.Context) error {
	content := c.FormValue("content")
	if len(content) == 0 {
		return c.String(http.StatusBadRequest, "content required\n")
	}
	if len(strings.TrimSpace(content)) == 0 {
		return c.String(http.StatusBadRequest, "content is empty\n")
	}
	if len(strings.TrimSpace(content)) > maxContentLen {
		return c.String(http.StatusBadRequest, "content is too long\n")
	}
	sha := s.g.sha1(content)
	code, err := s.g.save(sha, content)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "fail to save content")
	}
	result := fmt.Sprintf("%s/%s\n", c.Request().Host, NumberToChar(code))
	return c.String(http.StatusOK, result)
}

func (s *server) Start(addr string) error {
	return s.e.Start(addr)
}

type templateRender struct {
	templates *template.Template
}

func (t *templateRender) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewServer(redisAddr, dbAddr string) *server {
	server := &server{
		e: echo.New(),
		g: newCodeGenerator(redisAddr, dbAddr),
	}
	server.e.GET("/", server.index)
	server.e.Static("/js", "resources")
	server.e.Static("/css", "resources")
	server.e.File("/favicon.ico", "resources/favicon.ico")
	server.e.GET("/:code/:highlight", server.getContent)
	server.e.GET("/:code", server.getContent)
	server.e.POST("/new", server.newContent)
	render := &templateRender{
		templates: template.Must(template.ParseFiles("resources/content.html")),
	}
	server.e.Renderer = render
	return server
}

//MaxInt64 is 9223372036854775807, bigger than 62^10
//So when len(encodeChars) is 62, max code len is 10
func SetCodeLen(length int) {
	codeLen = length
}
