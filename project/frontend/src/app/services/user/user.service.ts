import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Observable } from 'rxjs';

export interface User {
  id: number;
  name: string;
  last_name: string;
  email: string;
  username: string;
  role: string;
  birth_date: string;
  gender: string;
}
export interface Absence {
  id: number;
  type: string;
  date: string;
  subject: string;
}
export interface AbsenceDTO {
  id: number;
  type: string;
  date: string;
}

export interface ClassDTO {
  id: number;
  title: string;
  year: number;
}
export interface TeacherClassesResponse {
  subject_name: string;
  classes: ClassDTO[];
}
interface StudentDTO {
  id: number;
  userId: number;
  name: string;
  lastName: string;
  numberOfAbsences: number;
  absences: AbsenceDTO[];
}


@Injectable({ providedIn: 'root' })
export class UserService {
  private http = inject(HttpClient);

  getUserProfile(): Observable<User> {
    const token = localStorage.getItem('jwt');
    let headers = new HttpHeaders();

    if (token) {
      headers = headers.set('Authorization', `Bearer ${token}`);
    }

    return this.http.get<User>(`${environment.apiBaseUrl}/profile`, { headers });
  }

  getStudentAbsences(studentId: number): Observable<{ student_id: number; count: number; absences: Absence[] }> {
    return this.http.get<{ student_id: number; count: number; absences: Absence[] }>(
      `${environment.schoolApiBaseUrl}/students/${studentId}/absences`
    );
  }

  updateAbsenceType(absenceId: number, newType: string): Observable<any> {
  return this.http.put(`${environment.schoolApiBaseUrl}/absences/${absenceId}/type`, {
    type: newType
  });
}


getTeacherClasses(teacherId: number): Observable<TeacherClassesResponse> {
  return this.http.get<TeacherClassesResponse>(
    `${environment.schoolApiBaseUrl}/teachers/${teacherId}/classes`
  );
}

getStudentsByClass(classId: number) {
  return this.http.get<StudentDTO[]>(
    `${environment.schoolApiBaseUrl}/api/classes/${classId}/students`
  );
}

getAbsenceCountForSubject(studentId: number, subjectId: number): Observable<{ count: number }> {
  return this.http.get<{ count: number }>(
    `${environment.schoolApiBaseUrl}/students/${studentId}/subjects/${subjectId}/absences/count`
  );
}

getStudentByUserId(userId: number): Observable<{ id: number; userId: number }> {
  return this.http.get<{ id: number; userId: number }>(
    `${environment.schoolApiBaseUrl}/students/by-user/${userId}`
  );
}


}
