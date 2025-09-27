import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ExaminationRequestComponent } from './examination-request.component';

describe('ExaminationRequestComponent', () => {
  let component: ExaminationRequestComponent;
  let fixture: ComponentFixture<ExaminationRequestComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ExaminationRequestComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ExaminationRequestComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
