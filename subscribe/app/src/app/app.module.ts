import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { SubscribeFormComponent } from "./subscribe-form/subscribe-form.component";
import { Client as MicroClient } from "@microhq/ng-client";
import { RouterModule } from "@angular/router";

@NgModule({
  declarations: [AppComponent, SubscribeFormComponent],
  imports: [BrowserModule, AppRoutingModule, RouterModule],
  providers: [MicroClient],
  bootstrap: [AppComponent]
})
export class AppModule {}
