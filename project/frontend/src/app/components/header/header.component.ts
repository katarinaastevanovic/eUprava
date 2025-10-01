import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { AuthService } from '../../services/auth/auth.service';
import { StudentNotificationService, StudentNotification } from '../../services/student-notification/student-notification.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './header.component.html'
})
export class HeaderComponent implements OnInit {
  showNotificationsSidebar = false;
  notifications: StudentNotification[] = [];
  studentId: number;

  constructor(private authService: AuthService, public router: Router, private notificationService: StudentNotificationService) {
    this.studentId = this.getUserIdFromToken();
  }

  ngOnInit(): void {
    this.loadNotifications();
  }

  logout() {
    this.authService.logout().then(() => {
      this.router.navigate(['/login']);
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

  loadNotifications() {
  this.notificationService.getNotifications(this.studentId).subscribe(
    res => {
      this.notifications = res.sort((a, b) => {
        if (a.read === b.read) return 0;
        return a.read ? 1 : -1; 
      });
    },
    err => console.error(err)
  );
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

  markAsRead(notif: StudentNotification) {
  if (!notif.read && notif.ID !== undefined) {
    this.notificationService.markAsRead(this.studentId, notif.ID).subscribe(() => {
      notif.read = true;
    });
  }
}

}
