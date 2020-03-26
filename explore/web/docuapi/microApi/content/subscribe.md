
---
weight: 11
title: subscribe
---
# subscribe

## Subscribe.ListSubscriptions
```go
package main
import (
  "github.com/micro/clients/go/client"
  subscribe_proto "github.com/micro/services/subscribe/proto"
)
func main() {
  c := client.NewClient(nil)
  req := subscribe_proto.ListSubscriptionsRequest{}
  rsp := subscribe_proto.ListSubscriptionsResponse{}
  if err := c.Call("go.micro.srv.subscribe", "Subscribe.ListSubscriptions", req, &rsp); err != nil {
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
      .call("go.micro.srv.subscribe", "Subscribe.ListSubscriptions", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
namespace | string | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
subscriptions | Subscription | 


### Message Subscription
Name |  Type | Description
--------- | --------- | ---------
email | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Subscribe.Subscribe
```go
package main
import (
  "github.com/micro/clients/go/client"
  subscribe_proto "github.com/micro/services/subscribe/proto"
)
func main() {
  c := client.NewClient(nil)
  req := subscribe_proto.SubscribeRequest{}
  rsp := subscribe_proto.SubscribeResponse{}
  if err := c.Call("go.micro.srv.subscribe", "Subscribe.Subscribe", req, &rsp); err != nil {
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
      .call("go.micro.srv.subscribe", "Subscribe.Subscribe", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
namespace | string |  namespaces are owned by the account that first creates a subscription in them
email | string | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

