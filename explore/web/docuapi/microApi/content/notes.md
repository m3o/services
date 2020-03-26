
---
weight: 11
title: notes
---
# notes

## Notes.CreateNote
```go
package main
import (
  "github.com/micro/clients/go/client"
  notes_proto "github.com/micro/services/notes/proto"
)
func main() {
  c := client.NewClient(nil)
  req := notes_proto.CreateNoteRequest{}
  rsp := notes_proto.CreateNoteResponse{}
  if err := c.Call("go.micro.srv.notes", "Notes.CreateNote", req, &rsp); err != nil {
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
      .call("go.micro.srv.notes", "Notes.CreateNote", {})
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

## Notes.DeleteNote
```go
package main
import (
  "github.com/micro/clients/go/client"
  notes_proto "github.com/micro/services/notes/proto"
)
func main() {
  c := client.NewClient(nil)
  req := notes_proto.DeleteNoteRequest{}
  rsp := notes_proto.DeleteNoteResponse{}
  if err := c.Call("go.micro.srv.notes", "Notes.DeleteNote", req, &rsp); err != nil {
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
      .call("go.micro.srv.notes", "Notes.DeleteNote", {})
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

## Notes.ListNotes
```go
package main
import (
  "github.com/micro/clients/go/client"
  notes_proto "github.com/micro/services/notes/proto"
)
func main() {
  c := client.NewClient(nil)
  req := notes_proto.ListNotesRequest{}
  rsp := notes_proto.ListNotesResponse{}
  if err := c.Call("go.micro.srv.notes", "Notes.ListNotes", req, &rsp); err != nil {
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
      .call("go.micro.srv.notes", "Notes.ListNotes", {})
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

## Notes.UpdateNote
```go
package main
import (
  "github.com/micro/clients/go/client"
  notes_proto "github.com/micro/services/notes/proto"
)
func main() {
  c := client.NewClient(nil)
  req := notes_proto.UpdateNoteRequest{}
  rsp := notes_proto.UpdateNoteResponse{}
  if err := c.Call("go.micro.srv.notes", "Notes.UpdateNote", req, &rsp); err != nil {
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
      .call("go.micro.srv.notes", "Notes.UpdateNote", {})
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

