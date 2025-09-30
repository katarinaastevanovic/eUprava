import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-welcome-page',
  standalone: true,           
  imports: [],                 
  templateUrl: './welcome-page.component.html',
  styleUrls: ['./welcome-page.component.css'] 
})
export class WelcomePageComponent {
  constructor(public router: Router) {}
}
