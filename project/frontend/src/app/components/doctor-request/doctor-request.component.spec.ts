import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DoctorRequestsComponent } from './doctor-request.component';

describe('DoctorRequestComponent', () => {
  let component: DoctorRequestsComponent;
  let fixture: ComponentFixture<DoctorRequestsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DoctorRequestsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DoctorRequestsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
