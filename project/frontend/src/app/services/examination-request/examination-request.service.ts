import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Request {
  id?: number;
  medicalRecordId: number;
  doctorId: number;
  type: 'REGULAR' | 'SPECIALIST' | 'URGENT';
  status?: 'REQUESTED' | 'APPROVED' | 'REJECTED';
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

  constructor(private http: HttpClient) {}

  createRequest(request: Request): Observable<Request> {
    return this.http.post<Request>(`${this.apiGatewayUrl}/medical/requests`, request);
  }

  getRequestsByPatient(patientId: number): Observable<Request[]> {
    return this.http.get<Request[]>(`${this.apiGatewayUrl}/medical/requests/patient/${patientId}`);
  }

  getRequestsByDoctor(doctorId: number): Observable<Request[]> {
    return this.http.get<Request[]>(`${this.apiGatewayUrl}/medical/requests/doctor/${doctorId}`);
  }

  approveRequest(requestId: number): Observable<void> {
    return this.http.patch<void>(`${this.apiGatewayUrl}/medical/requests/${requestId}/approve`, {});
  }

  rejectRequest(requestId: number): Observable<void> {
    return this.http.patch<void>(`${this.apiGatewayUrl}/medical/requests/${requestId}/reject`, {});
  }

  getDoctors(): Observable<Doctor[]> {
    return this.http.get<Doctor[]>(`${this.apiGatewayUrl}/auth/users/doctors`);
  }

  getAllStudents(): Observable<Student[]> {
    return this.http.get<Student[]>(`${this.apiGatewayUrl}/auth/users/students`);
  }
}
