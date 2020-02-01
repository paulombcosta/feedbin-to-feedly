package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
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

func main() {
	subscriptions := loadSubscriptions()
	fmt.Println(subscriptions)
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

func closeFile(f *os.File) {
	fmt.Println("closing file")
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		os.Exit(1)
	}
}
