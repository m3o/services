
---
weight: 11
title: helloworld
---
# helloworld

## Helloworld.Call
```go
package main
import (
  "github.com/micro/clients/go/client"
  helloworld_proto "github.com/micro/services/helloworld/proto"
)
func main() {
  c := client.NewClient(nil)
  req := helloworld_proto.Request{}
  rsp := helloworld_proto.Response{}
  if err := c.Call("go.micro.srv.helloworld", "Helloworld.Call", req, &rsp); err != nil {
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
      .call("go.micro.srv.helloworld", "Helloworld.Call", {})
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

