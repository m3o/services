import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";
import { NewProjectComponent } from "./new-project/new-project.component";
import { AuthGuard } from "./auth.guard";
import { WelcomeComponent } from "./welcome/welcome.component";
import { NotInvitedComponent } from "./not-invited/not-invited.component";

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
    path: "project/new",
    component: NewProjectComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "project/new/:id",
    component: NewProjectComponent,
    canActivate: [AuthGuard]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
