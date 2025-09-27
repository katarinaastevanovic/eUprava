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
  private apiUrl = 'http://localhost:8082';

  constructor(private http: HttpClient) {}

  createRequest(request: Request): Observable<Request> {
    return this.http.post<Request>(`${this.apiUrl}/requests`, request);
  }

  getRequestsByPatient(patientId: number): Observable<Request[]> {
    return this.http.get<Request[]>(`${this.apiUrl}/requests/patient/${patientId}`);
  }

  getDoctors(): Observable<Doctor[]> {
    return this.http.get<Doctor[]>('http://localhost:8080/users/doctors');
  }

  getAllStudents(): Observable<Student[]> {
    return this.http.get<Student[]>(`http://localhost:8080/users/students`);
  }

  getRequestsByDoctor(doctorId: number) {
  return this.http.get<Request[]>(`${this.apiUrl}/requests/doctor/${doctorId}`);
  }

  approveRequest(requestId: number) {
    return this.http.patch(`${this.apiUrl}/requests/${requestId}/approve`, {});
  }

  rejectRequest(requestId: number) {
    return this.http.patch(`${this.apiUrl}/requests/${requestId}/reject`, {});
  }

}
