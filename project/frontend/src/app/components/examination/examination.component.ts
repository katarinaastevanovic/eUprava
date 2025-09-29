import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
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
  medicalRecordId!: number; 
  examination: Examination = {
    requestId: 0,
    medicalRecordId: 0,
    diagnosis: '',
    therapy: '',
    note: ''
  };

  constructor(
    private route: ActivatedRoute,
    private examService: ExaminationService
  ) {}

  ngOnInit() {
  this.requestId = Number(this.route.snapshot.paramMap.get('id'));
  this.examination.requestId = this.requestId;

  console.log('Exam form for request ID:', this.requestId);

  this.examService.getMedicalRecordIdByRequest(this.requestId).subscribe({
    next: (recordId: number) => {
      this.examination.medicalRecordId = recordId;
      console.log('Loaded medicalRecordId:', recordId);
    },
    error: err => {
      console.error('Failed to load medicalRecordId:', err);
      alert('Cannot load medical record ID. Examination cannot be saved.');
    }
  });
}


  saveExamination() {
  console.log('Preparing to save examination:', this.examination);

  this.examService.createExamination(this.examination)
    .subscribe({
      next: res => {
        console.log('Examination saved successfully:', res);
        alert('Examination saved successfully!');
      },
      error: err => {
        console.error('Failed to save examination. Full error object:', err);

        if (err.error) {
          console.error('Error body from server:', err.error);
        }

        if (err.status) {
          console.error('HTTP status:', err.status);
        }

        if (err.message) {
          console.error('Error message:', err.message);
        }

        alert('Failed to save examination. Check console for details.');
      }
    });
}


  generateMedicalCertificate() {
    this.examService.generateMedicalCertificate(this.requestId)
      .subscribe(blob => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `medical_certificate_${this.requestId}.pdf`;
        a.click();
        window.URL.revokeObjectURL(url);
      }, err => {
        console.error('Failed to generate certificate:', err);
        alert('Failed to generate medical certificate.');
      });
  }
}
