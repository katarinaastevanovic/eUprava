import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export type TypeOfCertificate = 'REGULAR' | 'PE' | 'EXCURSION' | 'SICKNESS';

export interface MedicalCertificate {
  requestId: number;
  patientId: number;
  doctorId: number;
  date: string;
  type: TypeOfCertificate;
  note?: string;
}

@Injectable({
  providedIn: 'root'
})
export class MedicalCertificateService {
  private http = inject(HttpClient);
  private apiUrl = 'http://localhost:8080/api/medical/certificates';

  private getAuthHeaders(): { headers: HttpHeaders } {
    const token = localStorage.getItem('jwt');
    return { headers: new HttpHeaders({ 'Authorization': `Bearer ${token}` }) };
  }

  createCertificate(cert: MedicalCertificate): Observable<MedicalCertificate> {
    return this.http.post<MedicalCertificate>(this.apiUrl, cert, this.getAuthHeaders());
  }
}
