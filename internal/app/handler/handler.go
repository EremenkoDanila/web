package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/repository"
	"net/http"
	"strconv"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetOrders(ctx *gin.Context) {
	var orders []repository.Order
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		orders, err = h.Repository.GetOrders()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		orders, err = h.Repository.GetOrdersByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	// Получаем версии для каждого заказа
	type OrderWithVersions struct {
		Order    repository.Order
		Versions []repository.Version
		Latest   repository.Version
	}

	var ordersWithVersions []OrderWithVersions
	for _, order := range orders {
		versions, err := h.Repository.GetVersionsByOrderID(order.ID)
		if err != nil {
			logrus.Error(err)
			versions = []repository.Version{}
		}
		var latest repository.Version
		if len(versions) > 0 {
			latest = order.GetLatestVersion(versions)
		}
		ordersWithVersions = append(ordersWithVersions, OrderWithVersions{
			Order:    order,
			Versions: versions,
			Latest:   latest,
		})
	}

	ctx.HTML(http.StatusOK, "index_main.html", gin.H{
		"orders": ordersWithVersions,
		"query":  searchQuery,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	versions, err := h.Repository.GetVersionsByOrderID(order.ID)
	if err != nil {
		logrus.Error(err)
		versions = []repository.Version{}
	}

	ctx.HTML(http.StatusOK, "index_prod.html", gin.H{
		"order":    order,
		"versions": versions,
	})
}

func (h *Handler) GetCartItems(ctx *gin.Context) {
	cartOrders, err := h.Repository.GetCartOrders()
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка получения данных корзины",
		})
		return
	}

	type CartItemWithVersion struct {
		repository.Order
		FoundVersion repository.Version
		HasVersion   bool
		SearchQuery  string
	}

	var resultItems []CartItemWithVersion

	// Собираем поисковые запросы для всех карточек
	currentQueries := make(map[string]string)
	for _, item := range cartOrders {
		queryKey := "query_" + strconv.Itoa(item.ID)
		currentQueries[queryKey] = ctx.Query(queryKey)
	}

	for _, item := range cartOrders {
		queryKey := "query_" + strconv.Itoa(item.ID)
		searchQuery := currentQueries[queryKey]

		cartItem := CartItemWithVersion{
			Order:       item,
			SearchQuery: searchQuery,
		}

		// Ищем версию по запросу
		foundVersion, err := h.Repository.FindVersionByNumber(item.ID, searchQuery)
		if err == nil {
			cartItem.FoundVersion = foundVersion
			cartItem.HasVersion = true
		} else if searchQuery == "" {
			// Если запрос пустой - показываем последнюю версию
			versions, err := h.Repository.GetVersionsByOrderID(item.ID)
			if err != nil {
				logrus.Error(err)
			}
			if len(versions) > 0 {
				cartItem.FoundVersion = item.GetLatestVersion(versions)
				cartItem.HasVersion = true
			} else {
				cartItem.HasVersion = false
			}
		} else {
			// Если запрос не пустой, но версия не найдена
			cartItem.HasVersion = false
		}

		resultItems = append(resultItems, cartItem)
	}

	ctx.HTML(http.StatusOK, "index_rub.html", gin.H{
		"cartItems": resultItems,
	})
}