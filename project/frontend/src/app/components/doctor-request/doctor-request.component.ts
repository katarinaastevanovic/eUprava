import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { ExaminationRequestService, RequestWithStudent } from '../../services/examination-request/examination-request.service';

@Component({
  selector: 'app-doctor-requests',
  standalone: true,
  imports: [CommonModule, HttpClientModule],
  templateUrl: './doctor-request.component.html',
})
export class DoctorRequestsComponent implements OnInit {
  requests: RequestWithStudent[] = [];
  doctorId: number = 0;

  constructor(private requestService: ExaminationRequestService) {}

  ngOnInit() {
    this.doctorId = this.getUserIdFromToken();
    if (this.doctorId) {
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

  loadRequests() {
    this.requestService.getRequestsByDoctor(this.doctorId).subscribe(
      (res: RequestWithStudent[]) => {
        console.log('Requests with student names:', res);
        this.requests = res;
      },
      err => console.error('Failed to load requests:', err)
    );
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
