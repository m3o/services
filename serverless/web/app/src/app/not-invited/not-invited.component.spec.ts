import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NotInvitedComponent } from './not-invited.component';

describe('NotInvitedComponent', () => {
  let component: NotInvitedComponent;
  let fixture: ComponentFixture<NotInvitedComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NotInvitedComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NotInvitedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
