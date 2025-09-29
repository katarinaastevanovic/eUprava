import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DoctorApprovedRequestsComponent } from './doctor-approved-requests.component';

describe('DoctorApprovedRequestsComponent', () => {
  let component: DoctorApprovedRequestsComponent;
  let fixture: ComponentFixture<DoctorApprovedRequestsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DoctorApprovedRequestsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DoctorApprovedRequestsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
