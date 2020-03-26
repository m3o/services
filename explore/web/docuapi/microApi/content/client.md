
---
weight: 11
title: client
---
# client

## Client.Call
```go
package main
import (
  "github.com/micro/clients/go/client"
  client_proto "github.com/micro/services/client/proto"
)
func main() {
  c := client.NewClient(nil)
  req := client_proto.Request{}
  rsp := client_proto.Response{}
  if err := c.Call("go.micro.srv.client", "Client.Call", req, &rsp); err != nil {
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
      .call("go.micro.srv.client", "Client.Call", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```
 Call allows a single request to be made
### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
service | string | 
endpoint | string | 
content_type | string | 
body | bytes | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
body | bytes | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Client.Stream
```go
package main
import (
  "github.com/micro/clients/go/client"
  client_proto "github.com/micro/services/client/proto"
)
func main() {
  c := client.NewClient(nil)
  req := client_proto.Request{}
  rsp := client_proto.Response{}
  if err := c.Call("go.micro.srv.client", "Client.Stream", req, &rsp); err != nil {
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
      .call("go.micro.srv.client", "Client.Stream", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```
 Stream is a bidirectional stream
### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
service | string | 
endpoint | string | 
content_type | string | 
body | bytes | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
body | bytes | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

