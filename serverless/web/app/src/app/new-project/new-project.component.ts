import { Component, OnInit } from "@angular/core";
import * as types from "../types";
import { ProjectService } from "../project.service";
import { Router } from "@angular/router";
import { NotificationsService } from "angular2-notifications";

@Component({
  selector: "app-new-project",
  templateUrl: "./new-project.component.html",
  styleUrls: ["./new-project.component.css"]
})
export class NewProjectComponent implements OnInit {
  projectExists = false;
  loadingProjects = false;

  buildPacks: types.BuildPack[] = buildPacks;
  source = "";
  alias = "my-first-app";
  version = "";
  selectedBuildPackImageTag = "go";

  constructor(
    private ps: ProjectService,
    private router: Router,
    private notif: NotificationsService
  ) {}

  ngOnInit() {}

  keyPress($event) {}

  create() {
    const app = {
      name: this.alias,
      source: this.source,
      version: this.version,
      language: this.selectedBuildPackImageTag
    };
    console.log(app);
    this.ps
      .create(app)
      .then(() => {
        this.router.navigate(["/apps"]);
      })
      .catch(e => {
        this.notif.error("Error creating application", e);
      });
  }
}

const buildPacks: types.BuildPack[] = [
  {
    name: "Go",
    imageTag: "go"
  },
  {
    name: "Node.js",
    imageTag: "nodejs"
  },
  {
    name: "HTML",
    imageTag: "html"
  },
  {
    name: "Shell",
    imageTag: "shell"
  },
  {
    name: "PHP",
    imageTag: "php"
  },
  {
    name: "Python",
    imageTag: "python"
  },
  {
    name: "Ruby",
    imageTag: "ruby"
  },
  {
    name: "Rust",
    imageTag: "rust"
  },
  {
    name: "Java",
    imageTag: "java"
  }
];
