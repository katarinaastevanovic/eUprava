import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Notification } from '../../models/notification.model';

@Injectable({ providedIn: 'root' })
export class StudentNotificationService {
  private apiUrl = 'http://localhost:8080/api/school/'; 

  constructor(private http: HttpClient) {}

  getNotifications(userId: number): Observable<Notification[]> {
    const token = localStorage.getItem('jwt') || '';
    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);
    return this.http.get<Notification[]>(`${this.apiUrl}notifications/${userId}`, { headers });
  }

  markAsRead(userId: number, notifId: number): Observable<void> {
    const token = localStorage.getItem('jwt') || '';
    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);
    return this.http.put<void>(`${this.apiUrl}users/${userId}/notifications/${notifId}/read`, {}, { headers });
  }
}
