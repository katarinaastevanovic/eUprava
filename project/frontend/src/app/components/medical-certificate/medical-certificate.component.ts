import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';  
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MedicalCertificateService, MedicalCertificate, TypeOfCertificate } from '../../services/medical-certificate/medical-certificate.service';

@Component({
  selector: 'app-medical-certificate',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './medical-certificate.component.html'
})
export class MedicalCertificateComponent implements OnInit {
  requestId!: number;
  patientId!: number;
  doctorId!: number;

  certificate: MedicalCertificate = {
    requestId: 0,
    patientId: 0,
    doctorId: 0,
    date: new Date().toISOString().split('T')[0],
    type: 'REGULAR',
    note: ''
  };

  certificateTypes: TypeOfCertificate[] = ['REGULAR', 'PE', 'EXCURSION', 'SICKNESS'];

  constructor(
    private route: ActivatedRoute,
    private certService: MedicalCertificateService,
    private router: Router   
  ) {}

  ngOnInit() {
    this.requestId = Number(this.route.snapshot.paramMap.get('requestId'));
    this.certificate.requestId = this.requestId;

    this.certificate.patientId = 14;
    this.certificate.doctorId = 4;
  }

  saveCertificate() {
    this.certService.createCertificate(this.certificate).subscribe({
      next: res => {
        alert('Medical certificate saved successfully!');
        console.log('Saved certificate:', res);

        this.router.navigate(['/examination', this.requestId]);
      },
      error: err => {
        console.error('Failed to save certificate:', err);
        alert('Failed to save certificate. Check console for details.');
      }
    });
  }
}
