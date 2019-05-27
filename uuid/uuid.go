package uuid

import "github.com/satori/go.uuid"

type Version int

const (
	V1 Version = 1
	V4 Version = 4
)

// UUID return an uuid based on version.
func UUID(version Version) string {
	switch version {
	case V1:
		return uuid.NewV1().String()
	case V4:
		return uuid.NewV4().String()
	}
	return uuid.NewV4().String()
}
