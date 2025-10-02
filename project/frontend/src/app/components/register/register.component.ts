import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Component } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, ReactiveFormsModule, ValidationErrors, Validators } from '@angular/forms';
import { Observable, of, map, catchError } from 'rxjs';
import { Router } from '@angular/router';

function umcnValidator(control: AbstractControl): ValidationErrors | null {
  const value = control.value;

  if (!value || !/^\d{13}$/.test(value)) {
    return { invalidFormat: true };
  }

  const day = parseInt(value.substring(0, 2), 10);
  const month = parseInt(value.substring(2, 4), 10);
  let year = parseInt(value.substring(4, 7), 10);

  if (year >= 900) {
    year = 1000 + year; 
  } else {
    year = 2000 + year; 
  }

  const testDate = new Date(year, month - 1, day);
  if (
    testDate.getDate() !== day ||
    testDate.getMonth() + 1 !== month ||
    testDate.getFullYear() !== year
  ) {
    return { invalidDate: true };
  }

  return null;
}


@Component({
  selector: 'app-register',
  standalone: true,
  imports: [HttpClientModule,CommonModule,ReactiveFormsModule],
  templateUrl: './register.component.html',
styleUrl: './register.component.css'})
export class RegisterComponent {
  registerForm: FormGroup;

  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router) {
    this.registerForm = this.fb.group({
      umcn: ['', [Validators.required, umcnValidator]],
      name: ['', [Validators.required, Validators.pattern(/^[A-Z][a-z]+$/)]],
      lastName: ['', [Validators.required, Validators.pattern(/^[A-Z][a-z]+$/)]],
      email: ['', [Validators.required, Validators.email]],
      username: [
        '',
        [Validators.required, Validators.minLength(4), Validators.maxLength(20)],
        [this.usernameValidator.bind(this)] // 👈 async validator ovde
      ],
      password: [
        '',
        [
          Validators.required,
          Validators.minLength(8),
          Validators.pattern(/^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9]).+$/)
        ]
      ],
      role: ['', Validators.required],
    });
  }

  usernameValidator(control: AbstractControl): Observable<ValidationErrors | null> {
    if (!control.value) {
      return of(null);
    }

    return this.http.get(`/api/check-username?username=${control.value}`, { responseType: 'text' }).pipe(
      map(() => null), 
      catchError(err => {
        if (err.status === 409) {
          return of({ usernameTaken: true });
        }
        return of(null);
      })
    );
  }

  onSubmit() {
  if (this.registerForm.valid) {
    this.http.post('http://localhost:8080/api/auth/register', this.registerForm.value, { responseType: 'text' })
      .subscribe({
        next: (res) => {
          console.log('✅ Registracija uspela:', res);
          alert('Korisnik uspešno registrovan!');
          this.registerForm.reset();

           this.router.navigate(['/']);
        },
        error: (err) => {
          console.error('❌ Greška:', err);
          if (err.status === 409) {
            alert('Podaci moraju biti jedinstveni (UMCN, email ili username već postoje).');
          } else if (err.status === 400) {
            alert('Neispravan unos: ' + err.error);
          } else {
            alert('Došlo je do greške na serveru.');
          }
        }
      });
  } else {
    this.registerForm.markAllAsTouched();
  }
}



}
