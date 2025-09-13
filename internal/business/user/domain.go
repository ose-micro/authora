package user

import (
	"fmt"
	"time"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/utils"
	ose_error "github.com/ose-micro/error"
	"github.com/ose-micro/rid"
)

const CreatedEvent string = "events.authora.user_created"
const OnboardedEvent string = "events.authora.user_onboard"
const ChangeStateEvent string = "events.authora.user_change_state"

type Domain struct {
	*domain.Aggregate
	givenNames string
	familyName string
	email      string
	password   string
	metadata   map[string]interface{}
	status     *Status
}

type Params struct {
	Aggregate  *domain.Aggregate
	GivenNames string
	FamilyName string
	Email      string
	Password   string
	Metadata   map[string]interface{}
	Status     *Status
}

type Public struct {
	Id         string                 `json:"_id"`
	GivenNames string                 `json:"given_names"`
	FamilyName string                 `json:"family_name"`
	Email      string                 `json:"email"`
	Password   string                 `json:"password"`
	Metadata   map[string]interface{} `json:"metadata"`
	Version    int32                  `json:"version"`
	Status     *Status                `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	DeletedAt  *time.Time             `json:"deleted_at"`
	Events     []domain.Event         `json:"events"`
}

func (d *Domain) GivenNames() string {
	return d.givenNames
}

func (d *Domain) FamilyName() string {
	return d.familyName
}

func (d *Domain) Password() string {
	return d.password
}

func (d *Domain) Name() string {
	return fmt.Sprintf("%s %s", d.givenNames, d.familyName)
}

func (d *Domain) Email() string {
	return d.email
}

func (d *Domain) Status() *Status {
	return d.status
}

func (d *Domain) Metadata() map[string]interface{} {
	return d.metadata
}

func (d *Domain) Update(params Params) {
	if params.Metadata != nil {
		d.metadata = params.Metadata
		d.Touch()
	}

	if params.GivenNames != "" {
		d.givenNames = params.GivenNames
		d.Touch()
	}

	if params.FamilyName != "" {
		d.familyName = params.FamilyName
		d.Touch()
	}
}

func (d *Domain) ChangePassword(password string, oldPassword string) error {
	if !utils.CheckPasswordHash(oldPassword, d.password) {
		return ose_error.New(ose_error.ErrUnauthorized, "password does not match")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return ose_error.Wrap(err, ose_error.ErrUnauthorized, err.Error())
	}

	d.password = hash
	return nil
}

func (d *Domain) Public() *Public {
	return &Public{
		Id:         d.ID(),
		GivenNames: d.givenNames,
		FamilyName: d.familyName,
		Email:      d.email,
		Metadata:   d.metadata,
		Password:   d.password,
		Status:     d.status,
		Version:    d.Version(),
		CreatedAt:  d.CreatedAt(),
		UpdatedAt:  d.UpdatedAt(),
		DeletedAt:  d.DeletedAt(),
		Events:     d.Events(),
	}
}

func (p Public) Params() *Params {
	id := rid.Existing(p.Id)
	version := p.Version
	createdAt := p.CreatedAt
	updatedAt := p.UpdatedAt
	deletedAt := p.DeletedAt
	events := p.Events

	aggregate := domain.ExistingAggregate(*id, version, createdAt, updatedAt, deletedAt, events)

	return &Params{
		Aggregate:  aggregate,
		GivenNames: p.GivenNames,
		FamilyName: p.FamilyName,
		Email:      p.Email,
		Metadata:   p.Metadata,
		Password:   p.Password,
		Status:     p.Status,
	}
}
