package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/koron/go-dproxy"
)

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a // a < b && a < c
		}
	} else if b < c {
		return b // a >= b && b < c
	}
	return c // (a < b && a >= c) || (a >= b && b >= c)
}

func calcEditDistance(str1, str2 string) int {
	d := make([][]int, len(str1)+1)
	for i := 0; i <= len(str1); i++ {
		d[i] = make([]int, len(str2)+1)
	}

	for i := 0; i <= len(str1); i++ {
		d[i][0] = i
	}
	for j := 0; j <= len(str2); j++ {
		d[0][j] = j
	}

	for i := 1; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			cost := 0
			if str1[i-1] != str2[j-1] {
				cost = 1
			}
			d[i][j] = min3(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}

	return d[len(str1)][len(str2)]
}

type Publ struct {
	Score                    int
	Title, Type, Url, BibUrl string
}

func getPublList(title string) ([]*Publ, error) {
	u, err := url.Parse("https://dblp.org/search/publ/api")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("format", "json")
	q.Set("q", title)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var js interface{}
	err = json.Unmarshal(body, &js)
	if err != nil {
		return nil, err
	}

	hits, err := dproxy.New(js).M("result").M("hits").M("hit").Array()
	if err != nil {
		return nil, err
	}

	publList := make([]*Publ, 0)
	minScore := 1000
	for _, hit := range hits {
		info, err := dproxy.New(hit).M("info").Map()
		if err != nil {
			return nil, err
		}
		publ := &Publ{
			0,
			info["title"].(string),
			info["type"].(string),
			info["url"].(string),
			"",
		}
		publ.Score = calcEditDistance(publ.Title, title)
		publ.BibUrl = fmt.Sprintf("%s", publ.Url) + ".bib"
		publList = append(publList, publ)

		if minScore > publ.Score {
			minScore = publ.Score
		}
	}

	// Select entries that have the min score
	minScorePubList := make([]*Publ, 0)
	for _, publ := range publList {
		if publ.Score == minScore {
			minScorePubList = append(minScorePubList, publ)
		}
	}

	// If len > 1, then remove entries with type "Informal Publications"
	ret := make([]*Publ, 0)
	if len(minScorePubList) == 1 {
		ret = minScorePubList
	} else {
		for _, publ := range minScorePubList {
			if publ.Type != "Informal Publications" {
				ret = append(ret, publ)
			}
		}
	}

	return ret, nil
}

func fetchBib(publ *Publ) (string, error) {
	resp, err := http.Get(publ.BibUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	reader := os.Stdin
	writer := os.Stdout

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		time.Sleep(500 * time.Millisecond)
		title := scanner.Text()
		publList, err := getPublList(title)
		if err != nil || len(publList) == 0 || publList[0].Score > 10 {
			fmt.Fprintf(os.Stderr, "Can't get: %s\n", title)
			continue
		}
		if len(publList) > 1 {
			fmt.Fprintf(os.Stderr, ">1 candidates: %s\n", title)
		}
		for _, publ := range publList {
			time.Sleep(500 * time.Millisecond)
			body, err := fetchBib(publ)
			if err != nil {
				log.Fatal(err)
			}
			writer.WriteString(body)
		}
	}
}
