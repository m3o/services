
---
weight: 11
title: sumo
---
# sumo

## Sumo.Call
```go
package main
import (
  "github.com/micro/clients/go/client"
  sumo_proto "github.com/micro/services/sumo/proto"
)
func main() {
  c := client.NewClient(nil)
  req := sumo_proto.Request{}
  rsp := sumo_proto.Response{}
  if err := c.Call("go.micro.srv.sumo", "Sumo.Call", req, &rsp); err != nil {
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
      .call("go.micro.srv.sumo", "Sumo.Call", {})
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

## Sumo.PingPong
```go
package main
import (
  "github.com/micro/clients/go/client"
  sumo_proto "github.com/micro/services/sumo/proto"
)
func main() {
  c := client.NewClient(nil)
  req := sumo_proto.Ping{}
  rsp := sumo_proto.Pong{}
  if err := c.Call("go.micro.srv.sumo", "Sumo.PingPong", req, &rsp); err != nil {
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
      .call("go.micro.srv.sumo", "Sumo.PingPong", {})
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

## Sumo.Stream
```go
package main
import (
  "github.com/micro/clients/go/client"
  sumo_proto "github.com/micro/services/sumo/proto"
)
func main() {
  c := client.NewClient(nil)
  req := sumo_proto.StreamingRequest{}
  rsp := sumo_proto.StreamingResponse{}
  if err := c.Call("go.micro.srv.sumo", "Sumo.Stream", req, &rsp); err != nil {
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
      .call("go.micro.srv.sumo", "Sumo.Stream", {})
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

