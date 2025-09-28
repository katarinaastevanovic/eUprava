import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  standalone: true, 
  imports: [CommonModule] 
})

export class HeaderComponent {
  constructor(private authService: AuthService, private router: Router) {}

  logout() {
  this.authService.logout().then(() => {
    this.router.navigate(['/login']); 
  });
}

  isLoggedIn(): boolean {
    return !!localStorage.getItem('jwt');
  }
}
