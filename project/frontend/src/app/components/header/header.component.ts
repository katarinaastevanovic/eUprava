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
  role: 'STUDENT' | 'DOCTOR' | 'TEACHER' | null = null;

  constructor(
    private authService: AuthService,
    public router: Router,
    private studentNotificationService: StudentNotificationService,
    private doctorNotificationService: DoctorNotificationService
  ) { }

  ngOnInit(): void {
    this.authService.userRole$.subscribe(role => {
      this.role = role;
      this.userId = this.authService.getUserId();
      this.loadNotifications();
    });
  }

  logout() {
    this.authService.logout();
    this.router.navigate(['']);
  }

  isLoggedIn(): boolean {
    return !!this.authService.getRole();
  }

  loadNotifications() {
    if (!this.role) return;

    if (this.role === 'STUDENT') {
      this.studentNotificationService.getNotifications(this.userId).subscribe(
        res => this.notifications = this.sortNotifications(this.mapNotifications(res, 'STUDENT')),
        err => console.error('Error loading student notifications:', err)
      );
    } else if (this.role === 'DOCTOR') {
      this.doctorNotificationService.getNotifications(this.userId).subscribe(
        res => this.notifications = this.sortNotifications(this.mapNotifications(res, 'DOCTOR')),
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
      const service = this.role === 'STUDENT' ? this.studentNotificationService : this.doctorNotificationService;
      service.markAsRead(notif.userId, notif.id).subscribe({
        next: () => notif.read = true,
        error: err => console.error('Error marking notification as read:', err)
      });
    }
  }

  isStudent(): boolean {
    return this.role === 'STUDENT';
  }

  isDoctor(): boolean {
    return this.role === 'DOCTOR';
  }
}
