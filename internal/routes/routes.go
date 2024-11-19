package routes

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/lucasBiazon/botany-back/internal/middleware"
	"github.com/lucasBiazon/botany-back/internal/repositories"
	usecases_categoryplant "github.com/lucasBiazon/botany-back/internal/usecases/category_plant"
	usecases_categorytask "github.com/lucasBiazon/botany-back/internal/usecases/category_task"
	usecasesTasks "github.com/lucasBiazon/botany-back/internal/usecases/tasks"
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

	// Task Routes
	repositoryTask := repositories.NewTaskRepository(db, clientRedis)
	CreateTaskRoutes := usecasesTasks.NewCreateTaskUseCase(repositoryTask)
	FindAllTasksRoutes := usecasesTasks.NewFindAllTaskUseCase(repositoryTask)
	FindAllByStatusTaskRoutes := usecasesTasks.NewFindAllByStatusTaskUseCase(repositoryTask)
	FindAllByCategoryTaskRoutes := usecasesTasks.NewFindAllByCategoryTaskUseCase(repositoryTask)
	FindAllByDateTaskRoutes := usecasesTasks.NewFindAllByDateTaskUseCase(repositoryTask)
	FindAllByNameTaskRoutes := usecasesTasks.NewFindAllByNameTaskUseCase(repositoryTask)
	FindTasksNearDeadlineTaskRoutes := usecasesTasks.NewFindTaskNearDeadLineTaskUseCase(repositoryTask)
	FindTasksFarFromDeadlineTaskRoutes := usecasesTasks.NewFindTasksFarFromDeadlineUseCase(repositoryTask)
	FindTaskByIdTaskRoutes := usecasesTasks.NewFindByIDTaskUseCase(repositoryTask)
	UpdateTaskRoutes := usecasesTasks.NewUpdateTaskUseCase(repositoryTask)
	DeleteTaskRoutes := usecasesTasks.NewDeleteTaskUseCase(repositoryTask)

	taskHandlers := handlers.NewTaskHandlers(
		CreateTaskRoutes,
		DeleteTaskRoutes,
		FindAllByCategoryTaskRoutes,
		FindAllByDateTaskRoutes,
		FindAllByNameTaskRoutes,
		FindAllByStatusTaskRoutes,
		FindAllTasksRoutes,
		FindTaskByIdTaskRoutes,
		UpdateTaskRoutes,
		FindTasksNearDeadlineTaskRoutes,
		FindTasksFarFromDeadlineTaskRoutes,
	)

	// Category Task Routes
	repositoryCategoriesTask := repositories.NewCategoryTaskRepository(db, clientRedis)
	CreateCategoryTaskRoutes := usecases_categorytask.NewCreateCategoryTaskUseCase(repositoryCategoriesTask)
	FindAllCategoryTaskRoutes := usecases_categorytask.NewFindAllCategoryTaskUseCase(repositoryCategoriesTask)
	FindCategoryTaskByIdRoutes := usecases_categorytask.NewFindByIdCategoryTaskUseCase(repositoryCategoriesTask)
	FindByNameCategoryTaskRoutes := usecases_categorytask.NewFindByNameCategoryTaskUseCase(repositoryCategoriesTask)
	UpdateCategoryTaskRoutes := usecases_categorytask.NewUpdateCategoryTaskUseCase(repositoryCategoriesTask)
	DeleteCategoryTaskRoutes := usecases_categorytask.NewDeleteCategoryTaskUseCase(repositoryCategoriesTask)

	categoryTaskHandlers := handlers.NewCategoryTaskHandler(
		CreateCategoryTaskRoutes,
		DeleteCategoryTaskRoutes,
		FindByNameCategoryTaskRoutes,
		UpdateCategoryTaskRoutes,
		FindCategoryTaskByIdRoutes,
		FindAllCategoryTaskRoutes,
	)

	// Category Plant Routes
	repositoryCategoriesPlants := repositories.NewCategoryPlantRepository(db, clientRedis)
	CreateCategoryPlantRoutes := usecases_categoryplant.NewCreateCategoryPlantUseCase(repositoryCategoriesPlants)
	FindAllCategoryPlantRoutes := usecases_categoryplant.NewFindAllCategoryPlantUseCase(repositoryCategoriesPlants)
	FindByCategoryPlantRoutes := usecases_categoryplant.NewFindByIdCategoryPlantUseCase(repositoryCategoriesPlants)
	FindByNameCategoryPlantRoutes := usecases_categoryplant.NewCategoryPlantFindByNameUseCase(repositoryCategoriesPlants)
	UpdateCategoryPlantRoutes := usecases_categoryplant.NewUpdateCategoryPlantUseCase(repositoryCategoriesPlants)
	DeleteCategoryPlantRoutes := usecases_categoryplant.NewDeleteCategoryPlantUseCase(repositoryCategoriesPlants)

	categoryPlantHandlers := handlers.NewCategoryPlantHandler(
		CreateCategoryPlantRoutes,
		FindAllCategoryPlantRoutes,
		FindByCategoryPlantRoutes,
		FindByNameCategoryPlantRoutes,
		UpdateCategoryPlantRoutes,
		DeleteCategoryPlantRoutes,
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
	r.Route("/api/v1/task", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))
		r.Post("/", taskHandlers.CreateTaskHandler)
		r.Get("/", taskHandlers.FindAllTaskHandler)
		r.Get("/status", taskHandlers.FindAllByStatusTaskHandler)
		r.Get("/category", taskHandlers.FindAllByCategoryTaskHandler)
		r.Get("/date", taskHandlers.FindAllByDateTaskHandler)
		r.Get("/name", taskHandlers.FindAllByNameTaskHandler)
		r.Get("/deadline/near", taskHandlers.FindTaskNearDeadLineTaskHandler)
		r.Get("/deadline/far", taskHandlers.FindTasksFarFromDeadlineHandler)
		r.Get("/{id}", taskHandlers.FindByIDTaskHandler)
		r.Put("/{id}", taskHandlers.UpdateTaskHandler)
		r.Delete("/{id}", taskHandlers.DeleteTaskHandler)
	})
	r.Route("/api/v1/category-task", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))
		r.Post("/", categoryTaskHandlers.CreateCategoryTaskHandler)
		r.Get("/", categoryTaskHandlers.FindAllCategoryTaskHandler)
		r.Get("/id", categoryTaskHandlers.FindByIdCategoryTaskHandler)
		r.Get("/name", categoryTaskHandlers.FindByNameCategoryTaskHandler)
		r.Put("/", categoryTaskHandlers.UpdateCategoryTaskHandler)
		r.Delete("/", categoryTaskHandlers.DeleteCategoryTaskHandler)
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

	return r, nil
}
