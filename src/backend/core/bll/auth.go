package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/google/uuid"
	"time"
)

type (
	AuthBL interface {
		Login(ctx context.Context, userName string) (*model.User, error)
		SignUp(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error)
		//Invite(ctx context.Context, identity *security.Identity, param *parameters.InviteParam) (string, error)
		//JoinToTenant(ctx context.Context, state *parameters.OAuthParam, response *oauth.IdentityResponse) (*model.User, error)
		QuickLogin(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error)
	}
	authBL struct {
		userRepo repository.UserRepository
		cfg      *Config
	}
)

func NewAuthBL(userRepo repository.UserRepository,

	cfg *Config) AuthBL {
	return &authBL{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (a *authBL) Login(ctx context.Context, userName string) (*model.User, error) {
	user, err := a.userRepo.GetByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *authBL) SignUp(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error) {
	exists, err := a.userRepo.IsUserExists(ctx, identity.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrorBadRequest.New("user already exists")
	}
	userID := uuid.New()
	tenantID := uuid.New()

	// create user  and default connector and embedding model
	user := model.User{
		ID:         userID,
		TenantID:   tenantID,
		UserName:   identity.Email,
		FirstName:  identity.GivenName,
		LastName:   identity.FamilyName,
		ExternalID: identity.ID,
		Roles:      model.StringSlice{model.RoleSuperAdmin},
		Defaults: &model.Defaults{
			EmbeddingModel: &model.EmbeddingModel{
				TenantID:     tenantID,
				ModelID:      a.cfg.DefaultEmbeddingModel,
				ModelName:    a.cfg.DefaultEmbeddingModel,
				ModelDim:     a.cfg.DefaultEmbeddingVectorSize,
				IsActive:     true,
				CreationDate: time.Now().UTC(),
			},
		},
	}
	if user.FirstName == "" {
		user.FirstName = identity.Name
	}
	if err = a.userRepo.RegisterUser(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

//func (a *authBL) Invite(ctx context.Context, identity *security.Identity, param *parameters.InviteParam) (string, error) {
//
//	exists, err := a.userRepo.IsUserExists(ctx, param.Email)
//	if err != nil {
//		return "", err
//	}
//	if exists {
//		return "", utils.ErrorBadRequest.New("user already registered.")
//	}
//	//buf, err := json.Marshal(parameters.OAuthParam{Action: oauth.InviteState,
//	//	Role:     param.Role,
//	//	Email:    param.Email,
//	//	TenantID: identity.User.TenantID.String(),
//	//})
//	//if err != nil {
//	//	return "", utils.Internal.Wrap(err, "can not marshal payload")
//	//}
//	key := uuid.New()
//	//if err = a.storage.Save(key.String(), buf); err != nil {
//	//	return "", err
//	//}
//	state := base64.URLEncoding.EncodeToString([]byte(key.String()))
//
//	return fmt.Sprintf("%s/auth/google/invite?state=%s", a.redirectURL, state), nil
//
//}
//
//func (a *authBL) JoinToTenant(ctx context.Context, state *parameters.OAuthParam, response *oauth.IdentityResponse) (*model.User, error) {
//	//if state.Email != response.Email {
//	//	return nil, utils.ErrorPermission.New("email is not equals to invite")
//	//}
//	user := model.User{
//		ID:         uuid.New(),
//		TenantID:   uuid.MustParse(state.TenantID),
//		UserName:   state.Email,
//		FirstName:  response.GivenName,
//		LastName:   response.FamilyName,
//		ExternalID: response.ID,
//		Roles:      model.StringSlice{state.Role},
//	}
//	if err := a.userRepo.Create(ctx, &user); err != nil {
//		return nil, err
//	}
//	return &user, nil
//}

func (a *authBL) QuickLogin(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error) {
	exists, err := a.userRepo.IsUserExists(ctx, identity.Email)
	if err != nil {
		return nil, err
	}
	if !exists {
		return a.SignUp(ctx, identity)
	}
	user, err := a.userRepo.GetByUserName(ctx, identity.Email)
	if err != nil {
		return nil, err
	}
	if user.ExternalID == "" {
		if identity.GivenName == "" {
			user.FirstName = identity.Name
		} else {
			user.FirstName = identity.GivenName
		}
		user.LastName = identity.FamilyName
		user.ExternalID = identity.ID
	}
	if err = a.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
