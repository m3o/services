
---
weight: 11
title: payments
---
# payments

## Provider.CreatePaymentMethod
```go
package main
import (
  "github.com/micro/clients/go/client"
  payments_proto "github.com/micro/services/payments/proto"
)
func main() {
  c := client.NewClient(nil)
  req := payments_proto.CreatePaymentMethodRequest{}
  rsp := payments_proto.CreatePaymentMethodResponse{}
  if err := c.Call("go.micro.srv.payments", "Provider.CreatePaymentMethod", req, &rsp); err != nil {
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
      .call("go.micro.srv.payments", "Provider.CreatePaymentMethod", {})
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
user_id | string | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
payment_method | PaymentMethod | 


### Message PaymentMethod
Name |  Type | Description
--------- | --------- | ---------
id | string | 
created | int64 | 
user_id | string | 
type | string | 
card_brand | string | 
card_exp_month | string | 
card_exp_year | string | 
card_last_4 | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Provider.CreatePlan
```go
package main
import (
  "github.com/micro/clients/go/client"
  payments_proto "github.com/micro/services/payments/proto"
)
func main() {
  c := client.NewClient(nil)
  req := payments_proto.CreatePlanRequest{}
  rsp := payments_proto.CreatePlanResponse{}
  if err := c.Call("go.micro.srv.payments", "Provider.CreatePlan", req, &rsp); err != nil {
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
      .call("go.micro.srv.payments", "Provider.CreatePlan", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
plan | Plan | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### Message Plan
Name |  Type | Description
--------- | --------- | ---------
id | string | 
name | string | 
amount | int64 | 
currency | string | 
interval | PlanInterval | 
product_id | string | 


### Message PlanInterval
Name |  Type | Description
--------- | --------- | ---------


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Provider.CreateProduct
```go
package main
import (
  "github.com/micro/clients/go/client"
  payments_proto "github.com/micro/services/payments/proto"
)
func main() {
  c := client.NewClient(nil)
  req := payments_proto.CreateProductRequest{}
  rsp := payments_proto.CreateProductResponse{}
  if err := c.Call("go.micro.srv.payments", "Provider.CreateProduct", req, &rsp); err != nil {
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
      .call("go.micro.srv.payments", "Provider.CreateProduct", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
product | Product | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### Message Product
Name |  Type | Description
--------- | --------- | ---------
id | string | 
name | string | 
description | string | 
active | bool | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Provider.CreateSubscription
```go
package main
import (
  "github.com/micro/clients/go/client"
  payments_proto "github.com/micro/services/payments/proto"
)
func main() {
  c := client.NewClient(nil)
  req := payments_proto.CreateSubscriptionRequest{}
  rsp := payments_proto.CreateSubscriptionResponse{}
  if err := c.Call("go.micro.srv.payments", "Provider.CreateSubscription", req, &rsp); err != nil {
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
      .call("go.micro.srv.payments", "Provider.CreateSubscription", {})
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

## Provider.CreateUser
```go
package main
import (
  "github.com/micro/clients/go/client"
  payments_proto "github.com/micro/services/payments/proto"
)
func main() {
  c := client.NewClient(nil)
  req := payments_proto.CreateUserRequest{}
  rsp := payments_proto.CreateUserResponse{}
  if err := c.Call("go.micro.srv.payments", "Provider.CreateUser", req, &rsp); err != nil {
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
      .call("go.micro.srv.payments", "Provider.CreateUser", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
User | User | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------


### Message User
Name |  Type | Description
--------- | --------- | ---------
id | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Provider.DeletePaymentMethod
```go
package main
import (
  "github.com/micro/clients/go/client"
  payments_proto "github.com/micro/services/payments/proto"
)
func main() {
  c := client.NewClient(nil)
  req := payments_proto.DeletePaymentMethodRequest{}
  rsp := payments_proto.DeletePaymentMethodResponse{}
  if err := c.Call("go.micro.srv.payments", "Provider.DeletePaymentMethod", req, &rsp); err != nil {
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
      .call("go.micro.srv.payments", "Provider.DeletePaymentMethod", {})
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


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Provider.ListPaymentMethods
```go
package main
import (
  "github.com/micro/clients/go/client"
  payments_proto "github.com/micro/services/payments/proto"
)
func main() {
  c := client.NewClient(nil)
  req := payments_proto.ListPaymentMethodsRequest{}
  rsp := payments_proto.ListPaymentMethodsResponse{}
  if err := c.Call("go.micro.srv.payments", "Provider.ListPaymentMethods", req, &rsp); err != nil {
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
      .call("go.micro.srv.payments", "Provider.ListPaymentMethods", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
user_id | string | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
payment_methods | PaymentMethod | 


### Message PaymentMethod
Name |  Type | Description
--------- | --------- | ---------
id | string | 
created | int64 | 
user_id | string | 
type | string | 
card_brand | string | 
card_exp_month | string | 
card_exp_year | string | 
card_last_4 | string | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

