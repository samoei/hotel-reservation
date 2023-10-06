package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (b *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := b.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (b *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := b.store.Booking.GetBookingByID(c.Context(), c.Params(("id")))
	if err != nil {
		return err
	}
	return c.JSON(booking)
}
