
---
weight: 11
title: vtolstov
---
# vtolstov

## Vtolstov.Call
```go
package main
import (
  "github.com/micro/clients/go/client"
  vtolstov_proto "github.com/micro/services/vtolstov/proto"
)
func main() {
  c := client.NewClient(nil)
  req := vtolstov_proto.Request{}
  rsp := vtolstov_proto.Response{}
  if err := c.Call("go.micro.srv.vtolstov", "Vtolstov.Call", req, &rsp); err != nil {
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
      .call("go.micro.srv.vtolstov", "Vtolstov.Call", {})
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

## Vtolstov.PingPong
```go
package main
import (
  "github.com/micro/clients/go/client"
  vtolstov_proto "github.com/micro/services/vtolstov/proto"
)
func main() {
  c := client.NewClient(nil)
  req := vtolstov_proto.Ping{}
  rsp := vtolstov_proto.Pong{}
  if err := c.Call("go.micro.srv.vtolstov", "Vtolstov.PingPong", req, &rsp); err != nil {
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
      .call("go.micro.srv.vtolstov", "Vtolstov.PingPong", {})
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

## Vtolstov.Stream
```go
package main
import (
  "github.com/micro/clients/go/client"
  vtolstov_proto "github.com/micro/services/vtolstov/proto"
)
func main() {
  c := client.NewClient(nil)
  req := vtolstov_proto.StreamingRequest{}
  rsp := vtolstov_proto.StreamingResponse{}
  if err := c.Call("go.micro.srv.vtolstov", "Vtolstov.Stream", req, &rsp); err != nil {
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
      .call("go.micro.srv.vtolstov", "Vtolstov.Stream", {})
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

