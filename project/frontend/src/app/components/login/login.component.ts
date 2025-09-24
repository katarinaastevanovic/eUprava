import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './login.component.html',
})
export class LoginComponent {
  private authService = inject(AuthService);
  private router = inject(Router);

  email = '';
  password = '';
  loading = false;
  error = '';

  async loginWithFirebase() {
    this.loading = true;
    this.error = '';
    try {
      const idToken = await this.authService.loginWithGoogle();

      this.authService.loginBackend(idToken).subscribe({
        next: (res: any) => {
          localStorage.setItem('jwt', res.token);

          this.loading = false;
          this.router.navigateByUrl('/'); 
        },
        error: (err: any) => {
          if (err.status === 428 && err.error?.message === 'profile incomplete') {
            window.location.href = `/complete-profile?uid=${err.error.uid}&email=${err.error.email}`;
          } else {
            this.error = 'Firebase login failed';
            this.loading = false;
          }
        }
      });

    } catch (err) {
      console.error('Firebase login error:', err);
      this.error = 'Firebase login failed';
      this.loading = false;
    }
  }

  loginWithEmail() {
    this.loading = true;
    this.error = '';
    this.authService.loginWithEmail(this.email, this.password).subscribe({
      next: (res: any) => {
        localStorage.setItem('jwt', res.token);

        this.loading = false;
        this.router.navigateByUrl('/'); 
      },
      error: (err) => {
        console.error('Email login error:', err);
        this.error = 'Email login failed';
        this.loading = false;
      },
    });
  }
}
