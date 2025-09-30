import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { UserService, User, Absence } from '../../services/user/user.service';

@Component({
  selector: 'app-student-profile',
  standalone: true,
  imports: [CommonModule, FormsModule, HttpClientModule],
  templateUrl: './student-profile.component.html',
  styleUrls: ['./student-profile.component.css']
})
export class StudentProfileComponent implements OnInit {
  private userService = inject(UserService);
  private route = inject(ActivatedRoute);

  student: User | null = null;
  absences: Absence[] = [];
  loading = false;
  error = '';

  ngOnInit() {
    const studentId = Number(this.route.snapshot.paramMap.get('id'));
    if (studentId) {
      this.loadStudent(studentId);
      this.loadAbsences(studentId);
    }
  }

  loadStudent(id: number) {
  this.loading = true;
  this.userService.getUserById(id).subscribe({
    next: (data: User) => {  
      this.student = data;
      this.loading = false;
    },
    error: () => {
      this.error = 'Failed to load student';
      this.loading = false;
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
      next: () => console.log(`Absence ${absence.id} updated to ${absence.type}`),
      error: (err) => {
        console.error('Failed to update absence type', err);
        this.error = 'Failed to update absence type';
      }
    });
  }
}
