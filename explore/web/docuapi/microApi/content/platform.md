
---
weight: 11
title: platform
---
# platform

## Platform.CreateService
```go
package main
import (
  "github.com/micro/clients/go/client"
  platform_proto "github.com/micro/services/platform/proto"
)
func main() {
  c := client.NewClient(nil)
  req := platform_proto.CreateServiceRequest{}
  rsp := platform_proto.CreateServiceResponse{}
  if err := c.Call("go.micro.srv.platform", "Platform.CreateService", req, &rsp); err != nil {
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
      .call("go.micro.srv.platform", "Platform.CreateService", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
service | Service | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### Message Service
Name |  Type | Description
--------- | --------- | ---------
name | string | 
version | string | 
source | string | 
type | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Platform.DeleteService
```go
package main
import (
  "github.com/micro/clients/go/client"
  platform_proto "github.com/micro/services/platform/proto"
)
func main() {
  c := client.NewClient(nil)
  req := platform_proto.DeleteServiceRequest{}
  rsp := platform_proto.DeleteServiceResponse{}
  if err := c.Call("go.micro.srv.platform", "Platform.DeleteService", req, &rsp); err != nil {
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
      .call("go.micro.srv.platform", "Platform.DeleteService", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
service | Service | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### Message Service
Name |  Type | Description
--------- | --------- | ---------
name | string | 
version | string | 
source | string | 
type | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Platform.ListServices
```go
package main
import (
  "github.com/micro/clients/go/client"
  platform_proto "github.com/micro/services/platform/proto"
)
func main() {
  c := client.NewClient(nil)
  req := platform_proto.ListServicesRequest{}
  rsp := platform_proto.ListServicesResponse{}
  if err := c.Call("go.micro.srv.platform", "Platform.ListServices", req, &rsp); err != nil {
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
      .call("go.micro.srv.platform", "Platform.ListServices", {})
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
services | Service | 


### Message Service
Name |  Type | Description
--------- | --------- | ---------
name | string | 
version | string | 
source | string | 
type | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Platform.ReadService
```go
package main
import (
  "github.com/micro/clients/go/client"
  platform_proto "github.com/micro/services/platform/proto"
)
func main() {
  c := client.NewClient(nil)
  req := platform_proto.ReadServiceRequest{}
  rsp := platform_proto.ReadServiceResponse{}
  if err := c.Call("go.micro.srv.platform", "Platform.ReadService", req, &rsp); err != nil {
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
      .call("go.micro.srv.platform", "Platform.ReadService", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
service | Service | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
services | Service | 


### Message Service
Name |  Type | Description
--------- | --------- | ---------
name | string | 
version | string | 
source | string | 
type | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Platform.ReadUser
```go
package main
import (
  "github.com/micro/clients/go/client"
  platform_proto "github.com/micro/services/platform/proto"
)
func main() {
  c := client.NewClient(nil)
  req := platform_proto.ReadUserRequest{}
  rsp := platform_proto.ReadUserResponse{}
  if err := c.Call("go.micro.srv.platform", "Platform.ReadUser", req, &rsp); err != nil {
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
      .call("go.micro.srv.platform", "Platform.ReadUser", {})
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
name | string | 
email | string | 
login | string | 
avatar_url | string | 
team_name | string | 
team_url | string | 
organization_avatar_url | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Platform.UpdateService
```go
package main
import (
  "github.com/micro/clients/go/client"
  platform_proto "github.com/micro/services/platform/proto"
)
func main() {
  c := client.NewClient(nil)
  req := platform_proto.UpdateServiceRequest{}
  rsp := platform_proto.UpdateServiceResponse{}
  if err := c.Call("go.micro.srv.platform", "Platform.UpdateService", req, &rsp); err != nil {
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
      .call("go.micro.srv.platform", "Platform.UpdateService", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
service | Service | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### Message Service
Name |  Type | Description
--------- | --------- | ---------
name | string | 
version | string | 
source | string | 
type | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

