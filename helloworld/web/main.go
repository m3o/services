package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/micro/go-micro/v2/web"
	hello "github.com/micro/services/helloworld/proto/helloworld"
)

var (
	head = `<html>
<head>
  <style>
    body {
      margin: 25px;
      color: #0c2e47;
      font-family: medium-content-sans-serif-font,"Lucida Grande","Lucida Sans Unicode","Lucida Sans",Geneva,Arial,sans-serif;
    }
    .content {
      margin: 0 auto;
      max-width: 800px;
    }
    .content img {
      vertical-align: middle;
    }
    .title {
      font-weight: 600;
    }
    input {
      font-size: 20px;
      outline: 0;
      border: 0;
      border-bottom: 1px solid whitesmoke;
    }
  </style>
</head>
<body>
<div class="content">
<h1><img src="https://micro.mu/logo.png" width=50px height=auto /> <span class="title">Helloworld</span></h1>
`
	foot = `</div></body></html>`
	html = head + `<form method=post><input name=name type=text placeholder="Enter your name" autofocus></form>` + foot
)

func main() {
	service := web.NewService(
		web.Name("go.micro.web.helloworld"),
	)

	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()

			name := r.Form.Get("name")
			if len(name) == 0 {
				name = "World"
			}

			cli := service.Options().Service.Client()

			cl := hello.NewHelloworldService("go.micro.service.helloworld", cli)
			rsp, err := cl.Call(context.Background(), &hello.Request{
				Name: name,
			})

			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Write([]byte(head + `<h1>` + rsp.Msg + `</h1>` + foot))
			return
		}

		fmt.Fprint(w, html)
	})

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
