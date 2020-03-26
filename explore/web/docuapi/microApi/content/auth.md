
---
weight: 11
title: auth
---
# auth

## Auth.Login
```go
package main
import (
  "github.com/micro/clients/go/client"
  auth_proto "github.com/micro/services/auth/proto"
)
func main() {
  c := client.NewClient(nil)
  req := auth_proto.LoginRequest{}
  rsp := auth_proto.LoginResponse{}
  if err := c.Call("go.micro.srv.auth", "Auth.Login", req, &rsp); err != nil {
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
      .call("go.micro.srv.auth", "Auth.Login", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
email | string | 
password | string | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
token | string | 
user | User | 


### Message User
Name |  Type | Description
--------- | --------- | ---------
id | string | 
firstname | string | 
lastname | string | 
email | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Auth.Register
```go
package main
import (
  "github.com/micro/clients/go/client"
  auth_proto "github.com/micro/services/auth/proto"
)
func main() {
  c := client.NewClient(nil)
  req := auth_proto.RegisterRequest{}
  rsp := auth_proto.RegisterResponse{}
  if err := c.Call("go.micro.srv.auth", "Auth.Register", req, &rsp); err != nil {
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
      .call("go.micro.srv.auth", "Auth.Register", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
email | string | 
password | string | 
firstname | string | 
lastname | string | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
message | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

