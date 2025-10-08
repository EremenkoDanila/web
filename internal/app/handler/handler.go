package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/repository"
	"math/rand"
	"net/http"
	"time"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

func calculateInstallTime(size float32) int {
	rand.Seed(time.Now().UnixNano())
	speed := rand.Intn(10) + 1 // от 1 до 10 МБ/с
	timeMinutes := int((float32(size*1000) / float32(speed)) / 60)
	if timeMinutes < 1 {
		timeMinutes = 1
	}
	return timeMinutes
}

func (h *Handler) GetOrders(ctx *gin.Context) {
	var orders []repository.Order
	var err error

	searchQuery := ctx.Query("apps")
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
			continue
		}


		for i := range versions {
			_ = calculateInstallTime(versions[i].Size)
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

	cartCount := h.Repository.GetCartCount()

	ctx.HTML(http.StatusOK, "index_main.html", gin.H{
		"orders":    ordersWithVersions,
		"apps":      searchQuery,
		"CartCount": cartCount,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	title := ctx.Param("title")
	if title == "" {
		logrus.Error("Название приложения не указано")
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Название приложения не указано",
		})
		return
	}

	orders, err := h.Repository.GetOrdersByTitle(title)
	if err != nil || len(orders) == 0 {
		logrus.Error("Приложение не найдено: ", title)
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Приложение не найдено",
		})
		return
	}

	order := orders[0]

	versions, err := h.Repository.GetVersionsByOrderID(order.ID)
	if err != nil {
		logrus.Error(err)
	}

	cartCount := h.Repository.GetCartCount()

	ctx.HTML(http.StatusOK, "index_prod.html", gin.H{
		"order":     order,
		"versions":  versions,
		"CartCount": cartCount,
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
		FoundVersion  repository.Version
		HasVersion    bool
		SearchQuery   string
		TimeMinutes   int
		PresetVersion string 
	}

	var resultItems []CartItemWithVersion

	currentQueries := make(map[string]string)
	for _, item := range cartOrders {
		queryKey := "apps_" + item.Title
		currentQueries[queryKey] = ctx.Query(queryKey)
	}

	for _, item := range cartOrders {
		queryKey := "apps_" + item.Title
		searchQuery := currentQueries[queryKey]

		cartItem := CartItemWithVersion{
			Order:       item,
			SearchQuery: searchQuery,
			PresetVersion: h.Repository.GetCartDataVersion(item.ID), 
		}

		foundVersion, err := h.Repository.FindVersionByNumber(item.ID, searchQuery)
		if err == nil {
			cartItem.FoundVersion = foundVersion
			cartItem.HasVersion = true
			cartItem.TimeMinutes = calculateInstallTime(foundVersion.Size)
		} else if searchQuery == "" && cartItem.PresetVersion != "" {
			versions, err := h.Repository.GetVersionsByOrderID(item.ID)
			if err == nil && len(versions) > 0 {
				for _, v := range versions {
					if v.Number == cartItem.PresetVersion {
						cartItem.FoundVersion = v
						cartItem.HasVersion = true
						cartItem.TimeMinutes = calculateInstallTime(v.Size)
						break
					}
				}
			}
		}

		resultItems = append(resultItems, cartItem)
	}

	cartCount := h.Repository.GetCartCount()

	ctx.HTML(http.StatusOK, "index_rub.html", gin.H{
		"cartItems": resultItems,
		"CartCount": cartCount,
		"Phone":     h.Repository.GetCartPhone(),
	})
}


func (h *Handler) GetCartCount() int {
	return h.Repository.GetCartCount()
}
