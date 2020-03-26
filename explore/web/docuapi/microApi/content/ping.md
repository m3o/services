
---
weight: 11
title: ping
---
# ping

## Ping.Ping
```go
package main
import (
  "github.com/micro/clients/go/client"
  ping_proto "github.com/micro/services/ping/proto"
)
func main() {
  c := client.NewClient(nil)
  req := ping_proto.Request{}
  rsp := ping_proto.Response{}
  if err := c.Call("go.micro.srv.ping", "Ping.Ping", req, &rsp); err != nil {
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
      .call("go.micro.srv.ping", "Ping.Ping", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
ping | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

