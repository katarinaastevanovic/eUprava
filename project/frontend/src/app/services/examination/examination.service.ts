import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Examination {
  requestId: number;
  medicalRecordId: number;
  diagnosis: string;
  therapy: string;
  note: string;
}

@Injectable({
  providedIn: 'root'
})
export class ExaminationService {
  private http = inject(HttpClient);
  private apiUrl = 'http://localhost:8080/api/medical/examinations';

  private getAuthHeaders(): { headers: HttpHeaders } {
    const token = localStorage.getItem('jwt');
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
    return { headers };
  }

  createExamination(exam: Examination): Observable<Examination> {
    return this.http.post<Examination>(this.apiUrl, exam, this.getAuthHeaders());
  }

  getExaminationByRequest(requestId: number): Observable<Examination> {
  return this.http.get<Examination>(`${this.apiUrl}/${requestId}`, this.getAuthHeaders());
  }

  generateMedicalCertificate(requestId: number): Observable<Blob> {
    return this.http.get(`${this.apiUrl}/${requestId}/certificate`, {
      ...this.getAuthHeaders(),
      responseType: 'blob'
    });
  }

  getMedicalRecordIdByRequest(requestId: number): Observable<number> {
  return this.http.get<number>(`http://localhost:8080/api/medical/requests/${requestId}/medical-record-id`, this.getAuthHeaders());
  }

}
