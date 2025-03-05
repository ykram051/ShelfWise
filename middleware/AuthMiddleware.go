package middleware

import (
	"FinalProject/services"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket algorithm
type RateLimiter struct {
	tokens     map[string][]time.Time
	windowSize time.Duration
	maxTokens  int
	mu         sync.Mutex
}

func NewRateLimiter(windowSize time.Duration, maxTokens int) *RateLimiter {
	return &RateLimiter{
		tokens:     make(map[string][]time.Time),
		windowSize: windowSize,
		maxTokens:  maxTokens,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if _, exists := rl.tokens[key]; !exists {
		rl.tokens[key] = []time.Time{now}
		return true
	}

	// Remove expired tokens
	var validTokens []time.Time
	for _, t := range rl.tokens[key] {
		if now.Sub(t) <= rl.windowSize {
			validTokens = append(validTokens, t)
		}
	}

	if len(validTokens) >= rl.maxTokens {
		return false
	}

	rl.tokens[key] = append(validTokens, now)
	return true
}

// AuthMiddleware struct
type AuthMiddleware struct {
	AuthService  *services.AuthService
	OrderService *services.OrderService
	rateLimiter  *RateLimiter
}

// NewAuthMiddleware initializes middleware
func NewAuthMiddleware(authService *services.AuthService, orderService *services.OrderService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService:  authService,
		OrderService: orderService,
		rateLimiter:  NewRateLimiter(time.Minute, 60),
	}
}

// JWTAuthMiddleware ensures the request has a valid JWT token
func (m *AuthMiddleware) JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if !m.rateLimiter.Allow(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]
		claims, err := m.AuthService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract user info from token claims
		userID := int(claims["user_id"].(float64))
		role := claims["role"].(string)

		// Set user info in request headers
		r.Header.Set("X-User-ID", strconv.Itoa(userID))
		r.Header.Set("X-User-Role", role)

		// Check role-based access
		if !m.checkRoleAccess(r, userID, role) {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// checkRoleAccess enforces role-based permissions
func (m *AuthMiddleware) checkRoleAccess(r *http.Request, userID int, role string) bool {
	path := r.URL.Path
	method := r.Method

	// Public access (both customers and admins)
	if (strings.HasPrefix(path, "/api/books") && method == http.MethodGet) ||
		(strings.HasPrefix(path, "/api/authors") && method == http.MethodGet) {
		return true
	}

	// Admin-only actions
	if (strings.HasPrefix(path, "/api/books") || strings.HasPrefix(path, "/api/authors")) &&
		(method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete) {
		return role == "admin"
	}

	if strings.HasPrefix(path, "/api/customers") && !strings.HasPrefix(path, "/api/customers/"){
		return role == "admin"
	}

	// Customer profile access: Customers can only view/edit their own profiles
	if strings.HasPrefix(path, "/api/customers/") {
		if role == "admin" {
			return true // Admins can access any customer data
		}

		// Extract customer ID from URL
		customerIDStr := strings.TrimPrefix(path, "/api/customers/")
		customerID, err := strconv.Atoi(customerIDStr)
		if err != nil {
			return false // Invalid customer ID
		}

		log.Println("Authenticated User ID:", userID)
		log.Println("Requested Customer ID:", customerID)

		// Customers can only access their own profile
		return userID == customerID
	}

	// ✅ Fix: Allow customers to access only their own orders
	if strings.HasPrefix(path, "/api/orders") {
		if role == "admin" {
			return true // Admins can access all orders
		}

		// If it's a general list request (GET /api/orders), allow customers to see their own orders
		if path == "/api/orders" || path == "/api/orders/" {
			return true
		}

		// Extract order ID from URL
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 3 {
			return false // Invalid path (e.g., missing order ID)
		}

		orderID, err := strconv.Atoi(pathParts[3]) // Extract numeric order ID
		if err != nil {
			log.Println("Invalid order ID format:", pathParts[3])
			return false // Invalid order ID format
		}

		log.Println("Checking access for Order ID:", orderID, "User ID:", userID)

		// ✅ Ensure the order belongs to the customer
		if !m.isUserOrder(r.Context(), userID, orderID) {
			log.Println("Order does not belong to user:", userID)
			return false
		}

		return true // Allow access to the order if it belongs to the customer
	}

	return false // Default deny
}

// isUserOrder checks if the order belongs to the customer
func (m *AuthMiddleware) isUserOrder(ctx context.Context, userID int, orderID int) bool {
	order, err := m.OrderService.GetOrder(ctx, orderID)
	if err != nil {
		log.Println("Order not found or error retrieving order:", err)
		return false
	}

	log.Println("Order ID:", orderID, "Requested by User ID:", userID, "Actual Owner ID:", order.UserID)

	return order.UserID == userID // ✅ Only allow access if the order belongs to the customer
}
