import { Component, OnInit } from "@angular/core";
import { UserService } from "../user.service";
@Component({
  selector: "app-settings",
  templateUrl: "./settings.component.html",
  styleUrls: ["./settings.component.css"]
})
export class SettingsComponent implements OnInit {
  token = "";

  constructor(private us: UserService) {}

  ngOnInit() {
    this.token = this.us.token();
  }

  languages = ["bash"];
}
