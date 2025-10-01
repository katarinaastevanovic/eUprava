import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Request {
  id?: number;
  medicalRecordId: number;
  doctorId: number;
  type: 'REGULAR' | 'SPECIALIST' | 'URGENT';
  status?: 'REQUESTED' | 'APPROVED' | 'REJECTED' | 'FINISHED';
  needMedicalCertificate?: boolean;
}

export interface RequestWithStudent {
  id?: number;
  medicalRecordId: number;
  doctorId: number;
  type: 'REGULAR' | 'SPECIALIST' | 'URGENT';
  status?: 'REQUESTED' | 'APPROVED' | 'REJECTED' | 'FINISHED';
  studentName: string;
  needMedicalCertificate?: boolean;
}

export interface Doctor {
  id: number;
  name: string;
  lastName: string;
  role: string;
}

export interface Student {
  id: number;
  firstName: string;
  lastName: string;
}

@Injectable({
  providedIn: 'root'
})
export class ExaminationRequestService {
  private apiGatewayUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) { }

  private getAuthHeaders(): { headers: HttpHeaders } {
    const token = localStorage.getItem('jwt');
    const headers = new HttpHeaders({ 'Authorization': `Bearer ${token}` });
    return { headers };
  }

  createRequest(request: Request): Observable<Request> {
    return this.http.post<Request>(`${this.apiGatewayUrl}/school/requests`, request, this.getAuthHeaders());
  }

  getRequestsByPatient(patientId: number): Observable<Request[]> {
    return this.http.get<Request[]>(`${this.apiGatewayUrl}/medical/requests/patient/${patientId}`, this.getAuthHeaders());
  }

  getRequestsByDoctorPaginated(
    doctorId: number,
    page: number,
    pageSize: number,
    search: string = '',
    status: string = '',
    type: string = ''
  ): Observable<{ requests: RequestWithStudent[], totalPages: number }> {
    let url = `${this.apiGatewayUrl}/medical/requests/doctor/${doctorId}?page=${page}&pageSize=${pageSize}`;
    if (search) url += `&search=${encodeURIComponent(search)}`;
    if (status) url += `&status=${encodeURIComponent(status)}`;
    if (type) url += `&type=${encodeURIComponent(type)}`;
    return this.http.get<{ requests: RequestWithStudent[], totalPages: number }>(url, this.getAuthHeaders());
  }

  approveRequest(requestId: number): Observable<void> {
    return this.http.patch<void>(`${this.apiGatewayUrl}/medical/requests/${requestId}/approve`, {}, this.getAuthHeaders());
  }

  rejectRequest(requestId: number): Observable<void> {
    return this.http.patch<void>(`${this.apiGatewayUrl}/medical/requests/${requestId}/reject`, {}, this.getAuthHeaders());
  }

  getDoctors(): Observable<Doctor[]> {
    return this.http.get<Doctor[]>(`${this.apiGatewayUrl}/auth/users/doctors`, this.getAuthHeaders());
  }

  getAllStudents(): Observable<Student[]> {
    return this.http.get<Student[]>(`${this.apiGatewayUrl}/auth/users/students`, this.getAuthHeaders());
  }

  getApprovedRequestsByDoctor(doctorId: number, page: number = 1, search: string = ''): Observable<{ requests: RequestWithStudent[]; totalPages: number }> {
    let url = `${this.apiGatewayUrl}/medical/requests/doctor/${doctorId}/approved?page=${page}`;
    if (search) url += `&search=${encodeURIComponent(search)}`;
    return this.http.get<{ requests: RequestWithStudent[]; totalPages: number }>(url, this.getAuthHeaders());
  }
}
