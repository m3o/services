import {
  Component,
  OnInit,
  ViewEncapsulation,
  ViewChild,
  ElementRef
} from "@angular/core";
import { FormControl, Validators } from "@angular/forms";
import { Location } from "@angular/common";
import { UserService } from "../user.service";
import { ServiceService } from "../service.service";
import * as types from "../types";
import { Router, ActivatedRoute } from "@angular/router";
import * as _ from "lodash";
import * as rxjs from "rxjs";
import { debounceTime } from "rxjs/operators";
import { NotificationsService } from "angular2-notifications";

@Component({
  selector: "app-new-service",
  templateUrl: "./new-service.component.html",
  styleUrls: ["./new-service.component.css"],
  encapsulation: ViewEncapsulation.None
})
export class NewServiceComponent implements OnInit {
  serviceInput = new FormControl("", [Validators.required]);
  @ViewChild("sinput", { static: false }) sinput: ElementRef;
  alias = "";
  namespace = "go.micro";
  serviceType = "srv";
  serviceName = "";
  code: string = "";
  runCode: string = "";
  token = "";
  intervalId: any;
  buildTimerIntervalId: any;
  lastKeypress = new Date();
  loadingServices = false;
  loaded = false;
  events: types.Event[] = [];
  services: types.Service[] = [];
  lastInput;
  step = 0;
  // approximate time it will take to finisht the build
  maxBuildTimer = 60;
  serviceExists = false;
  minBuildTimer = 5;
  // no id for events so have to use timestamp
  lastBuildFailureTimestamp = 0;
  failureEvent: types.Event;
  buildTimer = this.maxBuildTimer;
  progressPercentage = 0;
  percentages = [0, 10, 20, 80];
  eventErrored = false;
  stepLabels = (): string[] => {
    return [
      "We are waiting for you to push your service...",
      "Found your service on GitHub. Waiting for the build to start...",
      "Build is in progress. Waiting for the build to finish...",
      "Build finished. Waiting for your service to start...",
      "Ready to roll! Redirecting you to your service page..."
    ];
  };

  constructor(
    private us: UserService,
    private ses: ServiceService,
    private router: Router,
    private location: Location,
    private activeRoute: ActivatedRoute,
    private notif: NotificationsService
  ) {}

  ngOnInit() {
    this.alias = this.us.user.login;
    this.serviceInput.markAsTouched();
    this.activeRoute.params.subscribe(p => {
      const id = <string>p["id"];
      if (id) {
        this.alias = _.last(id.split("."));
      }
      this.regen();
      this.serviceInput.markAsTouched();
    });

    this.token = this.us.token();
    this.serviceName =
      this.namespace + "." + this.serviceType + "." + this.alias;
    this.location.replaceState("/service/new/" + this.serviceName);

    this.loadAll(true);
    this.intervalId = setInterval(() => {
      this.loadAll();
    }, 1500);
    this.progressPercentage = this.percentages[this.step];
  }

  ngAfterViewInit() {
    const source = rxjs.fromEvent(this.sinput.nativeElement, "keyup");
    source.pipe(debounceTime(600)).subscribe(c => {
      this.loadAll(true);
    });
  }

  loadAll(setLoader?: boolean) {
    if (setLoader) {
      this.loadingServices = true;
    }
    this.ses
      .events(this.serviceName)
      .then(events => {
        this.events = events;
        this.checkEvents();
      })
      .catch(e => {
        if (this.eventErrored) {
          return;
        }
        this.eventErrored = true;
        let errMsg = "";
        try {
          errMsg = JSON.parse(e.error.error).detail;
        } catch (e) {}
        this.notif.error("Error listing events", errMsg);
      });
    this.ses.list().then(services => {
      this.services = services;
      this.checkServices(setLoader);
      if (setLoader) {
        this.loadingServices = false;
      }
      this.loaded = true;
    });
  }

  keyPress(event: any) {
    this.lastKeypress = new Date();
    this.location.replaceState("/service/new/" + this.serviceName);
  }

  checkEvents() {
    if (!this.events || this.events.length == 0) {
      return;
    }
    const e = _.last(_.orderBy(this.events, e => e.timestamp, "asc"));
    if (e.service.name != this.serviceName) {
      return;
    }
    // source updated
    if (e.type == 4) {
      this.step = 1;
      this.progressPercentage = this.percentages[1];
    }
    // build started and in progress
    if (e.type == 5) {
      this.step = 2;
      this.buildTimer = this.maxBuildTimer;
      this.progressPercentage = this.percentages[2];
      this.startBuildTimer(e);
    }
    // build finished
    if (e.type == 6) {
      this.step = 3;
      this.progressPercentage = this.percentages[3];
    }
    // build failure
    if (e.type == 7) {
      this.step = 1;
      if (this.buildTimerIntervalId) {
        this.stopBuildTimer();
      }
      if (this.progressPercentage == 0) {
        this.progressPercentage = this.percentages[1];
      }
      if (this.lastBuildFailureTimestamp == e.timestamp) {
        return;
      }
      this.lastBuildFailureTimestamp = e.timestamp;
      this.failureEvent = e;
      const buildNum =
        e.metadata && e.metadata["build"] ? e.metadata["build"] : "";
      this.notif.error(
        "Build failed",
        'Please see build <a target="_blank" href="' +
          this.buildUrl(e) +
          '">' +
          buildNum +
          "</a>"
      );
    }
  }

  // todo: this is copypasted from events-list.compinent.ts, fix that
  buildUrl(e: types.Event): string {
    if (!e.metadata) {
      return "";
    }
    const repo = e.metadata["repo"];
    const buildId = e.metadata["build"];
    // eg. https://github.com/micro/services/runs/466859781
    return "https://" + repo + "/actions/runs/" + buildId;
  }

  stopBuildTimer() {
    if (this.buildTimerIntervalId) {
      clearInterval(this.buildTimerIntervalId);
      this.buildTimerIntervalId = null;
    }
  }

  // the timer will only kick off after step 2
  startBuildTimer(e: types.Event) {
    const intervalSecs = 0.1;
    const secRange = this.maxBuildTimer - this.minBuildTimer;
    const secsSinceBuild =
      (new Date().getTime() - new Date(e.timestamp * 1000).getTime()) / 1000;
    if (secsSinceBuild > secRange) {
      this.buildTimer = this.minBuildTimer;
      this.progressPercentage = this.percentages[3];
      return;
    }

    const ratio = secsSinceBuild / secRange;
    this.buildTimer -= secRange * ratio;
    const percentageRange = this.percentages[3] - this.percentages[2];
    this.progressPercentage = this.percentages[2] + percentageRange * ratio;

    const percentageStep = secRange / intervalSecs;
    if (this.buildTimerIntervalId) {
      return;
    }
    this.buildTimerIntervalId = setInterval(() => {
      if (this.step !== 2) {
        return;
      }
      // the numbers below will depend heavily on the interval parameter of
      // the setInterval function

      if (this.buildTimer - intervalSecs <= this.minBuildTimer) {
        this.buildTimer = this.minBuildTimer;
        this.stopBuildTimer();
        return;
      }
      this.buildTimer -= intervalSecs;

      // calculating how much to add based on the difference in percentage between
      // the third and second step.
      this.progressPercentage +=
        (this.percentages[3] - this.percentages[2]) / percentageStep;
    }, intervalSecs * 1000);
  }

  checkServices(setExists?: boolean) {
    const inRegistry =
      this.services.filter(e => {
        return e.name == this.serviceName;
      }).length > 0;
    if (!inRegistry) {
      if (setExists) {
        this.serviceExists = false;
        this.serviceExists = false;
        this.serviceInput.markAsTouched();
      }
      return;
    }
    if (this.step == 0) {
      // Only checking for service exist on step 0
      // to support "service already exists "
      if (setExists) {
        this.serviceExists = true;
        this.serviceInput.setErrors({ incorrect: true });
        this.serviceInput.markAsTouched();
      }
    } else if (this.step < 4) {
      this.step = 4;
      this.progressPercentage = 100;
      setTimeout(() => {
        this.router.navigate(["/service/" + this.serviceName]);
      }, 3000);
    }
  }

  ngOnDestroy() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }

  regen() {
    this.serviceName =
      this.namespace + "." + this.serviceType + "." + this.alias;
    this.newCode();
    this.newRunCode();
  }

  languages = ["bash"];

  newCode() {
    this.code =
      `# Checkout the services repo
git clone https://github.com/micro/services && cd services
# Create new service
micro new ` +
      this.alias +
      `
cd ` +
      this.alias +
      `
# Build the service
make build
# Push to GitHub
git config --local core.hooksPath .githooks
git add . && git commit -m "Initialising service ` +
      this.alias +
      `" && git push`;
  }

  newRunCode() {
    this.runCode = `micro run --platform ` + this.alias;
  }
}
