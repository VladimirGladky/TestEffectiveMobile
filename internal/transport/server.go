package transport

import (
	_ "TestEffectiveMobile/docs"
	"TestEffectiveMobile/internal/config"
	"TestEffectiveMobile/internal/models"
	"TestEffectiveMobile/internal/service"
	"TestEffectiveMobile/pkg/logger"
	"TestEffectiveMobile/pkg/suberrors"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router.Run(s.cfg.Host + ":" + s.cfg.Port)
}

// @Summary Создаёт новую подписку
// @Tags Подписки
// @Accept json
// @Produce json
// @Param input body models.CreateSubscription true "Данные для создания подписки"
// @Success 200 {object} models.ID "Id созданной подписки"
// @Failure 400 {object} models.BadResponse "Неверный формат запроса"
// @Failure 405 {object} models.BadResponse "Метод не разрешён"
// @Failure 500 {object} models.BadResponse "Внутренняя ошибка сервера"
// @Router /create [post]
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

// @Summary Получает подписку по id
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path string true "ID подписки" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Success 200 {object} models.Subscription "Подписка"
// @Failure 400 {object} models.BadResponse "Неверный формат запроса"
// @Failure 404 {object} models.BadResponse "Подписка не найдена"
// @Failure 405 {object} models.BadResponse "Метод не разрешён"
// @Failure 500 {object} models.BadResponse "Внутренняя ошибка сервера"
// @Router /read/{id} [get]
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
			if errors.Is(err, suberrors.ErrIdSubscriptionNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Subscription id not found"})
				return
			}
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

// @Summary Обновляет подписку по id
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path string true "ID подписки" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Param input body models.UpdateSubscription true "Данные для обновления подписки"
// @Success 200 {object} models.GoodResponse "Подписка обновлена"
// @Failure 400 {object} models.BadResponse "Неверный формат запроса"
// @Failure 404 {object} models.BadResponse "Подписка не найдена"
// @Failure 405 {object} models.BadResponse "Метод не разрешён"
// @Failure 500 {object} models.BadResponse "Внутренняя ошибка сервера"
// @Router /update/{id} [put]
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
			if errors.Is(err, suberrors.ErrIdSubscriptionNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Subscription id not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error3", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.GoodResponse{Message: "Updated"})
	}
}

// @Summary Удаляет подписку по id
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path string true "ID подписки" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Success 200 {object} models.GoodResponse "Подписка удалена"
// @Failure 400 {object} models.BadResponse "Неверный формат запроса"
// @Failure 404 {object} models.BadResponse "Подписка не найдена"
// @Failure 405 {object} models.BadResponse "Метод не разрешён"
// @Failure 500 {object} models.BadResponse "Внутренняя ошибка сервера"
// @Router /delete/{id} [delete]
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
			if errors.Is(err, suberrors.ErrIdSubscriptionNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Subscription id not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
			return
		}
		c.JSON(http.StatusOK, models.GoodResponse{Message: "Deleted"})
	}
}

// @Summary Возвращает список подписок пользователя
// @Tags Подписки
// @Accept json
// @Produce json
// @Param user_id path string true "ID пользователя" example("user12345")
// @Success 200 {object} models.ListSubscriptionsResponse "Список подписок пользователя"
// @Failure 400 {object} models.BadResponse "Неверный формат запроса"
// @Failure 404 {object} models.BadResponse "Пользователь не найден"
// @Failure 405 {object} models.BadResponse "Метод не разрешён"
// @Failure 500 {object} models.BadResponse "Внутренняя ошибка сервера"
// @Router /list/{user_id} [get]
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
			if errors.Is(err, suberrors.ErrUserIdNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User id not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
			return
		}
		c.JSON(http.StatusOK, models.ListSubscriptionsResponse{Subscriptions: subs})
	}
}

// @Summary Возвращает сумму подписок пользователя
// @Tags Подписки
// @Accept json
// @Produce json
// @Param user_id query string true "ID пользователя" example("user12345")
// @Param start_date query string true "Дата начала периода" format(date) example(01-2006)
// @Param end_date query string true "Дата окончания периода" format(date) example(01-2006)
// @Param service_name query string true "Название сервиса" example("YouTube")
// @Success 200 {object} models.SumSubscriptionsResponse "Сумма подписок пользователя"
// @Failure 400 {object} models.BadResponse "Неверный формат запроса"
// @Failure 404 {object} models.BadResponse "Пользователь не найден"
// @Failure 405 {object} models.BadResponse "Метод не разрешён"
// @Failure 500 {object} models.BadResponse "Внутренняя ошибка сервера"
// @Router /sum [get]
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
			if errors.Is(err, suberrors.ErrUserIdNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User id not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.SumSubscriptionsResponse{Sum: sum})
	}
}
