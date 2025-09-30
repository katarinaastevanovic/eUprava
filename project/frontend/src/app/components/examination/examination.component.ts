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
  medicalRecordId!: number;
  showMedicalCertificateButton: boolean = false;

  examination: Examination = {
    requestId: 0,
    medicalRecordId: 0,
    diagnosis: '',
    therapy: '',
    note: ''
  };

  constructor(
    private route: ActivatedRoute,
    private examService: ExaminationService,
    private router: Router   // <-- dodan Router
  ) { }

  ngOnInit() {
    this.requestId = Number(this.route.snapshot.paramMap.get('id'));

    const saved = sessionStorage.getItem(`examination-${this.requestId}`);
    if (saved) {
      this.examination = JSON.parse(saved);
    } else {
      this.examination = {
        requestId: this.requestId,
        medicalRecordId: 0,
        diagnosis: '',
        therapy: '',
        note: ''
      };

      this.examService.getMedicalRecordIdByRequest(this.requestId).subscribe({
        next: recordId => this.examination.medicalRecordId = recordId
      });
    }

    this.examService.getRequestById(this.requestId).subscribe({
      next: req => this.showMedicalCertificateButton = !!req.needMedicalCertificate
    });
  }

  saveExamination() {
    this.examService.createExamination(this.examination).subscribe({
      next: res => {
        alert('Examination saved successfully!');
        sessionStorage.removeItem(`examination-${this.requestId}`);
        this.router.navigate(['/approved-requests']);
      },
      error: err => alert('Failed to save examination.')
    });
  }

  goToMedicalCertificate() {
    sessionStorage.setItem(
      `examination-${this.requestId}`,
      JSON.stringify(this.examination)
    );
    this.router.navigate(['/medical-certificate', this.requestId]);
  }

}
