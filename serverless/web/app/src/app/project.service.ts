import { Injectable } from "@angular/core";
import * as types from "./types";
import { HttpClient } from "@angular/common/http";
import { environment } from "../environments/environment";
import { UserService } from "./user.service";
import * as _ from "lodash";
@Injectable({
  providedIn: "root"
})
export class ProjectService {
  constructor(private us: UserService, private http: HttpClient) {}

  listOrganisations(): Promise<types.Organisation[]> {
    return new Promise<types.Organisation[]>((resolve, reject) => {
      return this.http
        .get<types.Organisation[]>(
          environment.backendUrl +
            "/v1/github/organisations?token=" +
            this.us.token()
        )
        .toPromise()
        .then(servs => {
          resolve(servs as types.Organisation[]);
        })
        .catch(e => {
          reject(e);
        });
    });
  }

  listRepositories(organisation: string): Promise<types.Repository[]> {
    return new Promise<types.Repository[]>((resolve, reject) => {
      return this.http
        .get<types.Repository[]>(
          environment.backendUrl +
            "/v1/github/repositories?token=" +
            this.us.token() +
            "&organisation=" +
            organisation
        )
        .toPromise()
        .then(servs => {
          resolve(servs as types.Repository[]);
        })
        .catch(e => {
          reject(e);
        });
    });
  }

  listContents(
    organisation: string,
    repository: string,
    path: string
  ): Promise<types.RepoContents[]> {
    return new Promise<types.RepoContents[]>((resolve, reject) => {
      return this.http
        .get<types.RepoContents[]>(
          environment.backendUrl +
            "/v1/github/folders?token=" +
            this.us.token() +
            "&organisation=" +
            organisation +
            "&repository=" +
            repository +
            "&path=" +
            path
        )
        .toPromise()
        .then(servs => {
          resolve(servs as types.RepoContents[]);
        })
        .catch(e => {
          reject(e);
        });
    });
  }
}
