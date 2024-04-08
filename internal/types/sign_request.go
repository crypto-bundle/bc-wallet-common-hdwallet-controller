package types

type SignRequestStatus uint8

const (
	SignRequestStatusPlaceholder SignRequestStatus = iota
	SignRequestStatusCreated
	SignRequestStatusPrepared
	SignRequestStatusSigned
	SignRequestStatusClosed
)

const (
	SignRequestStatusCreatedName  = "created"
	SignRequestStatusPreparedName = "prepared"
	SignRequestStatusSignedName   = "signed"
	SignRequestStatusClosedName   = "closed"
)

func (d SignRequestStatus) String() string {
	switch d {
	case SignRequestStatusCreated:
		return SignRequestStatusCreatedName
	case SignRequestStatusPrepared:
		return SignRequestStatusPreparedName
	case SignRequestStatusSigned:
		return SignRequestStatusSignedName
	case SignRequestStatusClosed:
		return SignRequestStatusClosedName

	case SignRequestStatusPlaceholder:
		fallthrough
	default:
		return ""
	}
}
