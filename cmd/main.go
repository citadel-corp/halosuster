package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/citadel-corp/halosuster/internal/common/db"
	"github.com/citadel-corp/halosuster/internal/common/middleware"
	"github.com/citadel-corp/halosuster/internal/image"
	"github.com/citadel-corp/halosuster/internal/medicalpatients"
	"github.com/citadel-corp/halosuster/internal/medicalrecords"
	"github.com/citadel-corp/halosuster/internal/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Connect to database
	// env := os.Getenv("ENV")
	// sslMode := "disable"
	// if env == "production" {
	// 	sslMode = "verify-full sslrootcert=ap-southeast-1-bundle.pem"
	// }
	// connStr := "postgres://[user]:[password]@[neon_hostname]/[dbname]?sslmode=require"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_PARAMS"))
	// dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	// 	os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), sslMode)
	db, err := db.Connect(connStr)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Cannot connect to database: %v", err))
		os.Exit(1)
	}

	// Create migrations
	// err = db.UpMigration()
	// if err != nil {
	// 	log.Error().Msg(fmt.Sprintf("Up migration failed: %v", err))
	// 	os.Exit(1)
	// }

	// initialize user domain
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	// initialize medical patient domain
	medicalPatientRepository := medicalpatients.NewRepository(db)
	medicalPatientService := medicalpatients.NewService(medicalPatientRepository)
	medicalPatientHandler := medicalpatients.NewHandler(medicalPatientService)

	// initialize medical record domain
	medicalRecordsRepository := medicalrecords.NewRepository(db)
	medicalRecordsService := medicalrecords.NewService(medicalRecordsRepository, medicalPatientRepository)
	medicalRecordsHandler := medicalrecords.NewHandler(medicalRecordsService)

	// initialize image domain
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Cannot create AWS session: %v", err))
		os.Exit(1)
	}
	imageService := image.NewService(sess)
	imageHandler := image.NewHandler(imageService)

	r := mux.NewRouter()
	r.Use(middleware.Logging)
	r.Use(middleware.PanicRecoverer)
	v1 := r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Service ready")
	})

	// user routes
	ur := v1.PathPrefix("/user").Subrouter()
	ur.HandleFunc("/it/register", userHandler.CreateITUser).Methods(http.MethodPost)
	ur.HandleFunc("/it/login", userHandler.LoginITUser).Methods(http.MethodPost)
	ur.HandleFunc("/nurse/register", middleware.AuthorizeITUser(userHandler.CreateNurseUser)).Methods(http.MethodPost)
	ur.HandleFunc("/nurse/login", userHandler.LoginNurseUser).Methods(http.MethodPost)
	ur.HandleFunc("", middleware.AuthorizeITUser(userHandler.ListUsers)).Methods(http.MethodGet)
	ur.HandleFunc("/nurse/{userId}", middleware.AuthorizeITUser(userHandler.UpdateNurse)).Methods(http.MethodPut)
	ur.HandleFunc("/nurse/{userId}", middleware.AuthorizeITUser(userHandler.DeleteNurse)).Methods(http.MethodDelete)
	ur.HandleFunc("/nurse/{userId}/access", middleware.AuthorizeITUser(userHandler.GrantNurseAccess)).Methods(http.MethodPost)

	// image routes
	ir := v1.PathPrefix("/image").Subrouter()
	ir.HandleFunc("", middleware.AuthorizeITAndNurseUser(imageHandler.UploadToS3)).Methods(http.MethodPost)

	// medical patient routes
	mpr := v1.PathPrefix("/medical/patient").Subrouter()
	mpr.HandleFunc("", middleware.AuthorizeITAndNurseUser(medicalPatientHandler.CreateMedicalPatient)).Methods(http.MethodPost)

	// medical record routes
	mr := v1.PathPrefix("/medical/record").Subrouter()
	mr.HandleFunc("", middleware.AuthorizeITAndNurseUser(medicalRecordsHandler.CreateMedicalRecord)).Methods(http.MethodPost)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Info().Msg(fmt.Sprintf("HTTP server listening on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error().Msg(fmt.Sprintf("HTTP server error: %v", err))
		}
		log.Info().Msg("Stopped serving new connections.")
	}()

	// Listen for the termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until termination signal received
	<-stop
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	log.Info().Msg(fmt.Sprintf("Shutting down HTTP server listening on %s", httpServer.Addr))
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Msg(fmt.Sprintf("HTTP server shutdown error: %v", err))
	}
	log.Info().Msg("Shutdown complete.")
}
