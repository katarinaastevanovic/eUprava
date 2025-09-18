import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Component } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, ReactiveFormsModule, ValidationErrors, Validators } from '@angular/forms';
import { Observable, of, map, catchError } from 'rxjs';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [HttpClientModule,CommonModule,ReactiveFormsModule],
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent {
  registerForm: FormGroup;

  constructor(private fb: FormBuilder, private http: HttpClient) {
    this.registerForm = this.fb.group({
      umcn: ['', [Validators.required, Validators.pattern(/^[0-9]{13}$/)]],
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

  // Inline async validator
  usernameValidator(control: AbstractControl): Observable<ValidationErrors | null> {
    if (!control.value) {
      return of(null);
    }

    return this.http.get(`/api/check-username?username=${control.value}`, { responseType: 'text' }).pipe(
      map(() => null), // OK → username slobodan
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
      console.log(this.registerForm.value);
      // HttpClient POST → backend
    } else {
      this.registerForm.markAllAsTouched();
    }
  }
}
