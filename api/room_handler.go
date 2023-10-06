package api

import (
	"context"
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

func (p BookingRoomParams) validateDates() error {
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

	if err := params.validateDates(); err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	//Get user from the context
	user, ok := c.Context().Value("user").(*types.User)

	// For some reason we could not fetch the user from the context
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(GenericResponse{
			Type: "error",
			Msg:  "internal server error",
		})
	}

	//Check if there are any other bookings associated with this room during the same period
	ok, err = r.roomCanBeBooked(c, roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(GenericResponse{
			Type: "error",
			Msg:  fmt.Sprintf("room with id %s is already booked", c.Params("id")),
		})
	}

	//prepare to store this booking in the database
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

	//update Room that its booked
	err = r.updateRoomBooked(c.Context(), roomID)

	if err != nil {
		fmt.Printf("Could not set room as booked %v", err)
	}

	return c.JSON(inserted)

}

func (r RoomHandler) updateRoomBooked(ctx context.Context, roomID primitive.ObjectID) error {
	filter := bson.M{"_id": roomID} // Specify the filter to match the document you want to update.
	update := bson.M{"$set": bson.M{"booked": true}}

	if err := r.store.Room.UpdateRoom(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (r RoomHandler) roomCanBeBooked(c *fiber.Ctx, roomID primitive.ObjectID, params BookingRoomParams) (bool, error) {
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
		return false, err
	}

	return len(bookings) == 0, nil
}
