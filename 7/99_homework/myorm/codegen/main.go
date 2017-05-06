package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Table struct {
	name    string
	columns []*Column
	pk      *Column
}

type Column struct {
	name      string
	fieldType string
	isNull    bool
	nameInDB  string
}

func WriteStringInFile(w *os.File, s string) {
	_, err := w.Write([]byte(s))
	PanicOnErr(err)
}

func GenerateHeaders(w *os.File) {
	s := "package user\n\n"
	WriteStringInFile(w, s)
	WriteStringInFile(w, "// Generated code!\n\n")
	s = "import \"database/sql\"\n\n"
	WriteStringInFile(w, s)
}

func GenerateDBInfoFunction(w *os.File) {
	WriteStringInFile(w, "var DB *sql.DB\n\n")
	s := "func SetDB(db *sql.DB) {\n"
	s += "\tDB = db\n}\n\n"
	WriteStringInFile(w, s)
}

func GenerateFindByPKFunction(w *os.File, table *Table, structName string, tableName string) {
	s := "func (data *" + structName + ") FindByPK(pk " + table.pk.fieldType + ") (err error) {\n"
	WriteStringInFile(w, s)
	for _, column := range table.columns {
		if column.isNull {
			s = "\tvar " + column.nameInDB + "_null *" + column.fieldType + "\n"
			WriteStringInFile(w, s)
		}
	}
	s = "\trow := DB.QueryRow(\"SELECT "
	for i, column := range table.columns {
		if i != 0 {
			s += ", "
		}
		s += column.nameInDB
	}
	s += " FROM " + tableName + " WHERE " + table.pk.nameInDB + "= ?\", pk)\n"
	WriteStringInFile(w, s)

	s = "\terr = row.Scan("
	for i, column := range table.columns {
		if i != 0 {
			s += ", "
		}
		if column.isNull {
			s += "&" + column.nameInDB + "_null"
		} else {
			s += "&data." + column.name
		}
	}
	s += ")\n"
	WriteStringInFile(w, s)

	for _, column := range table.columns {
		if column.isNull {
			s = "\tif " + column.nameInDB + "_null != nil {\n"
			s += "\t\tdata." + column.name + " = *" + column.nameInDB + "_null\n"
			s += "\t} else {\n"
			s += "\t\tdata." + column.name + " = \"\"\n\t}\n"
			WriteStringInFile(w, s)
		}
	}

	WriteStringInFile(w, "\treturn err\n}\n\n")
}

func GenerateUpdateFunction(w *os.File, table *Table, structName string, tableName string) {
	s := "func (data *" + structName + ") Update() (err error) {\n"
	WriteStringInFile(w, s)

	for _, column := range table.columns {
		if column.isNull {
			s = "\tvar " + column.nameInDB + "_null sql.NullString\n"
			WriteStringInFile(w, s)
			s = "\tif len(data." + column.name + ") == 0 {\n"
			s += "\t\t" + column.nameInDB + "_null = sql.NullString{}\n"
			s += "\t} else {\n"
			s += "\t\t" + column.nameInDB + "_null = sql.NullString{\n\t\t\tString: data." + column.name
			s += ",\n\t\t\tValid:  true,\n\t\t}\n\t}\n"
			WriteStringInFile(w, s)
		}
	}

	WriteStringInFile(w, "\t_, err = DB.Exec(\n")

	s = "\t\t\"UPDATE " + tableName + " SET "
	i := 0
	for _, column := range table.columns {
		if table.pk == column {
			continue
		}
		if i != 0 {
			s += ", "
		}
		s += column.nameInDB + " = ?"
		i++
	}
	s += " WHERE " + table.pk.nameInDB + " = ?\",\n"
	WriteStringInFile(w, s)

	s = "\t\t"
	i = 0
	for _, column := range table.columns {
		if table.pk == column {
			continue
		}
		if i != 0 {
			s += ", "
		}
		if column.isNull {
			s += column.nameInDB + "_null"
		} else {
			s += "data." + column.name
		}
		i++
	}
	s += ", data." + table.pk.name
	s += ",\n"
	WriteStringInFile(w, s)
	WriteStringInFile(w, "\t)\n\treturn err\n}\n\n")
}

func GenerateCreateFunction(w *os.File, table *Table, structName string, tableName string) {
	s := "func (data *" + structName + ") Create() (err error) {\n"
	WriteStringInFile(w, s)

	for _, column := range table.columns {
		if column.isNull {
			s = "\tvar " + column.nameInDB + "_null sql.NullString\n"
			WriteStringInFile(w, s)
			s = "\tif len(data." + column.name + ") == 0 {\n"
			s += "\t\t" + column.nameInDB + "_null = sql.NullString{}\n"
			s += "\t} else {\n"
			s += "\t\t" + column.nameInDB + "_null = sql.NullString{\n\t\t\tString: data." + column.name
			s += ",\n\t\t\tValid:  true,\n\t\t}\n\t}\n"
			WriteStringInFile(w, s)
		}
	}

	WriteStringInFile(w, "\tresult, err := DB.Exec(\n")
	s = "\t\t\"INSERT INTO " + tableName + "("
	i := 0
	for _, column := range table.columns {
		if i != 0 {
			s += ", "
		}
		if table.pk == column {
			continue
		}
		s += "`" + column.nameInDB + "`"
		i++
	}
	s += ") VALUES ("
	for i = 0; i < len(table.columns)-1; i++ {
		if i != 0 {
			s += ", "
		}
		s += "?"
	}
	s += ")\",\n"
	WriteStringInFile(w, s)
	s = "\t\t"
	i = 0
	for _, column := range table.columns {
		if i != 0 {
			s += ", "
		}
		if table.pk == column {
			continue
		}
		if column.isNull {
			s += column.nameInDB + "_null"
		} else {
			s += "data." + column.name
		}
		i++
	}
	s += ",\n\t)\n"
	WriteStringInFile(w, s)
	s = "\tif err != nil {\n\t\treturn\n\t}\n\n"
	s += "\tlastID, err := result.LastInsertId()\n"
	s += "\tif err != nil {\n\t\treturn\n\t}\n"
	s += "\tdata." + table.pk.name + " = " + table.pk.fieldType + "(lastID)\n"
	s += "\treturn nil\n}\n"
	WriteStringInFile(w, s)
}

func GenerateFile(path string) {
	fset := token.NewFileSet()
	pathSlice := strings.SplitAfter(path, "/")

	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	PanicOnErr(err)
	for decl_iter, decl := range f.Decls {
		typedecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		typeSpecDecl, ok := typedecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		structDecl, ok := typeSpecDecl.Type.(*ast.StructType)
		if !ok {
			continue
		}
		structName := typedecl.Specs[0].(*ast.TypeSpec).Name.Name

		comment := strings.Split(typedecl.Doc.Text(), "myorm:")
		if len(comment) != 2 {
			continue
		}
		tableName := comment[1][:len(comment[1])-1]

		outputFile := strings.Join(pathSlice[:len(pathSlice)-1], "") + tableName + "_myorm.go"
		w, err := os.Create(outputFile)
		PanicOnErr(err)
		defer w.Close()

		GenerateHeaders(w)

		if decl_iter == 0 {
			GenerateDBInfoFunction(w)
		}

		table := &Table{name: tableName}

		for _, field := range structDecl.Fields.List {
			column := Column{
				name:      field.Names[0].Name,
				nameInDB:  strings.ToLower(field.Names[0].Name),
				fieldType: field.Type.(*ast.Ident).Name,
			}
			ignored := false
			if field.Tag != nil {

				tagSlice := strings.Split(field.Tag.Value, "`myorm:")
				if len(tagSlice) != 2 {
					continue
				}

				tags := strings.Split(tagSlice[1][:len(tagSlice[1])-1], ",")
				for _, tag := range tags {
					switch tag {
					case "\"-\"":
						ignored = true
					case "\"null\"":
						column.isNull = true
					case "\"primary_key\"":
						if !ignored {
							table.pk = &column
						}
					default:
						columnNameSlice := strings.Split(tag, "column:")
						if len(columnNameSlice) == 2 {
							column.nameInDB = columnNameSlice[1][:len(columnNameSlice[1])-1]
						}
					}
				}
			}
			if !ignored {
				table.columns = append(table.columns, &column)
			}
		}

		// FindByPK generate
		GenerateFindByPKFunction(w, table, structName, tableName)
		// Update generate
		GenerateUpdateFunction(w, table, structName, tableName)
		// Create generate
		GenerateCreateFunction(w, table, structName, tableName)
	}
}

func main() {
	GenerateFile(os.Args[1])
}
