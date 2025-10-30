package handler

import (
	"fmt"
	"lab1/internal/app/ds"
	"net/http"
	"strconv"
	"time"
	"math/rand"
	"github.com/gin-gonic/gin"
)

// ===== Software =====
func (h *Handler) GetAllSoftware(ctx *gin.Context) {
	apps := ctx.Query("apps")
	var sw []ds.Software
	var err error

	if apps != "" {
		sw, err = h.Repository.SearchSoftwareByTitle(apps)
	} else {
		sw, err = h.Repository.GetAllSoftware()
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req, err := h.Repository.GetDraftRequest(2)
	requestID := 0
	cartCount := 0
	if err == nil {
		requestID = int(req.RequestID)
		cartCount = int(h.Repository.GetCartCount(2))
	}

	ctx.HTML(http.StatusOK, "software.page.tmpl", gin.H{
		"data":       sw,
		"cart_count": cartCount,
		"request_id": requestID,
		"apps":       apps,
	})
}

func (h *Handler) GetSoftwareByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	sw, err := h.Repository.GetSoftwareByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ПО не найдено"})
		return
	}

	req, err := h.Repository.GetDraftRequest(2)
	requestID := 0
	cartCount := 0
	if err == nil {
		requestID = int(req.RequestID)
		cartCount = int(h.Repository.GetCartCount(2))
	}

	ctx.HTML(http.StatusOK, "software_item.page.tmpl", gin.H{
		"order":      sw,
		"cart_count": cartCount,
		"request_id": requestID,
	})
}

// ===== Request / Cart =====
type CartItem struct {
	ds.Software
	Install string
	Final   string
}

func (h *Handler) GetCartByRequestID(ctx *gin.Context) {
	requestID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	req, err := h.Repository.GetRequestByID(uint(requestID), 2)
	if err != nil {
		ctx.HTML(http.StatusOK, "cart.page.tmpl", gin.H{"empty": true})
		return
	}

	installations, err := h.Repository.GetInstallationsByRequestID(uint(requestID))
	if err != nil || len(installations) == 0 {
		ctx.HTML(http.StatusOK, "cart.page.tmpl", gin.H{"empty": true})
		return
	}

	items := []CartItem{}
	for _, inst := range installations {
		install := ""
		final := ""

		if inst.InstallTime != nil {
			install = inst.InstallTime.Format("2006-01-02 15:04")
		}
		if inst.FinalTime != nil {
			final = inst.FinalTime.Format("2006-01-02 15:04")
		}

		items = append(items, CartItem{
			Software: inst.Software,
			Install:  install,
			Final:    final,
		})
	}

	phone := ""
	if req.Phone.Valid {
		phone = req.Phone.String
	}

	ctx.HTML(http.StatusOK, "cart.page.tmpl", gin.H{
		"items":      items,
		"cart_count": len(items),
		"Phone":      phone,
		"Count":      requestID,
		"Status":     req.Status,
	})
}

func (h *Handler) AddSoftwareToRequest(ctx *gin.Context) {
	userID := uint(2)

	softwareID, err := strconv.Atoi(ctx.PostForm("software_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid software id"})
		return
	}

	req, err := h.Repository.GetDraftRequest(userID)
	if err != nil {
		req, err = h.Repository.CreateDraftRequest(userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	count := h.Repository.GetInstallationCount(req.RequestID, uint(softwareID))
	if count > 0 {
		ctx.Redirect(http.StatusFound, "/software")
		return
	}

	err = h.Repository.AddSoftwareToRequest(req.RequestID, uint(softwareID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/software")
}

func (h *Handler) DeleteRequest(ctx *gin.Context) {
	userID := uint(2)
	requestID, err := strconv.Atoi(ctx.PostForm("request_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	err = h.Repository.DeleteRequestByCursor(uint(requestID), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/software")
}

// ===== Update Install Time =====
func (h *Handler) UpdateInstallTime(ctx *gin.Context) {
	userID := uint(2)

	softwareID, err := strconv.Atoi(ctx.PostForm("software_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid software id"})
		return
	}

	installTimeStr := ctx.PostForm("install_time")
	installTime, err := time.Parse("2006-01-02 15:04", installTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
		return
	}

	req, err := h.Repository.GetDraftRequest(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "draft request not found"})
		return
	}

	sw, err := h.Repository.GetSoftwareByID(uint(softwareID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "software not found"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	speed := float64(rand.Intn(10) + 1)

	durationMin := sw.Size / speed * 60
	finalTime := installTime.Add(time.Duration(durationMin) * time.Minute)

	err = h.Repository.UpdateInstallTime(req.RequestID, uint(softwareID), &installTime, &finalTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/software_request/%d", req.RequestID))
}

func (h *Handler) UpdatePhone(ctx *gin.Context) {
	userID := uint(2)

	requestID, err := strconv.Atoi(ctx.PostForm("request_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	phone := ctx.PostForm("phone")

	err = h.Repository.UpdatePhone(uint(requestID), userID, phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/software_request/%d", requestID))
}
