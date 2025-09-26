export interface FullMedicalRecord {
  patientId: number;
  name: string;
  lastName: string;
  jmbg: string;
  birthDate: string;
  gender: string;
  allergies: string;
  chronicDiseases: string;
  lastUpdate: string;
  examinations?: any[];
  requests?: any[];
}
