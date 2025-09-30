import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { AuthService } from '../../services/auth/auth.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  standalone: true, 
  imports: [CommonModule, RouterModule]
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
