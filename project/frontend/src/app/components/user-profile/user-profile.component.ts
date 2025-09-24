import { Component, OnInit, inject } from '@angular/core';
import { UserService, User } from '../../services/user/user.service';

@Component({
  selector: 'app-user-profile',
  imports: [],
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.css']
})
export class UserProfileComponent implements OnInit {
  private userService = inject(UserService);

  user: User | null = null;
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
      },
      error: (err) => {
        this.error = 'Failed to load user profile';
        this.loading = false;
      }
    });
  }
}
