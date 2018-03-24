package grsdb

import (
	"jcheng/grs/status"
)

type RStat_Json struct {
	Dir    status.Dirstat    `json:"dir"`
	Branch status.Branchstat `json:"branch"`
	Index  status.Indexstat  `json:"index"`
}

func (r *RStat_Json) Update(src status.RStat) {
	r.Branch = src.Branch
	r.Dir = src.Dir
	r.Index = src.Index
}
