import { TestBed } from '@angular/core/testing';

import { StudentNotificationService } from './student-notification.service';

describe('StudentNotificationService', () => {
  let service: StudentNotificationService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(StudentNotificationService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
