
---
weight: 11
title: pong
---
# pong

## Pong.Pong
```go
package main
import (
  "github.com/micro/clients/go/client"
  pong_proto "github.com/micro/services/pong/proto"
)
func main() {
  c := client.NewClient(nil)
  req := pong_proto.Request{}
  rsp := pong_proto.Response{}
  if err := c.Call("go.micro.srv.pong", "Pong.Pong", req, &rsp); err != nil {
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
      .call("go.micro.srv.pong", "Pong.Pong", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```
 Pong returns you a pong.
### Request Parameters
Name |  Type | Description
--------- | --------- | ---------

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
pong | string |  Pong will just contain the string literal "pong"


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

