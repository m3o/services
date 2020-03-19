import { Injectable } from "@angular/core";
import {
  CanActivate,
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  UrlTree,
  Router
} from "@angular/router";
import { Observable } from "rxjs";
import { UserService } from "./user.service";
import { environment } from "../environments/environment";

@Injectable({
  providedIn: "root"
})
export class AuthGuard implements CanActivate {
  constructor(
    private us: UserService,
    private r: Router) {}

  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ):
    | Observable<boolean | UrlTree>
    | Promise<boolean | UrlTree>
    | boolean
    | UrlTree {
    if (this.us.loggedIn()) {
      return true;
    }
    return new Observable<boolean>(observer => {
      this.us.isUserLoggedIn.subscribe(loggedIn => {
        if (loggedIn) {
          observer.next(true);
        } else {
          //confirm("redirect") ? window.location.href = environment.backendUrl + "/v1/github/login" : console.log("stopping")
          //window.location.href = environment.backendUrl + "/v1/github/login"
          this.r.navigate(['/'])
          observer.next(false);
        }
        observer.complete();
      });
    });
  }
}
