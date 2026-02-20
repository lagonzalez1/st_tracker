package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
	"tracker/app/cache"
	"tracker/app/config"
	"tracker/app/database"
	"tracker/app/middleware"
	"tracker/app/services"
	"tracker/app/sqs"
	"tracker/app/transport"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World end.!")
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	fmt.Println("POD STARTING...")

	db, db_err := database.ConnectDB(cfg.DB)
	if db_err != nil {
		log.Fatalf("Database connection failed: %v", db_err)
	}
	defer db.Close()

	s3Client, s3Error := config.ConnectS3()
	if s3Error != nil {
		fmt.Printf("s3 client connection failed: %v", s3Error)
		log.Fatalf("s3 client connection failed: %v", s3Error)
	}

	sqsClient, sqsError := config.ConnectSQS()
	if sqsError != nil {
		fmt.Printf("sqs connection failed: %v", sqsError)
		log.Fatalf("sqs connection failed: %v", sqsError)
	}

	// Interface
	valkeyClient, valError := config.LoadValKey(cfg)
	if valError != nil {
		fmt.Printf("valkey client connection failed: %v", valError)
		log.Fatalf("valkey client connection failed: %v", valError)
	}
	defer valkeyClient.Close()

	r := mux.NewRouter()

	cacheHandler := cache.New(valkeyClient)
	sqsHandler := sqs.New(sqsClient)

	authService := services.NewAuthService(db, s3Client, nil, valkeyClient, sqsClient)
	authHandler := transport.NewAuthHandler(authService, cacheHandler, sqsHandler, cfg)

	apiMiddleware := mux.NewRouter().PathPrefix("/api").Subrouter()
	apiMiddleware.Use(middleware.Middleware(authService, cacheHandler, cfg))

	// Testing allow different cors http://localhost:3000
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://presentifyclone.click", "http://localhost:3000", "https://checkout.stripe.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Stripe-Account", "Stripe-Signature", "X-District-ID", "X-Location-ID"},
		ExposedHeaders:   []string{"X-Access-Token"},
		AllowCredentials: true,
		Debug:            false,
	})

	r.HandleFunc("/health_check", authHandler.HealthCheck).Methods("GET")
	r.HandleFunc("/ready", authHandler.Ready).Methods("GET")
	r.HandleFunc("/hello", hello).Methods("GET")
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/get_assessment_questions_external", authHandler.GetAssessmentQuestionsExternal).Methods("GET")
	r.HandleFunc("/get_assessment_preassessment_external", authHandler.GetAssessmentQuestionsExternal).Methods("GET")
	r.HandleFunc("/get_assessment_session", authHandler.GetAssessmentSessionExternal).Methods("GET")
	r.HandleFunc("/get_student_details", authHandler.GetStudentDetails).Methods("GET")
	r.HandleFunc("/create_organization", authHandler.CreateOrganization).Methods("POST")
	r.HandleFunc("/create_student_assessment_response", authHandler.CreateStudentAssessmentResponse).Methods("POST")
	r.HandleFunc("/create_pre_assessment", authHandler.CreatePreAssessment).Methods("POST")
	r.HandleFunc("/stripe_webhook", authHandler.StripeWebhook).Methods("POST")

	apiMiddleware.HandleFunc("/create_student", authHandler.CreateStudent).Methods("POST")
	apiMiddleware.HandleFunc("/create_location", authHandler.CreateLocation).Methods("POST")
	apiMiddleware.HandleFunc("/create_admin_staff", authHandler.CreateAdminStaff).Methods("POST")
	apiMiddleware.HandleFunc("/create_district", authHandler.CreateDistrict).Methods("POST")
	apiMiddleware.HandleFunc("/create_program", authHandler.CreateProgram).Methods("POST")
	apiMiddleware.HandleFunc("/create_material", authHandler.CreateMaterial).Methods("POST")
	apiMiddleware.HandleFunc("/create_tutor", authHandler.CreateTutor).Methods("POST")
	apiMiddleware.HandleFunc("/create_semester", authHandler.CreateSemester).Methods("POST")

	apiMiddleware.HandleFunc("/create_semester_dates", authHandler.CreateSemesterDates).Methods("POST")
	apiMiddleware.HandleFunc("/delete_semester_dates", authHandler.DeleteSemesterDates).Methods("POST")
	apiMiddleware.HandleFunc("/get_semester_dates", authHandler.GetSemesterDates).Methods("GET")

	apiMiddleware.HandleFunc("/create_student_group", authHandler.CreateStudentGroup).Methods("POST")
	apiMiddleware.HandleFunc("/delete_student_group", authHandler.DeleteStudentGroup).Methods("POST")
	apiMiddleware.HandleFunc("/update_student_group", authHandler.UpdateStudentGroup).Methods("POST")

	apiMiddleware.HandleFunc("/create_semester_location", authHandler.CreateSemesterLocation).Methods("POST")
	apiMiddleware.HandleFunc("/update_semester_location", authHandler.UpdateSemesterLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_semester_location", authHandler.DeleteSemesterLocation).Methods("POST")

	apiMiddleware.HandleFunc("/create_admin_location", authHandler.CreateAdminLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_admin_location", authHandler.DeleteAdminLocation).Methods("POST")

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

	apiMiddleware.HandleFunc("/delete_schedule", authHandler.DeleteSchedule).Methods("POST")

	apiMiddleware.HandleFunc("/create_program_location", authHandler.CreateProgramLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_program_location", authHandler.DeleteProgramLocation).Methods("POST")

	apiMiddleware.HandleFunc("/session_search", authHandler.GetSessionSearch).Methods("GET")
	apiMiddleware.HandleFunc("/student_search", authHandler.GetStudentSessionSearch).Methods("GET")
	apiMiddleware.HandleFunc("/tutor_search", authHandler.GetTutorSearch).Methods("GET")
	apiMiddleware.HandleFunc("/student_assessment_search", authHandler.GetStudentAssesssmentSearch).Methods("GET")
	apiMiddleware.HandleFunc("/get_tutor_low_performance", authHandler.GetTutorLowPerformance).Methods("GET")
	apiMiddleware.HandleFunc("/get_absent_present", authHandler.GetAbsentPresent).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessment_completion", authHandler.GetAssessmentCompletion).Methods("GET")

	apiMiddleware.HandleFunc("/get_student_groups", authHandler.GetStudentGroups).Methods("GET")

	apiMiddleware.HandleFunc("/get_tutor_sessions", authHandler.GetTutorSessionAnalytics).Methods("GET")
	apiMiddleware.HandleFunc("/get_sessions", authHandler.GetTutorsSessions).Methods("GET")
	apiMiddleware.HandleFunc("/get_recent_sessions", authHandler.GetRecentSessionsTutors).Methods("GET")
	apiMiddleware.HandleFunc("/get_recent_location_sessions", authHandler.GetRecentLocationSessions).Methods("GET")
	apiMiddleware.HandleFunc("/get_location_session_average", authHandler.GetLocationSessionAverage).Methods("GET")

	apiMiddleware.HandleFunc("/create_s3_object", authHandler.CreateS3Object).Methods("POST")
	apiMiddleware.HandleFunc("/delete_s3_object", authHandler.DeleteS3Object).Methods("POST")
	apiMiddleware.HandleFunc("/create_teacher", authHandler.CreateTeacher).Methods("POST")

	apiMiddleware.HandleFunc("/announcement_ack", authHandler.CreateAckAnnouncements).Methods("POST")
	apiMiddleware.HandleFunc("/update_teacher", authHandler.UpdateTeacher).Methods("POST")
	apiMiddleware.HandleFunc("/delete_teacher", authHandler.DeleteTeacher).Methods("POST")
	apiMiddleware.HandleFunc("/get_teachers", authHandler.GetTeachers).Methods("GET")
	apiMiddleware.HandleFunc("/get_group_attendies", authHandler.GetGroupAttendies).Methods("GET")

	apiMiddleware.HandleFunc("/create_student_group_attendies", authHandler.CreateStudentGroupAttendies).Methods("POST")

	apiMiddleware.HandleFunc("/create_session", authHandler.CreateStudentSession).Methods("POST")
	apiMiddleware.HandleFunc("/create_survey_response", authHandler.CreateSurveyResponse).Methods("POST")
	apiMiddleware.HandleFunc("/create_assessment", authHandler.CreateAssessment).Methods("POST")
	apiMiddleware.HandleFunc("/update_assessment", authHandler.UpdateAssessment).Methods("POST")
	apiMiddleware.HandleFunc("/delete_assessment", authHandler.DeleteAssessment).Methods("POST")

	apiMiddleware.HandleFunc("/create_tutor_location", authHandler.CreateTutorLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_tutor_location", authHandler.DeleteTutorLocation).Methods("POST")

	apiMiddleware.HandleFunc("/delete_survey", authHandler.DeleteSurvey).Methods("POST")
	apiMiddleware.HandleFunc("/delete_session", authHandler.DeleteSession).Methods("POST")

	apiMiddleware.HandleFunc("/create_program_survey", authHandler.CreateProgramSurvey).Methods("POST")
	apiMiddleware.HandleFunc("/get_program_surveys", authHandler.GetProgramSurveys).Methods("GET")
	apiMiddleware.HandleFunc("/delete_program_survey", authHandler.DeleteProgramSurvey).Methods("POST")

	apiMiddleware.HandleFunc("/create_subject", authHandler.CreateSubject).Methods("POST")
	apiMiddleware.HandleFunc("/update_subject", authHandler.UpdateSubject).Methods("POST")
	apiMiddleware.HandleFunc("/delete_subject", authHandler.DeleteSubject).Methods("POST")

	// Implementation
	apiMiddleware.HandleFunc("/create_schedule_global", authHandler.CreateGlobalSchedule).Methods("POST")
	apiMiddleware.HandleFunc("/update_schedule_global", authHandler.UpdateGlobalSchedule).Methods("POST")
	apiMiddleware.HandleFunc("/delete_schedule_global", authHandler.DeleteGlobalSchedule).Methods("POST")
	apiMiddleware.HandleFunc("/get_schedule_global", authHandler.GetGlobalSchedule).Methods("GET")

	apiMiddleware.HandleFunc("/create_location_contact", authHandler.CreateLocationContact).Methods("POST")
	apiMiddleware.HandleFunc("/update_location_contact", authHandler.UpdateLocationContact).Methods("POST")
	apiMiddleware.HandleFunc("/delete_location_contact", authHandler.DeleteLocationContact).Methods("POST")

	apiMiddleware.HandleFunc("/get_location_contact", authHandler.GetLocationContact).Methods("GET")
	apiMiddleware.HandleFunc("/create_permission", authHandler.CreatePermission).Methods("POST")

	apiMiddleware.HandleFunc("/v3/create_schedule", authHandler.CreateScheduleV3).Methods("POST")
	apiMiddleware.HandleFunc("/v3/delete_entity_schedule", authHandler.DeleteEntitySchedule).Methods("POST")
	apiMiddleware.HandleFunc("/v3/get_entity_schedule", authHandler.GetEntitySchedule).Methods("GET")

	// This is not in use
	apiMiddleware.HandleFunc("/create_schedule", authHandler.CreateSchedule).Methods("POST")
	apiMiddleware.HandleFunc("/v2/create_schedule", authHandler.CreateSchedule).Methods("POST")
	apiMiddleware.HandleFunc("/v2/delete_schedule", authHandler.CreateSchedule).Methods("POST")
	apiMiddleware.HandleFunc("/v2/update_schedule", authHandler.CreateSchedule).Methods("POST")
	// This is not in use

	apiMiddleware.HandleFunc("/create_survey", authHandler.CreateSurvey).Methods("POST")
	apiMiddleware.HandleFunc("/update_survey", authHandler.UpdateSurvey).Methods("POST")

	apiMiddleware.HandleFunc("/create_subject_location", authHandler.CreateSubjectLocation).Methods("POST")
	apiMiddleware.HandleFunc("/delete_subject_location", authHandler.DeleteSubjectLocation).Methods("POST")
	apiMiddleware.HandleFunc("/location_subject_list", authHandler.GetSubjectLocations).Methods("GET")

	apiMiddleware.HandleFunc("/update_announcement", authHandler.UpdateAnnouncement).Methods("POST")
	apiMiddleware.HandleFunc("/create_announcement", authHandler.CreateAnnouncement).Methods("POST")
	apiMiddleware.HandleFunc("/delete_announcement", authHandler.DeleteAnnouncement).Methods("POST")

	apiMiddleware.HandleFunc("/get_object", authHandler.GetObject).Methods("GET")

	apiMiddleware.HandleFunc("/get_subjects", authHandler.GetSubjects).Methods("GET")
	apiMiddleware.HandleFunc("/get_surveys", authHandler.GetSurveys).Methods("GET")
	apiMiddleware.HandleFunc("/get_surveys_program_id", authHandler.GetSurveysByProgram).Methods("GET")
	apiMiddleware.HandleFunc("/get_session_info", authHandler.GetSessionInfo).Methods("GET")
	apiMiddleware.HandleFunc("/get_student_info", authHandler.GetStudentInfo).Methods("GET")

	apiMiddleware.HandleFunc("/get_announcements", authHandler.GetAnnouncements).Methods("GET")
	apiMiddleware.HandleFunc("/get_announcements_ack", authHandler.GetAnnouncementsAck).Methods("GET")
	apiMiddleware.HandleFunc("/get_tutor_file", authHandler.GetTutorFile).Methods("GET")
	apiMiddleware.HandleFunc("/get_student_file", authHandler.GetStudentFile).Methods("GET")
	apiMiddleware.HandleFunc("/get_locations", authHandler.GetLocations).Methods("GET")
	apiMiddleware.HandleFunc("/get_organization", authHandler.GetOrganization).Methods("GET")
	apiMiddleware.HandleFunc("/get_generation_usage", authHandler.GetGenerationUsage).Methods("GET")

	// Return session id as well for duplicate
	apiMiddleware.HandleFunc("/get_signed_url_materials", authHandler.GetSignedUrlMaterials).Methods("GET")
	apiMiddleware.HandleFunc("/get_session_accountability", authHandler.GetSessionAccountability).Methods("GET")
	apiMiddleware.HandleFunc("/get_sessions_scheduled", authHandler.GetSessionScheduled).Methods("GET")
	apiMiddleware.HandleFunc("/get_entity_scheduled_shift", authHandler.GetEntityScheduleShift).Methods("GET")
	apiMiddleware.HandleFunc("/get_tutors", authHandler.GetTutors).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessment_questions", authHandler.GetAssessmentQuestions).Methods("GET")
	apiMiddleware.HandleFunc("/get_schedule", authHandler.GetSchedules).Methods("GET")
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

	apiMiddleware.HandleFunc("/get_cycle_growth", authHandler.GetCycleGrowth).Methods("GET")
	apiMiddleware.HandleFunc("/get_cycle_growth_delimiters", authHandler.GetCycleGrowthDelim).Methods("GET")

	apiMiddleware.HandleFunc("/get_session_bchart", authHandler.GetSessionBChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessment_bchart", authHandler.GetAssessmentBChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_programs_bchart", authHandler.GetProgramsBChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_tutors_bchart", authHandler.GetTutorsBChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_session_analytics", authHandler.GetSessionAnalytics).Methods("GET")
	apiMiddleware.HandleFunc("/get_session_analytics_local", authHandler.GetSessionAnalyticsLocal).Methods("GET")
	apiMiddleware.HandleFunc("/get_tutor_session_analytics", authHandler.GetTutorSessionAnalytics).Methods("GET")

	apiMiddleware.HandleFunc("/get_assessment_trend", authHandler.GetAssessmentTrendLine).Methods("GET")
	apiMiddleware.HandleFunc("/get_session_trend", authHandler.GetSessionTrendLine).Methods("GET")
	apiMiddleware.HandleFunc("/get_semester_assessments_data", authHandler.GetSemestersVAssessmentChart).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessments_growth_data", authHandler.GetAssessmentGrowth).Methods("GET")
	apiMiddleware.HandleFunc("/get_session_vscore", authHandler.GetSessionVScore).Methods("GET")
	apiMiddleware.HandleFunc("/get_student_vassessment", authHandler.GetStudentVAssessments).Methods("GET")

	apiMiddleware.HandleFunc("/tutor_big_data", authHandler.UploadTutorBigData).Methods("POST")
	apiMiddleware.HandleFunc("/student_big_data", authHandler.UploadStudentBigData).Methods("POST")

	apiMiddleware.HandleFunc("/create_assessment_sessions", authHandler.CreateAssessmentSessions).Methods("POST")
	apiMiddleware.HandleFunc("/delete_assessment_sessions", authHandler.DeleteAssessmentSessions).Methods("POST")
	apiMiddleware.HandleFunc("/get_student_assessment_choices", authHandler.GetStudentAssessmentChoices).Methods("GET")
	apiMiddleware.HandleFunc("/delete_student_session", authHandler.DeleteStudentSession).Methods("POST")
	apiMiddleware.HandleFunc("/get_student_assessment_sessions", authHandler.GetStudentAssessmentSessions).Methods("GET")

	apiMiddleware.HandleFunc("/get_generated_questions", authHandler.GetGeneratedQuestion).Methods("GET")
	apiMiddleware.HandleFunc("/get_generated_materials", authHandler.GetGeneratedMaterials).Methods("GET")
	apiMiddleware.HandleFunc("/delete_generated_assessment", authHandler.MicroEventDeleteQuestions).Methods("POST")
	apiMiddleware.HandleFunc("/micro_generate", authHandler.MicroEventGenerate).Methods("POST")

	apiMiddleware.HandleFunc("/micro_student_report", authHandler.MicroEventStartStudentReport).Methods("POST")
	apiMiddleware.HandleFunc("/get_student_report", authHandler.MicroGetStudentReport).Methods("GET")

	apiMiddleware.HandleFunc("/get_tutor_file_url", authHandler.MicroGetTutorFile).Methods("GET")
	apiMiddleware.HandleFunc("/get_student_file_url", authHandler.MicroGetStudentFile).Methods("GET")

	apiMiddleware.HandleFunc("/get_sentiment", authHandler.GetSentiment).Methods("GET")
	apiMiddleware.HandleFunc("/get_sentiment_tutor", authHandler.GetSentimentByTutor).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessments_tutor", authHandler.GetAssessmentByTutor).Methods("GET")
	apiMiddleware.HandleFunc("/get_assessment_content", authHandler.GetAssessmentContent).Methods("GET")

	apiMiddleware.HandleFunc("/create_checkout_session", authHandler.CreateCheckoutSession).Methods("POST")
	apiMiddleware.HandleFunc("/create_portal_session", authHandler.CreatePortalSession).Methods("POST")
	apiMiddleware.HandleFunc("/get_subscriptions", authHandler.GetSubscriptions).Methods("GET")
	apiMiddleware.HandleFunc("/get_admin_locations", authHandler.GetAdminLocations).Methods("GET")

	r.PathPrefix("/api").Handler(apiMiddleware)

	fmt.Println("Number of cores avilable: ", runtime.GOMAXPROCS(0))

	handler := corsOptions.Handler(r)
	httpListen := &http.Server{
		Addr:           ":3333",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // MAX header size might be a issue with ZIP files ??
		Handler:        handler,
	}

	err = httpListen.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
