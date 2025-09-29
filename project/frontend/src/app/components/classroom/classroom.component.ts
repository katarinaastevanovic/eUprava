import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
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
  imports: [CommonModule,FormsModule],
  templateUrl: './classroom.component.html',
  styleUrls: ['./classroom.component.css']
})
export class ClassroomComponent implements OnInit {
  classId!: number;
  students: StudentDTO[] = [];
  loading = true;
  teacherSubjectId!: number;

  constructor(private route: ActivatedRoute, private userService: UserService) {}

  ngOnInit(): void {
    this.classId = +this.route.snapshot.paramMap.get('id')!;

    this.userService.getUserProfile().subscribe({
      next: (user) => {
        this.userService.getTeacherClasses(user.id).subscribe({
          next: (res) => {
            if (res.classes.length > 0) {
              this.teacherSubjectId = res.classes[0].id; 
              this.loadStudents();
            } else {
              console.error('Teacher has no classes assigned');
              this.loading = false;
            }
          },
          error: (err) => {
            console.error('Greška pri učitavanju nastavnikovih klasa', err);
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

  // Ako zatvorimo mod, resetujemo selekciju
  if (!this.selectMode) {
    this.students.forEach(s => s.selected = false);
  }
}

confirmSelection() {
  const selectedStudents = this.students.filter(s => s.selected);
  console.log('Selected students:', selectedStudents);
  alert(`Selected ${selectedStudents.length} student(s)`);

  // Trenutno ništa ne šalje na backend
}


  
}
