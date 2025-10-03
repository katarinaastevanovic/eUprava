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
  typeFilter: string = '';
  typesOfExamination: string[] = ['REGULAR', 'SPECIALIST', 'URGENT'];

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
      .getApprovedRequestsByDoctorFiltered(this.doctorId, page, this.searchTerm, this.typeFilter)
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
    if (page !== -1) this.loadApprovedRequests(page);
  }

  reviewRequest(requestId: number) {
    this.router.navigate(['/examination', requestId]);
  }

  onFilterChange() {
    this.loadApprovedRequests(1);
  }

  getPagesToShow(): number[] {
    const pages: number[] = [];
    const total = this.totalPages;
    const current = this.currentPage;

    if (total <= 3) {
      for (let i = 1; i <= total; i++) pages.push(i);
    } else {
      pages.push(1);

      if (current > 3) pages.push(-1);

      const start = Math.max(2, current - 1);
      const end = Math.min(total - 1, current + 1);

      for (let i = start; i <= end; i++) pages.push(i);

      if (current < total - 2) pages.push(-1);

      pages.push(total);
    }

    return pages;
  }
}
