
---
weight: 11
title: rex
---
# rex

## Rex.Call
```go
package main
import (
  "github.com/micro/clients/go/client"
  rex_proto "github.com/micro/services/rex/proto"
)
func main() {
  c := client.NewClient(nil)
  req := rex_proto.Request{}
  rsp := rex_proto.Response{}
  if err := c.Call("go.micro.srv.rex", "Rex.Call", req, &rsp); err != nil {
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
      .call("go.micro.srv.rex", "Rex.Call", {})
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

## Rex.PingPong
```go
package main
import (
  "github.com/micro/clients/go/client"
  rex_proto "github.com/micro/services/rex/proto"
)
func main() {
  c := client.NewClient(nil)
  req := rex_proto.Ping{}
  rsp := rex_proto.Pong{}
  if err := c.Call("go.micro.srv.rex", "Rex.PingPong", req, &rsp); err != nil {
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
      .call("go.micro.srv.rex", "Rex.PingPong", {})
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

## Rex.Stream
```go
package main
import (
  "github.com/micro/clients/go/client"
  rex_proto "github.com/micro/services/rex/proto"
)
func main() {
  c := client.NewClient(nil)
  req := rex_proto.StreamingRequest{}
  rsp := rex_proto.StreamingResponse{}
  if err := c.Call("go.micro.srv.rex", "Rex.Stream", req, &rsp); err != nil {
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
      .call("go.micro.srv.rex", "Rex.Stream", {})
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

