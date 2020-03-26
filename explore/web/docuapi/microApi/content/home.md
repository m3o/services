
---
weight: 11
title: home
---
# home

## Home.ListApps
```go
package main
import (
  "github.com/micro/clients/go/client"
  home_proto "github.com/micro/services/home/proto"
)
func main() {
  c := client.NewClient(nil)
  req := home_proto.ListAppsRequest{}
  rsp := home_proto.ListAppsResponse{}
  if err := c.Call("go.micro.srv.home", "Home.ListApps", req, &rsp); err != nil {
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
      .call("go.micro.srv.home", "Home.ListApps", {})
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
apps | App | 


### Message App
Name |  Type | Description
--------- | --------- | ---------
id | string | 
name | string | 
category | string | 
icon | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Home.ReadUser
```go
package main
import (
  "github.com/micro/clients/go/client"
  home_proto "github.com/micro/services/home/proto"
)
func main() {
  c := client.NewClient(nil)
  req := home_proto.ReadUserRequest{}
  rsp := home_proto.ReadUserResponse{}
  if err := c.Call("go.micro.srv.home", "Home.ReadUser", req, &rsp); err != nil {
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
      .call("go.micro.srv.home", "Home.ReadUser", {})
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
user | User | 


### Message User
Name |  Type | Description
--------- | --------- | ---------
first_name | string | 
last_name | string | 
profile_picture_url | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

