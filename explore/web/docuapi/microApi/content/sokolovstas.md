
---
weight: 11
title: sokolovstas
---
# sokolovstas

## Sokolovstas.Call
```go
package main
import (
  "github.com/micro/clients/go/client"
  sokolovstas_proto "github.com/micro/services/sokolovstas/proto"
)
func main() {
  c := client.NewClient(nil)
  req := sokolovstas_proto.Request{}
  rsp := sokolovstas_proto.Response{}
  if err := c.Call("go.micro.srv.sokolovstas", "Sokolovstas.Call", req, &rsp); err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(rsp)
}
```
```javascript
// To install "npm install --save @microhq/ng-client"
import { Component, OnInit } from "@angular/core";
import { ClientService } from "@microhq/ng-client";
@Component({
  selector: "app-example",
  templateUrl: "./example.component.html",
  styleUrls: ["./example.component.css"]
})
export class ExampleComponent implements OnInit {
  constructor(private mc: ClientService) {}
  ngOnInit() {
    this.mc
      .call("go.micro.srv.sokolovstas", "Sokolovstas.Call", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
name | string | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
msg | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Sokolovstas.PingPong
```go
package main
import (
  "github.com/micro/clients/go/client"
  sokolovstas_proto "github.com/micro/services/sokolovstas/proto"
)
func main() {
  c := client.NewClient(nil)
  req := sokolovstas_proto.Ping{}
  rsp := sokolovstas_proto.Pong{}
  if err := c.Call("go.micro.srv.sokolovstas", "Sokolovstas.PingPong", req, &rsp); err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(rsp)
}
```
```javascript
// To install "npm install --save @microhq/ng-client"
import { Component, OnInit } from "@angular/core";
import { ClientService } from "@microhq/ng-client";
@Component({
  selector: "app-example",
  templateUrl: "./example.component.html",
  styleUrls: ["./example.component.css"]
})
export class ExampleComponent implements OnInit {
  constructor(private mc: ClientService) {}
  ngOnInit() {
    this.mc
      .call("go.micro.srv.sokolovstas", "Sokolovstas.PingPong", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
stroke | int64 | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
stroke | int64 | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Sokolovstas.Stream
```go
package main
import (
  "github.com/micro/clients/go/client"
  sokolovstas_proto "github.com/micro/services/sokolovstas/proto"
)
func main() {
  c := client.NewClient(nil)
  req := sokolovstas_proto.StreamingRequest{}
  rsp := sokolovstas_proto.StreamingResponse{}
  if err := c.Call("go.micro.srv.sokolovstas", "Sokolovstas.Stream", req, &rsp); err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(rsp)
}
```
```javascript
// To install "npm install --save @microhq/ng-client"
import { Component, OnInit } from "@angular/core";
import { ClientService } from "@microhq/ng-client";
@Component({
  selector: "app-example",
  templateUrl: "./example.component.html",
  styleUrls: ["./example.component.css"]
})
export class ExampleComponent implements OnInit {
  constructor(private mc: ClientService) {}
  ngOnInit() {
    this.mc
      .call("go.micro.srv.sokolovstas", "Sokolovstas.Stream", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
count | int64 | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
count | int64 | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

