import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";
import { ServicesComponent } from "./services/services.component";
import { ServiceComponent } from "./service/service.component";
import { NewServiceComponent } from "./new-service/new-service.component";
import { AuthGuard } from "./auth.guard";
import { WelcomeComponent } from "./welcome/welcome.component";
import { NotInvitedComponent } from "./not-invited/not-invited.component";
import { SettingsComponent } from "./settings/settings.component";
import { EventsComponent } from "./events/events.component";

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
    path: "service/new",
    component: NewServiceComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "service/new/:id",
    component: NewServiceComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "service/:id/:tab",
    component: ServiceComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "service/:id",
    component: ServiceComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "settings/:id",
    component: SettingsComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "events",
    component: EventsComponent,
    canActivate: [AuthGuard]
  },
  { path: "services", component: ServicesComponent, canActivate: [AuthGuard] }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
