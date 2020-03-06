import { Component, OnInit } from "@angular/core";
import { ClientService } from "@microhq/ng-client";

@Component({
  selector: "app-subscribe-form",
  templateUrl: "./subscribe-form.component.html",
  styleUrls: ["./subscribe-form.component.css"],
  providers: []
})
export class SubscribeFormComponent implements OnInit {
  subscribed: boolean = false;

  constructor(private mc: ClientService) {}

  ngOnInit() {
    this.mc.call("go.micro.srv.greeter", "Say.Hello").then(response => {
      console.log(response);
    });
  }
}
