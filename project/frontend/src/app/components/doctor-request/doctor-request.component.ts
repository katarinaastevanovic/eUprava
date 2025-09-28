import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { ExaminationRequestService, Request, Student } from '../../services/examination-request/examination-request.service';

@Component({
  selector: 'app-doctor-requests',
  standalone: true,
  imports: [CommonModule, HttpClientModule],
  templateUrl: './doctor-request.component.html',
})
export class DoctorRequestsComponent implements OnInit {
  requests: Request[] = [];
  students: Student[] = [];
  doctorId: number = 0;

  constructor(private requestService: ExaminationRequestService) {}

  ngOnInit() {
    this.doctorId = this.getUserIdFromToken();
    if (this.doctorId) {
      this.loadStudents(); 
      this.loadRequests(); 
    }
  }

  getUserIdFromToken(): number {
    const token = localStorage.getItem('jwt');
    if (!token) return 0;
    const payload = token.split('.')[1];
    if (!payload) return 0;
    const decoded = JSON.parse(atob(payload));
    return Number(decoded.sub) || 0;
  }

  loadStudents() {
  this.requestService.getAllStudents().subscribe(
    res => {
      this.students = res.map((s: any) => ({
        id: s.ID,
        firstName: s.Name,
        lastName: s.LastName
      }));
      console.log('Mapped students:', this.students);
    },
    err => console.error('Failed to load students:', err)
  );
}

  loadRequests() {
    this.requestService.getRequestsByDoctor(this.doctorId).subscribe(
      res => {
        this.requests = res.map((r: any) => ({
          id: r.ID,
          medicalRecordId: r.MedicalRecordId,
          doctorId: r.DoctorId,
          type: r.Type,
          status: r.Status
        }));
        console.log('Doctor requests:', this.requests);
      },
      err => console.error(err)
    );
  }

  getStudentName(id: number): string {
    const student = this.students.find(s => s.id === id);
    return student ? `${student.firstName} ${student.lastName}` : `ID: ${id}`;
  }

  approveRequest(requestId: number) {
    this.requestService.approveRequest(requestId).subscribe(
      res => {
        console.log('Request approved:', res);
        this.loadRequests();
      },
      err => console.error(err)
    );
  }

  rejectRequest(requestId: number) {
    this.requestService.rejectRequest(requestId).subscribe(
      res => {
        console.log('Request rejected:', res);
        this.loadRequests();
      },
      err => console.error(err)
    );
  }
}
