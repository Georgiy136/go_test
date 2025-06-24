package postgres

import (
	"bytes"
	"strconv"
)

var (
	strCallProcedure = "CALL "
	strSelectAll     = "SELECT * FROM  "
	parenthesesLeft  = "("
	parenthesesRight = ")"
	specSymbol       = "$"
	strComma         = ","
)

type PgSpec struct {
	storedProcedure string
	params          *[]interface{}
	databaseKey     string
	useProcedure    bool
}

// SetParams - устанавливает параметры для вызова ХП.
func (s *PgSpec) SetParams(args ...interface{}) {
	s.params = &args
}

// SetStoredProcedure - устанавливает имя хранимой процедуры.
func (s *PgSpec) SetStoredProcedure(spName string) {
	s.storedProcedure = spName
}

func (s *PgSpec) SetUseProcedure() {
	s.useProcedure = true
}

func (s *PgSpec) SetUseFunction() {
	s.useProcedure = false
}

// GetStoredProcedureName - возвращает имя установленной хранимой процедуры.
func (s *PgSpec) GetStoredProcedureName() string { return s.storedProcedure }

// GetQuery - формирует вызов хранимой процедуры.
func (s *PgSpec) GetQuery() string {
	var buf bytes.Buffer
	isFirstParameter := true

	if s.useProcedure {
		addString(&buf, &strCallProcedure)
	} else {
		addString(&buf, &strSelectAll)
	}
	addString(&buf, &s.storedProcedure)
	addString(&buf, &parenthesesLeft)
	if s.params == nil {
		addString(&buf, &parenthesesRight)
		return buf.String()
	}

	for i := range *s.params {
		rnk := strconv.Itoa(i + 1)
		if !isFirstParameter {
			addString(&buf, &strComma)
		} else {
			isFirstParameter = false
		}

		addString(&buf, &specSymbol)
		addString(&buf, &rnk)
	}
	addString(&buf, &parenthesesRight)
	return buf.String()
}

func (s *PgSpec) GetParameters() []interface{} {
	if s.params == nil {
		return nil
	}

	return *s.params
}

func addString(buf *bytes.Buffer, str *string) {
	if _, err := buf.WriteString(*str); err != nil {
		panic(err)
	}
}
