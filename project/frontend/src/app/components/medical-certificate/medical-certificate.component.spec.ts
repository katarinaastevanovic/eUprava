import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MedicalCertificateComponent } from './medical-certificate.component';

describe('MedicalCertificateComponent', () => {
  let component: MedicalCertificateComponent;
  let fixture: ComponentFixture<MedicalCertificateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MedicalCertificateComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MedicalCertificateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
