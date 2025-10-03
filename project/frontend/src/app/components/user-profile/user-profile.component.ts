import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UserService, User, Absence, ClassDTO, GradeDTO, SubjectAverage } from '../../services/user/user.service';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { RouterLink } from '@angular/router';


@Component({
  selector: 'app-user-profile',
  imports: [HttpClientModule, CommonModule, FormsModule, RouterLink],
  standalone: true,
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.css']
})
export class UserProfileComponent implements OnInit {
  private userService = inject(UserService);

  user: User | null = null;
  absences: Absence[] = [];
  absenceStats: { excused: number; unexcused: number; pending: number; total: number } | null = null;
  showAbsenceStats = false;
  grades: GradeDTO[] | null = null;
  subjectAverages: SubjectAverage[] = [];
  showAverages = false;
  studentId: number | null = null;
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
          // ⬇️ absences i dalje koriste userId
          this.loadAbsences(data.id);

          // ⬇️ samo za ocene radimo mapiranje userId → studentId
          this.userService.getStudentByUserId(data.id).subscribe({
            next: (studentDto) => {
              this.studentId = studentDto.id;
              console.log("Mapiran userId", data.id, "na studentId", this.studentId);
              this.loadGrades(this.studentId);
            },
            error: (err) => {
              console.error("Greška pri mapiranju userId na studentId", err);
              this.error = 'Failed to map user to student';
            }
          });
        }

        if (data?.role?.toUpperCase() === 'TEACHER' && data?.id) {
          this.userService.getTeacherByUserId2(data.id).subscribe({
            next: (teacherDto) => {
              // Mapiramo backend polja na frontend-friendly
              const teacherId = teacherDto.ID ?? teacherDto.ID;
              const subjectId = teacherDto.SubjectID ?? teacherDto.SubjectID;
              console.log("Mapiran userId", data.id, "na teacherId", teacherId, "subjectId", subjectId);

              this.loadTeacherClasses(teacherId);
              this.teacherSubject = teacherDto.Title ?? ''; // ako želiš i title
            },
            error: (err) => {
              console.error("Greška pri mapiranju userId na teacherId", err);
              this.error = 'Failed to map user to teacher';
            }
          });
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

        // ⬇️ dodatno učitavanje statistike
        this.userService.getAbsenceStats(studentId).subscribe({
          next: (stats) => {
            this.absenceStats = {
              total: stats.total,
              excused: stats.excused,
              unexcused: stats.unexcused,
              pending: stats.pending
            };
          },
          error: () => {
            this.error = 'Failed to load absence stats';
          }
        });
      },
      error: () => {
        this.error = 'Failed to load absences';
      }
    });

  }

  loadGrades(studentId: number) {
    this.userService.getStudentGrades(studentId).subscribe({
      next: (grades: GradeDTO[]) => {  // 👈 dodaj tip
        this.grades = grades;
        console.log("Ocene:", grades);
      },
      error: (err: any) => {           // 👈 dodaj tip
        console.error('Greška pri dohvatanju ocena', err);
        this.error = 'Failed to load grades';
      }
    });
  }

  loadSubjectAverages(studentId: number) {
    this.showAverages = false;
    this.subjectAverages = [];
    this.userService.getStudentAveragesPerSubject(studentId).subscribe({
      next: (res) => {
        this.subjectAverages = res.subjects;
        this.showAverages = true;
        console.log("Proseci po predmetima:", res.subjects);
      },
      error: (err) => {
        console.error('Greška pri dohvatanju proseka', err);
        this.error = 'Failed to load averages';
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

  toggleAbsenceStats() {
    this.showAbsenceStats = !this.showAbsenceStats;
  }


}


