import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { RouterModule } from '@angular/router';
import { ExaminationRequestService, Request, Doctor } from '../../services/examination-request/examination-request.service';

@Component({
  selector: 'app-patient-requests',
  standalone: true,
  imports: [CommonModule, HttpClientModule, RouterModule],
  templateUrl: './patient-request.component.html',
})

export class PatientRequestsComponent implements OnInit {
  requests: Request[] = [];
  doctors: Doctor[] = [];
  patientId: number = 0;

  constructor(private requestService: ExaminationRequestService) {}

  ngOnInit() {
  this.loadDoctors();
  this.loadRequests();
}

  getUserIdFromToken(): number {
    const token = localStorage.getItem('jwt');
    if (!token) return 0;
    const payload = token.split('.')[1];
    if (!payload) return 0;
    const decoded = JSON.parse(atob(payload));
    return decoded.sub;
  }

  loadRequests() {
  this.requestService.getRequestsByPatient().subscribe(
    res => {
      this.requests = res.map((r: any) => ({
        id: r.id,
        medicalRecordId: r.medicalRecordId,
        doctorId: r.doctorId,
        type: r.type,
        status: r.status
      }));
      console.log('Mapped requests:', this.requests);
    },
    err => console.error(err)
  );
}

  loadDoctors() {
    this.requestService.getDoctors().subscribe(
      res => {
        this.doctors = res.map((d: any) => ({
          id: d.ID,
          name: d.Name,
          lastName: d.LastName,
          role: d.Role,
          email: d.Email
        }));
        console.log('Mapped doctors:', this.doctors);
      },
      err => console.error(err)
    );
  }

  getDoctorName(doctorId: number): string {
  const doctor = this.doctors.find(d => d.id === doctorId);
  return doctor ? `${doctor.name} ${doctor.lastName}` : 'Unknown Doctor';
}
}
