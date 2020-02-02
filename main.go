package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	BASE_URL = "https://cloud.feedly.com/v3/"
)

type Opml struct {
	XMLName xml.Name `xml:"opml"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

type Head struct {
	XMLName xml.Name `xml:"head"`
	Title   string   `xml:"title"`
}

type Body struct {
	XMLName       xml.Name       `xml:"body"`
	Subscriptions []Subscription `xml:"outline"`
}

type Subscription struct {
	XMLName xml.Name `xml:"outline"`
	Text    string   `xml:"text,attr"`
	Title   string   `xml:"title,attr"`
	Type    string   `xml:"type,attr"`
	XmlUrl  string   `xml:"xmlUrl,attr"`
	HtmlUrl string   `xml:"htmlUrl,attr"`
}

type CreateFeed struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func main() {
	subscriptions := loadSubscriptions()
	writeSubscriptions(subscriptions)
}

func loadSubscriptions() []Subscription {

	file, err := os.Open("subscriptions.xml")

	if err != nil {
		fmt.Println(err)
	}

	defer closeFile(file)

	byteArr, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		os.Exit(1)
	}

	var opml Opml

	xml.Unmarshal(byteArr, &opml)

	return opml.Body.Subscriptions
}

func writeSubscriptions(subs []Subscription) {
	client := &http.Client{}
	devKey := os.Getenv("FEEDLY_DEVELOPER_KEY")

	data, err := json.Marshal(&CreateFeed{ID: "feed/" + subs[0].XmlUrl, Title: subs[0].Title})

	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", BASE_URL+"subscriptions", bytes.NewReader(data))

	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", "OAuth "+devKey)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode == 200 {
		fmt.Printf("Sucessfully added feed %s", subs[0].Title)
	} else {
		fmt.Printf("Request failed with status code %d", resp.StatusCode)
	}

}

func closeFile(f *os.File) {
	fmt.Println("closing file")
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		os.Exit(1)
	}
}
