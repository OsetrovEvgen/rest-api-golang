package server

import "github.com/lithammer/shortuuid"

// EntityID ...
type EntityID string

const (
	projectID EntityID = "pj_"
	columnID  EntityID = "cl_"
	taskID    EntityID = "tk_"
	commentID EntityID = "cm_"
)

func genID(ent EntityID) *string {
	res := string(ent) + shortuuid.New()
	return &res
}
