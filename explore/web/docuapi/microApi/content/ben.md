
---
weight: 11
title: ben
---
# ben

## Ben.Call
```go
package main
import (
  "github.com/micro/clients/go/client"
  ben_proto "github.com/micro/services/ben/proto"
)
func main() {
  c := client.NewClient(nil)
  req := ben_proto.Request{}
  rsp := ben_proto.Response{}
  if err := c.Call("go.micro.srv.ben", "Ben.Call", req, &rsp); err != nil {
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
      .call("go.micro.srv.ben", "Ben.Call", {})
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

## Ben.PingPong
```go
package main
import (
  "github.com/micro/clients/go/client"
  ben_proto "github.com/micro/services/ben/proto"
)
func main() {
  c := client.NewClient(nil)
  req := ben_proto.Ping{}
  rsp := ben_proto.Pong{}
  if err := c.Call("go.micro.srv.ben", "Ben.PingPong", req, &rsp); err != nil {
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
      .call("go.micro.srv.ben", "Ben.PingPong", {})
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

## Ben.Stream
```go
package main
import (
  "github.com/micro/clients/go/client"
  ben_proto "github.com/micro/services/ben/proto"
)
func main() {
  c := client.NewClient(nil)
  req := ben_proto.StreamingRequest{}
  rsp := ben_proto.StreamingResponse{}
  if err := c.Call("go.micro.srv.ben", "Ben.Stream", req, &rsp); err != nil {
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
      .call("go.micro.srv.ben", "Ben.Stream", {})
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

