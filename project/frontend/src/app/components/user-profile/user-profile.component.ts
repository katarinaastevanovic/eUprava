import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UserService, User, Absence  } from '../../services/user/user.service';


@Component({
  selector: 'app-user-profile',
  imports: [CommonModule],
  standalone: true,
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.css']
})
export class UserProfileComponent implements OnInit {
  private userService = inject(UserService);

  user: User | null = null;
  absences: Absence[] = [];
  loading = false;
  error = '';

  ngOnInit() {
    this.loadUser();
  }

  loadUser() {
  this.loading = true;
  this.userService.getUserProfile().subscribe({
    next: (data) => {
      this.user = data;
      this.loading = false;

      if (data?.id) {
        this.loadAbsences(data.id);
      }
    },
    error: (err) => {
      this.error = 'Failed to load user profile';
      this.loading = false;
    }
  });
}


  loadAbsences(studentId: number) {
    this.userService.getStudentAbsences(studentId).subscribe({
      next: (data) => {
        this.absences = data.absences;
      },
      error: () => {
        this.error = 'Failed to load absences';
      }
    });
  }
}
