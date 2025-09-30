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
import { DoctorApprovedRequestsComponent } from './components/doctor-approved-requests/doctor-approved-requests.component';
import { ExaminationFormComponent } from './components/examination/examination.component';
import { ClassroomComponent } from './components/classroom/classroom.component';
import { StudentProfileComponent } from './components/student-profile/student-profile.component';
import { MedicalCertificateComponent } from './components/medical-certificate/medical-certificate.component';

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
    { path: 'approved-requests', component: DoctorApprovedRequestsComponent },
    { path: 'examination/:id', component: ExaminationFormComponent },
    { path: 'classroom/:id', component: ClassroomComponent },
    {path: 'student/:id/profile',component: StudentProfileComponent},
    { path: 'medical-certificate/:requestId', component: MedicalCertificateComponent },
    { path: '', redirectTo: '', pathMatch: 'full' },
    { path: '**', redirectTo: '' }
];
