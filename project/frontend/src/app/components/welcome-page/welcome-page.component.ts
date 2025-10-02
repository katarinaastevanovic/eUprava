import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common'; 

@Component({
  selector: 'app-welcome-page',
  standalone: true,           
  imports: [CommonModule],   
  templateUrl: './welcome-page.component.html',
  styleUrls: ['./welcome-page.component.css'] 
})
export class WelcomePageComponent implements OnInit {
  isLoggedIn = false;

  constructor(public router: Router) {}

  ngOnInit(): void {
    const token = localStorage.getItem('jwt');
    this.isLoggedIn = !!token;
  }
}
