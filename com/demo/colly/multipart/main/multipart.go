package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func generateFormDate() map[string][]byte{
	f,_:=os.Open("gocolly.jpg")
	defer f.Close()

	imgData,_:=ioutil.ReadAll(f)

	return map[string][]byte{
		"firstname":[]byte("one"),
		"lastname":[]byte("two"),
		"email":[]byte("onetwo@example.com"),
		"file":imgData,
	}
}

func setupServer()  {
	var handler http.HandlerFunc = func(w http.ResponseWriter,r *http.Request){
		fmt.Println("received request")
		err:=r.ParseMultipartForm(10000000)
		if err!=nil{
			fmt.Println("server:Error")
			w.WriteHeader(500)
			w.Write([]byte("<html><body>Internal Server Error</body></html>"))
			return
		}
		w.WriteHeader(200)
		fmt.Println("server:OK")
		w.Write([]byte("<html><body>Sucess</body></html>"))
	}

	go http.ListenAndServe(":8081",handler)
}

func main(){
	setupServer()

	c:=colly.NewCollector(colly.AllowURLRevisit(),colly.MaxDepth(5))

	c.OnHTML("html",func(e *colly.HTMLElement){
		fmt.Println(e.Text)
		time.Sleep(1*time.Second)
		e.Request.PostMultipart("http://localhost:8081",generateFormDate())
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Posting gocolly.jpg to ",request.URL.String())
	})

	c.PostMultipart("http://localhost:8081/",generateFormDate())
	c.Wait()
}
