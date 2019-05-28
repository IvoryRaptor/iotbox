package dataframe

import "github.com/IvoryRaptor/iotbox/streams"

type DataFrame struct {
	streams.Flow
	Schema *Schema
}

type InsertRow struct {
	Row Row
}

type RemoveRow struct {
	Row Row
}

type UpdateRow struct {
	Row Row
}
