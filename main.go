package main

import (
	"net/http"
    "github.com/gin-gonic/gin"
	"log"
	"fmt"
	"os"
	"io"
	"encoding/json"
)
type Person struct {
	Name string
	Age string
}
func main() {
	r := gin.Default()
	r.GET("ping", func(c *gin.Context) {
		name := c.Query("name")
		c.String(http.StatusOK,"Hello %s",name)
	})
	r.POST("upload",posting)
	r.GET("/someGet",getting)
	r.POST("/form_post",forPost)
	r.Run(":80")
}

func forPost(c *gin.Context) {

	form := c.DefaultPostForm("name", "yida")
	postForm := c.DefaultPostForm("age", "20")
	p := Person{form,postForm}
	fmt.Println("p",p)
	data, err := json.Marshal(p)
	if err != nil{
		log.Print(err)
	}
	fmt.Printf("data:%s",data)
	c.JSON(200,gin.H{
		"message":"ok",
		"data":data,
	})

}
func getting(c *gin.Context) {
	name :=c.Query("name")
	c.String(http.StatusOK,"hello %s",name)
}
func posting(c *gin.Context){
	file, header, _ := c.Request.FormFile("file")
	log.Printf(header.Filename)
	out, err := os.Create("C://myuser//gin//" + header.Filename)
	if err !=nil{
		panic(err)
	}
	defer out.Close()
	io.Copy(out,file)
	c.String(http.StatusOK,fmt.Sprintf("上传成功:%s",header.Filename))
}