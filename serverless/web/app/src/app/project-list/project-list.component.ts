import { Component, OnInit } from "@angular/core";
import * as types from "../types";
import { ProjectService } from "../project.service";

@Component({
  selector: "app-project-list",
  templateUrl: "./project-list.component.html",
  styleUrls: ["./project-list.component.css"]
})
export class ProjectListComponent implements OnInit {
  apps: types.App[];
  constructor(private ps: ProjectService) {}

  ngOnInit() {
    this.ps.list().then(apps => {
      this.apps = apps;
    });
  }
}
