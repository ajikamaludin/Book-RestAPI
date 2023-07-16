package handlers

import (
	"strconv"

	"grace/database"
	"grace/repositories"

	"github.com/gofiber/fiber/v2"
)

type CollectionHandler struct {
	Repository         *repositories.CollectionRepository
	BookRepository     *repositories.BookRepository
	CollectionBookRepo *repositories.CollectionBookRepository
}

// ...

func (h *CollectionHandler) GetCollection(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid collection ID"})
	}

	collection, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
	}

	return c.JSON(collection)
}

func (h *CollectionHandler) CreateCollection(c *fiber.Ctx) error {
	var collection database.Collection
	if err := c.BodyParser(&collection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := h.Repository.Create(&collection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create collection"})
	}

	return c.JSON(collection)
}

func (h *CollectionHandler) UpdateCollection(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid collection ID"})
	}

	collection, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
	}

	if err := c.BodyParser(&collection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := h.Repository.Update(collection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update collection"})
	}

	return c.JSON(collection)
}

func (h *CollectionHandler) DeleteCollection(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid collection ID"})
	}

	collection, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
	}

	if err := h.Repository.Delete(collection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete collection"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CollectionHandler) GetCollectionBooks(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid collection ID"})
	}

	collection, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
	}

	return c.JSON(collection.Books)
}

func (h *CollectionHandler) AddBookToCollection(c *fiber.Ctx) error {
	collectionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid collection ID"})
	}

	var bookIDPayload struct {
		BookID uint `json:"bookId"`
	}
	if err := c.BodyParser(&bookIDPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	book, err := h.BookRepository.FindByID(bookIDPayload.BookID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	collection, err := h.Repository.FindByID(uint(collectionID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
	}

	if err := h.CollectionBookRepo.AddBookToCollection(book, collection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add book to collection"})
	}

	return c.SendStatus(fiber.StatusAccepted)
}

func (h *CollectionHandler) RemoveBookFromCollection(c *fiber.Ctx) error {
	collectionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid collection ID"})
	}

	bookID, err := strconv.Atoi(c.Params("bookId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	collection, err := h.Repository.FindByID(uint(collectionID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
	}

	book, err := h.BookRepository.FindByID(uint(bookID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	if err := h.CollectionBookRepo.RemoveBookFromCollection(book, collection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove book from collection"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
