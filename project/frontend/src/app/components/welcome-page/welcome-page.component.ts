import { Component, OnInit, inject } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth/auth.service';

@Component({
  selector: 'app-welcome-page',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './welcome-page.component.html',
  styleUrls: ['./welcome-page.component.css']
})
export class WelcomePageComponent implements OnInit {
  isLoggedIn = false;
  authService = inject(AuthService);
  router = inject(Router);

  ngOnInit(): void {
    this.authService.userRole$.subscribe(role => {
      this.isLoggedIn = role !== null;
    });
  }
}
