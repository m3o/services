import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StatChartsComponent } from './stat-charts.component';

describe('StatChartsComponent', () => {
  let component: StatChartsComponent;
  let fixture: ComponentFixture<StatChartsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StatChartsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StatChartsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
