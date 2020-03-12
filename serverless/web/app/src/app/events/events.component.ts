import { Component, OnInit } from "@angular/core";
import { ServiceService } from "../service.service";
import * as types from "../types";

@Component({
  selector: "app-events",
  templateUrl: "./events.component.html",
  styleUrls: ["./events.component.css"]
})
export class EventsComponent implements OnInit {
  query: string = "";
  events: types.Event[] = [];

  constructor(private ses: ServiceService) {}

  ngOnInit() {
    this.ses.events().then(v => {
      this.events = v;
    });
  }
}
