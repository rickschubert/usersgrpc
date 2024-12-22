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

// This test adds a new user, gets it to assert that it exists, removes it, and
// gets it again to assert that the user cannot be retrieved any longer.
func TestRemoveUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addedUser, err := client.AddUser(ctx, &users.AddUserInput{
		FirstName: ptr.To(gofakeit.FirstName()),
		LastName:  ptr.To(gofakeit.LastName()),
	})
	require.NoError(t, err)
	require.NotEmpty(t, addedUser.GetId())

	userAfterAdding, err := client.GetUser(ctx, &users.GetUserInput{
		Id: ptr.To(addedUser.GetId()),
	})
	require.NoError(t, err)
	require.NotEmpty(t, userAfterAdding.GetId())

	removedUser, err := client.RemoveUser(ctx, &users.RemoveUserInput{
		Id: ptr.To(addedUser.GetId()),
	})
	require.NoError(t, err)
	require.Equal(t, addedUser.GetId(), removedUser.GetId())

	userAfterRemoving, err := client.GetUser(ctx, &users.GetUserInput{
		Id: ptr.To(addedUser.GetId()),
	})
	require.NoError(t, err)
	require.Empty(t, userAfterRemoving.GetId())
}
