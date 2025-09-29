import { Routes } from '@angular/router';
import { RegisterComponent } from './components/register/register.component';
import { AppComponent } from './app.component';
import { UserProfileComponent } from './components/user-profile/user-profile.component';
import { LoginComponent } from './components/login/login.component';
import { CompleteProfileComponent } from './components/complete-profile/complete-profile.component';
import { ClassroomComponent } from './components/classroom/classroom.component';
import { StudentProfileComponent } from './components/student-profile/student-profile.component';

export const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'register', component: RegisterComponent },
    { path: 'profile', component: UserProfileComponent },
    { path: 'login', component: LoginComponent },
    { path: 'complete-profile', component: CompleteProfileComponent },
    { path: 'classroom/:id', component: ClassroomComponent },
    {path: 'student/:id/profile',component: StudentProfileComponent},
    { path: '', redirectTo: '', pathMatch: 'full' },
    { path: '**', redirectTo: '' }
];
