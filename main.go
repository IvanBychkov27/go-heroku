package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

const doc = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Go-heroku</title>
`

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
	saveFileHTML(fileName)
	c.HTML(200, fileName, nil)
}

func saveFileHTML(fileName string) {
	t := time.Now().Format("02:04:06 15:04:05")
	res := doc + "<body><H1> Heroku time: " + t + " </H1></body></html>"

	fileName = "templates/" + fileName
	err := ioutil.WriteFile(fileName, []byte(res), 0777)
	if err != nil {
		fmt.Println("error write file ", fileName)
	}
}
