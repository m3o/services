import { Component, OnInit } from "@angular/core";
import { ClientService } from "@microhq/ng-client";
import { ActivatedRoute } from "@angular/router";

@Component({
  selector: "app-subscribe-form",
  templateUrl: "./subscribe-form.component.html",
  styleUrls: ["./subscribe-form.component.css"],
  providers: []
})
export class SubscribeFormComponent implements OnInit {
  email = "";
  subscribed: boolean = false;
  domain = "";
  error = "";

  constructor(private mc: ClientService, private route: ActivatedRoute) {}

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      this.domain = params["domain"];
      if (!this.domain || this.domain.length == 0) {
        this.error =
          "No domain parameter. Please embed this page with a domain query param.";
        return;
      }
    });
  }

  subscribe() {
    if (!this.email) {
      return
    }
    this.mc
      .call("go.micro.store", "Store.Write", {
        record: {
          key: this.domain + ":" + this.email,
          value: btoa(this.email)
        }
      })
      .then(response => {
        this.subscribed = true
        console.log(response);
      })
      .catch(e => {
        this.error = e;
      });
  }
}
