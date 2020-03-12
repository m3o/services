import { Component, OnInit, HostListener } from '@angular/core';
import { UserService } from '../user.service';
import { environment } from '../../environments/environment'

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.sass']
})
export class LoginComponent implements OnInit {

  constructor(private us: UserService) { }

  ngOnInit() {
  }

  @HostListener("click", ["$event"])
  public githubLogin(event: any) {
    this.us.logout()
    document.location.href = environment.backendUrl + '/v1/github/login';
    return false;
  }
}
