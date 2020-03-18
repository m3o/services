import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppFormComponent } from './app-form.component';

describe('AppFormComponent', () => {
  let component: AppFormComponent;
  let fixture: ComponentFixture<AppFormComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AppFormComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
