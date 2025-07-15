package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Govorov1705/ozon-test/config"
	"github.com/Govorov1705/ozon-test/graph"
	"github.com/Govorov1705/ozon-test/internal/logger"
	"github.com/Govorov1705/ozon-test/internal/middleware"
	"github.com/Govorov1705/ozon-test/internal/repositories"
	"github.com/Govorov1705/ozon-test/internal/services"
	"github.com/Govorov1705/ozon-test/internal/storages/inmemory"
	inmemRepos "github.com/Govorov1705/ozon-test/internal/storages/inmemory/repositories"
	"github.com/Govorov1705/ozon-test/internal/storages/postgresql"
	psqlRepos "github.com/Govorov1705/ozon-test/internal/storages/postgresql/repositories"
	"github.com/Govorov1705/ozon-test/internal/transactions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
)

func graphqlHandler(
	usersService *services.UsersService,
	postsService *services.PostsService,
	commentsService *services.CommentsService,
	allowedOrigins []string,
) gin.HandlerFunc {
	h := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(
		usersService,
		postsService,
		commentsService,
	)}))

	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" || origin == r.Header.Get("Host") {
					return true
				}
				return slices.Contains(allowedOrigins, origin)
			},
		},
	})
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	config.InitConfig()
	logger.InitLogger()

	var (
		txStarter    transactions.TxStarter
		usersRepo    repositories.UsersRepository
		postsRepo    repositories.PostsRepository
		commentsRepo repositories.CommentsRepository
	)

	switch config.Cfg.Storage {
	case config.StorageInmemory:
		logger.Logger.Info("Using in-memory storage")

		txStarter = &inmemory.InMemoryTxStarter{}

		usersRepo = inmemRepos.NewUsersRepository()
		postsRepo = inmemRepos.NewPostsRepository()
		commentsRepo = inmemRepos.NewCommentsRepository()
	case config.StoragePostgreSQL:
		logger.Logger.Info("Using PostgreSQL as a storage")

		storage := postgresql.NewStorage(config.Cfg.DBURL)
		txStarter = postgresql.NewPgxpoolTxStarter(storage.Pool)

		usersRepo = psqlRepos.NewUsersRepository(storage.Pool)
		postsRepo = psqlRepos.NewPostsRepository(storage.Pool)
		commentsRepo = psqlRepos.NewCommentsRepository(storage.Pool)
	default:
		logger.Logger.Fatal("Unsupported storage backend")
	}

	usersService := services.NewUsersService(usersRepo)
	postsService := services.NewPostsService(txStarter, postsRepo, commentsRepo)
	commentsService := services.NewCommentsService(txStarter, commentsRepo, postsRepo)

	if config.Cfg.Mode == config.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(middleware.Auth)
	r.Any("/query", graphqlHandler(
		usersService,
		postsService,
		commentsService,
		config.Cfg.AllowedOrigins,
	))
	r.GET("/", playgroundHandler())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}
	go func() {
		logger.Logger.Info("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("Error starting HTTP server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Logger.Info("Gracefully shutting down HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Error("HTTP server shutdown:", zap.Error(err))
	}
}
