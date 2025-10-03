import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MedicalRecordService } from '../../services/medical-record/medical-record.service';
import { FullMedicalRecord } from '../../models/medical-models/full-medical-record.model';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-medical-record',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './medical-record.component.html'
})
export class MedicalRecordComponent implements OnInit {
  record: FullMedicalRecord | null = null;
  userId!: number;

  allergies = '';
  chronicDiseases = '';
  error = '';
  success = '';
  loading = false;

  editMode = false; 

  constructor(private medicalRecordService: MedicalRecordService) {}

  ngOnInit(): void {
    this.setUserIdFromToken();
    if (this.userId) {
      this.getRecord();
    } else {
      this.error = 'User not authenticated';
    }
  }

private isFullMedicalRecord(obj: any): obj is FullMedicalRecord {
  return obj &&
         typeof obj.patientId === 'number' &&
         typeof obj.name === 'string' &&
         typeof obj.lastName === 'string' &&
         typeof obj.allergies === 'string' &&
         typeof obj.chronicDiseases === 'string';
}


  private setUserIdFromToken() {
    const token = localStorage.getItem('jwt');
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        this.userId = payload.sub;
      } catch (err) {
        console.error('Failed to decode token', err);
        this.userId = 0;
      }
    }
  }

  getRecord() {
  this.loading = true;
  this.error = '';
  this.success = '';

  this.medicalRecordService.getFullRecord(this.userId).subscribe({
    next: (data) => {
      console.log('Raw data from backend:', data);

      if (this.isFullMedicalRecord(data)) {
        const mappedExaminations = data.examinations?.map(exam => ({
          CreatedAt: exam.CreatedAt,
          Diagnosis: exam.diagnosis,
          Therapy: exam.therapy,
          Note: exam.note,
          ID: exam.ID,
          medicalRecordId: exam.medicalRecordId,
          requestId: exam.requestId,
          DeletedAt: exam.DeletedAt,
          UpdatedAt: exam.UpdatedAt
        })) || [];

        this.record = {
          ...data,
          examinations: mappedExaminations
        };

        this.allergies = data.allergies;
        this.chronicDiseases = data.chronicDiseases;
      } else {
        this.error = 'Invalid medical record type';
        console.log('Backend response does not match FullMedicalRecord');
      }

      this.loading = false;
    },
    error: (err) => {
      if (err.status === 404) {
        this.record = null;
      } else {
        this.error = 'Failed to fetch medical record';
      }
      this.loading = false;
    }
  });
}


  createRecord() {
    this.loading = true;
    this.error = '';
    this.success = '';

    this.medicalRecordService.createRecord(this.userId).subscribe({
      next: (data) => {
        this.record = data;
        this.allergies = data.allergies;
        this.chronicDiseases = data.chronicDiseases;
        this.loading = false;
        this.success = 'Medical record created successfully!';
      },
      error: () => {
        this.error = 'Failed to create record';
        this.loading = false;
      }
    });
  }

  toggleEdit() {
    this.editMode = !this.editMode;
    this.success = '';
    this.error = '';
  }

  updateRecord() {
    if (!this.record) {
      this.error = 'No medical record to update';
      return;
    }

    this.loading = true;
    this.error = '';
    this.success = '';

    const updatedRecord: FullMedicalRecord = {
      ...this.record,
      allergies: this.allergies,
      chronicDiseases: this.chronicDiseases
    };

    this.medicalRecordService.updateRecord(updatedRecord).subscribe({
      next: (data) => {
        this.getRecord(); 
        this.record = data;
        this.allergies = data.allergies;
        this.chronicDiseases = data.chronicDiseases;
        this.loading = false;
        this.editMode = false;
        this.success = 'Medical record updated successfully!';
      },
      error: () => {
        this.error = 'Failed to update medical record';
        this.loading = false;
      }
    });
  }

  getExaminations() {
    return this.record?.examinations || [];
  }

  getRequests() {
    return this.record?.requests || [];
  }
}
