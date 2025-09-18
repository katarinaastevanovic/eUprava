import { Routes } from '@angular/router';
import { RegisterComponent } from './components/register/register.component';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
{ path: 'register', component: RegisterComponent },  
  // { path: 'login', component: LoginComponent },
];
