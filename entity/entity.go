package entity

import "encoding/xml"

type Sdn struct {
	XMLName  xml.Name `xml:"sdnList"`
	SdnEntry []struct {
		Uid       string `xml:"uid"`
		LastName  string `xml:"lastName"`
		FirstName string `xml:"firstName"`
		SdnType   string `xml:"sdnType"`
	} `xml:"sdnEntry"`
}

type Person struct {
	Uid       string `xml:"uid"`
	LastName  string `xml:"lastName"`
	FirstName string `xml:"firstName"`
}
