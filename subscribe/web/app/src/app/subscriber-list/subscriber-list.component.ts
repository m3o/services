import { Component, OnInit } from "@angular/core";
import { ClientService } from "@microhq/ng-client";
import { ActivatedRoute } from "@angular/router";
import { environment } from '../../environments/environment';

interface Subscriber {
  email?: string;
}

@Component({
  selector: "app-subscriber-list",
  templateUrl: "./subscriber-list.component.html",
  styleUrls: ["./subscriber-list.component.css"]
})
export class SubscriberListComponent implements OnInit {
  subscribers: Subscriber[];
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

    this.mc.setOptions({ local: !environment.production })
    this.mc
      .call("go.micro.service.subscribe", "Subscribe.ListSubscriptions", {
        namespace: this.domain
      })
      .then((response: any) => {
        console.log(response)
        this.subscribers = response.subscriptions
      });
  }
}
