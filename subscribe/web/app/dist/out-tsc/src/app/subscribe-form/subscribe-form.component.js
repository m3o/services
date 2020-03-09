import * as tslib_1 from "tslib";
import { Component } from "@angular/core";
let SubscribeFormComponent = class SubscribeFormComponent {
    constructor(mc) {
        this.mc = mc;
    }
    ngOnInit() {
        this.mc.call("go.micro.srv.greeter", "Say.Hello");
    }
};
SubscribeFormComponent = tslib_1.__decorate([
    Component({
        selector: "app-subscribe-form",
        templateUrl: "./subscribe-form.component.html",
        styleUrls: ["./subscribe-form.component.css"],
        providers: []
    })
], SubscribeFormComponent);
export { SubscribeFormComponent };
//# sourceMappingURL=subscribe-form.component.js.map