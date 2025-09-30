import { TestBed } from '@angular/core/testing';

import { MedicalCertificateService } from './medical-certificate.service';

describe('MedicalCertificateService', () => {
  let service: MedicalCertificateService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(MedicalCertificateService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
