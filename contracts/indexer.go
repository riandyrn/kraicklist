package contracts

import "github.com/knightazura/data/model"

/**
/* A contract that can be applied to
/* any model that needed to indexed and searchable.
/*
/* Any model can implement various process which
/* kind of data that need to be indexed.
 */
type Indexer interface {
	ConvertToGeneralDocs() (out model.GeneralDocument, err error)
}
