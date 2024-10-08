package api

import (
	"e-commerce/internal/middleware"
	"e-commerce/internal/models"
	"e-commerce/internal/util"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Create User
func (u *HTTPHandler) CreateUser(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBind(&user); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	user.Password = hashedPassword

	err = u.Repository.CreateUser(user)
	if err != nil {
		util.Response(c, "User not created", 500, err.Error(), nil)
		return
	}
	util.Response(c, "User created", 200, nil, nil)

}

// Login User
func (u *HTTPHandler) LoginUser(c *gin.Context) {
	var loginRequest *models.LoginRequestUser
	err := c.ShouldBind(&loginRequest)
	if err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.Response(c, "Email and Password must not be empty", 400, nil, nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "Email does not exist", 404, err.Error(), nil)
		return
	}

	// Verify the password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		util.Response(c, "invalid email or password", 400, "invalid email or password", nil)
		return
	}

	accessClaims, refreshClaims := middleware.GenerateClaims(user.Email)

	secret := os.Getenv("JWT_SECRET")

	accessToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating access token", 500, err.Error(), nil)
		return
	}

	refreshToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating refresh token", 500, err.Error(), nil)
		return
	}

	c.Header("access_token", *accessToken)
	c.Header("refresh_token", *refreshToken)

	util.Response(c, "Login successful", 200, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// View Product listing
func (u *HTTPHandler) GetAllProducts(c *gin.Context) {
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	products, err := u.Repository.GetAllProducts()
	if err != nil {
		util.Response(c, "Error getting products", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Success", 200, products, nil)
}

// View Product by ID
func (u *HTTPHandler) GetProductByID(c *gin.Context) {
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	productID := c.Param("id")
	id, err := strconv.Atoi(productID)
	if err != nil {
		util.Response(c, "Error getting product", 500, err.Error(), nil)
		return
	}
	product, err := u.Repository.GetProductByID(uint(id))
	if err != nil {
		util.Response(c, "Error getting product", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Success", 200, product, nil)
}

// View cart
func (u *HTTPHandler) ViewCart(c *gin.Context) {
	// Get user ID from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	// Get cart items by user ID
	cartItems, err := u.Repository.GetCartsByUserID(user.ID)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}

	// If cart is empty, return early
	if len(cartItems) == 0 {
		util.Response(c, "Your cart is empty", 404, nil, nil)
		return
	}

	// Prepare the response structure
	var cartTotal models.CartTotal
	cartTotal.Cart = make([]*models.CartItem, len(cartItems))

	// Calculate the total price and prepare the cart items
	var total float64
	for i, cartItem := range cartItems {
		product, err := u.Repository.GetProductByID(cartItem.ProductID)
		if err != nil {
			util.Response(c, "Error fetching product details", 500, err.Error(), nil)
			return
		}
		cartTotal.Cart[i] = &models.CartItem{
			CartID:   cartItem.ID,
			Product:  product,
			Quantity: cartItem.Quantity,
		}
		total += float64(cartItem.Quantity) * product.Price
	}

	// Set the total price
	cartTotal.Total = total

	// Return the cart items and total price
	util.Response(c, "Cart fetched successfully", 200, gin.H{
		"cart":  cartTotal.Cart,
		"total": cartTotal.Total,
	}, nil)
}

func (u *HTTPHandler) AddToCart(c *gin.Context) {
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	var cart *models.IndividualItemInCart
	err = c.ShouldBind(&cart)
	if err != nil {
		util.Response(c, "invalid request", 401, err.Error(), nil)
		return
	}

	product, err := u.Repository.GetProductByID(cart.ProductID)
	if err != nil {
		util.Response(c, "product not found", 500, err.Error(), nil)
		return
	}

	if cart.Quantity > product.Quantity {
		util.Response(c, "product out of stock", 400, nil, nil)
		return
	}

	cart.UserID = user.ID

	err = u.Repository.AddToCart(cart)
	if err != nil {
		util.Response(c, "error adding product to cart", 500, err.Error(), nil)
		return
	}

	util.Response(c, "product added to cart", 200, product, nil)
}

func (u *HTTPHandler) EditCart(c *gin.Context) {
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	var cart *models.IndividualItemInCart
	err = c.ShouldBind(&cart)
	if err != nil {
		util.Response(c, "invalid request", 401, err.Error(), nil)
		return
	}

	shoppingCart, err := u.Repository.GetCartItemByProductID(cart.ProductID)
	if err != nil {
		util.Response(c, "Cart not found", 404, err.Error(), nil)
		return
	}

	product, err := u.Repository.GetProductByID(cart.ProductID)
	if err != nil {
		util.Response(c, "product not found", 500, err.Error(), nil)
		return
	}

	if cart.Quantity > product.Quantity {
		util.Response(c, "not enough products", 400, nil, nil)
		return

	}
	cart.UserID = user.ID
	cart.ID = shoppingCart.ID

	err = u.Repository.AddToCart(cart)
	if err != nil {
		util.Response(c, "Error editing product quantity.", 500, err.Error(), nil)
		return
	}

	util.Response(c, "product successfully added", 200, nil, nil)
}

// delete product from cart
func (u *HTTPHandler) DeleteProductFromCart(c *gin.Context) {
	// Get user id from context
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	// Get product by id
	productID := c.Param("id")

	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		util.Response(c, "Invalid product ID", 400, err.Error(), nil)
		return
	}

	// Validate request
	shoppingCart, err := u.Repository.GetCartItemByProductID(uint(productIDInt))
	if err != nil {
		util.Response(c, "Product not found", 404, err.Error(), nil)
		return
	}

	err = u.Repository.DeleteProductFromCart(shoppingCart)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Product deleted from cart", 200, nil, nil)
}

func (u *HTTPHandler) OrderHistory(c *gin.Context) {
	// check user is logged in
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error finding user in context", 500, err.Error(), nil)
		return

	}

	// go into order table to check user id to find order
	orders, err := u.Repository.OrderHistorybyUserID(user.ID)
	if err != nil {
		util.Response(c, "Order not found", 404, err.Error(), nil)
		return
	}
	
	var orderDetails []models.Order
	for _, order := range orders {
		orderItems, err := u.Repository.GetOrderItemsByOrderID(order.ID)
		if err != nil {
			util.Response(c, "Internal server error", 500, err.Error(), nil)
			return
		}
		order.Items = orderItems
		orderDetails = append(orderDetails, *order)
	}

	util.Response(c, "Orders fetched", 200, gin.H{
		"orders": orderDetails,
	}, nil)

}

func (u *HTTPHandler) PlaceOrder(c *gin.Context) {
	// Get user ID from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	// Get products in cart
	cartItems, err := u.Repository.GetCartsByUserID(user.ID)
	if err != nil {
		util.Response(c, "Error fetching cart items", 500, err.Error(), nil)
		return
	}

	if len(cartItems) == 0 {
		util.Response(c, "Cart is empty", 400, "No items in the cart", nil)
		return
	}

	// Calculate total and prepare order items
	var total float64
	var orderItems []*models.OrderItem
	for _, cartItem := range cartItems {
		product, err := u.Repository.GetProductByID(cartItem.ProductID)
		if err != nil {
			util.Response(c, "Error fetching product details", 500, err.Error(), nil)
			return
		}

		// Check if the product is out of stock
		if cartItem.Quantity > product.Quantity {
			util.Response(c, "Product out of stock", 400, "Product is out of stock", nil)
			return
		}

		// Calculate total price
		total += float64(cartItem.Quantity) * product.Price

		// Prepare order item
		orderItems = append(orderItems, &models.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
		})
	}

	// Prepare the order
	order := &models.Order{
		UserID: user.ID,
		Total:  total,
		Status: "PLACED",
		Items:  orderItems,
	}

	// Save the order and clear the cart within a transaction
	err = u.Repository.CreateOrder(order)
	if err != nil {
		util.Response(c, "Error creating order", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Order placed successfully", 200, nil, nil)
}
