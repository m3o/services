
---
weight: 11
title: account
---
# account

## Account.CreatePaymentMethod
```go
package main
import (
  "github.com/micro/clients/go/client"
  account_proto "github.com/micro/services/account/proto"
)
func main() {
  c := client.NewClient(nil)
  req := account_proto.CreatePaymentMethodRequest{}
  rsp := account_proto.CreatePaymentMethodResponse{}
  if err := c.Call("go.micro.srv.account", "Account.CreatePaymentMethod", req, &rsp); err != nil {
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
      .call("go.micro.srv.account", "Account.CreatePaymentMethod", {})
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

## Account.DeletePaymentMethod
```go
package main
import (
  "github.com/micro/clients/go/client"
  account_proto "github.com/micro/services/account/proto"
)
func main() {
  c := client.NewClient(nil)
  req := account_proto.DeletePaymentMethodRequest{}
  rsp := account_proto.DeletePaymentMethodResponse{}
  if err := c.Call("go.micro.srv.account", "Account.DeletePaymentMethod", req, &rsp); err != nil {
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
      .call("go.micro.srv.account", "Account.DeletePaymentMethod", {})
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

## Account.DeleteUser
```go
package main
import (
  "github.com/micro/clients/go/client"
  account_proto "github.com/micro/services/account/proto"
)
func main() {
  c := client.NewClient(nil)
  req := account_proto.DeleteUserRequest{}
  rsp := account_proto.DeleteUserResponse{}
  if err := c.Call("go.micro.srv.account", "Account.DeleteUser", req, &rsp); err != nil {
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
      .call("go.micro.srv.account", "Account.DeleteUser", {})
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

## Account.Login
```go
package main
import (
  "github.com/micro/clients/go/client"
  account_proto "github.com/micro/services/account/proto"
)
func main() {
  c := client.NewClient(nil)
  req := account_proto.LoginRequest{}
  rsp := account_proto.LoginResponse{}
  if err := c.Call("go.micro.srv.account", "Account.Login", req, &rsp); err != nil {
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
      .call("go.micro.srv.account", "Account.Login", {})
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
user | User | 
token | Token | 
secret | Token | 


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
payment_methods | PaymentMethod | 
subscriptions | Subscription | 


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


### Message Subscription
Name |  Type | Description
--------- | --------- | ---------


### Message Token
Name |  Type | Description
--------- | --------- | ---------
token | string | 
expires | int64 | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Account.ReadUser
```go
package main
import (
  "github.com/micro/clients/go/client"
  account_proto "github.com/micro/services/account/proto"
)
func main() {
  c := client.NewClient(nil)
  req := account_proto.ReadUserRequest{}
  rsp := account_proto.ReadUserResponse{}
  if err := c.Call("go.micro.srv.account", "Account.ReadUser", req, &rsp); err != nil {
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
      .call("go.micro.srv.account", "Account.ReadUser", {})
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
payment_methods | PaymentMethod | 
subscriptions | Subscription | 


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


### Message Subscription
Name |  Type | Description
--------- | --------- | ---------


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Account.RefreshToken
```go
package main
import (
  "github.com/micro/clients/go/client"
  account_proto "github.com/micro/services/account/proto"
)
func main() {
  c := client.NewClient(nil)
  req := account_proto.RefreshTokenRequest{}
  rsp := account_proto.RefreshTokenResponse{}
  if err := c.Call("go.micro.srv.account", "Account.RefreshToken", req, &rsp); err != nil {
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
      .call("go.micro.srv.account", "Account.RefreshToken", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
```

### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
secret | string | 

### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
token | Token | 


### Message Token
Name |  Type | Description
--------- | --------- | ---------
token | string | 
expires | int64 | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Account.Signup
```go
package main
import (
  "github.com/micro/clients/go/client"
  account_proto "github.com/micro/services/account/proto"
)
func main() {
  c := client.NewClient(nil)
  req := account_proto.SignupRequest{}
  rsp := account_proto.SignupResponse{}
  if err := c.Call("go.micro.srv.account", "Account.Signup", req, &rsp); err != nil {
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
      .call("go.micro.srv.account", "Account.Signup", {})
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
user | User | 
token | Token | 
secret | Token | 


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
payment_methods | PaymentMethod | 
subscriptions | Subscription | 


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


### Message Subscription
Name |  Type | Description
--------- | --------- | ---------


### Message Token
Name |  Type | Description
--------- | --------- | ---------
token | string | 
expires | int64 | 


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

## Account.UpdateUser
```go
package main
import (
  "github.com/micro/clients/go/client"
  account_proto "github.com/micro/services/account/proto"
)
func main() {
  c := client.NewClient(nil)
  req := account_proto.UpdateUserRequest{}
  rsp := account_proto.UpdateUserResponse{}
  if err := c.Call("go.micro.srv.account", "Account.UpdateUser", req, &rsp); err != nil {
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
      .call("go.micro.srv.account", "Account.UpdateUser", {})
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
payment_methods | PaymentMethod | 
subscriptions | Subscription | 


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


### Message Subscription
Name |  Type | Description
--------- | --------- | ---------


### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>

