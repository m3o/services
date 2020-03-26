
---
weight: 11
title: apps
---
# apps

## Apps.Import
```go
package main
import (
  "github.com/micro/clients/go/client"
  apps_proto "github.com/micro/services/apps/proto"
)
func main() {
  c := client.NewClient(nil)
  req := apps_proto.ImportRequest{}
  rsp := apps_proto.ImportResponse{}
  if err := c.Call("go.micro.srv.apps", "Apps.Import", req, &rsp); err != nil {
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
      .call("go.micro.srv.apps", "Apps.Import", {})
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


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

