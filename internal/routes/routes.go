package routes

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/lucasBiazon/botany-back/internal/middleware"
	"github.com/lucasBiazon/botany-back/internal/repositories"
	usecases_categoryplant "github.com/lucasBiazon/botany-back/internal/usecases/category-plant"
	usecases_plant "github.com/lucasBiazon/botany-back/internal/usecases/plant"
	usecases_specie "github.com/lucasBiazon/botany-back/internal/usecases/specie"
	usecases "github.com/lucasBiazon/botany-back/internal/usecases/user"

	services "github.com/lucasBiazon/botany-back/internal/service"
	handlers "github.com/lucasBiazon/botany-back/internal/web"
)

func InitializeRoutes(db *sql.DB, clientRedis *redis.Client, jwtService services.JWTService) (*chi.Mux, error) {

	// User Routes
	repository := repositories.NewUserRepository(db, clientRedis)
	RegisterUserRoutes := usecases.NewRegisterUserUseCase(repository)
	LoginUserRoutes := usecases.NewLoginUserUseCase(repository)
	FindUserRoutes := usecases.NewFindUserByIdUseCase(repository)
	DeleteUserRoutes := usecases.NewDeleteUserUseCase(repository)
	UpdateUserRoutes := usecases.NewUpdateUserUseCase(repository)
	RequestPasswordResetUserRoutes := usecases.NewRequestPasswordResetUseCase(repository, jwtService)
	ResetPasswordUserRoutes := usecases.NewResetPasswordUserUseCase(repository, jwtService)

	userHandlers := handlers.NewUserHandlers(
		RegisterUserRoutes,
		LoginUserRoutes,
		FindUserRoutes,
		DeleteUserRoutes,
		UpdateUserRoutes,
		RequestPasswordResetUserRoutes,
		ResetPasswordUserRoutes,
	)

	// Category Plant Routes
	repositoryCategoriesPlants := repositories.NewCategoryPlantRepository(db, clientRedis)
	CreateCategoryPlantRoutes := usecases_categoryplant.NewCreateCategoryPlantUseCase(repositoryCategoriesPlants)
	FindAllCategoryPlantRoutes := usecases_categoryplant.NewFindAllCategoryPlantUseCase(repositoryCategoriesPlants)
	FindByCategoryPlantRoutes := usecases_categoryplant.NewFindByIdCategoryPlantUseCase(repositoryCategoriesPlants)
	FindByNameCategoryPlantRoutes := usecases_categoryplant.NewFindByNameCategoryPlantUseCase(repositoryCategoriesPlants)
	UpdateCategoryPlantRoutes := usecases_categoryplant.NewUpdateCategoryPlantUseCase(repositoryCategoriesPlants)
	DeleteCategoryPlantRoutes := usecases_categoryplant.NewDeleteCategoryPlantUseCase(repositoryCategoriesPlants)

	categoryPlantHandlers := handlers.NewCategoryPlantHandler(
		CreateCategoryPlantRoutes,
		DeleteCategoryPlantRoutes,
		FindAllCategoryPlantRoutes,
		FindByCategoryPlantRoutes,
		FindByNameCategoryPlantRoutes,
		UpdateCategoryPlantRoutes,
	)

	// specie Routes
	repositorySpecies := repositories.NewSpeciesRepository(db, clientRedis)
	FindAllSpecieRoutes := usecases_specie.NewFindAllSpecieUseCase(repositorySpecies)
	FindByIdSpecieRoutes := usecases_specie.NewFindByIdSpecieUseCase(repositorySpecies)
	FindByNameSpecieRoutes := usecases_specie.NewFindByNameSpecieUseCase(repositorySpecies)

	specieHandlers := handlers.NewSpecieHandler(
		FindAllSpecieRoutes,
		FindByIdSpecieRoutes,
		FindByNameSpecieRoutes,
	)

	repositoryPlant := repositories.NewPlantRepositoryImpl(db, clientRedis)
	CreatePlantRoute := usecases_plant.NewCreatePlantUseCase(repositoryPlant, repositorySpecies)
	DeletePlantRoute := usecases_plant.NewDeletePlantUseCase(repositoryPlant, repositorySpecies)
	FindAllPlantRoute := usecases_plant.NewFindAllPlantUseCase(repositoryPlant)
	FindByCategoryNamePlantRoute := usecases_plant.NewFindByCategoryNamePlantUseCase(repositoryPlant)
	FindByIdPlantRoute := usecases_plant.NewFindByIdPlantUseCase(repositoryPlant)
	FindByNameCategoryPlantRoute := usecases_plant.NewFindByNameCategoryPlantUseCase(repositoryPlant)
	FindBySpecieNamePlantRoute := usecases_plant.NewFindBySpecieNamePlantUseCase(repositoryPlant)
	UpdatePlantRoute := usecases_plant.NewUpdatePlantUseCase(repositoryPlant)

	plantHandlers := handlers.NewPlantHandler(
		CreatePlantRoute,
		DeletePlantRoute,
		FindAllPlantRoute,
		FindByCategoryNamePlantRoute,
		FindByIdPlantRoute,
		FindByNameCategoryPlantRoute,
		FindBySpecieNamePlantRoute,
		UpdatePlantRoute,
		FindByIdSpecieRoutes,
	)

	// Routes
	r := chi.NewRouter()
	r.Use(middleware.ApiKeyMiddleware)
	r.Use(middleware.RateLimitMiddleware(clientRedis))
	r.Use(middleware.RetryMiddleware(3, 2))
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", userHandlers.RegisterUserHandler)
		r.Post("/register/confirm", userHandlers.ConfirmEmailHandler)
		r.Post("/register/resend-token", userHandlers.ResendTokenHandler)
		r.Post("/login", userHandlers.LoginUserHandler)
		r.Post("/password-reset/request", userHandlers.RequestPasswordResetUserHandler)
		r.Post("/password-reset", userHandlers.ResetPasswordUserHandler)
	})
	r.Route("/api/v1/user", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))
		r.Get("/", userHandlers.FindByIdUserHandler)
		r.Delete("/", userHandlers.DeleteUserHandler)
		r.Put("/", userHandlers.UpdateUserHandler)
	})

	r.Route("/api/v1/category-plant", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))
		r.Post("/", categoryPlantHandlers.CreateCategoryPlantHandler)
		r.Get("/", categoryPlantHandlers.FindAllCategoryPlantHandler)
		r.Get("/id", categoryPlantHandlers.FindByIdCategoryPlantHandler)
		r.Get("/name", categoryPlantHandlers.FindByNameCategoryPlantHandler)
		r.Put("/", categoryPlantHandlers.UpdateCategoryPlantHandler)
		r.Delete("/", categoryPlantHandlers.DeleteCategoryPlantHandler)
	})

	r.Route("/api/v1/specie", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))
		r.Get("/", specieHandlers.FindAllSpeciesHandler)
		r.Get("/id", specieHandlers.FindByIdSpecieHandler)
		r.Get("/name", specieHandlers.FindByNameSpecieHandler)
	})

	r.Route("/api/v1/plant", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))
		r.Post("/", plantHandlers.CreatePlantHandler)
		r.Delete("/", plantHandlers.DeletePlantHandler)
		r.Get("/", plantHandlers.FindAllPlantHandler)
		r.Get("/category-name", plantHandlers.FindByCategoryNamePlantHandler)
		r.Get("/id", plantHandlers.FindByIdPlantHandler)
		r.Get("/specie-plant-name", plantHandlers.FindBySpecieNamePlantHandler)
		r.Get("/name", plantHandlers.FindByNamePlantHandler)
		r.Put("/", plantHandlers.UpdatePlantHandler)
	})

	return r, nil
}
