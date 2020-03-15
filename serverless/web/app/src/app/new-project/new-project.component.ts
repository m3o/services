import { Component, OnInit } from "@angular/core";
import * as types from "../types";
import { ProjectService } from "../project.service";

@Component({
  selector: "app-new-project",
  templateUrl: "./new-project.component.html",
  styleUrls: ["./new-project.component.css"]
})
export class NewProjectComponent implements OnInit {
  buildPacks: types.BuildPack[] = buildPacks;
  organisations: types.Organisation[] = [];
  repositories: types.Repository[] = [];
  contents: types.RepoContents[] = [];
  step = 0;
  alias = "my-first-app";
  projectExists = false;
  loadingProjects = false;
  loaded = true;
  selectedOrg: string;
  selectedRepo: string;
  add = true;
  selectedBuildPack: types.BuildPack;
  path: string = "";

  constructor(private ps: ProjectService) {}

  ngOnInit() {
    this.ps.listOrganisations().then(orgs => {
      this.organisations = orgs;
    });
  }

  keyPress($event) {}

  orgSelected(v: string) {
    this.ps.listRepositories(v).then(repos => {
      this.repositories = repos;
    });
  }

  repoSelected(v: string) {
    this.loadFolders();
  }

  loadFolders() {
    this.ps
      .listContents(this.selectedOrg, this.selectedRepo, this.path)
      .then(contents => {
        this.contents = contents.filter(c => c.type == "dir");
      });
  }

  bpSelected(v: string) {}

  folderSelected(v: string) {
    this.add = false;
    this.path += "/" + v;
    this.loadFolders();
  }
}

const buildPacks: types.BuildPack[] = [
  {
    name: "go"
  },
  {
    name: "node.js"
  },
  {
    name: "shell"
  },
  {
    name: "php"
  },
  {
    name: "python"
  },
  {
    name: "ruby"
  },
  {
    name: "rust"
  },
  {
    name: "java"
  },
  {
    name: "html"
  }
];
