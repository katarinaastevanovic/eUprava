import { Routes } from '@angular/router';
import { RegisterComponent } from './components/register/register.component';
import { AppComponent } from './app.component';
import { UserProfileComponent } from './components/user-profile/user-profile.component';
import { LoginComponent } from './components/login/login.component';
import { CompleteProfileComponent } from './components/complete-profile/complete-profile.component';
import { MedicalRecordComponent } from './components/medical-record/medical-record.component';
import { ExaminationRequestComponent } from './components/examination-request/examination-request.component'; 
import { PatientRequestsComponent } from './components/patient-request/patient-request.component'; 
import { DoctorRequestsComponent } from './components/doctor-request/doctor-request.component'; 

export const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'register', component: RegisterComponent },  
    { path: 'profile', component: UserProfileComponent },  
    { path: 'login', component: LoginComponent },
    { path: 'complete-profile', component: CompleteProfileComponent },
    { path: 'medical-record', component: MedicalRecordComponent },
    { path: 'create-request', component: ExaminationRequestComponent },
    { path: 'requests', component: PatientRequestsComponent },
    { path: 'doctors-requests', component: DoctorRequestsComponent },
    { path: '', redirectTo: '', pathMatch: 'full' },
    { path: '**', redirectTo: '' }
];
