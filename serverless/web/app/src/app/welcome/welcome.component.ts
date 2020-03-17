import { Component, OnInit } from "@angular/core";
import { environment } from '../../environments/environment';
import { UserService } from '../user.service';
import { Router } from '@angular/router'

@Component({
  selector: "app-welcome",
  templateUrl: "./welcome.component.html",
  styleUrls: ["./welcome.component.css"]
})
export class WelcomeComponent implements OnInit {
  constructor(
    private us: UserService,
    private router: Router,
  ) {}

  ngOnInit() {
    if (this.us.loggedIn()) {
      this.router.navigate(['app/new'])
      return
    }
    this.us.isUserLoggedIn.subscribe(isIt => {
      if (isIt) {
        this.router.navigate(['app/new']);
      }
    })
  }

  login() {
    window.location.href = environment.backendUrl + "/v1/github/login"
  }
}
