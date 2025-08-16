package transport

import (
	"TestEffectiveMobile/internal/config"
	"TestEffectiveMobile/internal/models"
	"TestEffectiveMobile/internal/service"
	"TestEffectiveMobile/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SubscriptionServer struct {
	Service service.SubscriptionServiceInterface
	cfg     *config.Config
	ctx     context.Context
}

func New(srv service.SubscriptionServiceInterface, cfg *config.Config, ctx context.Context) *SubscriptionServer {
	return &SubscriptionServer{
		Service: srv,
		cfg:     cfg,
		ctx:     ctx,
	}
}

func (s *SubscriptionServer) Run() error {
	router := gin.Default()
	logger.GetLoggerFromCtx(s.ctx).Info("gin framework is running")
	api := router.Group("/api/v1")
	{
		api.POST("/create", CreateSubscriptionHandler(s))
		api.GET("/read/:id", ReadSubscriptionHandler(s))
		api.PUT("/update/:id", UpdateSubscriptionHandler(s))
		api.DELETE("/delete/:id", DeleteSubscriptionHandler(s))
		api.GET("/list/:user_id", ListSubscriptionsHandler(s))
		api.GET("/sum", CalculateSumSubscriptionsHandler(s))
	}
	return router.Run(s.cfg.Host + ":" + s.cfg.Port)
}

func CreateSubscriptionHandler(s *SubscriptionServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
				return
			}
		}()
		if c.Request.Method != http.MethodPost {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}
		var request *models.Subscription
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
			return
		}
		id, err := s.Service.Create(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error3"})
			return
		}
		c.JSON(http.StatusOK, models.ID{Id: id})
	}
}

func ReadSubscriptionHandler(s *SubscriptionServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
				return
			}
		}()
		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}
		id := c.Param("id")
		sub, err := s.Service.Read(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
			return
		}
		c.JSON(http.StatusOK, models.Subscription{
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			Id:          id,
			UserId:      sub.UserId,
			StartDate:   sub.StartDate,
			EndDate:     sub.EndDate,
		})
	}
}

func UpdateSubscriptionHandler(s *SubscriptionServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
				return
			}
		}()
		if c.Request.Method != http.MethodPut {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}
		id := c.Param("id")
		var request *models.UpdateSubscription
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
			return
		}
		err := s.Service.Update(id, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error3"})
			return
		}
		c.JSON(http.StatusOK, models.GoodResponse{Message: "Updated"})
	}
}

func DeleteSubscriptionHandler(s *SubscriptionServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
				return
			}
		}()
		if c.Request.Method != http.MethodDelete {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}
		id := c.Param("id")
		err := s.Service.Delete(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
			return
		}
		c.JSON(http.StatusOK, models.GoodResponse{Message: "Deleted"})
	}
}

func ListSubscriptionsHandler(s *SubscriptionServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
				return
			}
		}()
		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}
		userId := c.Param("user_id")
		subs, err := s.Service.ListSubscriptions(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
			return
		}
		c.JSON(http.StatusOK, models.ListSubscriptionsResponse{Subscriptions: subs})
	}
}

func CalculateSumSubscriptionsHandler(s *SubscriptionServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
				return
			}
		}()
		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		userID := c.Query("user_id")
		nameService := c.Query("service_name")
		sum, err := s.Service.CalculateSumSubscriptions(userID, startDate, endDate, nameService)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
			return
		}
		c.JSON(http.StatusOK, models.SumSubscriptionsResponse{Sum: sum})
	}
}
