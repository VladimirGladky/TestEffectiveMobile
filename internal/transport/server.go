package transport

import (
	"TestEffectiveMobile/internal/config"
	"TestEffectiveMobile/internal/models"
	"TestEffectiveMobile/internal/service"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
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
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/create", CreateSubscriptionHandler(s))
	router.HandleFunc("/api/v1/read/{id}", ReadSubscriptionHandler(s))
	router.HandleFunc("/api/v1/update/{id}", UpdateSubscriptionHandler(s))
	router.HandleFunc("/api/v1/delete/{id}", DeleteSubscriptionHandler(s))
	router.HandleFunc("/api/v1/list/{user_id}", ListSubscriptionsHandler(s))
	router.HandleFunc("/api/v1/sum", CalculateSumSubscriptionsHandler(s))
	return http.ListenAndServe(s.cfg.Host+":"+s.cfg.Port, router)
}

func CreateSubscriptionHandler(s *SubscriptionServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "Method not allowed"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		request := new(models.Subscription)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error4"}`))
			}
		}(r.Body)
		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error5"}`))
			return
		}
		id, err := s.Service.Create(request)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error6"}`))
			}
			return
		}
		err = json.NewEncoder(w).Encode(models.ID{Id: id})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error7"}`))
		}
	}
}

func ReadSubscriptionHandler(s *SubscriptionServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "Method not allowed"}`))
			return
		}
		id := mux.Vars(r)["id"]
		sub, err := s.Service.Read(id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error6"}`))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(models.Subscription{
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			UserId:      sub.UserId,
			StartDate:   sub.StartDate,
			EndDate:     sub.EndDate,
			Id:          sub.Id,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error7"}`))
		}
	}
}

func UpdateSubscriptionHandler(s *SubscriptionServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "Method not allowed"}`))
			return
		}
		id := mux.Vars(r)["id"]
		request := new(models.UpdateSubscription)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error4"}`))
			}
		}(r.Body)
		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error5"}`))
			return
		}
		err = s.Service.Update(id, request)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error6"}`))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(models.GoodResponse{Message: "Updated"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error7"}`))
		}
	}
}

func DeleteSubscriptionHandler(s *SubscriptionServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "Method not allowed"}`))
			return
		}
		id := mux.Vars(r)["id"]
		err := s.Service.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error6"}`))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(models.GoodResponse{Message: "Deleted"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error7"}`))
		}
	}
}

func ListSubscriptionsHandler(s *SubscriptionServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "Method not allowed"}`))
			return
		}
		userId := mux.Vars(r)["user_id"]
		subs, err := s.Service.ListSubscriptions(userId)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error6"}`))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(models.ListSubscriptionsResponse{Subscriptions: subs})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error7"}`))
		}
	}
}

func CalculateSumSubscriptionsHandler(s *SubscriptionServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "Method not allowed"}`))
			return
		}
		query := r.URL.Query()
		startDate := query.Get("start_date")
		endDate := query.Get("end_date")
		userID := query.Get("user_id")
		nameService := query.Get("service_name")
		sum, err := s.Service.CalculateSumSubscriptions(userID, startDate, endDate, nameService)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error6"}`))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(models.SumSubscriptionsResponse{Sum: sum})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error7"}`))
		}
	}
}
