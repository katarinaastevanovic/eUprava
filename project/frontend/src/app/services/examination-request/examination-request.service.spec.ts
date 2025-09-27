import { TestBed } from '@angular/core/testing';

import { ExaminationRequestService } from '../examination-request/examination-request.service';

describe('ExaminationRequestService', () => {
  let service: ExaminationRequestService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ExaminationRequestService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
