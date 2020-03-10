import * as tslib_1 from "tslib";
import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { SubscribeFormComponent } from "./subscribe-form/subscribe-form.component";
import { Client as MicroClient } from "@microhq/ng-client";
import { RouterModule } from "@angular/router";
let AppModule = class AppModule {
};
AppModule = tslib_1.__decorate([
    NgModule({
        declarations: [AppComponent, SubscribeFormComponent],
        imports: [BrowserModule, AppRoutingModule, RouterModule],
        providers: [MicroClient],
        bootstrap: [AppComponent]
    })
], AppModule);
export { AppModule };
//# sourceMappingURL=app.module.js.map