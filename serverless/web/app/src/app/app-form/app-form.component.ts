import { Component, OnInit, Input } from "@angular/core";
import * as types from "../types";
import { ProjectService } from "../project.service";
import { Router } from "@angular/router";
import { NotificationsService } from "angular2-notifications";

@Component({
  selector: "app-app-form",
  templateUrl: "./app-form.component.html",
  styleUrls: ["./app-form.component.css"]
})
export class AppFormComponent implements OnInit {
  @Input() app: types.App;
  create = false;

  projectExists = false;
  loadingApps = false;

  buildPacks: types.BuildPack[] = buildPacks;
  selectedBuildPackImageTag = "go";

  constructor(
    private ps: ProjectService,
    private router: Router,
    private notif: NotificationsService
  ) {
    if (!this.app || !this.app.name) {
      this.create = true;
      this.app = {
        name: "my-app-" + makeid(6)
      };
    }
  }

  ngOnInit() {}

  keyPress($event) {}

  createApp() {
    const app = {
      name: this.app.name,
      source: this.app.source,
      version: this.app.version,
      language: this.selectedBuildPackImageTag
    };
    this.ps
      .create(app)
      .then(() => {
        this.router.navigate(["/apps"]);
      })
      .catch(e => {
        this.notif.error("Error creating application", e);
      });
  }

  saveApp() {
    this.notif.alert("App update is not implemented yet");
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

function makeid(length) {
  var result = "";
  var characters = "abcdefghijklmnopqrstuvwxyz0123456789";
  var charactersLength = characters.length;
  for (var i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}
