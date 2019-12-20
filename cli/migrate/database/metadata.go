package database

type MetadataDriver interface {
	ExportMetadata() (interface{}, error)

	ResetMetadata() error

	ReloadMetadata() error

	GetInconsistentMetadata() (interface{}, error)

	DropInconsistentMetadata() error

	ApplyMetadata(data interface{}) error

	Query(data []interface{}) error
}
