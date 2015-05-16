// shorty.go
package main

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/render"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var siteurl = "4pr.es/"
var db *sql.DB
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var coder Coder

type Coder struct {
	Length    uint
	Shrt, Url string
}

type Web struct {
	Get                    bool
	Proto, Banner, Content string
	Err                    error
}

func init() {
	var err error
	db, err = sql.Open("mysql", "shortener:passwd@/short")
	if err != nil {
		panic(err)
	}
	coder.Length = 6
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout:     "index",
		Extensions: []string{".tmpl", ".html"},
		Charset:    "UTF-8",
	}))
	//Landing page
	m.Get("/", func(req *http.Request, r render.Render) {
		s, ex := url.Parse(req.URL.String())
		if ex != nil {
			panic(ex)
		}
		fmt.Printf("URL SCHEME %s", s.Scheme)
		var land Web
		land.Proto = req.URL.Scheme
		land.Get = true
		land.Banner = "Get short URL for"
		land.Content = ""
		r.HTML(200, "get", land)
	})
	//Create entry for shortened URL
	m.Post("/", func(req *http.Request, r render.Render) {
		short, err := createUrl(req.FormValue("url"))
		var post Web
		post.Proto = req.URL.Scheme
		post.Get = false
		if err != nil {
			post.Banner = "Error :("
			post.Content = "Something did not work while trying to shorten URL " +
				req.FormValue("url") + "\n" + err.Error()
			post.Err = err
			r.HTML(500, "error", post)
		}
		post.Content = post.Proto + "://" + siteurl + short
		r.HTML(200, "post", post)
	})
	//Redirection to original URL
	m.Get("/:short", func(params martini.Params, w http.ResponseWriter, req *http.Request, r render.Render) {
		err := getUrl(params["short"], w, req)
		if err != nil {
			var e404 Web
			e404.Proto = req.URL.Scheme
			e404.Get = false
			e404.Err = err
			e404.Banner = "404 :("
			e404.Content = "URL " + params["short"] + " not found!"
			r.HTML(404, "error", e404)
		}
	})
	fmt.Println("Shortening URLS on localhost:1337 (HTTPS) and localhost:4337 (HTTPS)")
	http.Handle("/", m)
	http.ListenAndServe(":1337")
	http.ListenAndServeTLS(":4337", "/home/inge/dev/go/src/s4pres/tls/server.crt", "/home/inge/dev/go/src/s4pres/tls/server.key")
}

func createUrl(input string) (string, error) {
	coder.Url = input
	coder.Shrt = shorten(coder.Length)
	for urlPresent(coder.Shrt) {
		coder.Shrt = shorten(coder.Length)
	}
	_, err := db.Exec("INSERT INTO urls VALUES (null, ?, ?, null)", coder.Url, coder.Shrt)
	if err != nil {
		return "", err
	}
	return coder.Shrt, nil
}

func urlPresent(url string) bool {
	var s sql.NullString
	db.QueryRow("SELECT short from urls WHERE short = ?", coder.Shrt).Scan(&s)
	if s.Valid {
		return true
	}
	return false
}

func getUrl(short string, w http.ResponseWriter, req *http.Request) error {
	var redir string
	err := db.QueryRow("SELECT url FROM urls WHERE short = ?", short).Scan(&redir)
	if err != nil {
		return err
	}
	http.Redirect(w, req, redir, 301)
	return nil
}

func shorten(c uint) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, c)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
