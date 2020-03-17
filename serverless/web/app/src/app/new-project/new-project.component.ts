import { Component, OnInit } from "@angular/core";
import * as types from "../types";
import { ProjectService } from "../project.service";
import { Router } from "@angular/router";

@Component({
  selector: "app-new-project",
  templateUrl: "./new-project.component.html",
  styleUrls: ["./new-project.component.css"]
})
export class NewProjectComponent implements OnInit {
  buildPacks: types.BuildPack[] = buildPacks;
  organisations: types.Organisation[] = [];
  source = "";
  repositories: types.Repository[] = [];
  contents: types.RepoContents[] = [];
  step = 0;
  alias = "my-first-app";
  version = "";

  projectExists = false;
  loadingProjects = false;
  loaded = true;
  selectedOrg: string;
  selectedRepo: string;
  add = true;
  selectedBuildPack: types.BuildPack;
  path: string = "";

  constructor(private ps: ProjectService, private router: Router) {}

  ngOnInit() {}

  keyPress($event) {}

  create() {
    this.ps
      .create({
        name: this.alias,
        source: this.source,
        version: this.version
      })
      .then(() => {
        this.router.navigate(["/"]);
      });
  }
}

const buildPacks: types.BuildPack[] = [
  {
    name: "Go"
  },
  {
    name: "Node.js"
  },
  {
    name: "HTML"
  },
  {
    name: "Shell"
  },
  {
    name: "PHP"
  },
  {
    name: "Python"
  },
  {
    name: "Ruby"
  },
  {
    name: "Rust"
  },
  {
    name: "Java"
  }
];
