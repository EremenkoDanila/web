package handler

import (
	"net/http"
	"strconv"

	"lab1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

func (h *Handler) GetOrders(ctx *gin.Context) {
	searchQuery := ctx.Query("apps")

	var orders []repository.Order
	var err error
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

	if err != nil {
		logrus.Error(err)
	}

	cartCount := h.Repository.GetCartCount()
	cartID := h.Repository.GetCartID()

	ctx.HTML(http.StatusOK, "index_main.html", gin.H{
		"orders":    orders,
		"apps":      searchQuery,
		"CartCount": cartCount,
		"CartID":    cartID,
	})
}

func (h *Handler) GetOrderByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		logrus.Error("ID приложения не указан")
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "ID приложения не указан",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("Некорректный ID приложения: ", idStr)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Некорректный ID приложения",
		})
		return
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error("Приложение не найдено: ", id)
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Приложение не найдено",
		})
		return
	}

	cartID := h.Repository.GetCartID()

	ctx.HTML(http.StatusOK, "index_prod.html", gin.H{
		"order":     order,
		"CartCount": h.Repository.GetCartCount(),
		"CartID":    cartID,
	})
}

func (h *Handler) GetCartItems(ctx *gin.Context) {
	countParam := ctx.Param("count")
	idParam := ctx.Param("id")

	cartOrders, err := h.Repository.GetCartOrders()
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка получения корзины",
		})
		return
	}

	type CartItem struct {
		repository.Order
		InstallTime string
		FinishTime  string
	}

	var resultItems []CartItem
	for _, order := range cartOrders {
		chosenTime := h.Repository.GetCartInstallTime(order.ID)
		end := ""
		if chosenTime != "" {
			end = repository.CalculateEndTime(chosenTime, order.Size)
		}
		resultItems = append(resultItems, CartItem{
			Order:       order,
			InstallTime: chosenTime,
			FinishTime:  end,
		})
	}

	ctx.HTML(http.StatusOK, "index_rub.html", gin.H{
		"cartItems": resultItems,
		"CartCount": h.Repository.GetCartCount(),
		"Phone":     h.Repository.GetCartPhone(),
		"Count":     countParam,
		"CartID":    idParam,
	})
}
