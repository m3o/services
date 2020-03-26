
---
weight: 11
title: distributed
---
# distributed

## DistributedNotes.CreateNote
```go
package main
import (
  "github.com/micro/clients/go/client"
  distributed_proto "github.com/micro/services/distributed/proto"
)
func main() {
  c := client.NewClient(nil)
  req := distributed_proto.CreateNoteRequest{}
  rsp := distributed_proto.CreateNoteResponse{}
  if err := c.Call("go.micro.srv.distributed", "DistributedNotes.CreateNote", req, &rsp); err != nil {
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
      .call("go.micro.srv.distributed", "DistributedNotes.CreateNote", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
note | Note | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
note | Note | 


### Message Note
Name |  Type | Description
--------- | --------- | ---------
id | string | 
created | int64 | 
title | string | 
text | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## DistributedNotes.DeleteNote
```go
package main
import (
  "github.com/micro/clients/go/client"
  distributed_proto "github.com/micro/services/distributed/proto"
)
func main() {
  c := client.NewClient(nil)
  req := distributed_proto.DeleteNoteRequest{}
  rsp := distributed_proto.DeleteNoteResponse{}
  if err := c.Call("go.micro.srv.distributed", "DistributedNotes.DeleteNote", req, &rsp); err != nil {
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
      .call("go.micro.srv.distributed", "DistributedNotes.DeleteNote", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
note | Note | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### Message Note
Name |  Type | Description
--------- | --------- | ---------
id | string | 
created | int64 | 
title | string | 
text | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## DistributedNotes.ListNotes
```go
package main
import (
  "github.com/micro/clients/go/client"
  distributed_proto "github.com/micro/services/distributed/proto"
)
func main() {
  c := client.NewClient(nil)
  req := distributed_proto.ListNotesRequest{}
  rsp := distributed_proto.ListNotesResponse{}
  if err := c.Call("go.micro.srv.distributed", "DistributedNotes.ListNotes", req, &rsp); err != nil {
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
      .call("go.micro.srv.distributed", "DistributedNotes.ListNotes", {})
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
notes | Note | 


### Message Note
Name |  Type | Description
--------- | --------- | ---------
id | string | 
created | int64 | 
title | string | 
text | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## DistributedNotes.UpdateNote
```go
package main
import (
  "github.com/micro/clients/go/client"
  distributed_proto "github.com/micro/services/distributed/proto"
)
func main() {
  c := client.NewClient(nil)
  req := distributed_proto.UpdateNoteRequest{}
  rsp := distributed_proto.UpdateNoteResponse{}
  if err := c.Call("go.micro.srv.distributed", "DistributedNotes.UpdateNote", req, &rsp); err != nil {
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
      .call("go.micro.srv.distributed", "DistributedNotes.UpdateNote", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
note | Note | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### Message Note
Name |  Type | Description
--------- | --------- | ---------
id | string | 
created | int64 | 
title | string | 
text | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

