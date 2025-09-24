import { Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { CompleteProfileComponent } from './components/complete-profile/complete-profile.component';

export const routes: Routes = [
  { path: 'login', component: LoginComponent },
  { path: 'complete-profile', component: CompleteProfileComponent },
  { path: '', redirectTo: '', pathMatch: 'full' },
  { path: '**', redirectTo: '' }
];
