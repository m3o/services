
---
weight: 11
title: graphql
---
# graphql

## Graphql.Call
```go
package main
import (
  "github.com/micro/clients/go/client"
  graphql_proto "github.com/micro/services/graphql/proto"
)
func main() {
  c := client.NewClient(nil)
  req := graphql_proto.go.api.Request{}
  rsp := graphql_proto.go.api.Response{}
  if err := c.Call("go.micro.srv.graphql", "Graphql.Call", req, &rsp); err != nil {
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
      .call("go.micro.srv.graphql", "Graphql.Call", {})
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

