package types

type ErrorType int

const (
	NoError         ErrorType = 1
	ConversionError           = 1001
	DatabaseError             = 4001
	CacheError                = 4002

	JsonParsingError      = 5001
	UriParsingError       = 5001
	ClaimsExtractingError = 5002
	FileExtractingError   = 5003

	FileTypeError   = 6001
	FileSavingError = 6002
)