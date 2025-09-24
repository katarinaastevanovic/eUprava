import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, NgForm } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../services/auth/auth.service';

@Component({
  selector: 'app-complete-profile',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './complete-profile.component.html'
})
export class CompleteProfileComponent {
  uid = '';
  email = '';
  username = '';
  name = '';
  lastName = '';
  umcn = '';
  role = ''; 
  roles = ['STUDENT', 'TEACHER', 'DOCTOR', 'PARENT']; 
  error = '';

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private authService: AuthService
  ) {
    this.uid = this.route.snapshot.queryParamMap.get('uid') || '';
    this.email = this.route.snapshot.queryParamMap.get('email') || '';
  }

  complete(form: NgForm) {
    if (!form.valid) {
      this.error = 'Please fill out all required fields';
      return;
    }

    const profile = {
      uid: this.uid,
      email: this.email,
      username: this.username,
      name: this.name,
      lastName: this.lastName,
      umcn: this.umcn,
      role: this.role 
    };

    this.authService.completeProfile(profile).subscribe({
      next: (res: any) => {
        localStorage.setItem('jwt', res.token); 
        this.router.navigate(['/']); 
      },
      error: () => {
        this.error = 'Greška pri kompletiranju profila';
      }
    });
  }
}
