import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UserService, User, Absence, ClassDTO  } from '../../services/user/user.service';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';


@Component({
  selector: 'app-user-profile',
  imports: [HttpClientModule, CommonModule, FormsModule],
  standalone: true,
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.css']
})
export class UserProfileComponent implements OnInit {
  private userService = inject(UserService);

  user: User | null = null;
  absences: Absence[] = [];
  teacherClasses: ClassDTO[] = [];
  teacherSubject: string = '';
  loading = false;
  error = '';

  ngOnInit() {
    this.loadUser();
  }

  loadUser() {
  this.loading = true;
  this.userService.getUserProfile().subscribe({
    next: (data) => {
      this.user = data;
      this.loading = false;

      if (data?.role?.toUpperCase() === 'STUDENT' && data?.id) {
      this.loadAbsences(data.id); 
    }
     if (data?.role?.toUpperCase() === 'TEACHER' && data?.id) {
        this.loadTeacherClasses(data.id);  
      }
    },
    error: (err) => {
      this.error = 'Failed to load user profile';
      this.loading = false;
    }
  });
}

loadTeacherClasses(teacherId: number) {
  console.log('TeacherId za razrede:', teacherId);
  this.userService.getTeacherClasses(teacherId).subscribe({
    next: (response) => {
      console.log('Dobijen odgovor:', response);
      this.teacherSubject = response.subject_name;
      this.teacherClasses = response.classes;
    },
    error: (err) => {
      console.error('Greška pri dohvatanju razreda', err);
      this.error = 'Failed to load classes';
    }
  });
}




  loadAbsences(studentId: number) {
    this.userService.getStudentAbsences(studentId).subscribe({
      next: (data) => {
        this.absences = data.absences;
      },
      error: () => {
        this.error = 'Failed to load absences';
      }
    });
  }

  onAbsenceTypeChange(absence: Absence) {
  this.userService.updateAbsenceType(absence.id, absence.type).subscribe({
    next: () => {
      console.log(`Absence ${absence.id} updated to ${absence.type}`);
    },
    error: (err) => {
      console.error('Failed to update absence type', err);
      this.error = 'Failed to update absence type';
    }
  });
}

}


