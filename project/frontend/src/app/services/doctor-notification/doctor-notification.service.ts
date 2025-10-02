import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Notification } from '../../models/notification.model';

@Injectable({ providedIn: 'root' })
export class DoctorNotificationService {
  private baseUrl = 'http://localhost:8080/api/medical/notifications';

  constructor(private http: HttpClient) {}

  getNotifications(userId: number): Observable<Notification[]> {
  const token = localStorage.getItem('jwt') || '';
  console.log('Fetching doctor notifications for userId:', userId, 'with token:', token);

  const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);
  return this.http.get<Notification[]>(`${this.baseUrl}/${userId}`, { headers });
}

markAsRead(userId: number, notifId: number): Observable<void> {
  const token = localStorage.getItem('jwt') || '';
  console.log(`Marking doctor notification as read. userId: ${userId}, notifId: ${notifId}, token: ${token}`);

  const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);
  return this.http.patch<void>(`${this.baseUrl}/${userId}/${notifId}/read`, {}, { headers });
}

}
