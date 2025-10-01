import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ExaminationRequestService, Request, Doctor } from '../../services/examination-request/examination-request.service';

@Component({
  selector: 'app-request',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './examination-request.component.html',
  styleUrls: ['./examination-request.component.css']
})
export class ExaminationRequestComponent implements OnInit {
  requests: Request[] = [];
  doctors: Doctor[] = [];
  patientId: number;
  needMedicalCertificate: boolean = false;
  errorMessage: string = '';
  successMessage: string = '';

  constructor(private requestService: ExaminationRequestService) {
    this.patientId = this.getUserIdFromToken();
  }

  ngOnInit(): void {
    this.loadDoctors();
  }

  getUserIdFromToken(): number {
    const token = localStorage.getItem('jwt');
    if (!token) return 0;
    const payload = token.split('.')[1];
    if (!payload) return 0;
    const decoded = JSON.parse(atob(payload));
    return decoded.sub;
  }

  loadDoctors() {
    this.requestService.getDoctors().subscribe(
      (res: any[]) => {
        this.doctors = res.map((d: any) => ({
          id: d.ID,
          name: d.Name,
          lastName: d.LastName,
          role: d.Role,
          email: d.Email
        }));
      },
      err => console.error(err)
    );
  }

  submitRequestFromSelect(doctorIdString: string, typeString: string) {
    const doctorId = +doctorIdString;
    const type = typeString as 'REGULAR' | 'SPECIALIST' | 'URGENT';
    this.submitRequest(doctorId, type, this.needMedicalCertificate);
  }

  submitRequest(doctorId: number, type: 'REGULAR' | 'SPECIALIST' | 'URGENT', needCertificate: boolean) {
    this.errorMessage = '';
    this.successMessage = '';

    if (!doctorId || !type) {
      this.errorMessage = 'Please select doctor and type';
      setTimeout(() => this.errorMessage = '', 5000);
      return;
    }

    const request: Request = {
      medicalRecordId: this.patientId,
      doctorId,
      type,
      needMedicalCertificate: needCertificate
    };

    this.requestService.createRequest(request).subscribe(
      res => {
        this.successMessage = 'Request successfully created!';

        setTimeout(() => this.successMessage = '', 2000);

        const doctorSelect = (document.querySelector('#doctorSelect') as HTMLSelectElement);
        const typeSelect = (document.querySelector('#typeSelect') as HTMLSelectElement);
        if (doctorSelect) doctorSelect.value = '';
        if (typeSelect) typeSelect.value = '';
        this.needMedicalCertificate = false;

      },
      err => {
        if (err.status === 429) {
          this.errorMessage = 'Too many requests! Please wait a minute before trying again.';
        } else if (err.error) {
          this.errorMessage = 'Error: ' + err.error;
        } else {
          this.errorMessage = 'An unexpected error occurred.';
        }

        setTimeout(() => this.errorMessage = '', 5000);
      }
    );
  }
}

