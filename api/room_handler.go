package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api/db"
	"github.com/samoei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

type BookingRoomParams struct {
	FromDate   time.Time `json:"fromDate,omitempty"`
	TillDate   time.Time `json:"tillDate,omitempty"`
	NumPersons int       `json:"numPersons,omitempty"`
}

func (p BookingRoomParams) validate() error {
	now := time.Now()

	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("you can not book a room in the past")
	}

	return nil

}

func (r RoomHandler) HandleBookingRoom(c *fiber.Ctx) error {
	var params BookingRoomParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	err := params.validate()

	if err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(GenericResponse{
			Type: "error",
			Msg:  "internal server error",
		})
	}
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := r.store.Booking.GetBookings(c.Context(), where)

	if err != nil {
		return err
	}

	if len(bookings) > 0 {
		return c.Status(http.StatusBadRequest).JSON(GenericResponse{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", roomID.String()),
		})
	}

	booking := types.Booking{
		RoomID:     roomID,
		UserID:     user.ID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	inserted, err := r.store.Booking.InsertBookings(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)

}
