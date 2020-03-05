import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { SubscribeFormComponent } from "./subscribe-form/subscribe-form.component";
import { ClientModule, ClientService } from "@microhq/ng-client";
import { RouterModule } from "@angular/router";

@NgModule({
  declarations: [AppComponent, SubscribeFormComponent],
  imports: [BrowserModule, AppRoutingModule, RouterModule, ClientModule],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {}
