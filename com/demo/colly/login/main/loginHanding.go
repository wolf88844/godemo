package main

import (
	"github.com/gocolly/colly"
	"log"
)

func main(){
	c:= colly.NewCollector()

	err:=c.Post("http://www.pupedu.cn/app/login/j_spring_security_check",map[string]string{"username":"admin","password":"admin"})
	if err != nil{
		log.Fatal(err)
	}

	c.OnResponse(func(r * colly.Response){
		log.Println("response received",r.StatusCode)
		log.Println(r.Body)
	})

	c.Visit("https://example.com")
}
