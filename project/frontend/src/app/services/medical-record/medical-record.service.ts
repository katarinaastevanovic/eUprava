import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { FullMedicalRecord } from '../../models/medical-models/full-medical-record.model';

@Injectable({
  providedIn: 'root'
})
export class MedicalRecordService {
  private apiUrl = 'http://localhost:8080/api/medical/medical-record';

  constructor(private http: HttpClient) {}

  getFullRecord(userId: number): Observable<FullMedicalRecord> {
  const token = localStorage.getItem('jwt');
  const headers = { Authorization: `Bearer ${token}` };
  return this.http.get<FullMedicalRecord>(`${this.apiUrl}/full/${userId}`, { headers });
}

createRecord(userId: number): Observable<FullMedicalRecord> {
  const token = localStorage.getItem('jwt');
  const headers = { Authorization: `Bearer ${token}` };
  return this.http.post<FullMedicalRecord>(
    `${this.apiUrl}/medical-records`,
    { userId },
    { headers }
  );
}

updateRecord(record: FullMedicalRecord) {
  const token = localStorage.getItem('jwt');
  const headers = { Authorization: `Bearer ${token}` };
  return this.http.put<FullMedicalRecord>(
    `${this.apiUrl}/${record.patientId}`,
    {
      allergies: record.allergies,
      chronicDiseases: record.chronicDiseases
    },
    { headers }
  );
}
}
