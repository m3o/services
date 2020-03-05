import { Component, OnInit } from "@angular/core";
import { Client } from "@microhq/ng-client";

@Component({
  selector: "app-subscribe-form",
  templateUrl: "./subscribe-form.component.html",
  styleUrls: ["./subscribe-form.component.css"],
  providers: []
})
export class SubscribeFormComponent implements OnInit {
  constructor(private mc: Client) {

  }

  ngOnInit() {
    this.mc.call("go.micro.srv.greeter", "Say.Hello");
  }
}
