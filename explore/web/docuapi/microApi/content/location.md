
---
weight: 11
title: location
---
# location

## Location.Read
```go
package main
import (
  "github.com/micro/clients/go/client"
  location_proto "github.com/micro/services/location/proto"
)
func main() {
  c := client.NewClient(nil)
  req := location_proto.ReadRequest{}
  rsp := location_proto.ReadResponse{}
  if err := c.Call("go.micro.srv.location", "Location.Read", req, &rsp); err != nil {
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
      .call("go.micro.srv.location", "Location.Read", {})
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

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
entity | Entity | 


### Message Entity
Name |  Type | Description
--------- | --------- | ---------
id | string | 
type | string | 
location | Point | 


### Message Point
Name |  Type | Description
--------- | --------- | ---------
latitude | double | 
longitude | double | 
timestamp | int64 | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Location.Save
```go
package main
import (
  "github.com/micro/clients/go/client"
  location_proto "github.com/micro/services/location/proto"
)
func main() {
  c := client.NewClient(nil)
  req := location_proto.SaveRequest{}
  rsp := location_proto.SaveResponse{}
  if err := c.Call("go.micro.srv.location", "Location.Save", req, &rsp); err != nil {
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
      .call("go.micro.srv.location", "Location.Save", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
entity | Entity | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### Message Entity
Name |  Type | Description
--------- | --------- | ---------
id | string | 
type | string | 
location | Point | 


### Message Point
Name |  Type | Description
--------- | --------- | ---------
latitude | double | 
longitude | double | 
timestamp | int64 | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Location.Search
```go
package main
import (
  "github.com/micro/clients/go/client"
  location_proto "github.com/micro/services/location/proto"
)
func main() {
  c := client.NewClient(nil)
  req := location_proto.SearchRequest{}
  rsp := location_proto.SearchResponse{}
  if err := c.Call("go.micro.srv.location", "Location.Search", req, &rsp); err != nil {
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
      .call("go.micro.srv.location", "Location.Search", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
center | Point | 
radius | double | 
type | string | 
numEntities | int64 | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
entities | Entity | 


### Message Point
Name |  Type | Description
--------- | --------- | ---------
latitude | double | 
longitude | double | 
timestamp | int64 | 


### Message Entity
Name |  Type | Description
--------- | --------- | ---------
id | string | 
type | string | 
location | Point | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

