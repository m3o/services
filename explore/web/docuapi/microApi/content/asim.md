
---
weight: 11
title: asim
---
# asim

## Asim.Call
```go
package main
import (
  "github.com/micro/clients/go/client"
  asim_proto "github.com/micro/services/asim/proto"
)
func main() {
  c := client.NewClient(nil)
  req := asim_proto.Request{}
  rsp := asim_proto.Response{}
  if err := c.Call("go.micro.srv.asim", "Asim.Call", req, &rsp); err != nil {
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
      .call("go.micro.srv.asim", "Asim.Call", {})
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

## Asim.PingPong
```go
package main
import (
  "github.com/micro/clients/go/client"
  asim_proto "github.com/micro/services/asim/proto"
)
func main() {
  c := client.NewClient(nil)
  req := asim_proto.Ping{}
  rsp := asim_proto.Pong{}
  if err := c.Call("go.micro.srv.asim", "Asim.PingPong", req, &rsp); err != nil {
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
      .call("go.micro.srv.asim", "Asim.PingPong", {})
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

## Asim.Stream
```go
package main
import (
  "github.com/micro/clients/go/client"
  asim_proto "github.com/micro/services/asim/proto"
)
func main() {
  c := client.NewClient(nil)
  req := asim_proto.StreamingRequest{}
  rsp := asim_proto.StreamingResponse{}
  if err := c.Call("go.micro.srv.asim", "Asim.Stream", req, &rsp); err != nil {
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
      .call("go.micro.srv.asim", "Asim.Stream", {})
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

