package authdatamanager

import "github.com/lfq7413/tomato/types"

type facebook struct{}

func (a facebook) ValidateAuthData(authData types.M) error {
	return nil
}
