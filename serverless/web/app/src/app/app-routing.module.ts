import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";
import { NewAppComponent } from "./new-app/new-app.component";
import { AuthGuard } from "./auth.guard";
import { WelcomeComponent } from "./welcome/welcome.component";
import { NotInvitedComponent } from "./not-invited/not-invited.component";
import { AppListComponent } from "./app-list/app-list.component";
import { AppComponent } from "./app/app.component";

const routes: Routes = [
  {
    path: "",
    component: WelcomeComponent,
    pathMatch: "full"
  },
  {
    path: "not-invited",
    component: NotInvitedComponent
  },
  {
    path: "app/new",
    component: NewAppComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "app/new/:id",
    component: NewAppComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "app/:id",
    component: AppComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "apps",
    component: AppListComponent,
    canActivate: [AuthGuard]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
