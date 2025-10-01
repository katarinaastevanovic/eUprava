import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface StudentNotification {
  ID?: number;
  userId: number;
  message: string;
  read: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class StudentNotificationService {
  private apiUrl = 'http://localhost:8080/api/school/'; 

  constructor(private http: HttpClient) {}

  getNotifications(studentId: number): Observable<StudentNotification[]> {
    return this.http.get<StudentNotification[]>(`${this.apiUrl}notifications/${studentId}`);
  }

  markAsRead(userId: number, notifId: number): Observable<any> {
  return this.http.put(`${this.apiUrl}users/${userId}/notifications/${notifId}/read`, {});
}
}
