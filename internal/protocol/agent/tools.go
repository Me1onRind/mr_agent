package agent

type ExcuteSQLInput struct {
	DBLabel string `json:"db_label"`
	SQL     string `json:"sql"`
}

type ExcuteSQLOutput struct {
	RowsAffected int64            `json:"rows_affected,omitempty"`
	Rows         []map[string]any `json:"rows,omitempty"`
	Error        string           `json:"error,omitempty"`
}
