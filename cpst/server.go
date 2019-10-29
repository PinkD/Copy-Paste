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

type server struct {
	g *codeGenerator
	e *echo.Echo
}

func (s *server) index(c echo.Context) error {
	if strings.Contains(c.Request().UserAgent(), "curl") {
		return c.String(http.StatusOK, fmt.Sprintf("cat filename | curl -F \"content=<-\" %s/new\n", c.Request().Host))
	}
	return c.Render(http.StatusOK, "index", map[string]string{
		"host": c.Request().Host,
	})
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
		return c.String(http.StatusInternalServerError, fmt.Sprintf("fail to save content: %s", err.Error()))
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
	s := &server{
		e: echo.New(),
		g: newCodeGenerator(redisAddr, dbAddr),
	}
	s.e.GET("/", s.index)
	s.e.Static("/js", "resources")
	s.e.Static("/css", "resources")
	s.e.File("/favicon.ico", "resources/favicon.ico")
	s.e.GET("/:code/:highlight", s.getContent)
	s.e.GET("/:code", s.getContent)
	s.e.POST("/new", s.newContent)
	render := &templateRender{
		templates: template.Must(template.ParseGlob("resources/*.html")),
	}
	s.e.Renderer = render
	return s
}
