package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/google/uuid"
)

type (
	TenantBL interface {
		GetUsers(ctx context.Context, user *model.User) ([]*model.User, error)
		AddUser(ctx context.Context, user *model.User, email, role string) (*model.User, error)
		UpdateUser(ctx context.Context, user *model.User, id uuid.UUID, role string) (*model.User, error)
	}
	tenantBL struct {
		tenantRepo repository.TenantRepository
		userRepo   repository.UserRepository
	}
)

func (b *tenantBL) GetUsers(ctx context.Context, user *model.User) ([]*model.User, error) {
	if len(user.Roles) == 0 || user.Roles[0] == model.RoleUser {
		return nil, utils.ErrorPermission.New("access denied")
	}
	return b.tenantRepo.GetUsers(ctx, user.TenantID)
}

func NewTenantBL(tenantRepo repository.TenantRepository, userRepo repository.UserRepository) TenantBL {
	return &tenantBL{tenantRepo: tenantRepo, userRepo: userRepo}
}

func (b *tenantBL) AddUser(ctx context.Context, user *model.User, email, role string) (*model.User, error) {
	exists, err := b.userRepo.IsUserExists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrorBadRequest.New("user already exists")
	}
	newUser := &model.User{
		ID:         uuid.New(),
		TenantID:   user.TenantID,
		UserName:   email,
		FirstName:  "",
		LastName:   "",
		ExternalID: "",
		Roles:      model.StringSlice{role},
	}
	if err := b.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (b *tenantBL) UpdateUser(ctx context.Context, user *model.User, id uuid.UUID, role string) (*model.User, error) {
	user, err := b.userRepo.GetByIDAndTenantID(ctx, id, user.TenantID)
	if err != nil {
		return nil, err
	}
	user.Roles = model.StringSlice{role}
	if err = b.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
