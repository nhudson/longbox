package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
}

func Search(log logrus.FieldLogger, url string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "POST" {
			http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
			return
		}

		c := colly.NewCollector()

		p := &pageInfo{Links: make(map[string]int)}

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Request.AbsoluteURL(e.Attr("href"))
			if link != "" {
				p.Links[link]++
			}
		})

		c.OnResponse(func(r *colly.Response) {
			log.Println("response received", r.StatusCode)
			p.StatusCode = r.StatusCode
		})

		c.OnError(func(r *colly.Response, err error) {
			log.Println("error:", r.StatusCode, err)
			p.StatusCode = r.StatusCode
		})

		c.Visit(url)

		b, err := json.Marshal(p)
		if err != nil {
			log.Errorln("failed to serialize response:", err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)

	})
}
