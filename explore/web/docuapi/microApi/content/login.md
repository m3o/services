
---
weight: 11
title: login
---
# login

## Login.CreateLogin
```go
package main
import (
  "github.com/micro/clients/go/client"
  login_proto "github.com/micro/services/login/proto"
)
func main() {
  c := client.NewClient(nil)
  req := login_proto.CreateLoginRequest{}
  rsp := login_proto.CreateLoginResponse{}
  if err := c.Call("go.micro.srv.login", "Login.CreateLogin", req, &rsp); err != nil {
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
      .call("go.micro.srv.login", "Login.CreateLogin", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
id | string | 
email | string | 
password | string | 
validate_only | bool | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Login.UpdateEmail
```go
package main
import (
  "github.com/micro/clients/go/client"
  login_proto "github.com/micro/services/login/proto"
)
func main() {
  c := client.NewClient(nil)
  req := login_proto.UpdateEmailRequest{}
  rsp := login_proto.UpdateEmailResponse{}
  if err := c.Call("go.micro.srv.login", "Login.UpdateEmail", req, &rsp); err != nil {
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
      .call("go.micro.srv.login", "Login.UpdateEmail", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```
 TODO: Remove the RPC and replace with consuming users update events
 once we have update diff implemented.
### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
old_email | string | 
new_email | string | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Login.VerifyLogin
```go
package main
import (
  "github.com/micro/clients/go/client"
  login_proto "github.com/micro/services/login/proto"
)
func main() {
  c := client.NewClient(nil)
  req := login_proto.VerifyLoginRequest{}
  rsp := login_proto.VerifyLoginResponse{}
  if err := c.Call("go.micro.srv.login", "Login.VerifyLogin", req, &rsp); err != nil {
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
      .call("go.micro.srv.login", "Login.VerifyLogin", {})
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
id | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

