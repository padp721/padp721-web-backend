package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/padp721/padp721-web-backend/models"
	"github.com/padp721/padp721-web-backend/utilities"
	"golang.org/x/crypto/bcrypt"
)

func UserCreate(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	var newUser models.UserRegister
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: err.Error(),
		})
	}

	//* BEGIN TRANSACTION
	dbTx, err := db.Begin(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}
	defer dbTx.Rollback(c.Context())

	id := uuid.New()

	passwordString := utilities.GeneratePasswordString(newUser.Password, newUser.Username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	sql := "INSERT INTO public.users(id, username, name, email, phone) VALUES($1, $2, $3, $4, $5)"
	_, err = dbTx.Exec(c.Context(), sql, id, newUser.Username, newUser.Name, newUser.Email, newUser.Phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	sql = "INSERT INTO secret.users(user_id, password) VALUES($1, $2)"
	_, err = dbTx.Exec(c.Context(), sql, id, string(hashedPassword))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	//* COMMIT TRANSACTION
	err = dbTx.Commit(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Message: "User Created.",
	})
}

func UsersGet(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	sql := "SELECT id, username, name, email, phone FROM public.users"
	rows, err := db.Query(c.Context(), sql)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Name, &user.Email, &user.Phone); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
				Message: err.Error(),
			})
		}
		users = append(users, user)
	}

	return c.JSON(models.ResponseData{
		Message: "Data fetch success!",
		Data:    users,
	})
}

func UserGet(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	id := c.Params("id")

	var user models.User
	sql := "SELECT id, username, name, email, phone FROM public.users WHERE id=$1"
	err := db.QueryRow(c.Context(), sql, id).Scan(&user.Id, &user.Username, &user.Name, &user.Email, &user.Phone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.Response{
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(models.ResponseData{
		Message: "Data fetch success!",
		Data:    user,
	})
}

func UserUpdate(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	var updatedUser models.UserUpdate
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: err.Error(),
		})
	}

	id := c.Params("id")
	sql := "UPDATE public.users SET username=$1, name=$2, email=$3, phone=$4 WHERE id=$5"
	commandTag, err := db.Exec(c.Context(), sql, updatedUser.Username, updatedUser.Name, updatedUser.Email, updatedUser.Phone, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: "User not found!",
		})
	}

	return c.JSON(models.Response{
		Message: "User Updated!",
	})
}

func UserDelete(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	//* BEGIN TRANSACTION
	dbTx, err := db.Begin(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}
	defer dbTx.Rollback(c.Context())

	var userCount int
	sql := "SELECT COUNT(id) FROM public.users"
	err = dbTx.QueryRow(c.Context(), sql).Scan(&userCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	if userCount <= 1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: "You can't delete more user. This is the last user available.",
		})
	}

	id := c.Params("id")

	sql = "DELETE FROM secret.users WHERE user_id=$1"
	commandTag, err := dbTx.Exec(c.Context(), sql, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Message: "User not found!",
		})
	}

	sql = "DELETE FROM public.users WHERE id=$1"
	commandTag, err = dbTx.Exec(c.Context(), sql, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Message: "User not found!",
		})
	}

	//* COMMIT TRANSACTION
	err = dbTx.Commit(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(models.Response{
		Message: "User Deleted.",
	})
}

func UserChangePassword(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	var inputPassword models.UserChangePassword
	if err := c.BodyParser(&inputPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: fmt.Sprintf("Error parsing request body: %v", err),
		})
	}

	//* CHECK IF BOTH NEW PASSWORD FIELD IS EQUAL
	if inputPassword.NewPassword != inputPassword.ReNewPassword {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: "Make sure new password in both inputs are equal!",
		})
	}

	userId := c.Locals("userId").(string)

	var (
		username                string
		hashedOldPasswordString string
	)

	sql := `
		SELECT a.username, b.password
		FROM public.users AS a
		INNER JOIN secret.users AS b ON a.id = b.user_id
		WHERE a.id=$1
	`
	err := db.QueryRow(c.Context(), sql, userId).Scan(&username, &hashedOldPasswordString)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.Response{
				Message: "User tidak ditemukan!",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: fmt.Sprintf("Error getting user data: %v", err),
		})
	}

	newPasswordString := utilities.GeneratePasswordString(inputPassword.NewPassword, username)
	oldPasswordString := utilities.GeneratePasswordString(inputPassword.OldPassword, username)

	//* CHECK IF OLD PASSWORD INPUT IS SAME AS OLD PASSWORD DB
	err = bcrypt.CompareHashAndPassword([]byte(hashedOldPasswordString), []byte(oldPasswordString))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: "Old password is wrong.",
		})
	}

	//* CHECK IF NEW PASSWORD IS SAME AS OLD PASSWORD
	err = bcrypt.CompareHashAndPassword([]byte(hashedOldPasswordString), []byte(newPasswordString))
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: "New password is the same as old password!",
		})
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPasswordString), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	sql = "UPDATE secret.users SET password=$1 WHERE user_id=$2"
	commandTag, err := db.Exec(c.Context(), sql, string(hashedNewPassword), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: fmt.Sprintf("Error updating password to database: %v", err),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: "User not found!",
		})
	}

	return c.JSON(models.Response{
		Message: "Password changed. You need to re-login to app.",
	})
}
