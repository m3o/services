import { Component, OnInit } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import * as types from "../types";
import { ProjectService } from "../project.service";
import { NotificationsService } from "angular2-notifications";

@Component({
  selector: "app-app",
  templateUrl: "./app-single.component.html",
  styleUrls: ["./app-single.component.css"]
})
export class AppSingleComponent implements OnInit {
  app: types.App;

  constructor(
    private ps: ProjectService,
    private activeRoute: ActivatedRoute,
    private notif: NotificationsService
  ) {}

  ngOnInit() {
    this.activeRoute.params.subscribe(p => {
      this.ps
        .list()
        .then(apps => {
          this.app = apps.apps.filter(a => a.name == <string>p["id"])[0];
        })
        .catch(e => {
          this.notif.error(e);
        });
    });
  }
}
