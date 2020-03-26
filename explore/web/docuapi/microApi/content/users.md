
---
weight: 11
title: users
---
# users

## Users.Delete
```go
package main
import (
  "github.com/micro/clients/go/client"
  users_proto "github.com/micro/services/users/proto"
)
func main() {
  c := client.NewClient(nil)
  req := users_proto.DeleteRequest{}
  rsp := users_proto.DeleteResponse{}
  if err := c.Call("go.micro.srv.users", "Users.Delete", req, &rsp); err != nil {
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
      .call("go.micro.srv.users", "Users.Delete", {})
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
id | string | 
created | int64 | 
updated | int64 | 
first_name | string | 
last_name | string | 
email | string | 
username | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Users.Read
```go
package main
import (
  "github.com/micro/clients/go/client"
  users_proto "github.com/micro/services/users/proto"
)
func main() {
  c := client.NewClient(nil)
  req := users_proto.ReadRequest{}
  rsp := users_proto.ReadResponse{}
  if err := c.Call("go.micro.srv.users", "Users.Read", req, &rsp); err != nil {
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
      .call("go.micro.srv.users", "Users.Read", {})
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
id | string | 
created | int64 | 
updated | int64 | 
first_name | string | 
last_name | string | 
email | string | 
username | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Users.Update
```go
package main
import (
  "github.com/micro/clients/go/client"
  users_proto "github.com/micro/services/users/proto"
)
func main() {
  c := client.NewClient(nil)
  req := users_proto.UpdateRequest{}
  rsp := users_proto.UpdateResponse{}
  if err := c.Call("go.micro.srv.users", "Users.Update", req, &rsp); err != nil {
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
      .call("go.micro.srv.users", "Users.Update", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
user | User | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
user | User | 


### Message User
Name |  Type | Description
--------- | --------- | ---------
id | string | 
created | int64 | 
updated | int64 | 
first_name | string | 
last_name | string | 
email | string | 
username | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

