import { Injectable } from "@angular/core";
import * as types from "./types";
import { environment } from "../environments/environment";
import { ClientService } from "@microhq/ng-client";
import * as _ from "lodash";

@Injectable({
  providedIn: "root"
})
export class ProjectService {
  constructor(
    private mc: ClientService,
  ) {
    this.mc.setOptions({ local: !environment.production });
  }

  list(): Promise<types.App[]> {
    return this.mc.call<types.App[]>(
      "go.micro.service.serverless",
      "Apps.List",
      {}
    );
  }

  create(app: types.App): Promise<void> {
    return this.mc.call("go.micro.service.serverless", "Apps.Create", {
      app: app
    })
  }
}
