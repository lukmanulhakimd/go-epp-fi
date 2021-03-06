package epp

import (
	"encoding/xml"
	"time"
)

type APIGreeting struct {
	XMLName  xml.Name `xml:"epp"`
	Host     string   `xml:"host,attr"`
	Domain   string   `xml:"domain,attr"`
	Contact  string   `xml:"contact,attr"`
	Obj      string   `xml:"obj,attr"`
	Xmlns    string   `xml:"xmlns,attr"`
	Greeting Greeting `xml:"greeting"`
}

type Greeting struct {
	SvID    string `xml:"svID"`
	RawDate string `xml:"svDate"`
	SvDate  time.Time
	SvcMenu struct {
		Version      string   `xml:"version"`
		Lang         string   `xml:"lang"`
		ObjURI       []string `xml:"objURI"`
		SvcExtension struct {
			ExtURI []string `xml:"extURI,omitempty"`
		} `xml:"svcExtension"`
	} `xml:"svcMenu"`
}
