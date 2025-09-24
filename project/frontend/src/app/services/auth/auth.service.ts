import { Injectable, inject } from '@angular/core';
import { Auth, signInWithPopup, GoogleAuthProvider } from '@angular/fire/auth';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private auth = inject(Auth);
  private http = inject(HttpClient);

  async loginWithGoogle(): Promise<string> {
    const provider = new GoogleAuthProvider();
    provider.setCustomParameters({
      client_id: environment.firebase.webClientId
    });

    console.log('Pokrecem Google popup login...');
    const credential = await signInWithPopup(this.auth, provider);

    if (!credential.user) throw new Error('No user returned');

    const idToken = await credential.user.getIdToken();
    console.log('Firebase ID Token dobijen:', idToken);
    return idToken;
  }

  loginBackend(idToken: string) {
    console.log('Saljem ID token na backend:', idToken);
    return this.http.post(`${environment.apiBaseUrl}/firebase-login`, { idToken });
  }

  completeProfile(profile: any) {
    console.log('Kompletiranje profila:', profile);
    return this.http.post(`${environment.apiBaseUrl}/complete-profile`, profile);
  }

  loginWithEmail(email: string, password: string) {
  return this.http.post(`${environment.apiBaseUrl}/login`, { email, password });
}

}