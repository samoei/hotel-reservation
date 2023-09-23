package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

type HotelQueryParams struct {
	Rooms bool
	Stars int
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	fmt.Println(qparams)

	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)

	if err != nil {
		return err
	}

	return c.JSON(&hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	var (
		hotelID = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)

	if err != nil {
		return err
	}

	return c.JSON(&hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	var (
		hotelID = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelID": oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filter)

	if err != nil {
		return err
	}

	return c.JSON(&rooms)

}
