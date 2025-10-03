import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ExaminationService, Examination } from '../../services/examination/examination.service';

@Component({
  selector: 'app-examination',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './examination.component.html',
})
export class ExaminationFormComponent implements OnInit {
  requestId!: number;
  showMedicalCertificateButton: boolean = false;

  examination: Examination = {
    requestId: 0,
    medicalRecordId: 0,
    diagnosis: '',
    therapy: '',
    note: ''
  };

  medicalRecord: any;

  constructor(
    private route: ActivatedRoute,
    private examService: ExaminationService,
    private router: Router
  ) { }

  ngOnInit() {
    this.requestId = Number(this.route.snapshot.paramMap.get('id'));
    console.log('Request ID from route:', this.requestId);

    const saved = sessionStorage.getItem(`examination-${this.requestId}`);
    if (saved) {
      this.examination = JSON.parse(saved);
      if (this.examination.medicalRecordId) {
        this.loadFullMedicalRecord(this.examination.medicalRecordId);
      }
    } else {
      this.examination = {
        requestId: this.requestId,
        medicalRecordId: 0,
        diagnosis: '',
        therapy: '',
        note: ''
      };

      this.examService.getMedicalRecordIdByRequest(this.requestId).subscribe({
        next: recordId => {
          console.log('MedicalRecordId fetched:', recordId);
          this.examination.medicalRecordId = recordId;
          this.loadFullMedicalRecord(recordId);
        },
        error: err => console.error('Failed to get medicalRecordId:', err)
      });
    }

    this.examService.getRequestById(this.requestId).subscribe({
      next: req => {
        console.log('Request fetched:', req);
        this.showMedicalCertificateButton = !!req.needMedicalCertificate;
        console.log('showMedicalCertificateButton value:', this.showMedicalCertificateButton);
      },
      error: err => console.error('Failed to load request info:', err)
    });
  }


  private loadFullMedicalRecord(medicalRecordId: number) {
    this.examService.getFullMedicalRecordById(medicalRecordId).subscribe({
      next: record => this.medicalRecord = record,
      error: err => console.error('Failed to load medical record:', err)
    });
  }

  saveExamination() {
    this.examService.createExamination(this.examination).subscribe({
      next: res => {
        alert('Examination saved successfully!');
        sessionStorage.removeItem(`examination-${this.requestId}`);
        window.location.href = 'http://localhost:4200/approved-requests';
      },
      error: err => {
        console.error('Failed to save examination:', err);
        alert('Failed to save examination.');
      }
    });
  }

  goToMedicalCertificate() {
    sessionStorage.setItem(`examination-${this.requestId}`, JSON.stringify(this.examination));
    this.router.navigate(['/medical-certificate', this.requestId]);
  }
}
