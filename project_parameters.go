package shuffler

import "bitbucket.org/yasnikoff/shuffler/parameters"

func (pr *Project)DefaultParams() *parameters.Collection {
	return pr.Signature().NewCollection()
}

func (pr *Project)Signature() *parameters.Signature {
	return pr.SigData
}

