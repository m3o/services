
---
weight: 11
title: platform-test
---
# platform-test

## PlatformTest.GetHealth
```go
package main
import (
  "github.com/micro/clients/go/client"
  platform_test_proto "github.com/micro/services/platform-test/proto"
)
func main() {
  c := client.NewClient(nil)
  req := platform_test_proto.GetHealthRequest{}
  rsp := platform_test_proto.GetHealthResponse{}
  if err := c.Call("go.micro.srv.platform-test", "PlatformTest.GetHealth", req, &rsp); err != nil {
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
      .call("go.micro.srv.platform-test", "PlatformTest.GetHealth", {})
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
auth | HealthStatus | 
broker | HealthStatus | 
config | HealthStatus | 
registry | HealthStatus | 
runtime | HealthStatus | 
store | HealthStatus | 


### Message HealthStatus
Name |  Type | Description
--------- | --------- | ---------


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

