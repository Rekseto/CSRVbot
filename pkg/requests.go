package pkg

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

func GetCSRVCode() (string, error) {
	req, err := http.NewRequest("POST", "https://craftserve.pl/api/generate_voucher", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("csrvbot", os.Getenv("CSRV_SECRET"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("getCSRVCode http.DefaultClient.Do(req) " + err.Error())
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Code string `json:"code"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.Code, nil
}

func GetDocs(prefix string) ([]string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/craftserve/docs/contents", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("getDocs http.DefaultClient.Do(req) " + err.Error())
		return nil, err
	}

	var data []struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println("getDocs resp.Body.Close() ", err)
	}

	var docs []string
	hiddenDocs := []string{"README.md", "todo.md"}

DOCS:
	for _, doc := range data {
		if !strings.HasSuffix(doc.Name, ".md") {
			continue
		}
		if !strings.HasPrefix(doc.Name, prefix) {
			continue
		}
		for _, hiddenDoc := range hiddenDocs {
			if doc.Name == hiddenDoc {
				continue DOCS
			}
		}

		docs = append(docs, doc.Name[:len(doc.Name)-3])
	}

	return docs, nil
}

func GetDocExists(name string) (bool, error) {
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
