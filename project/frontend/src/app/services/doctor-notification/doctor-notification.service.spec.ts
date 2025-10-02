import { TestBed } from '@angular/core/testing';

import { DoctorNotificationService } from './doctor-notification.service';

describe('DoctorNotificationService', () => {
  let service: DoctorNotificationService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(DoctorNotificationService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
