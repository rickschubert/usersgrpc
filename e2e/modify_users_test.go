package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/rickschubert/usersgrpc/users"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

// This test adds a new user, modifies it, and then gets it again using the same
// ID to check that the modifications happened.
func TestModifyUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addedUser, err := client.AddUser(ctx, &users.AddUserInput{
		FirstName: ptr.To(gofakeit.FirstName()),
		LastName:  ptr.To(gofakeit.LastName()),
		Nickname:  ptr.To(gofakeit.Username()),
		Email:     ptr.To(gofakeit.Email()),
		Country:   ptr.To(gofakeit.Country()),
	})
	require.NoError(t, err)
	require.NotEmpty(t, addedUser.GetId())

	// Properties to modify
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	nickname := gofakeit.Username()
	email := gofakeit.Email()
	country := gofakeit.Country()

	timeBeforeUpdate := time.Now()

	modifiedUser, err := client.ModifyUser(ctx, &users.ModifyUserInput{
		Id:        ptr.To(addedUser.GetId()),
		FirstName: ptr.To(firstName),
		LastName:  ptr.To(lastName),
		Nickname:  ptr.To(nickname),
		Email:     ptr.To(email),
		Country:   ptr.To(country),
	})

	require.NoError(t, err)
	require.Equal(t, modifiedUser.GetId(), addedUser.GetId())
	require.Equal(t, modifiedUser.GetFirstName(), firstName)
	require.Equal(t, modifiedUser.GetLastName(), lastName)
	require.Equal(t, modifiedUser.GetNickname(), nickname)
	require.Equal(t, modifiedUser.GetEmail(), email)
	require.Equal(t, modifiedUser.GetCountry(), country)

	// Check that the new updatedAt value is less than 1 second ago
	updatedAt, err := time.Parse(time.RFC3339, modifiedUser.GetUpdatedAt())
	require.NoError(t, err)
	diff := updatedAt.Sub(timeBeforeUpdate)
	require.Less(t, diff.Abs().Milliseconds(), int64(1000))
	require.Greater(t, diff.Abs().Milliseconds(), int64(0))

	gotUser, err := client.GetUser(ctx, &users.GetUserInput{
		Id: ptr.To(addedUser.GetId()),
	})
	require.NoError(t, err)
	require.Equal(t, gotUser.GetId(), addedUser.GetId())
	require.Equal(t, gotUser.GetFirstName(), firstName)
	require.Equal(t, gotUser.GetLastName(), lastName)
	require.Equal(t, gotUser.GetNickname(), nickname)
	require.Equal(t, gotUser.GetEmail(), email)
	require.Equal(t, gotUser.GetCountry(), country)
}
