import { Component, OnInit, Input, ViewChild, ElementRef } from "@angular/core";
import * as types from "../types";

@Component({
  selector: "app-logs",
  templateUrl: "./logs.component.html",
  styleUrls: ["./logs.component.css"]
})
export class LogsComponent implements OnInit {
  @ViewChild("scrollMe", { static: false })
  private myScrollContainer: ElementRef;
  @Input() logs: types.LogRecord[] = [];

  constructor() {}

  ngOnInit() {
    // todo this is a disgusting hack and we need to find something better
    setTimeout(() => {
      this.scrollToBottom();
    }, 500);
  }

  ngOnChange() {}

  scrollToBottom(): void {
    try {
      this.myScrollContainer.nativeElement.scrollTop = this.myScrollContainer.nativeElement.scrollHeight;
    } catch (err) {}
  }
}
