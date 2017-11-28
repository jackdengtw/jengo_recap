package model

import (
	"github.com/qetuantuan/jengo_recap/api"
)

type InnerRun struct {
	api.Run
	Job     string `json:"job"`
	HashId  []byte `bson:"_id"`
	BuildId string `json:"build_id" bson:"build_id"`
}

func (r *InnerRun) ToApiObj() *api.Run {
	return &r.Run
}

func NewInnerRun(r *api.Run) *InnerRun {
	return &InnerRun{
		Run: *r,
	}
}

type InnerRuns []*InnerRun

func (rs InnerRuns) ToApiObj() (ar api.Runs) {
	for _, r := range rs {
		ar = append(ar, *r.ToApiObj())
	}
	return
}

func NewInnerRuns(ar api.Runs) (rs InnerRuns) {
	for _, r := range ar {
		rs = append(rs, NewInnerRun(&r))
	}
	return
}
