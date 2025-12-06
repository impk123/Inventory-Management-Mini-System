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

// ProductHandler holds dependencies for product handling
type ProductHandler struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewProductHandler(db *gorm.DB, cache *redis.Client) *ProductHandler {
	return &ProductHandler{
		db:    db,
		cache: cache,
	}
}

// CreateProduct creates a new product
type CreateProductRequest struct {
	SKU         string  `json:"sku" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
	CostPrice   float64 `json:"cost_price"`
	Quantity    int     `json:"quantity" binding:"min=0"`
	MinQuantity int     `json:"min_quantity" binding:"min=0"`
	MaxQuantity int     `json:"max_quantity" binding:"min=0"`
	Location    string  `json:"location"`
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the given details
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body CreateProductRequest true "Product details"
// @Success 200 {object} MessageResponse "Product created successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if SKU exists
	var existingProduct models.Product
	if err := h.db.Where("sku = ?", req.SKU).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product SKU already exists"})
		return
	}

	product := models.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		UnitPrice:   req.UnitPrice,
		CostPrice:   req.CostPrice,
		Quantity:    req.Quantity,
		MinQuantity: req.MinQuantity,
		MaxQuantity: req.MaxQuantity,
		Location:    req.Location,
		IsActive:    true,
		CreatedBy:   userID.(uint),
		UpdatedBy:   userID.(uint),
	}

	if err := h.db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Invalidate List Cache
	cache.InvalidateProductList(h.cache, c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})

}

// GetProducts all products

// GetProducts godoc
// @Summary Get all products
// @Description Get a list of all products with pagination and filtering
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param search query string false "Search term"
// @Param category query string false "Category filter"
// @Success 200 {object} ProductListResponse
// @Router /products [get]

func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	category := c.Query("category")
	offset := (page - 1) * limit

	// 1. Try Cache
	cacheKey := cache.GenerateProductListKey(h.cache, c.Request.Context(), c.Request.URL.Query())
	var cachedResp ProductListResponse
	if err := cache.GetCached(h.cache, c.Request.Context(), cacheKey, &cachedResp); err == nil {
		c.JSON(http.StatusOK, cachedResp)
		return
	}

	var products []models.Product
	var total int64

	query := h.db.Model(&models.Product{}).Where("is_active = ?", true)

	if search != "" {
		query = query.Where("sku LIKE ? OR name LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&products)

	response := ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		Limit:    limit,
		Pages:    (int(total) + limit - 1) / limit,
	}

	// 2. Set Cache (30 Minutes)
	// We cache the FULL response, so next time it's instant return
	cache.SetCached(h.cache, c.Request.Context(), cacheKey, response, 30*time.Minute)

	c.JSON(http.StatusOK, response)

}

// GetProduct by id
// GetProductByID godoc
// @Summary Get a product by ID
// @Description Get detailed information of a specific product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// หาใน cache ก่อน
	cacheKey := cache.GenerateProductKey(uint(id))
	var product models.Product
	if err := cache.GetCached(h.cache, c.Request.Context(), cacheKey, &product); err == nil {
		c.JSON(http.StatusOK, product)
		return
	}

	// ถ้าไม่เจอ หาใน database
	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Cache the result
	cache.SetCached(h.cache, c.Request.Context(), cacheKey, product, 5*time.Minute) // 5 minutes

	c.JSON(http.StatusOK, product)
}

// Update Product
// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param product body CreateProductRequest true "Product details"
// @Success 200 {object} MessageResponse "Product updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Category = req.Category
	product.UnitPrice = req.UnitPrice
	product.CostPrice = req.CostPrice
	product.MinQuantity = req.MinQuantity
	product.MaxQuantity = req.MaxQuantity
	product.Location = req.Location
	product.UpdatedBy = userID.(uint)

	if err := h.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Invalidate List Cache
	cache.InvalidateProductList(h.cache, c.Request.Context())
	cacheKey := cache.GenerateProductKey(uint(id))
	h.cache.Del(c.Request.Context(), cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// Delete Product
// DeleteProduct godoc
// @Summary Delete a product
// @Description Soft delete a product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} MessageResponse "Product deleted successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Soft delete
	product.IsActive = false
	if err := h.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	// Invalidate List Cache
	cache.InvalidateProductList(h.cache, c.Request.Context())
	cacheKey := cache.GenerateProductKey(uint(id))
	h.cache.Del(c.Request.Context(), cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})

}
