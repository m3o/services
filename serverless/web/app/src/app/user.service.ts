import { Injectable } from "@angular/core";
import * as types from "./types";
import { HttpClient, HttpParams } from "@angular/common/http";
import { Subject } from "rxjs";
import { environment } from "../environments/environment";
import { CookieService } from "ngx-cookie-service";
import { NotificationsService } from "angular2-notifications";
import { Router } from "@angular/router";

@Injectable()
export class UserService {
  public user: types.User = {} as types.User;
  public isUserLoggedIn = new Subject<boolean>();

  constructor(
    private http: HttpClient,
    private cookie: CookieService,
    private notif: NotificationsService,
    private router: Router
  ) {
    this.get()
      .then(user => {
        for (const k of Object.keys(user)) {
          this.user[k] = user[k];
        }
        this.isUserLoggedIn.next(true);
      })
      .catch(e => {
        this.isUserLoggedIn.next(false);
      });
  }

  loggedIn(): boolean {
    return this.user && this.user.name != undefined;
  }

  logout() {
    // todo We are nulling out the name here because that's what we use
    // for user existence checks.
    this.user.name = "";
    this.cookie.delete("micro_token", "/")
    document.location.href = "/";
  }

  token(): string {
    return this.cookie.get("micro_token")
  }

  // gets current user
  get(): Promise<types.User> {
    if (!this.token() || this.token().length === 0) {
      return Promise.reject("Not logged in")
    }
    return this.http
      .get<types.User>(
        environment.backendUrl + "/v1/user?token=" + this.token()
      )
      .toPromise();
  }
}
