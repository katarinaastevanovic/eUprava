import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ExaminationRequestService, RequestWithStudent } from '../../services/examination-request/examination-request.service';

@Component({
  selector: 'app-doctor-approved-requests',
  standalone: true,
  imports: [CommonModule, HttpClientModule, FormsModule],
  templateUrl: './doctor-approved-requests.component.html',
})
export class DoctorApprovedRequestsComponent implements OnInit {
  approvedRequests: RequestWithStudent[] = [];
  totalPages: number = 0;
  currentPage: number = 1;
  doctorId: number = 0;
  searchTerm: string = ''; 

  constructor(
    private requestService: ExaminationRequestService,
    private router: Router
  ) { }

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

  loadApprovedRequests(page: number = 1) {
    this.currentPage = page;
    this.requestService
      .getApprovedRequestsByDoctor(this.doctorId, page, this.searchTerm)
      .subscribe(
        (res: any) => {
          this.approvedRequests = res.requests;
          this.totalPages = res.totalPages;
        },
        err => console.error('Failed to load approved requests:', err)
      );
  }

  onSearch() {
    this.loadApprovedRequests(1);
  }

  prevPage() {
    if (this.currentPage > 1) {
      this.loadApprovedRequests(this.currentPage - 1);
    }
  }

  nextPage() {
    if (this.currentPage < this.totalPages) {
      this.loadApprovedRequests(this.currentPage + 1);
    }
  }

  goToPage(page: number) {
    this.loadApprovedRequests(page);
  }

  reviewRequest(requestId: number) {
    this.router.navigate(['/examination', requestId]);
  }
}
