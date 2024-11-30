//package dictionary_api

package main

import (
	"errors"
	"fmt"
	"github.com/LaFei/dictionary-api/utils"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "ok")
}

func main() {
	router := gin.Default()
	//router.GET("/:word", SearchAWord)

	//SearchWord()

	err2 := utils.ReadSrtFilesInInputDirAndWriteWordCountsToCSV()
	if err2 != nil {
		return
	}

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}

func SearchAWord(c *gin.Context) {

	word := c.Param("word")

	GetWordFromDB(word)

	wikiHtml, err := GetWordFromWiki(word)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "not found from Wiki")
	}

	camDicHtml, err := GetWordFromCamDic(word)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "not found from Cam Dic")
	}

	SaveWordIntoDB(word, wikiHtml, camDicHtml)
	//CreateSearchRecord()
	//ParseSourcesForApiResponse()

	c.IndentedJSON(http.StatusOK, "ok")
}

func GetWordFromDB(word string) {

	result := utils.GetWordInNeo4jDb(word)

	if result == nil {
		return
	}
	fmt.Println(result["word"].(struct{}))

	//wikiIoReader := io.NopCloser(strings.NewReader(.(string)))
	//
	//ParseWordFromWiki(wikiIoReader)

}

func GetWordFromCamDic(word string) (string, error) {

	client := &http.Client{
		Transport: &http.Transport{},
	}

	// create HTTP request
	req, err := http.NewRequest("GET", "https://dictionary.cambridge.org/dictionary/english/"+word, nil)
	if err != nil {
		// Handle error
	}

	// set User-Agent header
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	)

	// make HTTP request
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// close the response body
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
	}

	camDicHtmlBuf := new(strings.Builder)
	n2, err2 := io.Copy(camDicHtmlBuf, response.Body)
	if err2 != nil {
		fmt.Println("Error:", err2)
	}

	fmt.Println("n2:", n2)

	camDicHtmlString := camDicHtmlBuf.String()

	return camDicHtmlString, nil
}

func GetWordFromWiki(word string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{},
	}

	// create HTTP request
	req, err := http.NewRequest("GET", "https://simple.wiktionary.org/wiki/"+word, nil)
	if err != nil {
		// Handle error
	}

	// set User-Agent header
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	)

	// make HTTP request
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// close the response body
	defer response.Body.Close()

	fmt.Printf("response.StatusCode: %v", response.StatusCode)

	if response.StatusCode != 200 {
		log.Println("status code error: %d %s", response.StatusCode, response.Status)
		return "", errors.New("Wiki record not found for: " + word)
	}

	wikiHtmlBuf := new(strings.Builder)
	n, err := io.Copy(wikiHtmlBuf, response.Body)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("n:", n)

	wikiHtmlString := wikiHtmlBuf.String()

	return wikiHtmlString, nil
}

func SaveWordIntoDB(word string, wikiHtml string, camDicHtml string) {

	fmt.Println("camDicHtmlString len:", len(wikiHtml), len(camDicHtml))
	utils.UpdateWordInNeo4jDb(word, wikiHtml, camDicHtml)
}

func CreateSearchRecord(word string) {

}

func ParseWordFromCamDic(html io.ReadCloser) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".pos.dpos").
		Each(func(i int, s *goquery.Selection) {
			content := s.Text()
			fmt.Printf("%d: %s\n", i, content)
		})
}

func ParseWordFromWiki(htmlReader io.ReadCloser) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("tr > td > p ").
		Each(func(i int, s *goquery.Selection) {
			content := s.Text()
			fmt.Printf("%d: %s\n", i, content)
		})
}

func ParseSourcesForApiResponse(html io.ReadCloser) {
	//ParseWordFromCamDic(html io.ReadCloser)
	//ParseWordFromWiki(html io.ReadCloser)
}
