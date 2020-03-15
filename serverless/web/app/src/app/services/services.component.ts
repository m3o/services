import { Component, OnInit } from "@angular/core";
import { ServiceService } from "../service.service";
import * as types from "../types";
import { NotificationsService } from "angular2-notifications";

var groupBy = function(xs, key) {
  return xs.reduce(function(rv, x) {
    (rv[x[key]] = rv[x[key]] || []).push(x);
    return rv;
  }, {});
};

@Component({
  selector: "app-services",
  templateUrl: "./services.component.html",
  styleUrls: ["./services.component.scss"]
})
export class ServicesComponent implements OnInit {
  services: Map<string, types.Service[]>;
  query: string;

  constructor(
    private ses: ServiceService,
    private notif: NotificationsService
  ) {}

  ngOnInit() {
    this.ses
      .list()
      .then(servs => {
        this.services = groupBy(servs, "name");
      })
      .catch(e => {
        console.log(e);
        this.notif.error(
          "Error listing services",
          JSON.parse(e.error.error).detail
        );
      });
  }
}
