package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/thegeorgenikhil/hackernews-go-graphql/graph/generated"
	"github.com/thegeorgenikhil/hackernews-go-graphql/graph/model"
	"github.com/thegeorgenikhil/hackernews-go-graphql/internal/auth"
	"github.com/thegeorgenikhil/hackernews-go-graphql/internal/links"
	"github.com/thegeorgenikhil/hackernews-go-graphql/internal/users"
	"github.com/thegeorgenikhil/hackernews-go-graphql/pkg/jwt"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("access denied")
	}
	//check if links.Title is empty or link.Address is empty
	if input.Title == "" || input.Address == "" {
		return &model.Link{}, fmt.Errorf("Title or Address cannot be empty")
	}

	var link links.Link
	link.Title = input.Title
	
	// checking whether the link has https:// in it
	if !strings.HasPrefix(input.Address, "https://") {
		link.Address = "https://" + input.Address
	} else {
		link.Address = input.Address
	}

	link.User = user
	linkID := link.Save()
	graphqlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}

	return &model.Link{ID: strconv.FormatInt(linkID, 10), Title: link.Title, Address: link.Address, User: graphqlUser}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password

	if user.Username == "" || user.Password == "" {
		return "", fmt.Errorf("Username or Password cannot be empty")
	}

	// check if user already exists
	userExists, _ := users.GetUserIdByUsername(user.Username)

	// if err is nil, that means we have got the user id, duplicate user
	if userExists == 0 {
		// if user doesn't exist, create a new user
		user.Create()
		token, err := jwt.GenerateToken(user.Username)
		if err != nil {
			return "", err
		}
		return token, nil
	}
	return "", fmt.Errorf("user already exists")
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()

	if !correct {
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, err
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	var resultLinks []*model.Link
	var dbLinks []links.Link = links.GetAll()
	for _, link := range dbLinks {
		graphqlUser := &model.User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		resultLinks = append(resultLinks, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address, User: graphqlUser})
	}
	return resultLinks, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
