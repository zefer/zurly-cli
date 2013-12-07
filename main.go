package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"os"
	"io/ioutil"
)

const VERSION = "0.0.0"
const API_HOST = "zurly.herokuapp.com"

type shortenRequest struct {
	Url string
}

type shortenResponse struct {
	Id string
	LongUrl string
}

type errorResponse struct {
	Message string
}

func main() {
	app := cli.NewApp()
	app.Name = "zurly"
	app.Usage = "Shortens a url"
	app.Version = VERSION
	app.Action = func(c *cli.Context) {
		if len(c.Args()) < 1 {
			fmt.Println("Correct usage is: zurl http://long.url")
			return
		}
		longUrl := c.Args()[0]
		resp, err := shorten(longUrl)
		if err != nil {
			fmt.Println(err)
		} else {
			shortUrl := fmt.Sprintf("http://%s/%s", API_HOST, resp.Id)
			fmt.Println(shortUrl)
		}
	}

	app.Run(os.Args)
}

func shorten(longUrl string) (r shortenResponse, err error) {
	url := fmt.Sprintf("http://%s/", API_HOST)

	data := shortenRequest{Url: longUrl}
	urlJson, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(urlJson))
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", fmt.Sprintf("github.com/zefer/zurly-cli v%s", VERSION))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error shortening url: [%v]", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusCreated {
		var e errorResponse
		err = json.Unmarshal(body, &e)
		if err != nil {
			return
		}
		err = fmt.Errorf("Status: %d, Error: %v", resp.StatusCode, e.Message)
		return
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		err = fmt.Errorf("Error parsing API response: [%v]", err)
	}

	return
}
