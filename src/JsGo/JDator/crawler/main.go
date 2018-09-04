package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Frame 43: 836 bytes on wire (6688 bits), 836 bytes captured (6688 bits) on interface 0
// Ethernet II, Src: 9e:b6:aa:f7:72:b1 (9e:b6:aa:f7:72:b1), Dst: Raisecom_2c:bb:de (c8:50:e9:2c:bb:de)
// Internet Protocol Version 4, Src: 192.168.1.4, Dst: 203.76.216.1
// Transmission Control Protocol, Src Port: 61906, Dst Port: 80, Seq: 1, Ack: 1, Len: 782
// Hypertext Transfer Protocol
//     GET /shop/92678772 HTTP/1.1\r\n
//     Host: www.dianping.com\r\n
//     Connection: keep-alive\r\n
//     Cache-Control: max-age=0\r\n
//     Upgrade-Insecure-Requests: 1\r\n
//     User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36\r\n
//     Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\n
// Accept-Encoding: gzip, deflate\r\n
// Accept-Language: zh-CN,zh;q=0.9\r\n

func ExampleScrape() {
	// Request the HTML page.
	c1 := &http.Cookie{Name: "cy", Value: "1"}
	c2 := &http.Cookie{Name: "cityid", Value: "1"}
	c3 := &http.Cookie{Name: "cye", Value: "shanghai"}
	c5 := &http.Cookie{Name: "_lxsdk_cuid", Value: "162a4355c83c8-08fcce581305e2-b34356b-1fa400-162a4355c83c8"}
	c6 := &http.Cookie{Name: "_lxsdk", Value: "162a4355c83c8-08fcce581305e2-b34356b-1fa400-162a4355c83c8"}
	c7 := &http.Cookie{Name: "_hc.v", Value: "f925a499-dca1-077f-bad7-adcb6dbb0ef5.1523173383"}
	c8 := &http.Cookie{Name: "s_ViewType", Value: "10"}
	c9 := &http.Cookie{Name: "_lx_utm", Value: "utm_source%3DBaidu%26utm_medium%3Dorganic"}
	c10 := &http.Cookie{Name: "_lxsdk_s", Value: "162a46e50b7-aaf-299-ea4%7C%7C11"}

	reqest, err := http.NewRequest("GET", "http://www.dianping.com/shop/92678772", nil)

	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	reqest.Header.Set("Accept-Encoding", "gzip, deflate")
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	reqest.Header.Set("Upgrade-Insecure-Requests", "1")
	reqest.Header.Set("Cache-Control", "max-age=0")
	reqest.Header.Set("Connection", "keep-alive")

	reqest.AddCookie(c1)
	reqest.AddCookie(c2)
	reqest.AddCookie(c3)
	reqest.AddCookie(c5)
	reqest.AddCookie(c6)
	reqest.AddCookie(c7)
	reqest.AddCookie(c8)
	reqest.AddCookie(c9)
	reqest.AddCookie(c10)

	client := &http.Client{}
	res, err := client.Do(reqest)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("rep = %v\n", res)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)

	}

	fmt.Printf("%v", res.Body)

	// Find the review items
	doc.Find(".brief-info").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		// band := s.Find("a").Text()
		// title := s.Find("i").Text()
		// fmt.Printf("Review %d: %s - %s\n", i, band, title)
		fmt.Printf("%s\n", s.Find("div").Text())
	})
}

func main() {
	ExampleScrape()
}
