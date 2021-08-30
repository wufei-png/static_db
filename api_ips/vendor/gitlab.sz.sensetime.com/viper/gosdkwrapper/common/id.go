package common

//go:generate go run ./kestrel_id_gen/main.go
//go:generate goimports -w ../common

type ObjectID uint64
type ObjectMainID uint64

// nolint: golint
const UNKNOWN_OBJECTID ObjectID = 0

var (
	string2ObjectIDMap map[string]ObjectID
)

func (o ObjectID) String() string {
	return objectID2StringMap[o]
}

func ObjectIDFromString(s string) ObjectID {
	return string2ObjectIDMap[s]
}

func GetObjectMainID(t ObjectID) ObjectMainID {
	return objectMainIDMap[t]
}

const (
	ObjectMainIDUnknown ObjectMainID = iota
	ObjectMainIDFacial
	ObjectMainIDBody
	ObjectMainIDTrafficParticipant
	ObjectMainIDTrafficInfrastructure
	ObjectMainIDClothing
	ObjectMainIDCostumeAccessory
)

var objectMainIDMap = map[ObjectID]ObjectMainID{}

func init() {
	string2ObjectIDMap = make(map[string]ObjectID)
	for k, v := range objectID2StringMap {
		string2ObjectIDMap[v] = k
	}
}
