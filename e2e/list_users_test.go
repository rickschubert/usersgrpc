package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/rickschubert/usersgrpc/config"
	"github.com/rickschubert/usersgrpc/users"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

// Works based on the assumption that there are 3 users in the database with a
// country of "Deathstar". The results will be split over 2 pages of maximum 2
// results.
func TestListUsersWithFilterByCountry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deathStarUsersPageOne, err := client.ListUsers(ctx, &users.ListUsersInput{
		Country: ptr.To("Deathstar"),
	})
	require.NoError(t, err)
	require.Equal(t, config.PaginationSize(), len(deathStarUsersPageOne.Users))
	require.NotEmpty(t, deathStarUsersPageOne.GetNextPageToken())

	deathStarUsersPageTwo, err := client.ListUsers(ctx, &users.ListUsersInput{
		Country:       ptr.To("Deathstar"),
		NextPageToken: ptr.To(deathStarUsersPageOne.GetNextPageToken()),
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(deathStarUsersPageTwo.Users))
	require.Empty(t, deathStarUsersPageTwo.GetNextPageToken())
}

// // Works based on the assumption that there are 4 users in the database in
// // total. The results will be split over 2 pages with 2 items each; and the last
// // page would be empty.
func TestListUsersWithoutFilter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deathStarUsersPageOne, err := client.ListUsers(ctx, &users.ListUsersInput{})
	require.NoError(t, err)
	require.Equal(t, config.PaginationSize(), len(deathStarUsersPageOne.Users))
	require.NotEmpty(t, deathStarUsersPageOne.GetNextPageToken())

	deathStarUsersPageTwo, err := client.ListUsers(ctx, &users.ListUsersInput{
		NextPageToken: ptr.To(deathStarUsersPageOne.GetNextPageToken()),
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(deathStarUsersPageTwo.Users))
	require.NotEmpty(t, deathStarUsersPageTwo.GetNextPageToken())

	deathStarUsersPageThree, err := client.ListUsers(ctx, &users.ListUsersInput{
		NextPageToken: ptr.To(deathStarUsersPageTwo.GetNextPageToken()),
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(deathStarUsersPageThree.Users))
	require.Empty(t, deathStarUsersPageThree.GetNextPageToken())
}
