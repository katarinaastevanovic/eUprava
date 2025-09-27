import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PatientRequestComponent } from './patient-request.component';

describe('PatientRequestComponent', () => {
  let component: PatientRequestComponent;
  let fixture: ComponentFixture<PatientRequestComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PatientRequestComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PatientRequestComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
