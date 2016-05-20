package util

import (
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func GetResponseHeader(url, key string, params map[string]string) (string, error) {
	client := httpClient()
	response, err := client.PostForm(url, mapToUrlValues(params))

	if err != nil {
		return "", err
	}

	header := response.Header.Get("Location")

	if header == "" {
		return "", errors.New("Header: " + key + " not found")
	}

	return header, nil
}

func getElementBy(attname string, id string, n *html.Node) (element *html.Node, ok bool) {
	for _, a := range n.Attr {
		if a.Key == attname && a.Val == id {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = getElementBy(attname, id, c); ok {
			return
		}
	}
	return
}

func GetResponseForm(urlbase string, params map[string]string, authparams map[string]string) (*http.Client, string, error) {
	client := httpClient()
	loginurl, _ := url.Parse(urlbase + "/login")
	parameters := url.Values{}
	parameters.Add("service", params["service"])
	loginurl.RawQuery = parameters.Encode()

	// First: GET login page to catch LT and execution
	response, err := client.Get(loginurl.String())
	//fmt.Println(response)
	if err != nil {
		return client, "", err
	}

	if response.StatusCode != 200 {
		errMsg := fmt.Sprintf("response should be 200 but is: %d", response.StatusCode)
		return client, "", errors.New(errMsg)
	}

	body := response.Body
	defer response.Body.Close()

	LT, exec := "", ""
	doc, _ := html.Parse(body)

	element, ok := getElementBy("name", "lt", doc)
	if ok {
		for _, a := range element.Attr {
			if a.Key == "value" {
				LT = a.Val
				continue
			}
		}
	}

	element, ok = getElementBy("name", "execution", doc)
	if ok {
		for _, a := range element.Attr {
			if a.Key == "value" {
				exec = a.Val
				continue
			}
		}
	}

	// Second: POST auth form
	authparams["lt"] = LT
	authparams["service"] = params["service"]
	authparams["auto"] = "true"
	authparams["_eventId"] = "submit"
	authparams["execution"] = exec
	response2, err := client.PostForm(loginurl.String(), mapToUrlValues(authparams))
	if err != nil {
		return client, "", err
	}
	if response2.StatusCode != 200 {
		errMsg := fmt.Sprintf("response should be 200 but is: %d", response2.StatusCode)
		return client, "", errors.New(errMsg)
	}

	//	fmt.Println(response2)

	// Third: on success, client is redirected to test service which put User header
	header := response2.Header.Get("User")
	//	fmt.Println(client.Jar.Cookies(loginurl))
	/*	body2, _ := ioutil.ReadAll(response2.Body)
		defer response2.Body.Close()
		fmt.Println(string(body2))*/
	if header == "" {
		return client, "", errors.New("Header: User not found")
	}

	return client, header, nil
}

func GetResponseBody(url string, params map[string]string) (string, error) {
	client := httpClient()
	response, err := client.PostForm(url, mapToUrlValues(params))
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		errMsg := fmt.Sprintf("response should be 200 but is: %d", response.StatusCode)
		return "", errors.New(errMsg)
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	return string(body), nil
}

func httpClient() *http.Client {
	cookieJar, _ := cookiejar.New(nil)
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	return &http.Client{Transport: transport, Jar: cookieJar}
}

func mapToUrlValues(hash map[string]string) url.Values {
	values := url.Values{}

	for key, value := range hash {
		values.Add(key, value)
	}

	return values
}
