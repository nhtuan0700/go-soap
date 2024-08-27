package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Example api: https://www.postman.com/cs-demo/public-soap-apis/request/rlue7u2/numbertowords?tab=body

type NumberToWordsRequest struct {
	NumberToWords NumberToWords `xml:"NumberToWords"`
}

type NumberToWords struct {
	UbiNum uint   `xml:"ubiNum"`
	Ns     string `xml:"xmlns,attr"`
}

func GetXmlBodyRequest(v any) (string, error) {
	type Envelope struct {
		XMLName xml.Name `xml:"soap:Envelope"`
		Soap    string   `xml:"xmlns:soap,attr"`
		Body    any      `xml:"soap:Body"`
	}
	body := Envelope{
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: v,
	}

	out, err := xml.MarshalIndent(body, " ", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal xml: %v", err)
	}

	return string(out), nil
}

type NumberToWordsResponse struct {
	Body struct {
		Response struct {
			NumberToWordsResult string `xml:"NumberToWordsResult"`
		} `xml:"NumberToWordsResponse"`
	}
}

func SoapCall(service string, request interface{}) string {
	body, err := GetXmlBodyRequest(request)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	response, err := client.Post(service, "text/xml", bytes.NewBufferString(body))

	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	content, _ := io.ReadAll(response.Body)
	s := strings.TrimSpace(string(content))
	return s
}

func main() {
	request := NumberToWordsRequest{
		NumberToWords: NumberToWords{
			Ns:     "http://www.dataaccess.com/webservicesserver/",
			UbiNum: 500,
		},
	}

	body, err := GetXmlBodyRequest(request)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("xml:\n", body)

	responseString := SoapCall("https://www.dataaccess.com/webservicesserver/NumberConversion.wso", request)
	fmt.Println("------")
	fmt.Println("Response: \n", responseString)

	var res NumberToWordsResponse
	err = xml.Unmarshal([]byte(responseString), &res)
	if err != nil {
		fmt.Printf("failed to unmarshal xml: %v", err)
		return
	}
	fmt.Println("-----")
	fmt.Printf("\n%+v\n", res)
}
