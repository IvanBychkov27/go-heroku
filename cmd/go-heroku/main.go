package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

const (
	docStart = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Go-heroku</title>
<body>
`
	form = `
<form action="/" method="GET">
	<br>
	<p>Количество соперников:
	    <input type="number" name="nPlayers" value="{nPlayers}"  min="1" max="15" size="1" step="1">
	</p>
	<br>
	<p> <button type="submit">Раздать карты</button></p>
</form>
<br><br>
`
	button = `<p> <button type="submit">Отправить</button>`
	docEnd = `</body></html>`
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("$PORT must be set")
		return
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.GET("/", mainHandler)
	router.Run(":" + port)
}

func mainHandler(c *gin.Context) {
	nPlayers := c.Request.URL.Query().Get("nPlayers")
	data := buildResultData(nPlayers)

	fileName := "index.tmpl.html"
	saveFileHTML("templates/"+fileName, data)
	c.HTML(200, fileName, nil)
}

func buildResultData(nPlayers string) []byte {
	if nPlayers == "" {
		nPlayers = "1"
	}
	dForm := strings.Replace(form, "{nPlayers}", nPlayers, 1)
	res := docStart + "<H1> Heroku time: " + timeNow(3) + " </H1> <br><br> <H3>Кол-во игроков: " + nPlayers + "</H3> <br> " + dForm
	res += cardsImage(5) // вывод 5ти карт без повторов
	res += docEnd

	return []byte(res)
}

func saveFileHTML(fileName string, data []byte) {
	err := ioutil.WriteFile(fileName, data, 0777)
	if err != nil {
		fmt.Println("error write file ", fileName)
	}
}

func cardsImage(n int) string {
	var res string
	temp := make(map[string]struct{})
	for i := 0; i < n; i++ {
		card := randomCard()
		_, ok := temp[card]
		if !ok {
			res += imageData(card)
			temp[card] = struct{}{}
		} else {
			i--
		}
	}
	return res
}

func imageData(cardName string) string {
	openFileName := "cardImage/" + cardName + ".jpg"
	fileData, err := ioutil.ReadFile(openFileName)
	if err != nil {
		fmt.Println("error open file", err.Error())
		return openFileName
	}
	return `<img src="data:image/jpg; base64,` + base64.StdEncoding.EncodeToString(fileData) + `">`
}

func randomCard() string {
	cards := []string{
		"21", "22", "23", "24",
		"31", "32", "33", "34",
		"41", "42", "43", "44",
		"51", "52", "53", "54",
		"61", "62", "63", "64",
		"71", "72", "73", "74",
		"81", "82", "83", "84",
		"91", "92", "93", "94",
		"101", "102", "103", "104",
		"111", "112", "113", "114",
		"121", "122", "123", "124",
		"131", "132", "133", "134",
		"141", "142", "143", "144",
	}
	return cards[rand.Intn(len(cards))]
}

func timeNow(addHour int) string {
	y := time.Now().Year()
	mec := time.Now().Month()
	d := time.Now().Day()
	h := time.Now().Hour() + addHour
	m := time.Now().Minute()
	s := time.Now().Second()
	return time.Date(y, mec, d, h, m, s, 0, time.UTC).Format("02.01.2006  15:04:05")
}
