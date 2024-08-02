package storage

var mariaDriver = "mysql"

type MariaDBErrors string

const (
	MariaDBErrorsBeginTx          = MariaDBErrors("Begin Tx")
	MariaDBErrorsCommitTx         = MariaDBErrors("Commit Tx")
	MariaDBErrorsPrepareStatement = MariaDBErrors("Prepare Statement")
	MariaDBErrorsExecStatement    = MariaDBErrors("Execute Statement")
	MariaDBErrorsExecQuery        = MariaDBErrors("Execute Query")
	MariaDBErrorsQueryRow         = MariaDBErrors("Query Row")
	MariaDBErrorsQuery            = MariaDBErrors("Query Results")
	MariaDBErrorsScanResult       = MariaDBErrors("Scanning Query Rows")
)
