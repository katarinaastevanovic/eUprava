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
      console.log('Doctors loaded:', this.doctors);
    },
    err => console.error(err)
  );
}


  submitRequestFromSelect(doctorIdString: string, typeString: string) {
    const doctorId = +doctorIdString; 
    const type = typeString as 'REGULAR' | 'SPECIALIST' | 'URGENT';
    this.submitRequest(doctorId, type);
  }

  submitRequest(doctorId: number, type: 'REGULAR' | 'SPECIALIST' | 'URGENT') {
    if (!doctorId || !type) {
      alert('Please select doctor and type');
      return;
    }

    const request: Request = {
      medicalRecordId: this.patientId,
      doctorId,
      type
    };

    this.requestService.createRequest(request).subscribe(
      res => {
        alert('Request created!');
      },
      err => console.error(err)
    );
  }

  getDoctorName(doctorId: number): string {
    const doctor = this.doctors.find(d => d.id === doctorId);
    return doctor ? `${doctor.name} ${doctor.lastName}` : '';
  }
}
