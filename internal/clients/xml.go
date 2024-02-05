package clients

import (
	"context"
	"encoding/xml"
	"net/http"
	"sdn_list/internal/services"
)

type XmlClient struct {
	url string
}

func NewXmlClient(url string) *XmlClient {
	return &XmlClient{url: url}
}

func (c *XmlClient) GetAll(ctx context.Context) (<-chan services.Person, <-chan services.MetaData, error) {

	response, err := http.Get(c.url)

	if err != nil {
		return nil, nil, err
	}

	resCh := make(chan services.Person, 10)
	metaCh := make(chan services.MetaData, 1)

	decoder := xml.NewDecoder(response.Body)

	go func(decoder *xml.Decoder) {
		defer func() {
			response.Body.Close()
			close(resCh)
			close(metaCh)
		}()

		for token, err := decoder.Token(); token != nil; token, err = decoder.Token() {
			if err != nil {
				continue
			}
			switch se := token.(type) {
			case xml.StartElement:
				if se.Name.Local == "sdnEntry" {
					sdnEntry := SdnEntry{}
					decoder.DecodeElement(&sdnEntry, &se)
					if sdnEntry.SdnType == "Individual" {
						person := Person{
							Uid:       sdnEntry.Uid,
							FirstName: sdnEntry.FirstName,
							LastName:  sdnEntry.LastName,
						}
						resCh <- convertPerson(person)
						for _, akaPerson := range sdnEntry.AkaList.Persons {
							akaPerson.Uid = sdnEntry.Uid
							resCh <- convertPerson(akaPerson)
						}
					}
					continue
				}
				if se.Name.Local == "publshInformation" {
					metaData := MetaData{}
					decoder.DecodeElement(&metaData, &se)
					metaCh <- convertMetaData(metaData)
				}
			}
		}
	}(decoder)

	return resCh, metaCh, nil
}

func convertPerson(person Person) services.Person {
	return services.Person{
		Uid:       person.Uid,
		FirstName: person.FirstName,
		LastName:  person.LastName,
	}
}

func convertMetaData(metaData MetaData) services.MetaData {
	return services.MetaData{
		PublishDate: metaData.PublishDate,
	}
}

type SdnEntry struct {
	Uid       int     `xml:"uid"`
	FirstName string  `xml:"firstName"`
	LastName  string  `xml:"lastName"`
	SdnType   string  `xml:"sdnType"`
	AkaList   AkaList `xml:"akaList"`
}

type AkaList struct {
	Persons []Person `xml:"aka"`
}

type Person struct {
	Uid       int    `xml:"uid"`
	FirstName string `xml:"firstName"`
	LastName  string `xml:"lastName"`
}

type MetaData struct {
	PublishDate string `xml:"Publish_Date"`
}
