import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LogUserInComponent } from './log-user-in.component';

describe('LogUserInComponent', () => {
  let component: LogUserInComponent;
  let fixture: ComponentFixture<LogUserInComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LogUserInComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LogUserInComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
