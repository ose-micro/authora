package user

import (
	"fmt"
	"time"
)

type State int32

const (
	StateUnspecified         State = iota
	StateInvited                   // User invited, not yet registered
	StatePendingVerification       // Registered but email/phone not confirmed
	StateInactive                  // Exists but not enabled (admin or system hold)
	StateActive                    // Fully verified and allowed to use the system
	StateDormant                   // Inactive for a long period, auto-marked dormant
	StateSuspended                 // Temporarily blocked by admin or system rule
	StateLocked                    // Auto-locked (e.g. too many failed logins)
	StateBanned                    // Permanently blocked for policy violation
	StateDeactivated               // User voluntarily deactivated account
	StateArchived                  // Retained for records, no login allowed
	StateDeleted                   // Permanently removed from system
)

type Status struct {
	State    State     `json:"state" bson:"state"`
	Previous *State    `json:"previous" bson:"previous"`
	OccurOn  time.Time `json:"occur_on" bson:"occur_on"`
}

var allowedTransitions = map[State][]State{
	StateInvited:             {StatePendingVerification, StateDeleted},
	StatePendingVerification: {StateActive, StateInactive, StateDeleted},
	StateActive:              {StateSuspended, StateDormant, StateDeactivated, StateDeleted},
	StateSuspended:           {StateActive, StateBanned, StateDeleted},
	StateDormant:             {StateActive, StateDeleted},
	StateDeactivated:         {StateActive, StateDeleted},
}

func (s *Status) IsActive() bool {
	return s.State == StateActive
}

func (s *Status) IsPendingVerification() bool {
	return s.State == StatePendingVerification
}

func (s *Status) IsInvited() bool {
	return s.State == StateInvited
}

func (s *Status) IsDormant() bool {
	return s.State == StateDormant
}

func (s *Status) IsLocked() bool {
	return s.State == StateLocked
}

func (s *Status) IsBanned() bool {
	return s.State == StateBanned
}

func (s *Status) IsArchived() bool {
	return s.State == StateArchived
}

func (s *Status) IsSuspended() bool {
	return s.State == StateSuspended
}

func (s *Status) IsDeactivated() bool {
	return s.State == StateDeactivated
}

func (s *Status) IsDeleted() bool {
	return s.State == StateDeleted
}

func (s *Status) IsInactive() bool {
	return s.State == StateInactive
}

func Active() Status {
	return NewStatus(StateActive)
}

func Invited() Status {
	return NewStatus(StateInvited)
}

func Dormant() Status {
	return NewStatus(StateDormant)
}

func Locked() Status {
	return NewStatus(StateLocked)
}

func Banned() Status {
	return NewStatus(StateBanned)
}

func Archived() Status {
	return NewStatus(StateArchived)
}

func Suspended() Status {
	return NewStatus(StateSuspended)
}

func Deactivated() Status {
	return NewStatus(StateDeactivated)
}

func Deleted() Status {
	return NewStatus(StateDeleted)
}

func Inactive() Status {
	return NewStatus(StateInactive)
}

func (s State) String() string {
	switch s {
	case StateInvited:
		return "Invited"
	case StatePendingVerification:
		return "PendingVerification"
	case StateInactive:
		return "Inactive"
	case StateActive:
		return "Active"
	case StateDormant:
		return "Dormant"
	case StateSuspended:
		return "Suspended"
	case StateLocked:
		return "Locked"
	case StateBanned:
		return "Banned"
	case StateDeactivated:
		return "Deactivated"
	case StateArchived:
		return "Archived"
	case StateDeleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}

func (s *Status) ChangeState(next State) error {
	if s.State == next {
		return fmt.Errorf("status is already %v", next)
	}

	allowed := false
	for _, st := range allowedTransitions[s.State] {
		if st == next {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("transition from %v to %v not allowed", s.State, next)
	}

	prev := s.State
	s.Previous = &prev
	s.State = next
	s.OccurOn = time.Now()
	return nil
}

func NewStatus(state State) Status {
	return Status{State: state, Previous: nil, OccurOn: time.Now()}
}
