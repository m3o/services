import { async, ComponentFixture, TestBed } from "@angular/core/testing";

import { AppSingleComponent } from "./app-single.component";

describe("AppSingleComponent", () => {
  let component: AppSingleComponent;
  let fixture: ComponentFixture<AppSingleComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [AppSingleComponent]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppSingleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
