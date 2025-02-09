package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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
	apiMiddleware.Use(middleware.Middleware)

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	//api.ConnectSheetsAPI()

	r.HandleFunc("/hello", hello).Methods("GET")
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	apiMiddleware.HandleFunc("/create_student", authHandler.CreateStudent).Methods("POST")
	apiMiddleware.HandleFunc("/create_location", authHandler.CreateLocation).Methods("POST")
	apiMiddleware.HandleFunc("/create_admin_staff", authHandler.CreateAdminStaff).Methods("POST")
	apiMiddleware.HandleFunc("/create_district", authHandler.CreateDistrict).Methods("POST")
	apiMiddleware.HandleFunc("/create_program", authHandler.CreateProgram).Methods("POST")
	apiMiddleware.HandleFunc("/create_materials", authHandler.CreateMaterial).Methods("POST")
	apiMiddleware.HandleFunc("/create_tutor", authHandler.CreateTutor).Methods("POST")
	apiMiddleware.HandleFunc("/create_semester", authHandler.CreateSemester).Methods("POST")

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

	apiMiddleware.HandleFunc("/session_search", authHandler.GetSessionSearch).Methods("POST")

	apiMiddleware.HandleFunc("/create_session", authHandler.CreateStudentSession).Methods("POST")
	apiMiddleware.HandleFunc("/create_assessment", authHandler.CreateAssessment).Methods("POST")

	apiMiddleware.HandleFunc("/create_subject", authHandler.CreateSubject).Methods("POST")
	apiMiddleware.HandleFunc("/update_subject", authHandler.UpdateSubject).Methods("POST")
	apiMiddleware.HandleFunc("/delete_subject", authHandler.DeleteSubject).Methods("POST")

	apiMiddleware.HandleFunc("/get_subjects", authHandler.GetSubjects).Methods("GET")

	apiMiddleware.HandleFunc("/get_locations", authHandler.GetLocations).Methods("GET")
	apiMiddleware.HandleFunc("/get_tutors", authHandler.GetTutors).Methods("GET")
	apiMiddleware.HandleFunc("/get_materials", authHandler.GetMaterials).Methods("GET")
	apiMiddleware.HandleFunc("/get_programs", authHandler.GetPrograms).Methods("GET")
	apiMiddleware.HandleFunc("/get_districts", authHandler.GetDistricts).Methods("GET")
	apiMiddleware.HandleFunc("/get_students", authHandler.GetStudents).Methods("GET")
	apiMiddleware.HandleFunc("/get_admins", authHandler.GetAdmins).Methods("GET")
	apiMiddleware.HandleFunc("/get_semesters", authHandler.GetSemesters).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessments", authHandler.GetAssessments).Methods("GET")

	apiMiddleware.HandleFunc("/location_program_list", authHandler.GetLocationPrograms).Methods("GET")

	r.PathPrefix("/api").Handler(apiMiddleware)

	handler := corsOptions.Handler(r)

	err := http.ListenAndServe(":3333", handler)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
