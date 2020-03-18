import { Component, OnInit } from "@angular/core";
import * as types from "../types";
import { ProjectService } from "../project.service";
import { NotificationsService } from "angular2-notifications";

var groupBy = function(xs, key) {
  return xs.reduce(function(rv, x) {
    (rv[x[key]] = rv[x[key]] || []).push(x);
    return rv;
  }, {});
};

@Component({
  selector: "app-project-list",
  templateUrl: "./app-list.component.html",
  styleUrls: ["./app-list.component.css"]
})
export class AppListComponent implements OnInit {
  apps: types.App[];
  query = "";

  constructor(
    private ps: ProjectService,
    private notif: NotificationsService
  ) {}

  ngOnInit() {
    this.ps
      .list()
      .then(apps => {
        this.apps = groupBy(apps.apps, "name");
      })
      .catch(e => {
        console.log(e);
        this.notif.error(
          "Error listing services",
          e
        );
      });
  }
}
