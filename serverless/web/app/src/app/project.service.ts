import { Injectable } from "@angular/core";
import * as types from "./types";
import { environment } from "../environments/environment";
import { ClientService } from "@microhq/ng-client";
import * as _ from "lodash";

interface AppListResponse {
  apps: types.App[];
}

@Injectable({
  providedIn: "root"
})
export class ProjectService {
  constructor(private mc: ClientService) {
    this.mc.setOptions({ local: !environment.production });
  }

  list(): Promise<AppListResponse> {
    return this.mc
      .call<AppListResponse>("go.micro.service.serverless", "Apps.List", {})
      .then(rsp => {
        return {
          apps: rsp.apps.map(app => {
            app.name = app.name.split("/")[1];
            return app;
          })
        };
      });
  }

  create(app: types.App): Promise<void> {
    return this.mc.call("go.micro.service.serverless", "Apps.Create", {
      app: app
    });
  }
}
