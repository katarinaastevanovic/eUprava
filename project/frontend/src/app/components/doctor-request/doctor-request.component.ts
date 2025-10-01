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

  constructor(private requestService: ExaminationRequestService) { }

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

  currentPage = 1;
  totalPages = 1;
  pageSize = 10;

  loadRequests() {
    this.requestService.getRequestsByDoctorPaginated(this.doctorId, this.currentPage, this.pageSize)
      .subscribe(res => {
        this.requests = res.requests;
        this.totalPages = res.totalPages;
      }, err => console.error(err));
  }

  prevPage() {
    if (this.currentPage > 1) {
      this.currentPage--;
      this.loadRequests();
    }
  }

  nextPage() {
    if (this.currentPage < this.totalPages) {
      this.currentPage++;
      this.loadRequests();
    }
  }

  goToPage(page: number) {
  this.currentPage = page;
  this.loadRequests();
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
