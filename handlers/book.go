package handlers

import (
	"fmt"
	"strconv"
	"time"

	"grace/database"
	"grace/repositories"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	Repository *repositories.BookRepository
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var book database.Book

	if err := c.BodyParser(&book); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := h.Repository.Create(&book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create book"})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	book, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := h.Repository.Update(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update book"})
	}

	return c.JSON(book)
}

func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	book, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	if err := h.Repository.Delete(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete book"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *BookHandler) GetBooks(c *fiber.Ctx) (err error) {
	author := c.Query("author")
	genre := c.Query("genre")
	var start time.Time
	var end time.Time

	rstart, rend := c.Query("start"), c.Query("end")
	if rstart != "" {
		start, err = time.Parse("2006-01-02", c.Query("start"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start date"})
		}
	}

	if rend != "" {
		end, err = time.Parse("2006-01-02", c.Query("end"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end date"})
		}
	}

	books, err := h.Repository.Find(author, genre, start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get books"})
	}

	return c.JSON(books)
}

func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	book, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	return c.JSON(book)
}
