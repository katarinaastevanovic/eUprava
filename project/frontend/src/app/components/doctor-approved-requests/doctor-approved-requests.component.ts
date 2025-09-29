import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { ExaminationRequestService, RequestWithStudent } from '../../services/examination-request/examination-request.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-doctor-approved-requests',
  standalone: true,
  imports: [CommonModule, HttpClientModule],
  templateUrl: './doctor-approved-requests.component.html',
})
export class DoctorApprovedRequestsComponent implements OnInit {
  approvedRequests: RequestWithStudent[] = [];
  doctorId: number = 0;

  constructor(
    private requestService: ExaminationRequestService, 
    private router: Router
  ) {}

  ngOnInit() {
    this.doctorId = this.getUserIdFromToken();
    if (this.doctorId) {
      this.loadApprovedRequests();
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

  loadApprovedRequests() {
    this.requestService.getApprovedRequestsByDoctor(this.doctorId).subscribe(
      (res: RequestWithStudent[]) => {
        console.log('Approved requests:', res);
        this.approvedRequests = res;
      },
      err => console.error('Failed to load approved requests:', err)
    );
  }

  reviewRequest(requestId: number) {
    this.router.navigate(['/examination', requestId]);
  }
}
