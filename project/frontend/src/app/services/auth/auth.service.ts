import { Injectable, inject } from '@angular/core';
import { Auth, signInWithPopup, GoogleAuthProvider, signOut } from '@angular/fire/auth';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { BehaviorSubject } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private auth = inject(Auth);
  private http = inject(HttpClient);

  private apiGatewayUrl = 'http://localhost:8080/api/auth';

  private userRoleSubject = new BehaviorSubject<'STUDENT' | 'DOCTOR' | null>(this.getRoleFromToken());
  userRole$ = this.userRoleSubject.asObservable();

  private userIdSubject = new BehaviorSubject<number>(this.getUserIdFromToken());
  userId$ = this.userIdSubject.asObservable();

  async loginWithGoogle(): Promise<string> {
    const provider = new GoogleAuthProvider();
    provider.setCustomParameters({
      prompt: 'select_account',
      client_id: environment.firebase.webClientId
    });

    const credential = await signInWithPopup(this.auth, provider);
    if (!credential.user) throw new Error('No user returned');

    return await credential.user.getIdToken();
  }

  loginBackend(idToken: string) {
    return this.http.post(`${this.apiGatewayUrl}/firebase-login`, { idToken });
  }

  loginWithEmail(email: string, password: string) {
    return this.http.post(`${this.apiGatewayUrl}/login`, { email, password });
  }

  completeProfile(profile: any) {
    return this.http.post(`${this.apiGatewayUrl}/complete-profile`, profile);
  }

  async logout(): Promise<void> {
    try {
      await signOut(this.auth);
      localStorage.removeItem('jwt');
      this.updateAuthState(); 
    } catch (err) {
      console.error('Firebase logout error', err);
    }
  }

  updateAuthState() {
    this.userRoleSubject.next(this.getRoleFromToken());
    this.userIdSubject.next(this.getUserIdFromToken());
  }

  getRole(): 'STUDENT' | 'DOCTOR' | null {
    return this.userRoleSubject.value;
  }

  getUserId(): number {
    return this.userIdSubject.value;
  }

  private getRoleFromToken(): 'STUDENT' | 'DOCTOR' | null {
    const token = localStorage.getItem('jwt');
    if (!token) return null;
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      if (Array.isArray(payload.roles)) {
        if (payload.roles.includes('STUDENT')) return 'STUDENT';
        if (payload.roles.includes('DOCTOR')) return 'DOCTOR';
      }
      return payload.role?.toUpperCase() === 'STUDENT' ? 'STUDENT' :
             payload.role?.toUpperCase() === 'DOCTOR' ? 'DOCTOR' : null;
    } catch {
      return null;
    }
  }

  private getUserIdFromToken(): number {
    const token = localStorage.getItem('jwt');
    if (!token) return 0;
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.sub || 0;
    } catch {
      return 0;
    }
  }
}
