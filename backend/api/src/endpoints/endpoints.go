package endpoints

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	"soul-connection.com/api/src/endpoints/clothes"
	"soul-connection.com/api/src/endpoints/customers"
	"soul-connection.com/api/src/endpoints/employees"
	"soul-connection.com/api/src/endpoints/encounters"
	"soul-connection.com/api/src/endpoints/events"
	"soul-connection.com/api/src/endpoints/payments"
	"soul-connection.com/api/src/endpoints/tips"
	"soul-connection.com/api/src/middleware"
)

type Endpoint struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

type ModelRoutes struct {
	BasePath string
	Routes   []Endpoint
}

func CreateRouter(database *sql.DB, fileStorage *mongo.Database) (*mux.Router, error) {
	apiKey := os.Getenv("API_KEY")

	employeesBucket, err := gridfs.NewBucket(fileStorage, options.GridFSBucket().SetName("employeesBucket"))
	customersBucket, err := gridfs.NewBucket(fileStorage, options.GridFSBucket().SetName("customersBucket"))
	clothesBucket, err := gridfs.NewBucket(fileStorage, options.GridFSBucket().SetName("clothesBucket"))
	if err != nil {
		return nil, err
	}

	// authModel := auth.AuthModel{Auth: auth.ApiKeyAuth{ApiKey: apiKey}}
	employeeModel := employees.EmployeesModel{Employees: employees.EmployeesDB{DB: database, Bucket: employeesBucket}}
	customerModel := customers.CustomersModel{Customers: customers.CustomersDB{DB: database, Bucket: customersBucket}}
	eventModel := events.EventModel{Events: events.EventsDB{DB: database}}
	paymentModel := payments.PaymentModel{Payments: payments.PaymentsDB{DB: database}}
	encounterModel := encounters.EncounterModel{Encounters: encounters.EncountersDB{DB: database}}
	clotheModel := clothes.ClothesModel{Clothes: clothes.ClothesDB{DB: database, Bucket: clothesBucket}}
	tipModel := tips.TipModel{Tips: tips.TipsDB{DB: database}}

	// publicRoutes := []ModelRoutes{
	// 	{
	// 		BasePath: "/api/auth",
	// 		Routes: []Endpoint{
	// 			{Path: "/login", Handler: authModel.Login, Method: http.MethodPost},
	// 		},
	// 	},
	// }

	publicRoutes := []ModelRoutes{
		{
			BasePath: "/api/employees",
			Routes: []Endpoint{
				{Path: "", Handler: employeeModel.GetAllEmployees, Method: http.MethodGet},
				{Path: "", Handler: employeeModel.AddEmployee, Method: http.MethodPost},
				{Path: "/{employee_id}", Handler: employeeModel.GetEmployeeById, Method: http.MethodGet},
				{Path: "/{employee_id}", Handler: employeeModel.DeleteEmployee, Method: http.MethodDelete},
				{Path: "/{employee_id}", Handler: employeeModel.PatchEmployee, Method: http.MethodPatch},
				{Path: "/{employee_id}/image", Handler: employeeModel.GetImage, Method: http.MethodGet},
			},
		},
		{
			BasePath: "/api/customers",
			Routes: []Endpoint{
				{Path: "", Handler: customerModel.GetAllCustomers, Method: http.MethodGet},
				{Path: "", Handler: customerModel.AddCustomer, Method: http.MethodPost},
				{Path: "/{customer_id}", Handler: customerModel.GetCustomerById, Method: http.MethodGet},
				{Path: "/{customer_id}", Handler: customerModel.DeleteCustomer, Method: http.MethodDelete},
				{Path: "/{customer_id}", Handler: customerModel.PatchCustomer, Method: http.MethodPatch},
				{Path: "/{customer_id}/image", Handler: customerModel.GetImage, Method: http.MethodGet},
				{Path: "/employees/{employee_id}", Handler: customerModel.GetCustomerById, Method: http.MethodGet},
			},
		},
		{
			BasePath: "/api/events",
			Routes: []Endpoint{
				{Path: "", Handler: eventModel.GetAllEvents, Method: http.MethodGet},
				{Path: "", Handler: eventModel.AddEvent, Method: http.MethodPost},
				{Path: "/{event_id}", Handler: eventModel.GetEventsById, Method: http.MethodGet},
				{Path: "/{event_id}", Handler: eventModel.DeleteEvent, Method: http.MethodDelete},
				{Path: "/{event_id}", Handler: eventModel.PatchEvent, Method: http.MethodPatch},
			},
		},
		{
			BasePath: "/api/payments",
			Routes: []Endpoint{
				{Path: "", Handler: paymentModel.GetAllPayments, Method: http.MethodGet},
				{Path: "", Handler: paymentModel.AddPayment, Method: http.MethodPost},
				{Path: "/{payment_id}", Handler: paymentModel.GetPaymentsById, Method: http.MethodGet},
				{Path: "/{payment_id}", Handler: paymentModel.DeletePayment, Method: http.MethodDelete},
				{Path: "/{payment_id}", Handler: paymentModel.PatchPayment, Method: http.MethodPatch},
				{Path: "/customer/{customer_id}", Handler: paymentModel.GetPaymentsByCustomerId, Method: http.MethodGet},
			},
		},
		{
			BasePath: "/api/encounters",
			Routes: []Endpoint{
				{Path: "", Handler: encounterModel.GetAllEncounters, Method: http.MethodGet},
				{Path: "", Handler: encounterModel.AddEncounter, Method: http.MethodPost},
				{Path: "/{encounter_id}", Handler: encounterModel.GetEncounterById, Method: http.MethodGet},
				{Path: "/{encounter_id}", Handler: encounterModel.DeleteEncounter, Method: http.MethodDelete},
				{Path: "/{encounter_id}", Handler: encounterModel.PatchEncounter, Method: http.MethodPatch},
				{Path: "/customer/{customer_id}", Handler: encounterModel.GetEncounterByCustomerId, Method: http.MethodGet},
			},
		},
		{
			BasePath: "/api/clothes",
			Routes: []Endpoint{
				{Path: "", Handler: clotheModel.GetAllClothes, Method: http.MethodGet},
				{Path: "", Handler: clotheModel.AddClothe, Method: http.MethodPost},
				{Path: "/{clothe_id}", Handler: clotheModel.GetClotheById, Method: http.MethodGet},
				{Path: "/{clothe_id}", Handler: clotheModel.DeleteClothe, Method: http.MethodDelete},
				{Path: "/{clothe_id}", Handler: clotheModel.PatchClothes, Method: http.MethodPatch},
				{Path: "/{clothe_id}/image", Handler: clotheModel.GetImage, Method: http.MethodGet},
				{Path: "/customer/{customer_id}", Handler: clotheModel.GetClotheByCustomerId, Method: http.MethodGet},
			},
		},
		{
			BasePath: "/api/tips",
			Routes: []Endpoint{
				{Path: "", Handler: tipModel.GetAllTips, Method: http.MethodGet},
				{Path: "", Handler: tipModel.AddTip, Method: http.MethodPost},
				{Path: "/{tip_id}", Handler: tipModel.GetTipById, Method: http.MethodGet},
				{Path: "/{tip_id}", Handler: tipModel.DeleteTip, Method: http.MethodDelete},
				{Path: "/{tip_id}", Handler: tipModel.PatchTips, Method: http.MethodPatch},
			},
		},
	}

	router := mux.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("WEB_URL")},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler)
	router.Use(middleware.Logging)

	publicRouter := router.PathPrefix("").Subrouter()

	authProvider := middleware.AuthProvider{ApiKey: apiKey}
	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(authProvider.Auth)

	attachRoutes(publicRouter, publicRoutes)
	// attachRoutes(protectedRouter, protectedRoutes)

	return router, nil
}

func attachRoutes(router *mux.Router, endpoints []ModelRoutes) {
	for _, endpoint := range endpoints {
		for _, route := range endpoint.Routes {
			fullPath := endpoint.BasePath + route.Path
			router.HandleFunc(fullPath, route.Handler).Methods(route.Method)
		}
	}
}
