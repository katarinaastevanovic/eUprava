package handlers

import (
	"encoding/json"
	"net/http"
	"school-service/models"
	"school-service/services"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type SchoolHandler struct {
	Service *services.SchoolService
}

func NewSchoolHandler(s *services.SchoolService) *SchoolHandler {
	return &SchoolHandler{Service: s}
}

type AbsenceResponse struct {
	ID      uint   `json:"id"`
	Type    string `json:"type"`
	Date    string `json:"date"`
	Subject string `json:"subject"`
}

type CreateAbsenceRequest struct {
	Date      time.Time `json:"date"`
	StudentID uint      `json:"studentId"`
	SubjectID uint      `json:"subjectId"`
}

type AbsenceResponse2 struct {
	ID         uint      `json:"id"`
	Type       string    `json:"type"`
	Date       time.Time `json:"date"`
	StudentIDs []uint    `json:"studentIds"`
	SubjectID  uint      `json:"subjectId"`
}

type ClassDTO struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

type TeacherClassesResponse struct {
	SubjectName string     `json:"subject_name"`
	Classes     []ClassDTO `json:"classes"`
}

func (h *SchoolHandler) GetStudentAbsences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	absences, count, err := h.Service.GetAbsencesByUserID(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responseAbsences []AbsenceResponse
	for _, a := range absences {
		responseAbsences = append(responseAbsences, AbsenceResponse{
			ID:      a.ID,
			Type:    string(a.Type),
			Date:    a.Date.Format("2006-01-02 15:04"),
			Subject: a.Subject.Name,
		})
	}

	resp := map[string]interface{}{
		"user_id":  userID,
		"count":    count,
		"absences": responseAbsences,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

type UpdateAbsenceRequest struct {
	Type models.AbsenceType `json:"type"`
}

func (h *SchoolHandler) UpdateAbsenceType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	absenceID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid absence ID", http.StatusBadRequest)
		return
	}

	var req UpdateAbsenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Type != models.Excused && req.Type != models.Unexcused && req.Type != models.Pending {
		http.Error(w, "Invalid absence type", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateAbsenceType(uint(absenceID), req.Type); err != nil {
		http.Error(w, "Failed to update absence type", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Absence type updated successfully",
	})
}

func (h *SchoolHandler) CreateAbsences(w http.ResponseWriter, r *http.Request) {
	var req AbsenceResponse2
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.StudentIDs) == 0 {
		http.Error(w, "studentIds are required", http.StatusBadRequest)
		return
	}
	if req.SubjectID == 0 {
		http.Error(w, "subjectId is required", http.StatusBadRequest)
		return
	}

	date := req.Date
	if date.IsZero() {
		date = time.Now()
	}

	var absences []models.Absence
	for _, sid := range req.StudentIDs {
		absences = append(absences, models.Absence{
			Type:      models.Pending,
			Date:      date,
			StudentID: sid,
			SubjectID: req.SubjectID,
		})
	}

	if err := h.Service.CreateAbsences(absences); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *SchoolHandler) GetClassesForTeacher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	teacherID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	var teacher models.Teacher
	if err := h.Service.DB.First(&teacher, teacherID).Error; err != nil {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	var subject models.Subject
	if err := h.Service.DB.First(&subject, teacher.SubjectID).Error; err != nil {
		http.Error(w, "Subject not found", http.StatusNotFound)
		return
	}

	classes, err := h.Service.GetClassesByTeacherID(uint(teacherID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dtos := make([]ClassDTO, 0, len(classes))
	for _, c := range classes {
		dtos = append(dtos, ClassDTO{
			ID:    c.ID,
			Title: c.Title,
			Year:  c.Year,
		})
	}

	response := TeacherClassesResponse{
		SubjectName: subject.Name,
		Classes:     dtos,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SchoolHandler) GetStudentsByClass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classIDStr := vars["classID"]

	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	students, err := h.Service.GetStudentsByClassID(uint(classID))
	if err != nil {
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (h *SchoolHandler) GetAbsenceCountForSubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	studentID, err := strconv.ParseUint(vars["studentID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	subjectID, err := strconv.ParseUint(vars["subjectID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	count, err := h.Service.GetAbsenceCountForSubject(uint(studentID), uint(subjectID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"count": count})
}

func (h *SchoolHandler) GetStudentByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	studentDTO, err := h.Service.GetStudentDTOByUserID(uint(userID))
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentDTO)
}

func (h *SchoolHandler) GetStudentFullProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	dto, err := h.Service.GetFullStudentProfileByUserID(uint(userID))
	if err != nil {
		http.Error(w, "Failed to load full student profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}

type GradeItem struct {
	Value int       `json:"value"`
	Date  time.Time `json:"date,omitempty"`
}

type GradesResponse struct {
	StudentID   uint        `json:"student_id"`
	SubjectID   uint        `json:"subject_id"`
	TeacherID   uint        `json:"teacher_id"`
	SubjectName string      `json:"subject_name"`
	Grades      []GradeItem `json:"grades"`
}

func (h *SchoolHandler) GetGradesByStudentSubjectAndTeacherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	studentID, err := strconv.Atoi(vars["studentID"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	subjectID, err := strconv.Atoi(vars["subjectID"])
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	teacherID, err := strconv.Atoi(vars["teacherID"])
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	grades, err := h.Service.GetGradesByStudentSubjectAndTeacher(uint(studentID), uint(subjectID), uint(teacherID))
	if err != nil {
		http.Error(w, "Failed to fetch grades", http.StatusInternalServerError)
		return
	}

	var gradeItems []GradeItem
	for _, g := range grades {
		gradeItems = append(gradeItems, GradeItem{
			Value: g.Value,
			Date:  g.Date,
		})
	}

	subjectName := ""
	if len(grades) > 0 {
		subjectName = grades[0].Subject.Name
	}

	response := GradesResponse{
		StudentID:   uint(studentID),
		SubjectID:   uint(subjectID),
		TeacherID:   uint(teacherID),
		SubjectName: subjectName,
		Grades:      gradeItems,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type StudentGradeResponse struct {
	SubjectID   uint      `json:"subject_id"`
	SubjectName string    `json:"subject_name"`
	Value       int       `json:"value"`
	Date        time.Time `json:"date,omitempty"`
}

func (h *SchoolHandler) GetAllGradesByStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID, err := strconv.Atoi(vars["studentID"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	grades, err := h.Service.GetAllGradesByStudent(uint(studentID))
	if err != nil {
		http.Error(w, "Failed to fetch grades", http.StatusInternalServerError)
		return
	}

	var response []StudentGradeResponse
	for _, g := range grades {
		response = append(response, StudentGradeResponse{
			SubjectID:   g.SubjectID,
			SubjectName: g.Subject.Name,
			Value:       g.Value,
			Date:        g.Date,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SchoolHandler) GetAverageByTeacherAndSubjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID, _ := strconv.Atoi(vars["studentID"])
	subjectID, _ := strconv.Atoi(vars["subjectID"])
	teacherID, _ := strconv.Atoi(vars["teacherID"])

	avg, err := h.Service.GetAverageGradeByTeacherAndSubject(uint(studentID), uint(subjectID), uint(teacherID))
	if err != nil {
		http.Error(w, "Failed to fetch average", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"student_id": studentID,
		"subject_id": subjectID,
		"teacher_id": teacherID,
		"average":    avg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *SchoolHandler) GetAverageByStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID, _ := strconv.Atoi(vars["studentID"])

	avg, err := h.Service.GetAverageGradeByStudent(uint(studentID))
	if err != nil {
		http.Error(w, "Failed to fetch average", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"student_id": studentID,
		"average":    avg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *SchoolHandler) GetAverageByStudentPerSubjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID, _ := strconv.Atoi(vars["studentID"])

	averages, err := h.Service.GetAverageGradeByStudentPerSubject(uint(studentID))
	if err != nil {
		http.Error(w, "Failed to fetch averages", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"student_id": studentID,
		"subjects":   averages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *SchoolHandler) GetTeacherByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	teacher, err := h.Service.GetTeacherByUserID(uint(userID))
	if err != nil {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func (h *SchoolHandler) SearchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classIDStr := vars["classId"]
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	students, err := h.Service.SearchStudentsByName(uint(classID), query)
	if err != nil {
		http.Error(w, "Failed to search students: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(students)
}

func (h *SchoolHandler) SortStudentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classIDStr := vars["classId"]
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	order := r.URL.Query().Get("order")
	if order != "asc" && order != "desc" {
		order = "asc" // default
	}

	students, err := h.Service.SortStudentsByLastName(uint(classID), order)
	if err != nil {
		http.Error(w, "Failed to sort students: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(students)
}

func (h *SchoolHandler) CheckStudentMedicalCertificate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	// Uzmi Authorization header od profesora koji šalje zahtev
	token := r.Header.Get("Authorization")
	if token != "" {
		// Ako stiže u formatu "Bearer xyz", skini prefix da ne dupliraš
		token = strings.TrimPrefix(token, "Bearer ")
	}

	hasCert, err := h.Service.CheckStudentCertificate(uint(userId), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"userId":         userId,
		"hasCertificate": hasCert,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
