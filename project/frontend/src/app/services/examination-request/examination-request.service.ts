import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
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

  private getAuthHeaders(): { headers: HttpHeaders } {
    const token = localStorage.getItem('jwt');
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
    return { headers };
  }

  createRequest(request: Request): Observable<Request> {
    return this.http.post<Request>(
      `${this.apiGatewayUrl}/medical/requests`, 
      request, 
      this.getAuthHeaders()
    );
  }

  getRequestsByPatient(patientId: number): Observable<Request[]> {
    return this.http.get<Request[]>(
      `${this.apiGatewayUrl}/medical/requests/patient/${patientId}`,
      this.getAuthHeaders()
    );
  }

  getRequestsByDoctor(doctorId: number): Observable<Request[]> {
    return this.http.get<Request[]>(
      `${this.apiGatewayUrl}/medical/requests/doctor/${doctorId}`,
      this.getAuthHeaders()
    );
  }

  approveRequest(requestId: number): Observable<void> {
    return this.http.patch<void>(
      `${this.apiGatewayUrl}/medical/requests/${requestId}/approve`,
      {},
      this.getAuthHeaders()
    );
  }

  rejectRequest(requestId: number): Observable<void> {
    return this.http.patch<void>(
      `${this.apiGatewayUrl}/medical/requests/${requestId}/reject`,
      {},
      this.getAuthHeaders()
    );
  }

  getDoctors(): Observable<Doctor[]> {
    return this.http.get<Doctor[]>(
      `${this.apiGatewayUrl}/auth/users/doctors`,
      this.getAuthHeaders()
    );
  }

  getAllStudents(): Observable<Student[]> {
    return this.http.get<Student[]>(
      `${this.apiGatewayUrl}/auth/users/students`,
      this.getAuthHeaders()
    );
  }
}
