package types

type SignRequestStatus uint8

const (
	SignRequestStatusPlaceholder SignRequestStatus = iota
	SignRequestStatusCreated
	SignRequestStatusPrepared
	SignRequestStatusSigned
	SignRequestStatusFailed
)

const (
	SignRequestStatusCreatedName  = "created"
	SignRequestStatusPreparedName = "prepared"
	SignRequestStatusSignedName   = "signed"
	SignRequestStatusFailedName   = "failed"
)

func (d SignRequestStatus) String() string {
	switch d {
	case SignRequestStatusCreated:
		return SignRequestStatusCreatedName
	case SignRequestStatusPrepared:
		return SignRequestStatusPreparedName
	case SignRequestStatusSigned:
		return SignRequestStatusSignedName
	case SignRequestStatusFailed:
		return SignRequestStatusFailedName

	case SignRequestStatusPlaceholder:
		fallthrough
	default:
		return ""
	}
}
