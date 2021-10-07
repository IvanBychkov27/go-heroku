package main

import (
	"fmt"
	"io/ioutil"
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
	<p> <button type="submit">Отправить</button></p>
</form>
`
	button = `<p> <button type="submit">Отправить</button>`
	docEnd = `</body></html>`
)

func main() {
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
	fileName := "index.tmpl.html"

	nPlayers := c.Request.URL.Query().Get("nPlayers")

	fmt.Println("nPlayers =", nPlayers)

	saveFileHTML(fileName, nPlayers)
	c.HTML(200, fileName, nil)
}

func saveFileHTML(fileName, nPlayers string) {
	t := time.Now().Format("02.01.2006  15:04:05")
	if nPlayers == "" {
		nPlayers = "1"
	}

	dForm := strings.Replace(form, "{nPlayers}", nPlayers, 1)

	res := docStart + "<H1> Heroku time: " + t + " </H1> <br><br> <H3>Кол-во игроков: " + nPlayers + "</H3> <br> " + dForm + docEnd

	fileName = "templates/" + fileName
	err := ioutil.WriteFile(fileName, []byte(res), 0777)
	if err != nil {
		fmt.Println("error write file ", fileName)
	}
}
