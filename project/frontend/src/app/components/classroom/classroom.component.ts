import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { UserService } from '../../services/user/user.service';
import { FormsModule } from '@angular/forms';

interface StudentDTO {
  id: number;
  userId?: number;
  name: string;
  lastName: string;
  numberOfAbsences: number;
  absences?: any[];
  selected?: boolean; 
}


@Component({
  selector: 'app-classroom',
  standalone: true,
  imports: [CommonModule,FormsModule,RouterLink],
  templateUrl: './classroom.component.html',
  styleUrls: ['./classroom.component.css']
})
export class ClassroomComponent implements OnInit {
  classId!: number;
  students: StudentDTO[] = [];
  loading = true;
  teacherSubjectId!: number;
  searchQuery = '';
  sortOrder: string = 'asc';

  constructor(private route: ActivatedRoute, private userService: UserService) {}

  ngOnInit(): void {
  this.classId = +this.route.snapshot.paramMap.get('id')!;

  this.userService.getUserProfile().subscribe({
    next: (user) => {
      // Dohvati teacherDto preko postojeće metode
      this.userService.getTeacherByUserId2(user.id).subscribe({
        next: (teacherDto) => {
          // Mapiramo polja iz backend-a na frontend-friendly
          const teacherId = teacherDto.ID;
          const subjectId = teacherDto.SubjectID;

          console.log("Mapiran teacher:", teacherId, "subjectId:", subjectId);
          this.teacherSubjectId = subjectId;

          this.loadStudents();
        },
        error: (err) => {
          console.error('Greška pri mapiranju userId na teacherId', err);
          this.loading = false;
        }
      });
    },
    error: (err) => {
      console.error('Greška pri učitavanju profila', err);
      this.loading = false;
    }
  });
}

  private loadStudents() {
    this.userService.getStudentsByClass(this.classId).subscribe({
      next: (students: StudentDTO[]) => {
        console.log('API response:', students);
        this.students = students;

        this.students.forEach(student => {
          this.userService.getAbsenceCountForSubject(student.id, this.teacherSubjectId)
            .subscribe({
              next: (res) => student.numberOfAbsences = res.count,
              error: (err) => console.error('Failed to load absence count', err)
            });
        });

        this.loading = false;
      },
      error: (err: any) => {
        console.error('Greška prilikom učitavanja učenika', err);
        this.loading = false;
      }
    });
  }

  selectMode = false;

toggleSelectMode() {
  this.selectMode = !this.selectMode;

  if (!this.selectMode) {
    this.students.forEach(s => s.selected = false);
  }
}

confirmSelection() {
  const selectedStudents = this.students.filter(s => s.selected).map(s => s.id);

  if (selectedStudents.length === 0) {
    alert("You didn't select any student!");
    return;
  }

  this.userService.createAbsences(selectedStudents, this.teacherSubjectId).subscribe({
    next: (res) => {
      console.log("Absences created:", res);
      alert(`Sucessfully created ${selectedStudents.length} absences!`);

      this.selectMode = false;
      this.loadStudents();
    },
    error: (err) => {
      console.error("Greška pri kreiranju izostanaka:", err);
      alert("Došlo je do greške pri slanju izostanaka.");
    }
  });
}

onSearch() {
  if (!this.searchQuery) {
    this.loadStudents(); 
    return;
  }

  this.userService.searchStudents(this.classId, this.searchQuery).subscribe({
    next: (students: StudentDTO[]) => {
      console.log('Search result:', students);
      this.students = students;
    },
    error: (err) => {
      console.error('Greška pri pretrazi učenika', err);
    }
  });
}

onSort() {
  this.userService.sortStudents(this.classId, this.sortOrder).subscribe({
    next: (students) => {
      this.students = students;
    },
    error: (err) => {
      console.error('Greška pri sortiranju učenika', err);
    }
  });
}
  
}