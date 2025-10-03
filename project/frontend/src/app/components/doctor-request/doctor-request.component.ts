import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { ExaminationRequestService, RequestWithStudent } from '../../services/examination-request/examination-request.service';

@Component({
  selector: 'app-doctor-requests',
  standalone: true,
  imports: [CommonModule, HttpClientModule, FormsModule],
  templateUrl: './doctor-request.component.html',
})
export class DoctorRequestsComponent implements OnInit {
  requests: RequestWithStudent[] = [];
  doctorId: number = 0;

  currentPage = 1;
  totalPages = 1;
  pageSize = 10;
  searchTerm = '';
  sortPending = false;

  statusFilter = '';
  typeFilter = '';
  typesOfExamination: string[] = ['REGULAR', 'SPECIALIST', 'URGENT'];

  constructor(private requestService: ExaminationRequestService) { }

  ngOnInit() {
    this.doctorId = this.getUserIdFromToken();
    if (this.doctorId) this.loadRequests();
  }

  getUserIdFromToken(): number {
    const token = localStorage.getItem('jwt');
    if (!token) return 0;
    const payload = token.split('.')[1];
    if (!payload) return 0;
    const decoded = JSON.parse(atob(payload));
    return Number(decoded.sub) || 0;
  }

  loadRequests(page: number = 1) {
    this.currentPage = page;
    const sortParam = this.sortPending ? 'requestedFirst' : '';

    console.log('Loading requests:', {
      page,
      sortParam,
      search: this.searchTerm,
      status: this.statusFilter,
      type: this.typeFilter
    });

    this.requestService.getRequestsByDoctorPaginated(
      this.doctorId,
      page,
      this.pageSize,
      this.searchTerm,
      this.statusFilter,
      this.typeFilter,
      sortParam
    ).subscribe(res => {
      this.requests = res.requests;
      this.totalPages = res.totalPages;
    }, err => console.error(err));
  }

  onSearch() {
    this.loadRequests(1);
  }

  onFilterChange() {
    this.loadRequests(1);
  }

  prevPage() {
    if (this.currentPage > 1) this.loadRequests(this.currentPage - 1);
  }

  nextPage() {
    if (this.currentPage < this.totalPages) this.loadRequests(this.currentPage + 1);
  }

  goToPage(page: number) {
    if (page !== -1) this.loadRequests(page);
  }

  approveRequest(requestId: number) {
    this.requestService.approveRequest(requestId).subscribe(() => this.loadRequests(this.currentPage));
  }

  rejectRequest(requestId: number) {
    this.requestService.rejectRequest(requestId).subscribe(() => this.loadRequests(this.currentPage));
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
