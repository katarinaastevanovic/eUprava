import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { AuthService } from '../../services/auth/auth.service';
import { StudentNotificationService } from '../../services/student-notification/student-notification.service';
import { DoctorNotificationService } from '../../services/doctor-notification/doctor-notification.service';
import { Notification } from '../../models/notification.model';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './header.component.html'
})
export class HeaderComponent implements OnInit {
  showNotificationsSidebar = false;
  notifications: Notification[] = [];
  userId: number = 0;
  role: 'STUDENT' | 'DOCTOR' | null = null;

  constructor(
    private authService: AuthService,
    public router: Router,
    private studentNotificationService: StudentNotificationService,
    private doctorNotificationService: DoctorNotificationService
  ) { }

  ngOnInit(): void {
    this.userId = this.getUserIdFromToken();
    this.role = this.getRoleFromToken();
    console.log('UserId:', this.userId, 'Role:', this.role);
    this.loadNotifications();
  }

  logout() {
    this.authService.logout().then(() => {
      this.router.navigate(['']).then(() => window.location.reload());
    });
  }

  isLoggedIn(): boolean {
    return !!localStorage.getItem('jwt');
  }

  getUserIdFromToken(): number {
    const token = localStorage.getItem('jwt');
    if (!token) return 0;
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.sub;
  }

  getRoleFromToken(): 'STUDENT' | 'DOCTOR' | null {
    const token = localStorage.getItem('jwt');
    if (!token) return null;
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      if (Array.isArray(payload.roles)) {
        if (payload.roles.includes('STUDENT')) return 'STUDENT';
        if (payload.roles.includes('DOCTOR')) return 'DOCTOR';
        return null;
      }
      return payload.role?.toUpperCase() === 'STUDENT' ? 'STUDENT' :
        payload.role?.toUpperCase() === 'DOCTOR' ? 'DOCTOR' : null;
    } catch (err) {
      console.error('Failed to decode token', err);
      return null;
    }
  }

  loadNotifications() {
    if (!this.role) return;

    console.log('Loading notifications for userId:', this.userId, 'role:', this.role);

    if (this.role === 'STUDENT') {
      this.studentNotificationService.getNotifications(this.userId).subscribe(
        res => {
          console.log('Student notifications received:', res);
          this.notifications = this.sortNotifications(this.mapNotifications(res, 'STUDENT'));
        },
        err => console.error('Error loading student notifications:', err)
      );
    } else if (this.role === 'DOCTOR') {
      this.doctorNotificationService.getNotifications(this.userId).subscribe(
        res => {
          console.log('Doctor notifications received:', res);
          this.notifications = this.sortNotifications(this.mapNotifications(res, 'DOCTOR'));
        },
        err => console.error('Error loading doctor notifications:', err)
      );
    }
  }

  private mapNotifications(res: any[], role: 'STUDENT' | 'DOCTOR'): Notification[] {
    return res.map(n => ({
      id: role === 'STUDENT' ? n.ID : n.id,
      message: n.message,
      read: n.read,
      userId: n.userId,
      createdAt: n.CreatedAt || n.createdAt,
      updatedAt: n.UpdatedAt || n.updatedAt
    }));
  }

  private sortNotifications(notifs: Notification[]): Notification[] {
    return notifs.sort((a, b) => (a.read === b.read ? 0 : a.read ? 1 : -1));
  }

  get unreadCount(): number {
    return this.notifications.filter(n => !n.read).length;
  }

  openNotificationPanel() {
    this.showNotificationsSidebar = true;
  }

  closeNotificationPanel() {
    this.showNotificationsSidebar = false;
  }

  markAsRead(notif: Notification) {
    if (!notif.read && notif.id !== undefined) {
      if (this.role === 'STUDENT') {
        this.studentNotificationService.markAsRead(notif.userId, notif.id).subscribe({
          next: () => notif.read = true,
          error: err => console.error('Error marking student notification as read:', err)
        });
      } else if (this.role === 'DOCTOR') {
        this.doctorNotificationService.markAsRead(notif.userId, notif.id).subscribe({
          next: () => notif.read = true,
          error: err => console.error('Error marking doctor notification as read:', err)
        });
      }
    }
  }

  isStudent(): boolean {
    return this.role === 'STUDENT';
  }

  isDoctor(): boolean {
    return this.role === 'DOCTOR';
  }
}
