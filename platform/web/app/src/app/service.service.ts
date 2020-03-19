import { Injectable } from "@angular/core";
import * as types from "./types";
import { HttpClient } from "@angular/common/http";
import { environment } from "../environments/environment";
import { UserService } from "./user.service";
import * as _ from "lodash";

export interface RPCRequest {
  service: string;
  endpoint: string;
  method?: string;
  address?: string;
  request: any;
}

@Injectable({
  providedIn: "root"
})
export class ServiceService {
  constructor(private us: UserService, private http: HttpClient) {}

  list(): Promise<types.Service[]> {
    return new Promise<types.Service[]>((resolve, reject) => {
      return this.http
        .get<types.Service[]>(
          environment.backendUrl + "/v1/services?token=" + this.us.token()
        )
        .toPromise()
        .then(servs => {
          resolve(servs as types.Service[]);
        })
        .catch(e => {
          reject(e);
        });
    });
  }

  logs(service: string): Promise<types.LogRecord[]> {
    return new Promise<types.LogRecord[]>((resolve, reject) => {
      return this.http
        .get<types.LogRecord[]>(
          environment.backendUrl +
            "/v1/service/logs?service=" +
            service +
            "&token=" +
            this.us.token()
        )
        .toPromise()
        .then(servs => {
          resolve(servs as types.LogRecord[]);
        })
        .catch(e => {
          reject(e);
        });
    });
  }

  stats(service: string, version?: string): Promise<types.DebugSnapshot[]> {
    return new Promise<types.DebugSnapshot[]>((resolve, reject) => {
      return this.http
        .get<types.DebugSnapshot[]>(
          environment.backendUrl +
            "/v1/service/stats?service=" +
            service +
            "&token=" +
            this.us.token()
        )
        .toPromise()
        .then(servs => {
          resolve(servs as types.DebugSnapshot[]);
        })
        .catch(e => {
          reject(e);
        });
    });
  }

  trace(service?: string): Promise<types.Span[]> {
    const qs = service ? "service=" + service + "&" : "";
    return new Promise<types.Span[]>((resolve, reject) => {
      return this.http
        .get<types.Span[]>(
          environment.backendUrl +
            "/v1/service/trace?" +
            qs +
            "token=" +
            this.us.token() +
            "&limit=1000"
        )
        .toPromise()
        .then(servs => {
          resolve(servs as types.Span[]);
        })
        .catch(e => {
          reject(e);
        });
    });
  }

  call(rpc: RPCRequest): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      return this.http
        .post<string>(environment.backendUrl + "/v1/service/call", rpc)
        .toPromise()
        .then(response => {
          resolve(JSON.stringify(response, null, "  "));
        })
        .catch(e => {
          reject(e);
        });
    });
  }

  events(service?: string): Promise<types.Event[]> {
    const serviceQuery = service ? "?service=" + service : "";
    return new Promise<types.Event[]>((resolve, reject) => {
      return this.http
        .get<types.Event[]>(
          environment.backendUrl + "/v1/events" + serviceQuery
        )
        .toPromise()
        .then(events => {
          resolve(_.orderBy(events, e => e.timestamp, ["desc"]));
        })
        .catch(e => {
          reject(e);
        });
    });
  }
}
