import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
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
  // dodaj ostale polja po potrebi
}

@Injectable({ providedIn: 'root' })
export class UserService {
  private http = inject(HttpClient);

  getUserProfile(): Observable<User> {
    return this.http.get<User>(`${environment.apiBaseUrl}/profile`);
    // backend mora da koristi token iz localStorage ili cookie za identifikaciju
  }
}
