package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type GithubClient struct {
	HiddenDocs []string
}

func NewGithubClient() *GithubClient {
	return &GithubClient{
		HiddenDocs: []string{"README.md", "todo.md"},
	}
}

type ContentsResponse []struct {
	Name string `json:"name"`
}

func (g *GithubClient) GetDocs(prefix string) ([]string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/craftserve/docs/contents", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("getDocs http.DefaultClient.Do(req) " + err.Error())
		return nil, err
	}

	var contents ContentsResponse
	err = json.NewDecoder(resp.Body).Decode(&contents)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println("getDocs resp.Body.Close() ", err)
	}

	var docs []string

DOCS:
	for _, doc := range contents {
		if !strings.HasSuffix(doc.Name, ".md") {
			continue
		}
		if !strings.HasPrefix(doc.Name, prefix) {
			continue
		}
		for _, hiddenDoc := range g.HiddenDocs {
			if doc.Name == hiddenDoc {
				continue DOCS
			}
		}

		docs = append(docs, doc.Name[:len(doc.Name)-3])
	}

	return docs, nil
}

func (g *GithubClient) GetDocExists(name string) (bool, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/craftserve/docs/contents/"+name+".md", nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("getDocExists http.DefaultClient.Do(req) " + err.Error())
		return false, err
	}

	return resp.StatusCode == 200, nil
}
