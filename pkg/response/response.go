package response

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PaginationMeta holds pagination metadata.
type PaginationMeta struct {
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	ItemsPerPage int   `json:"items_per_page"`
}

// APIResponse is the common structure for all API responses.
type APIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
	Errors  interface{}     `json:"errors,omitempty"`
}

// ValidationError defines the structure for a single validation error.
type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// Success sends a standardized success response.
func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a standardized success response with pagination metadata.
func SuccessWithMeta(c *fiber.Ctx, statusCode int, message string, data interface{}, meta *PaginationMeta) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends a standardized error response.
func Error(c *fiber.Ctx, statusCode int, message string, errs ...interface{}) error {
	var errorsData interface{}
	if len(errs) > 0 {
		errorsData = errs
		// If there's only one error and it's a single item, don't wrap it in another array
		if len(errs) == 1 {
			errorsData = errs[0]
		}
	}

	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Message: message,
		Errors:  errorsData,
	})
}

// ValidationErrorResponse handles validation errors from the validator package.
func ValidationErrorResponse(c *fiber.Ctx, err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ValidationError, len(ve))
		for i, fe := range ve {
			out[i] = ValidationError{Field: strings.ToLower(fe.Field()), Error: msgForTag(fe)}
		}
		return Error(c, http.StatusBadRequest, "Validation failed", out)
	}

	return Error(c, http.StatusBadRequest, "Invalid input", err.Error())
}

// msgForTag returns a user-friendly message for a given validation tag.
func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "This field must be at least " + fe.Param() + " characters long"
	case "len":
		return "This field must be exactly " + fe.Param() + " characters long"
	case "oneof":
		return "This field must be one of " + fe.Param()
	case "strong_password":
		return "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character."
	}
	return fe.Error() // default error
}
