import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { UserService, User, Absence, GradesResponse } from '../../services/user/user.service';

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
  gradesResponse: GradesResponse | null = null;
  singleAverage: number | null = null;
  loading = false;
  error = '';

  ngOnInit() {
    const routeStudentId = Number(this.route.snapshot.paramMap.get('id'));
    console.log("Route student id:", routeStudentId);

    if (routeStudentId) {
      this.loadStudent(routeStudentId);
      this.loadAbsences(routeStudentId);

      // Dohvati ulogovanog korisnika
      this.userService.getUserProfile().subscribe({
        next: (user) => {
          console.log("Ulogovani korisnik:", user);

          // Ako je korisnik nastavnik, dohvatimo njegov teacherId
          // Ako je korisnik nastavnik, dohvatimo njegov teacherId
          this.userService.getTeacherByUserId(user.id).subscribe({
            next: (teacher: any) => {
              console.log("Teacher pronađen:", teacher);

              const teacherId = teacher.id ?? teacher.ID;
              const subjectId = teacher.subject_id ?? teacher.SubjectID;

              // ✅ Umesto routeStudentId koristimo pravi student.id iz profila
              this.userService.getUserById(routeStudentId).subscribe({
                next: (studentProfile) => {
                  console.log("Student profil:", studentProfile);
                  const realStudentId = studentProfile.id; // backend-ov ID studenta

                  console.log("Pozivam loadGrades sa:", realStudentId, subjectId, teacherId);
                  this.loadGrades(realStudentId, subjectId, teacherId);
                },
                error: (err) => {
                  console.error("Greška pri dohvatanju studenta", err);
                  this.error = "Failed to load student profile";
                }
              });
            },
            error: (err) => {
              console.error("Greška pri dohvatanju nastavnika", err);
              this.error = "Failed to load teacher info";
            }
          });


        },
        error: (err) => {
          console.error("Greška pri dohvatanju ulogovanog korisnika", err);
          this.error = "Failed to get logged-in user info";
        }
      });
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

  loadGrades(studentId: number, subjectId: number, teacherId: number) {
    console.log("Pozivam loadGrades sa:", studentId, subjectId, teacherId);

    this.userService.getGradesByStudentSubjectAndTeacher(studentId, subjectId, teacherId)
      .subscribe({
        next: (res) => {
          console.log("RAW Dohvaćene ocene:", res);
          const safeGrades = res.grades ?? [];
          console.log("Safe grades:", safeGrades);

          this.gradesResponse = { ...res, grades: safeGrades };
          console.log("Postavljen gradesResponse:", this.gradesResponse);
        },
        error: (err) => {
          console.error('Greška pri dohvatanju ocena', err);
          this.error = 'Failed to load grades';
        }
      });
  }

  loadSingleAverage(studentId: number, subjectId: number, teacherId: number) {
  this.userService.getStudentSubjectTeacherAverage(studentId, subjectId, teacherId).subscribe({
    next: (res) => {
      this.singleAverage = res.average;
      console.log('Single average:', res);
    },
    error: (err) => {
      console.error('Greška pri dohvatanju proseka za jednog profesora', err);
      this.error = 'Failed to load average';
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

