package pkg

import "time"

type Timeframe struct {
	Since time.Time
	Until time.Time
}

type Algorithm uint16

const (
	None Algorithm = iota
	SHA256
	SHA512
	SHA3_256
	SHA3_512
	Base64
)

func (a Algorithm) String() string {
	return [...]string{"None", "SHA256", "SHA512", "SHA3_256", "SHA3_512", "Base64"}[a]
}

type TokenType uint16

const (
	AccessToken TokenType = iota
	RefreshToken
	PasswordResetToken
	EmailVerificationToken
)

func (t TokenType) Name() string {
	return [...]string{"AccessToken", "RefreshToken", "PasswordResetToken", "EmailVerificationToken"}[t]
}

type FriendRelationType uint8

const (
	Contact FriendRelationType = iota
	Family
	Colleague
	Acquaintance
	Nakama
	Custom1
	Custom2
	Blocked
)

func (f FriendRelationType) String() string {
	return [...]string{"Contact", "Family", "Colleague", "Acquaintance", "Friend", "Custom1", "Custom2", "Blocked"}[f]
}

type SocialRequestStatus uint8

const (
	Pending SocialRequestStatus = iota
	Accepted
	Declined
	Invited
	Left
)

func (f SocialRequestStatus) String() string {
	return [...]string{"Pending", "Accepted", "Declined", "Invited", "Left"}[f]
}
