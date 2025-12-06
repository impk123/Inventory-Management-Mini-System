package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/impk123/Inventory-Management-Mini-System/internal/models"
	"github.com/impk123/Inventory-Management-Mini-System/pkg/cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// StockHandler holds dependencies for stock handling
type StockHandler struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewStockHandler(db *gorm.DB, cache *redis.Client) *StockHandler {
	return &StockHandler{
		db:    db,
		cache: cache,
	}
}

// CreateStockMovement creates a new stock movement
// POST /api/stock
// CreateStockMovement godoc
// @Summary Create a stock movement
// @Description Record a new stock movement (IN, OUT, ADJUST)
// @Tags stocks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param movement body models.StockUpdateRequest true "Stock movement details"
// @Success 200 {object} MessageResponse "Stock updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /stocks [post]
func (h *StockHandler) CreateStockMovement(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req models.StockUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get product info
	var product models.Product
	if err := h.db.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Calculate new quantity
	oldQty := product.Quantity
	newQty := oldQty

	switch req.Type {
	case "IN":
		newQty = oldQty + req.Quantity
	case "OUT":
		if oldQty < req.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Insufficient stock",
				"available": oldQty,
				"requested": req.Quantity,
			})
			return
		}
		newQty = oldQty - req.Quantity
	case "ADJUST":
		newQty = req.Quantity
	}

	// Start transaction
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// Update product quantity
		product.Quantity = newQty
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// Create stock movement record
		movement := models.Stock{
			ProductID:   req.ProductID,
			Type:        req.Type,
			Quantity:    req.Quantity,
			OldQuantity: oldQty,
			NewQuantity: newQty,
			Notes:       req.Notes,
			CreatedBy:   userID.(uint),
			CreatedAt:   time.Now(),
		}

		return tx.Create(&movement).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock movement"})
		return
	}

	// Invalidate Caches
	// 1. Stock Lists (History & Current Stock) --> เพราะมี Movement ใหม่
	cache.InvalidateStockList(h.cache, c.Request.Context())

	// 2. Product Lists --> เพราะ Quantity เปลี่ยน
	cache.InvalidateProductList(h.cache, c.Request.Context())

	// 3. Specific Product --> เพราะ Quantity เปลี่ยน
	productCacheKey := cache.GenerateProductKey(req.ProductID)
	h.cache.Del(c.Request.Context(), productCacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})

}

// GetStockHistory gets stock history for a product
// GET /api/stock/product/:id
// GetStockHistory godoc
// @Summary Get stock history
// @Description Get stock movement history for a specific product
// @Tags stocks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} StockHistoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stocks/product/{id} [get]
func (h *StockHandler) GetStockHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	// 1. Try Cache
	// ใช้ GenerateStockHistoryKey แทน เพื่อให้ Key มี ID ของสินค้า และไม่ชนกับ GetCurrentStock
	cacheKey := cache.GenerateStockHistoryKey(h.cache, c.Request.Context(), id, c.Request.URL.Query())
	type HistoryResponse struct {
		Stock     []models.Stock `json:"stock"`
		Movements []models.Stock `json:"movements"`
	}
	var cachedResp HistoryResponse
	if err := cache.GetCached(h.cache, c.Request.Context(), cacheKey, &cachedResp); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"stock":     cachedResp.Stock,
			"movements": cachedResp.Movements,
		})
		return
	}

	var stock []models.Stock
	if err := h.db.Where("product_id = ?", id).Find(&stock).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	var movements []models.Stock
	h.db.Where("product_id = ?", id).
		Order("created_at DESC").
		Limit(100).
		Find(&movements)

	response := HistoryResponse{
		Stock:     stock,
		Movements: movements,
	}

	// 2. Set Cache
	cache.SetCached(h.cache, c.Request.Context(), cacheKey, response, 30*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"stock":     stock,
		"movements": movements,
	})
}

// GetCurrentStock gets current stock for all products
// GET /api/stock
// GetCurrentStock godoc
// @Summary Get current stock
// @Description Get current stock levels for all products
// @Tags stocks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} CurrentStockResponse
// @Router /stocks [get]
func (h *StockHandler) GetCurrentStock(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	// 1. Try Cache
	cacheKey := cache.GenerateStockListKey(h.cache, c.Request.Context(), c.Request.URL.Query())
	type CurrentStockResponse struct {
		Products      []models.Product `json:"products"`
		Total         int64            `json:"total"`
		Page          int              `json:"page"`
		Limit         int              `json:"limit"`
		LowStock      []models.Product `json:"low_stock"`
		LowStockCount int              `json:"low_stock_count"`
	}
	var cachedResp CurrentStockResponse
	if err := cache.GetCached(h.cache, c.Request.Context(), cacheKey, &cachedResp); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"products":        cachedResp.Products,
			"total":           cachedResp.Total,
			"page":            cachedResp.Page,
			"limit":           cachedResp.Limit,
			"low_stock":       cachedResp.LowStock,
			"low_stock_count": cachedResp.LowStockCount,
		})
		return
	}

	var products []models.Product
	var total int64

	// Get paginated products
	h.db.Model(&models.Product{}).Count(&total)
	h.db.Offset(offset).Limit(limit).Find(&products)

	// Get low stock items
	var lowStock []models.Product
	h.db.Where("quantity <= min_quantity").Find(&lowStock)

	response := CurrentStockResponse{
		Products:      products,
		Total:         total,
		Page:          page,
		Limit:         limit,
		LowStock:      lowStock,
		LowStockCount: len(lowStock),
	}

	// 2. Set Cache
	cache.SetCached(h.cache, c.Request.Context(), cacheKey, response, 30*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"products":        products,
		"total":           total,
		"page":            page,
		"limit":           limit,
		"low_stock":       lowStock,
		"low_stock_count": len(lowStock),
	})
}
