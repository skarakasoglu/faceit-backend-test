package user

import (
	"context"
)

type GetManyParameters struct {
	Page    int
	PerPage int
	Filter  Entity
}

type Repository interface {
	Create(ctx context.Context, entity Entity) (Entity, error)
	Update(ctx context.Context, entity Entity) (Entity, error)
	DeleteById(ctx context.Context, id string) error
	GetMany(ctx context.Context, parameters GetManyParameters) ([]Entity, error)
}

type service struct {
	repo Repository
}

var _ Service = (*service)(nil)

type ServiceOpts func(*service)

func NewService(opts ...ServiceOpts) Service {
	s := &service{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithRepository(repo Repository) ServiceOpts {
	return func(s *service) {
		s.repo = repo
	}
}

func (s *service) Create(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error) {
	entity, err := s.repo.Create(ctx, Entity{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Nickname:  request.Nickname,
		Password:  request.Password,
		Email:     request.Email,
		Country:   request.Country,
	})
	if err != nil {
		return CreateUserResponse{}, repositoryError(err)
	}

	return CreateUserResponse{
		User{
			Id:        entity.Id,
			FirstName: entity.FirstName,
			LastName:  entity.LastName,
			Nickname:  entity.Nickname,
			Password:  entity.Password,
			Email:     entity.Email,
			Country:   entity.Country,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
	}, nil
}

func (s *service) Update(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error) {
	entity, err := s.repo.Update(ctx, Entity{
		Id:        request.Id,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Nickname:  request.Nickname,
		Password:  request.Password,
		Email:     request.Email,
		Country:   request.Country,
	})
	if err != nil {
		return UpdateUserResponse{}, repositoryError(err)
	}

	return UpdateUserResponse{
		User{
			Id:        entity.Id,
			FirstName: entity.FirstName,
			LastName:  entity.LastName,
			Nickname:  entity.Nickname,
			Password:  entity.Password,
			Email:     entity.Email,
			Country:   entity.Country,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
	}, nil
}

func (s *service) DeleteById(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error) {
	err := s.repo.DeleteById(ctx, request.Id)
	if err != nil {
		return DeleteUserResponse{}, repositoryError(err)
	}

	return DeleteUserResponse{
		Id: request.Id,
	}, nil
}

func (s *service) GetMany(ctx context.Context, request GetUsersManyRequest) (GetUsersManyResponse, error) {
	params := GetManyParameters{
		Page:    request.Page,
		PerPage: request.PerPage,
		Filter: Entity{
			Id:        request.Filter.Id,
			FirstName: request.Filter.FirstName,
			LastName:  request.Filter.LastName,
			Nickname:  request.Filter.Nickname,
			Password:  request.Filter.Password,
			Email:     request.Filter.Email,
			Country:   request.Filter.Country,
			CreatedAt: request.Filter.CreatedAt,
			UpdatedAt: request.Filter.UpdatedAt,
		},
	}

	entities, err := s.repo.GetMany(ctx, params)
	if err != nil {
		return GetUsersManyResponse{}, repositoryError(err)
	}

	users := make([]User, len(entities))
	for i, entity := range entities {
		users[i] = User{
			Id:        entity.Id,
			FirstName: entity.FirstName,
			LastName:  entity.LastName,
			Nickname:  entity.Nickname,
			Password:  entity.Password,
			Email:     entity.Email,
			Country:   entity.Country,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		}
	}

	return GetUsersManyResponse{Users: users}, nil
}
