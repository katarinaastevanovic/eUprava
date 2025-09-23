import { Routes } from '@angular/router';
import { RegisterComponent } from './components/register/register.component';
import { AppComponent } from './app.component';
import { UserProfileComponent } from './components/user-profile/user-profile.component';

export const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'register', component: RegisterComponent },  
    { path: 'profile', component: UserProfileComponent },  
  // { path: 'login', component: LoginComponent },
];
