package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"tracker/app/database"
	"tracker/app/middleware"
	"tracker/app/services"
	"tracker/app/transport"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	db, db_err := database.ConnectDB()
	if db_err != nil {
		log.Fatalf("Database connection failed: %v", db_err)
	}
	defer db.Close()

	r := mux.NewRouter()

	authService := services.NewAuthService(db)
	authHandler := transport.NewAuthHandler(authService)

	apiMiddleware := mux.NewRouter().PathPrefix("/api").Subrouter()
	apiMiddleware.Use(middleware.Middleware(authService))

	// Testing allow different cors http://localhost:3000
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://presentifyclone.click"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true, // Log CORS issues
	})
	//api.ConnectSheetsAPI()

	r.HandleFunc("/hello", hello).Methods("GET")
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/create_organization", authHandler.CreateOrganization).Methods("POST")
	r.HandleFunc("/health_check", authHandler.HealthCheck).Methods("GET") // FOR APPLICATIO LOAD BALANCER

	apiMiddleware.HandleFunc("/create_student", authHandler.CreateStudent).Methods("POST")
	apiMiddleware.HandleFunc("/create_location", authHandler.CreateLocation).Methods("POST")
	apiMiddleware.HandleFunc("/create_admin_staff", authHandler.CreateAdminStaff).Methods("POST")
	apiMiddleware.HandleFunc("/create_district", authHandler.CreateDistrict).Methods("POST")
	apiMiddleware.HandleFunc("/create_program", authHandler.CreateProgram).Methods("POST")
	apiMiddleware.HandleFunc("/create_material", authHandler.CreateMaterial).Methods("POST")
	apiMiddleware.HandleFunc("/create_tutor", authHandler.CreateTutor).Methods("POST")
	apiMiddleware.HandleFunc("/create_semester", authHandler.CreateSemester).Methods("POST")

	apiMiddleware.HandleFunc("/create_semester_location", authHandler.CreateSemesterLocation).Methods("POST")
	apiMiddleware.HandleFunc("/update_semester_location", authHandler.UpdateSemesterLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_semester_location", authHandler.DeleteSemesterLocation).Methods("POST")

	apiMiddleware.HandleFunc("/update_admin_staff", authHandler.UpdateAdminStaff).Methods("POST")
	apiMiddleware.HandleFunc("/delete_admin_staff", authHandler.DeleteAdminStaff).Methods("POST")

	apiMiddleware.HandleFunc("/update_program", authHandler.UpdateProgram).Methods("POST")
	apiMiddleware.HandleFunc("/delete_program", authHandler.DeleteProgram).Methods("POST")

	apiMiddleware.HandleFunc("/update_semester", authHandler.UpdateSemester).Methods("POST")
	apiMiddleware.HandleFunc("/delete_semester", authHandler.DeleteSemester).Methods("POST")

	apiMiddleware.HandleFunc("/update_district", authHandler.UpdateDistrict).Methods("POST")
	apiMiddleware.HandleFunc("/delete_district", authHandler.DeleteDistrict).Methods("POST")

	apiMiddleware.HandleFunc("/update_location", authHandler.UpdateLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_location", authHandler.DeleteLocation).Methods("POST")

	apiMiddleware.HandleFunc("/update_student", authHandler.UpdateStudent).Methods("POST")
	apiMiddleware.HandleFunc("/delete_student", authHandler.DeleteStudent).Methods("POST")

	apiMiddleware.HandleFunc("/update_material", authHandler.UpdateMaterial).Methods("POST")
	apiMiddleware.HandleFunc("/delete_material", authHandler.DeleteMaterial).Methods("POST")

	apiMiddleware.HandleFunc("/update_tutor", authHandler.UpdateTutor).Methods("POST")
	apiMiddleware.HandleFunc("/delete_tutor", authHandler.DeleteTutor).Methods("POST")

	apiMiddleware.HandleFunc("/create_program_location", authHandler.CreateProgramLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_program_location", authHandler.DeleteProgramLocation).Methods("POST")

	apiMiddleware.HandleFunc("/session_search", authHandler.GetSessionSearch).Methods("GET")
	apiMiddleware.HandleFunc("/student_search", authHandler.GetStudentSessionSearch).Methods("GET")

	apiMiddleware.HandleFunc("/get_tutor_sessions", authHandler.GetTutorSessionAnalytics).Methods("GET")

	apiMiddleware.HandleFunc("/create_session", authHandler.CreateStudentSession).Methods("POST")
	apiMiddleware.HandleFunc("/create_assessment", authHandler.CreateAssessment).Methods("POST")
	apiMiddleware.HandleFunc("/update_assessment", authHandler.UpdateAssessment).Methods("POST")
	apiMiddleware.HandleFunc("/delete_assessment", authHandler.DeleteAssessment).Methods("POST")

	apiMiddleware.HandleFunc("/create_tutor_location", authHandler.CreateTutorLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_tutor_location", authHandler.DeleteTutorLocation).Methods("POST")

	apiMiddleware.HandleFunc("/create_subject", authHandler.CreateSubject).Methods("POST")
	apiMiddleware.HandleFunc("/update_subject", authHandler.UpdateSubject).Methods("POST")
	apiMiddleware.HandleFunc("/delete_subject", authHandler.DeleteSubject).Methods("POST")

	apiMiddleware.HandleFunc("/create_permission", authHandler.CreatePermission).Methods("POST")

	apiMiddleware.HandleFunc("/create_subject_location", authHandler.CreateSubjectLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_subject_location", authHandler.DeleteSubjectLocation).Methods("POST")
	apiMiddleware.HandleFunc("/location_subject_list", authHandler.GetSubjectLocations).Methods("GET")

	apiMiddleware.HandleFunc("/update_announcement", authHandler.UpdateAnnouncement).Methods("POST")
	apiMiddleware.HandleFunc("/create_announcement", authHandler.CreateAnnouncement).Methods("POST")
	apiMiddleware.HandleFunc("/delete_announcement", authHandler.DeleteAnnouncement).Methods("POST")

	apiMiddleware.HandleFunc("/get_subjects", authHandler.GetSubjects).Methods("GET")
	apiMiddleware.HandleFunc("/get_session_info", authHandler.GetSessionInfo).Methods("GET")
	apiMiddleware.HandleFunc("/get_student_info", authHandler.GetStudentInfo).Methods("GET")

	apiMiddleware.HandleFunc("/get_announcements", authHandler.GetAnnouncements).Methods("GET")

	apiMiddleware.HandleFunc("/get_locations", authHandler.GetLocations).Methods("GET")

	apiMiddleware.HandleFunc("/get_tutors", authHandler.GetTutors).Methods("GET")
	apiMiddleware.HandleFunc("/get_materials", authHandler.GetMaterials).Methods("GET")
	apiMiddleware.HandleFunc("/get_programs", authHandler.GetPrograms).Methods("GET")
	apiMiddleware.HandleFunc("/get_districts", authHandler.GetDistricts).Methods("GET")
	apiMiddleware.HandleFunc("/get_students", authHandler.GetStudents).Methods("GET")
	apiMiddleware.HandleFunc("/get_admins", authHandler.GetAdmins).Methods("GET")
	apiMiddleware.HandleFunc("/get_semesters", authHandler.GetSemesters).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessments", authHandler.GetAssessments).Methods("GET")
	apiMiddleware.HandleFunc("/get_tutor_locations", authHandler.GetTutorLocations).Methods("GET")
	apiMiddleware.HandleFunc("/location_program_list", authHandler.GetLocationPrograms).Methods("GET")
	apiMiddleware.HandleFunc("/get_permissions", authHandler.GetPermissions).Methods("GET")
	apiMiddleware.HandleFunc("/get_org_permissions", authHandler.GetOrganizationPermissions).Methods("GET")
	apiMiddleware.HandleFunc("/semester_location_list", authHandler.GetSemesterLocations).Methods("GET")

	apiMiddleware.HandleFunc("/get_session_bchart", authHandler.GetSessionBChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessment_bchart", authHandler.GetAssessmentBChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_programs_bchart", authHandler.GetProgramsBChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_tutors_bchart", authHandler.GetTutorsBChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_session_analytics", authHandler.GetSessionAnalytics).Methods("GET")
	apiMiddleware.HandleFunc("/get_tutor_session_analytics", authHandler.GetSessionAnalytics).Methods("GET")

	apiMiddleware.HandleFunc("/get_assessment_trend", authHandler.GetAssessmentTrendLine).Methods("GET")
	apiMiddleware.HandleFunc("/get_session_trend", authHandler.GetSessionTrendLine).Methods("GET")
	apiMiddleware.HandleFunc("/get_semester_assessments_data", authHandler.GetSemestersVAssessmentChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessments_growth_data", authHandler.GetAssessmentGrowth).Methods("GET")

	apiMiddleware.HandleFunc("/tutor_big_data", authHandler.UploadTutorBigData).Methods("POST")
	apiMiddleware.HandleFunc("/student_big_data", authHandler.UploadStudentBigData).Methods("POST")

	r.PathPrefix("/api").Handler(apiMiddleware)

	handler := corsOptions.Handler(r)
	httpListen := &http.Server{
		Addr:           ":3333",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // MAX header size might be a issue with ZIP files ??
		Handler:        handler,
	}

	err := httpListen.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
