package bebop

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// GenerateSettings holds customization options for what Generate should do.
type GenerateSettings struct {
	// PackageName is optional if the target bebop file defines a go_package
	// constant. If both are provided, PackageName will take precedence.
	PackageName string

	typeByters        map[string]string
	typeByteReaders   map[string]string
	typeMarshallers   map[string]string
	typeUnmarshallers map[string]string
	typeLengthers     map[string]string
	customRecordTypes map[string]struct{}

	ImportGenerationMode
	imported          []File
	importTypeAliases map[string]string

	GenerateUnsafeMethods bool
	SharedMemoryStrings   bool
}

type ImportGenerationMode uint8

const (
	// ImportGenerationModeSeparate will generate separate go files for
	// every bebop file, and will assume that imported files are
	// already generated. If imported file types are used and their containing
	// files do not contain a go_package constant, this mode will fail.
	ImportGenerationModeSeparate ImportGenerationMode = iota

	// ImportGenerationModeCombined will generate one go file including
	// all definitions from all imports. This does not require go_package
	// is defined anywhere, and maintains compatibility with files compilable
	// by the original bebopc compiler.
	ImportGenerationModeCombined ImportGenerationMode = iota
)

var allImportModes = map[ImportGenerationMode]struct{}{
	ImportGenerationModeSeparate: {},
	ImportGenerationModeCombined: {},
}

var reservedWords = map[string]struct{}{
	"map":      {},
	"array":    {},
	"struct":   {},
	"message":  {},
	"enum":     {},
	"readonly": {},
}

func (gs GenerateSettings) Validate() error {
	if _, ok := allImportModes[gs.ImportGenerationMode]; !ok {
		return fmt.Errorf("unknown import mode: %d", gs.ImportGenerationMode)
	}
	return nil
}

// Validate verifies a File can be successfully generated.
func (f File) Validate() error {
	allConsts := map[string]struct{}{}
	for _, c := range f.Consts {
		if _, ok := allConsts[c.Name]; ok {
			return fmt.Errorf("const has duplicated name %s", c.Name)
		}
		allConsts[c.Name] = struct{}{}
	}
	customTypes := map[string]struct{}{}
	structTypeUsage := map[string]map[string]bool{}
	for _, en := range f.Enums {
		if _, ok := primitiveTypes[en.Name]; ok {
			return fmt.Errorf("enum shares primitive type name %s", en.Name)
		}
		if _, ok := reservedWords[en.Name]; ok {
			return fmt.Errorf("enum shares reserved word name %s", en.Name)
		}
		if _, ok := customTypes[en.Name]; ok {
			return fmt.Errorf("enum has duplicated name %s", en.Name)
		}
		if en.Namespace != "" {
			nameWithoutPrefix := strings.TrimPrefix(en.Name, en.Namespace+".")
			customTypes[nameWithoutPrefix] = struct{}{}
		}
		customTypes[en.Name] = struct{}{}
		optionNames := map[string]struct{}{}
		optionValues := map[int32]struct{}{}
		for _, opt := range en.Options {
			if _, ok := optionNames[opt.Name]; ok {
				return fmt.Errorf("enum %s has duplicate option name %s", en.Name, opt.Name)
			}
			if _, ok := optionValues[opt.Value]; ok {
				return fmt.Errorf("enum %s has duplicate option value %d", en.Name, opt.Value)
			}
			optionNames[opt.Name] = struct{}{}
			optionValues[opt.Value] = struct{}{}
		}
	}
	for _, st := range f.Structs {
		if _, ok := primitiveTypes[st.Name]; ok {
			return fmt.Errorf("struct shares primitive type name %s", st.Name)
		}
		if _, ok := reservedWords[st.Name]; ok {
			return fmt.Errorf("struct shares reserved word name %s", st.Name)
		}
		if _, ok := customTypes[st.Name]; ok {
			return fmt.Errorf("struct has duplicated name %s", st.Name)
		}
		if st.Namespace != "" {
			nameWithoutPrefix := strings.TrimPrefix(st.Name, st.Namespace+".")
			customTypes[nameWithoutPrefix] = struct{}{}
			structTypeUsage[nameWithoutPrefix] = st.usedTypes()
		}
		customTypes[st.Name] = struct{}{}
		structTypeUsage[st.Name] = st.usedTypes()
		stNames := map[string]struct{}{}
		for _, fd := range st.Fields {
			if _, ok := stNames[fd.Name]; ok {
				return fmt.Errorf("struct %s has duplicate field name %s", st.Name, fd.Name)
			}
			stNames[fd.Name] = struct{}{}
		}
	}
	for _, msg := range f.Messages {
		if _, ok := primitiveTypes[msg.Name]; ok {
			return fmt.Errorf("message shares primitive type name %s", msg.Name)
		}
		if _, ok := reservedWords[msg.Name]; ok {
			return fmt.Errorf("message shares reserved word name %s", msg.Name)
		}
		if _, ok := customTypes[msg.Name]; ok {
			return fmt.Errorf("message has duplicated name %s", msg.Name)
		}
		if msg.Namespace != "" {
			nameWithoutPrefix := strings.TrimPrefix(msg.Name, msg.Namespace+".")
			customTypes[nameWithoutPrefix] = struct{}{}
		}
		customTypes[msg.Name] = struct{}{}
		msgNames := map[string]struct{}{}
		for _, fd := range msg.Fields {
			if _, ok := msgNames[fd.Name]; ok {
				return fmt.Errorf("message %s has duplicate field name %s", msg.Name, fd.Name)
			}

			msgNames[fd.Name] = struct{}{}
		}
	}
	allTypes := customTypes
	for typ := range primitiveTypes {
		allTypes[typ] = struct{}{}
	}
	for _, st := range f.Structs {
		for _, fd := range st.Fields {
			if err := typeDefined(fd.FieldType, allTypes); err != nil {
				return err
			}
		}
	}
	for _, msg := range f.Messages {
		for _, fd := range msg.Fields {
			if err := typeDefined(fd.FieldType, allTypes); err != nil {
				return err
			}
		}
	}

	// Todo: within a given struct, enum, or message, a field / option cannot
	// have a duplicate name

	// Determine which structs include themselves as required fields (which would lead
	// to the struct taking up infinite size)
	delta := true
	// This is an unbounded loop because for this case:
	// fooStruct:{barStruct}, barStruct{bizzStruct}, bizzStruct{bazStruct}, bazStruct{fooStruct}
	// After each iteration we have these updated usages:
	// 1. fooStruct:{barStruct, bizzStruct}, barStruct{bizzStruct, bazStruct}, bizzStruct{bazStruct, fooStruct}, bazStruct{fooStruct, barStruct}
	// 2. fooStruct:{barStruct, bizzStruct, bazStruct, fooStruct}, barStruct{bizzStruct, bazStruct, fooStruct, barStruct}, bizzStruct{bazStruct, fooStruct, barStruct, bizzStruct}, bazStruct{fooStruct, barStruct, bizzStruct, bazStruct}
	// ... and as the chain of structs gets longer the required iterations also increases.
	for delta {
		delta = false
		for stName, usage := range structTypeUsage {
			for stName2, usage2 := range structTypeUsage {
				if stName == stName2 {
					continue
				}
				// If struct1 includes struct2, it also includes
				// all fields that struct2 includes
				if usage[stName2] {
					for k, v := range usage2 {
						if !usage[k] {
							delta = true
						}
						usage[k] = v
					}
				}
			}
			structTypeUsage[stName] = usage
		}
	}
	for stName, usage := range structTypeUsage {
		if usage[stName] {
			return fmt.Errorf("struct %s recursively includes itself as a required field", stName)
		}
	}

	// Todo: union validation

	return nil
}

func typeDefined(ft FieldType, allTypes map[string]struct{}) error {
	if ft.Array != nil {
		return typeDefined(*ft.Array, allTypes)
	}
	if ft.Map != nil {
		if _, ok := allTypes[ft.Map.Key]; !ok {
			return fmt.Errorf("map key type %s undefined", ft.Map.Key)
		}
		return typeDefined(ft.Map.Value, allTypes)
	}
	if _, ok := allTypes[ft.Simple]; !ok {
		return fmt.Errorf("type %s undefined", ft.Simple)
	}
	return nil
}

// Generate writes a .go file out to w. If f has imports, it will
// parse the files at those imports and generate them according to
// settings.
func (f File) Generate(w io.Writer, settings GenerateSettings) error {
	if err := settings.Validate(); err != nil {
		return fmt.Errorf("invalid generation settings: %w", err)
	}

	if len(f.Imports) != 0 {
		thisFilePath := f.FileName
		if !path.IsAbs(f.FileName) {
			// assume relative to our current directory
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot determine working directory for import lookup: %v", err)
			}
			thisFilePath = filepath.Join(wd, thisFilePath)
		}
		thisDir := filepath.Dir(thisFilePath)
		// imports cause all elements of the imported file to be exported, so we
		// must follow imports infinitely deep.
		imports := make([]string, len(f.Imports))
		copy(imports, f.Imports)
		// TODO: should we forbid recursive imports? Yes, in non-combined mode
		// TODO: why are imports not scoped to a namespace?
		imported := map[string]struct{}{}
		for i := 0; i < len(imports); i++ {
			imp := imports[i]
			impPath := filepath.Join(thisDir, imp)
			if _, ok := imported[impPath]; ok {
				continue
			}
			impF, err := os.Open(impPath)
			if err != nil {
				return fmt.Errorf("failed to open imported file %s: %w", imp, err)
			}
			impFile, err := ReadFile(impF)
			if err != nil {
				impF.Close()
				return fmt.Errorf("failed to parse imported file %s: %w", imp, err)
			}
			impF.Close()
			settings.imported = append(settings.imported, impFile)
			imports = append(imports, impFile.Imports...)
			imported[impPath] = struct{}{}
		}
	}
	imports := []string{}
	settings.importTypeAliases = make(map[string]string)
	switch settings.ImportGenerationMode {
	case ImportGenerationModeSeparate:
		for _, imp := range settings.imported {
			if imp.GoPackage == "" {
				return fmt.Errorf("cannot import %s: file has no %s const", imp.FileName, goPackage)
			}
			packageName := imp.GoPackage
			packageNamespace := path.Base(packageName)
			for _, st := range imp.Structs {
				st.Namespace = packageNamespace
				namespacedName := packageNamespace + "." + st.Name
				settings.importTypeAliases[st.Name] = namespacedName
				st.Name = namespacedName
				f.Structs = append(f.Structs, st)
			}
			for _, un := range imp.Unions {
				un.Namespace = packageNamespace
				namespacedName := packageNamespace + "." + un.Name
				settings.importTypeAliases[un.Name] = namespacedName
				un.Name = namespacedName
				f.Unions = append(f.Unions, un)
			}
			for _, msg := range imp.Messages {
				msg.Namespace = packageNamespace
				namespacedName := packageNamespace + "." + msg.Name
				settings.importTypeAliases[msg.Name] = namespacedName
				msg.Name = namespacedName
				f.Messages = append(f.Messages, msg)
			}
			for _, enm := range imp.Enums {
				enm.Namespace = packageNamespace
				namespacedName := packageNamespace + "." + enm.Name
				settings.importTypeAliases[enm.Name] = namespacedName
				enm.Name = namespacedName
				f.Enums = append(f.Enums, enm)
			}
			// TODO: only import if we actually use a type from this package
			imports = append(imports, packageName)
		}
	case ImportGenerationModeCombined:
		// treat all imported files as a part of this file. Do not observe GoPackage.
		for _, imp := range settings.imported {
			f.Consts = append(f.Consts, imp.Consts...)
			f.Structs = append(f.Structs, imp.Structs...)
			f.Unions = append(f.Unions, imp.Unions...)
			f.Messages = append(f.Messages, imp.Messages...)
			f.Enums = append(f.Enums, imp.Enums...)
		}
	}

	if err := f.Validate(); err != nil {
		return fmt.Errorf("cannot generate file: %w", err)
	}
	settings.typeMarshallers = f.typeMarshallers()
	settings.typeByters = f.typeByters()
	settings.typeByteReaders = f.typeByteReaders(settings)
	settings.typeUnmarshallers = f.typeUnmarshallers()
	settings.typeLengthers = f.typeLengthers()
	settings.customRecordTypes = f.customRecordTypes()

	usedTypes := f.usedTypes()
	if settings.PackageName == "" && f.GoPackage != "" {
		settings.PackageName = path.Base(f.GoPackage)
	} else if settings.PackageName == "" {
		return fmt.Errorf("no package name is defined, provide a %s const or an explicit package name setting", goPackage)
	}

	writeLine(w, "// Code generated by bebopc-go; DO NOT EDIT.")
	writeLine(w, "")
	writeLine(w, "package %s", settings.PackageName)
	writeLine(w, "")
	if len(f.Messages)+len(f.Unions) != 0 {
		imports = append(imports, "bytes")
	}
	if len(f.Messages)+len(f.Structs)+len(f.Unions) != 0 {
		imports = append(imports, "io")
	}
	for _, c := range f.Consts {
		if c.impossibleGoConst() {
			imports = append(imports, "math")
			break
		}
	}
	if usedTypes["date"] {
		imports = append(imports, "time")
	}
	if len(f.Messages)+len(f.Structs)+len(f.Unions) != 0 {
		imports = append(imports, "github.com/200sc/bebop")
		imports = append(imports, "github.com/200sc/bebop/iohelp")
	}
	if len(imports) != 0 {
		writeLine(w, "import (")
		for _, i := range imports {
			writeLine(w, "\t%q", i)
		}
		writeLine(w, ")")
		writeLine(w, "")
	}

	impossibleGoConsts := []Const{}

	if len(f.Consts) != 0 {
		writeLine(w, "const (")
		for _, con := range f.Consts {
			if con.impossibleGoConst() {
				impossibleGoConsts = append(impossibleGoConsts, con)
			} else {
				con.Generate(w, settings)
			}
		}
		writeLine(w, ")")
		writeLine(w, "")
	}
	if len(impossibleGoConsts) != 0 {
		writeLine(w, "var (")
		for _, con := range impossibleGoConsts {
			con.Generate(w, settings)
		}
		writeLine(w, ")")
		writeLine(w, "")
	}
	for _, en := range f.Enums {
		if en.Namespace != "" {
			continue
		}
		en.Generate(w, settings)
	}
	for _, st := range f.Structs {
		if st.Namespace != "" {
			continue
		}
		st.Generate(w, settings)
	}
	for _, msg := range f.Messages {
		if msg.Namespace != "" {
			continue
		}
		msg.Generate(w, settings)
	}
	for _, union := range f.Unions {
		if union.Namespace != "" {
			continue
		}
		union.Generate(w, settings)
	}
	return nil
}

func writeLine(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format+"\n", args...)
}

func writeLineWithTabs(w io.Writer, format string, depth int, args ...string) {
	var assigner string
	var receiver string
	var typename string
	var namespace string
	var baretype string
	if len(args) > 0 {
		assigner = args[0]
		if assigner[0] == '&' || assigner[0] == '*' {
			receiver = assigner[1:]
		} else {
			if assigner[0] == '(' {
				receiver = "(*" + assigner[1:]
			} else {
				receiver = "*" + assigner
			}
		}
	}
	if len(args) > 1 {
		typename = args[1]
		if strings.Contains(typename, ".") {
			splitType := strings.Split(typename, ".")
			if len(splitType) == 2 {
				namespace = splitType[0]
				baretype = splitType[1]
			}
		}
	}
	// add tabs
	tbs := strings.Repeat("\t", depth)
	format = tbs + format
	format = strings.Replace(format, "\n", "\n"+tbs, -1)

	// %RECV = receiver
	// %ASGN = assigner
	// %TYPE = typename
	// %NAMESPACE = namespace
	// %BARETYPE = type name without namespace
	format = strings.Replace(format, "%RECV", receiver, -1)
	format = strings.Replace(format, "%ASGN", assigner, -1)
	format = strings.Replace(format, "%TYPE", typename, -1)
	format = strings.Replace(format, "%BARETYPE", baretype, -1)
	format = strings.Replace(format, "%NAMESPACE", namespace, -1)
	format = strings.Replace(format, "%KNAME", depthName("k", depth), -1)
	format = strings.Replace(format, "%VNAME", depthName("v", depth), -1)

	fmt.Fprint(w, format+"\n")
}

func getLineWithTabs(format string, depth int, args ...string) string {
	var b = new(bytes.Buffer)
	writeLineWithTabs(b, format, depth, args...)
	return b.String()
}

func writeComment(w io.Writer, depth int, comment string) {
	if comment == "" {
		return
	}
	tbs := strings.Repeat("\t", depth)

	commentLines := strings.Split(comment, "\n")
	for _, cm := range commentLines {
		writeLine(w, tbs+"//%s", cm)
	}
}

// Generate writes a .go enum definition out to w.
func (en Enum) Generate(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(en.Name)
	writeComment(w, 0, en.Comment)
	writeLine(w, "type %s uint32", exposedName)
	writeLine(w, "")
	if len(en.Options) != 0 {
		writeLine(w, "const (")
		for _, opt := range en.Options {
			writeComment(w, 1, opt.Comment)
			if opt.Deprecated {
				writeLine(w, "\t// Deprecated: %s", opt.DeprecatedMessage)
			}
			writeLine(w, "\t%s_%s %s = %d", exposedName, opt.Name, exposedName, opt.Value)
		}
		writeLine(w, ")")
		writeLine(w, "")
	}
}

// Generate writes a .go struct definition out to w.
func (st Struct) Generate(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(st.Name)
	if st.OpCode != 0 {
		writeLine(w, "const %sOpCode = 0x%x", exposedName, st.OpCode)
		writeLine(w, "")
	}
	writeLine(w, "var _ bebop.Record = &%s{}", exposedName)
	writeLine(w, "")
	writeComment(w, 0, st.Comment)
	writeLine(w, "type %s struct {", exposedName)
	for _, fd := range st.Fields {
		writeFieldDefinition(fd, w, st.ReadOnly, false, settings)
	}
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp %s) MarshalBebop() []byte {", exposedName)
	writeLine(w, "\tbuf := make([]byte, bbp.Size())")
	writeLine(w, "\tbbp.MarshalBebopTo(buf)")
	writeLine(w, "\treturn buf")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp %s) MarshalBebopTo(buf []byte) int {", exposedName)
	if st.OpCode != 0 {
		writeLine(w, "\tiohelp.WriteUint32Bytes(buf, uint32(%sOpCode))", exposedName)
		if len(st.Fields) > 0 {
			writeLine(w, "\tat := 4")
		}
	} else {
		if len(st.Fields) > 0 {
			writeLine(w, "\tat := 0")
		}
	}
	for _, fd := range st.Fields {
		name := exposeName(fd.Name)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeFieldByter("bbp."+name, fd.FieldType, w, settings, 1)
	}
	if len(st.Fields) == 0 {
		writeLine(w, "\treturn 0")
	} else {
		writeLine(w, "\treturn at")
	}
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp *%s) UnmarshalBebop(buf []byte) (err error) {", exposedName)
	if st.OpCode != 0 {
		if len(st.Fields) > 0 {
			writeLine(w, "\tat := 4")
		}
	} else {
		if len(st.Fields) > 0 {
			writeLine(w, "\tat := 0")
		}
	}
	for _, fd := range st.Fields {
		name := exposeName(fd.Name)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeFieldReadByter("bbp."+name, fd.FieldType, w, settings, 1, true)
	}
	writeLine(w, "\treturn nil")
	writeLine(w, "}")
	writeLine(w, "")
	if settings.GenerateUnsafeMethods {
		writeLine(w, "func (bbp *%s) MustUnmarshalBebop(buf []byte) {", exposedName)
		if st.OpCode != 0 {
			if len(st.Fields) > 0 {
				writeLine(w, "\tat := 4")
			}
		} else {
			if len(st.Fields) > 0 {
				writeLine(w, "\tat := 0")
			}
		}
		for _, fd := range st.Fields {
			name := exposeName(fd.Name)
			if st.ReadOnly {
				name = unexposeName(fd.Name)
			}
			writeFieldReadByter("bbp."+name, fd.FieldType, w, settings, 1, false)
		}
		writeLine(w, "}")
		writeLine(w, "")
	}
	isFirstTopLength = true
	writeLine(w, "func (bbp %s) EncodeBebop(iow io.Writer) (err error) {", exposedName)
	writeLine(w, "\tw := iohelp.NewErrorWriter(iow)")
	if st.OpCode != 0 {
		writeLine(w, "\tiohelp.WriteUint32(w, uint32(%sOpCode))", exposedName)
	}
	for _, fd := range st.Fields {
		name := exposeName(fd.Name)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeFieldMarshaller("bbp."+name, fd.FieldType, w, settings, 1)
	}
	writeLine(w, "\treturn w.Err")
	writeLine(w, "}")
	writeLine(w, "")
	isFirstTopLength = true
	writeLine(w, "func (bbp *%s) DecodeBebop(ior io.Reader) (err error) {", exposedName)
	writeLine(w, "\tr := iohelp.NewErrorReader(ior)")
	if st.OpCode != 0 {
		writeLine(w, "\tr.Read(make([]byte, 4))")
	}
	for _, fd := range st.Fields {
		name := exposeName(fd.Name)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeStructFieldUnmarshaller("&bbp."+name, fd.FieldType, w, settings, 1)
	}
	writeLine(w, "\treturn r.Err")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp %s) Size() int {", exposedName)
	writeLine(w, "\tbodyLen := 0")
	if st.OpCode != 0 {
		writeLine(w, "\tbodyLen += 4")
	}
	for _, fd := range st.Fields {
		name := exposeName(fd.Name)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		name = "bbp." + name
		writeMessageFieldBodyCount(name, fd.FieldType, w, settings, 1)
	}
	writeLine(w, "\treturn bodyLen")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func Make%[1]s(r iohelp.ErrorReader) (%[1]s, error) {", exposedName)
	writeLine(w, "\tv := %s{}", exposedName)
	writeLine(w, "\terr := v.DecodeBebop(r)")
	writeLine(w, "\treturn v, err")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func Make%[1]sFromBytes(buf []byte) (%[1]s, error) {", exposedName)
	writeLine(w, "\tv := %s{}", exposedName)
	writeLine(w, "\terr := v.UnmarshalBebop(buf)")
	writeLine(w, "\treturn v, err")
	writeLine(w, "}")
	writeLine(w, "")
	if settings.GenerateUnsafeMethods {
		writeLine(w, "func MustMake%[1]sFromBytes(buf []byte) %[1]s {", exposedName)
		writeLine(w, "\tv := %s{}", exposedName)
		writeLine(w, "\tv.MustUnmarshalBebop(buf)")
		writeLine(w, "\treturn v")
		writeLine(w, "}")
		writeLine(w, "")
	}
	// TODO: slices are not really readonly, we need to return a copy.
	if st.ReadOnly {
		for _, fd := range st.Fields {
			writeLine(w, "func (bbp %s) Get%s() %s {", exposedName, exposeName(fd.Name), fd.FieldType.goString(settings))
			writeLine(w, "\treturn bbp.%s", unexposeName(fd.Name))
			writeLine(w, "}")
			writeLine(w, "")
		}
		writeLine(w, "func New%s(", exposedName)
		for _, fd := range st.Fields {
			writeLine(w, "\t\t%s %s,", unexposeName(fd.Name), fd.FieldType.goString(settings))
		}
		writeLine(w, "\t) %s {", exposedName)
		writeLine(w, "\treturn %s{", exposedName)
		for _, fd := range st.Fields {
			writeLine(w, "\t\t%s: %s,", unexposeName(fd.Name), unexposeName(fd.Name))
		}
		writeLine(w, "\t}")
		writeLine(w, "}")
		writeLine(w, "")
	}
}

type fieldWithNumber struct {
	UnionField UnionField
	Field
	num uint8
}

// Generate writes a .go message definition out to w.
func (msg Message) Generate(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(msg.Name)
	if msg.OpCode != 0 {
		writeLine(w, "const %sOpCode = 0x%x", exposedName, msg.OpCode)
		writeLine(w, "")
	}
	writeLine(w, "var _ bebop.Record = &%s{}", exposedName)
	writeLine(w, "")
	writeComment(w, 0, msg.Comment)
	writeLine(w, "type %s struct {", exposedName)
	fields := make([]fieldWithNumber, 0, len(msg.Fields))
	for i, fd := range msg.Fields {
		fields = append(fields, fieldWithNumber{
			Field: fd,
			num:   i,
		})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].num < fields[j].num
	})
	for _, fd := range fields {
		writeFieldDefinition(fd.Field, w, false, true, settings)
	}
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp %s) MarshalBebop() []byte {", exposedName)
	writeLine(w, "\tbuf := make([]byte, bbp.Size())")
	writeLine(w, "\tbbp.MarshalBebopTo(buf)")
	writeLine(w, "\treturn buf")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp %s) MarshalBebopTo(buf []byte) int {", exposedName)
	if msg.OpCode != 0 {
		writeLine(w, "\tiohelp.WriteUint32Bytes(buf, uint32(%sOpCode))", exposedName)
		writeLine(w, "\tat := 4")
		writeLine(w, "\tiohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-8))")
	} else {
		writeLine(w, "\tat := 0")
		writeLine(w, "\tiohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))")
	}
	writeLine(w, "\tat += 4")
	for _, fd := range fields {
		name := exposeName(fd.Name)
		num := strconv.Itoa(int(fd.num))
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "buf[at] = %ASGN", 2, num)
		writeLineWithTabs(w, "at++", 2)
		writeFieldByter(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn at")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp *%s) UnmarshalBebop(buf []byte) (err error) {", exposedName)
	if msg.OpCode != 0 {
		writeLine(w, "\tat := 4")
	} else {
		writeLine(w, "\tat := 0")
	}
	writeLine(w, "\t_ = iohelp.ReadUint32Bytes(buf[at:])")
	writeLine(w, "\tbuf = buf[4:]")
	writeLine(w, "\tfor {")
	writeLine(w, "\t\tswitch buf[at] {")
	for _, fd := range fields {
		name := exposeName(fd.Name)
		writeLine(w, "\t\tcase %d:", fd.num)
		writeLine(w, "\t\t\tat += 1")
		writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
		writeFieldReadByter("(*bbp."+name+")", fd.FieldType, w, settings, 3, true)
	}
	writeLine(w, "\t\tdefault:")
	writeLine(w, "\t\t\treturn nil")
	writeLine(w, "\t\t}")
	writeLine(w, "\t}")
	writeLine(w, "}")
	writeLine(w, "")
	if settings.GenerateUnsafeMethods {
		writeLine(w, "func (bbp *%s) MustUnmarshalBebop(buf []byte) {", exposedName)
		if msg.OpCode != 0 {
			writeLine(w, "\tat := 4")
		} else {
			writeLine(w, "\tat := 0")
		}
		writeLine(w, "\t_ = iohelp.ReadUint32Bytes(buf[at:])")
		writeLine(w, "\tbuf = buf[4:]")
		writeLine(w, "\tfor {")
		writeLine(w, "\t\tswitch buf[at] {")
		for _, fd := range fields {
			name := exposeName(fd.Name)
			writeLine(w, "\t\tcase %d:", fd.num)
			writeLine(w, "\t\t\tat += 1")
			writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
			writeFieldReadByter("(*bbp."+name+")", fd.FieldType, w, settings, 3, false)
		}
		writeLine(w, "\t\tdefault:")
		writeLine(w, "\t\t\treturn")
		writeLine(w, "\t\t}")
		writeLine(w, "\t}")
		writeLine(w, "}")
		writeLine(w, "")
	}
	isFirstTopLength = true
	writeLine(w, "func (bbp %s) EncodeBebop(iow io.Writer) (err error) {", exposedName)
	writeLine(w, "\tw := iohelp.NewErrorWriter(iow)")
	if msg.OpCode != 0 {
		writeLine(w, "\tiohelp.WriteUint32(w, uint32(%sOpCode))", exposedName)
		writeLine(w, "\tiohelp.WriteUint32(w, uint32(bbp.Size()-8))")
	} else {
		writeLine(w, "\tiohelp.WriteUint32(w, uint32(bbp.Size()-4))")
	}
	for _, fd := range fields {
		name := exposeName(fd.Name)
		num := strconv.Itoa(int(fd.num))
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "w.Write([]byte{%ASGN})", 2, num)
		writeFieldMarshaller(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\tw.Write([]byte{0})")
	writeLine(w, "\treturn w.Err")
	writeLine(w, "}")
	writeLine(w, "")
	isFirstTopLength = true
	writeLine(w, "func (bbp *%s) DecodeBebop(ior io.Reader) (err error) {", exposedName)
	writeLine(w, "\ter := iohelp.NewErrorReader(ior)")
	if msg.OpCode != 0 {
		writeLine(w, "\tiohelp.ReadUint32(er)")
	}
	writeLine(w, "\tbodyLen := iohelp.ReadUint32(er)")
	// why read the entire body upfront? Because we're allowed
	// to exit early, and if we do exit early and this message
	// is a field of another record we need that record to resume
	// reading at the byte after this entire body.
	writeLine(w, "\tbody := make([]byte, bodyLen)")
	writeLine(w, "\ter.Read(body)")
	writeLine(w, "\tr := iohelp.NewErrorReader(bytes.NewReader(body))")
	writeLine(w, "\tfor {")
	writeLine(w, "\t\tswitch iohelp.ReadByte(r) {")
	for _, fd := range fields {
		writeLine(w, "\t\tcase %d:", fd.num)
		name := exposeName(fd.Name)
		writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
		writeMessageFieldUnmarshaller("bbp."+name, fd.FieldType, w, settings, 3)
	}
	// ref: https://github.com/RainwayApp/bebop/wiki/Wire-format#messages, final paragraph
	// for some reason we're allowed to skip parsing all remaining fields if we see one
	// that we don't know about.
	writeLine(w, "\t\tdefault:")
	writeLine(w, "\t\t\treturn er.Err")
	writeLine(w, "\t\t}")
	writeLine(w, "\t}")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp %s) Size() int {", exposedName)
	// size at front (4) + 0 byte (1)
	// q: why do messages end in a 0 byte?
	// a: (I think) because then we can loop reading a single byte for each field, and if we read 0
	// we know we're done and don't have to unread the byte
	writeLine(w, "\tbodyLen := 5")
	if msg.OpCode != 0 {
		writeLine(w, "\tbodyLen += 4")
	}
	for _, fd := range fields {
		name := exposeName(fd.Name)
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "bodyLen += 1", 2)
		writeMessageFieldBodyCount(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn bodyLen")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func Make%[1]s(r iohelp.ErrorReader) (%[1]s, error) {", exposedName)
	writeLine(w, "\tv := %s{}", exposedName)
	writeLine(w, "\terr := v.DecodeBebop(r)")
	writeLine(w, "\treturn v, err")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func Make%[1]sFromBytes(buf []byte) (%[1]s, error) {", exposedName)
	writeLine(w, "\tv := %s{}", exposedName)
	writeLine(w, "\terr := v.UnmarshalBebop(buf)")
	writeLine(w, "\treturn v, err")
	writeLine(w, "}")
	writeLine(w, "")
	if settings.GenerateUnsafeMethods {
		writeLine(w, "func MustMake%[1]sFromBytes(buf []byte) %[1]s {", exposedName)
		writeLine(w, "\tv := %s{}", exposedName)
		writeLine(w, "\tv.MustUnmarshalBebop(buf)")
		writeLine(w, "\treturn v")
		writeLine(w, "}")
		writeLine(w, "")
	}
}

// Generate writes a .go union definition out to w.
func (u Union) Generate(w io.Writer, settings GenerateSettings) {
	fields := make([]fieldWithNumber, 0, len(u.Fields))
	for i, ufd := range u.Fields {
		var fd Field
		if ufd.Struct != nil {
			fd.FieldType.Simple = ufd.Struct.Name
		}
		if ufd.Message != nil {
			fd.FieldType.Simple = ufd.Message.Name
		}
		fd.Name = fd.FieldType.Simple
		fields = append(fields, fieldWithNumber{
			UnionField: ufd,
			Field:      fd,
			num:        i,
		})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].num < fields[j].num
	})
	for _, field := range fields {
		if field.UnionField.Struct != nil {
			field.UnionField.Struct.Generate(w, settings)
		}
		if field.UnionField.Message != nil {
			field.UnionField.Message.Generate(w, settings)
		}
	}
	exposedName := exposeName(u.Name)
	if u.OpCode != 0 {
		writeLine(w, "const %sOpCode = 0x%x", exposedName, u.OpCode)
		writeLine(w, "")
	}
	writeLine(w, "var _ bebop.Record = &%s{}", exposedName)
	writeLine(w, "")
	writeComment(w, 0, u.Comment)
	writeLine(w, "type %s struct {", exposedName)
	for _, fd := range fields {
		writeFieldDefinition(fd.Field, w, false, true, settings)
	}
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp %s) MarshalBebop() []byte {", exposedName)
	writeLine(w, "\tbuf := make([]byte, bbp.Size())")
	writeLine(w, "\tbbp.MarshalBebopTo(buf)")
	writeLine(w, "\treturn buf")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp %s) MarshalBebopTo(buf []byte) int {", exposedName)
	if u.OpCode != 0 {
		writeLine(w, "\tiohelp.WriteUint32Bytes(buf, uint32(%sOpCode))", exposedName)
		writeLine(w, "\tat := 4")
		writeLine(w, "\tiohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-8))")
	} else {
		writeLine(w, "\tat := 0")
		writeLine(w, "\tiohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))")
	}
	writeLine(w, "\tat += 4")
	for _, fd := range fields {
		name := exposeName(fd.Name)
		num := strconv.Itoa(int(fd.num))
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "buf[at] = %ASGN", 2, num)
		writeLineWithTabs(w, "at++", 2)
		writeFieldByter(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "return at", 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn at")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp *%s) UnmarshalBebop(buf []byte) (err error) {", exposedName)
	if u.OpCode != 0 {
		writeLine(w, "\tat := 4")
	} else {
		writeLine(w, "\tat := 0")
	}
	writeLine(w, "\t_ = iohelp.ReadUint32Bytes(buf[at:])")
	writeLine(w, "\tbuf = buf[4:]")
	writeLine(w, "\tif len(buf) == 0 {")
	writeLine(w, "\t\treturn iohelp.UnpopulatedUnion")
	writeLine(w, "\t}")
	writeLine(w, "\tfor {")
	writeLine(w, "\t\tswitch buf[at] {")
	for _, fd := range fields {
		name := exposeName(fd.Name)
		writeLine(w, "\t\tcase %d:", fd.num)
		writeLine(w, "\t\t\tat += 1")
		writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
		writeFieldReadByter("(*bbp."+name+")", fd.FieldType, w, settings, 3, true)
		writeLine(w, "\t\t\treturn nil")
	}
	writeLine(w, "\t\tdefault:")
	writeLine(w, "\t\t\treturn nil")
	writeLine(w, "\t\t}")
	writeLine(w, "\t}")
	writeLine(w, "}")
	writeLine(w, "")
	if settings.GenerateUnsafeMethods {
		writeLine(w, "func (bbp *%s) MustUnmarshalBebop(buf []byte) {", exposedName)
		if u.OpCode != 0 {
			writeLine(w, "\tat := 4")
		} else {
			writeLine(w, "\tat := 0")
		}
		writeLine(w, "\t_ = iohelp.ReadUint32Bytes(buf[at:])")
		writeLine(w, "\tbuf = buf[4:]")
		writeLine(w, "\tfor {")
		writeLine(w, "\t\tswitch buf[at] {")
		for _, fd := range fields {
			name := exposeName(fd.Name)
			writeLine(w, "\t\tcase %d:", fd.num)
			writeLine(w, "\t\t\tat += 1")
			writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
			writeFieldReadByter("(*bbp."+name+")", fd.FieldType, w, settings, 3, false)
			writeLine(w, "\t\t\treturn")
		}
		writeLine(w, "\t\tdefault:")
		writeLine(w, "\t\t\treturn")
		writeLine(w, "\t\t}")
		writeLine(w, "\t}")
		writeLine(w, "}")
		writeLine(w, "")
	}
	isFirstTopLength = true
	writeLine(w, "func (bbp %s) EncodeBebop(iow io.Writer) (err error) {", exposedName)
	writeLine(w, "\tw := iohelp.NewErrorWriter(iow)")
	if u.OpCode != 0 {
		writeLine(w, "\tiohelp.WriteUint32(w, uint32(%sOpCode))", exposedName)
		writeLine(w, "\tiohelp.WriteUint32(w, uint32(bbp.Size()-8))")
	} else {
		writeLine(w, "\tiohelp.WriteUint32(w, uint32(bbp.Size()-4))")
	}
	for _, fd := range fields {
		name := exposeName(fd.Name)
		num := strconv.Itoa(int(fd.num))
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "w.Write([]byte{%ASGN})", 2, num)
		writeFieldMarshaller(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "return w.Err", 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn w.Err")
	writeLine(w, "}")
	writeLine(w, "")
	isFirstTopLength = true
	writeLine(w, "func (bbp *%s) DecodeBebop(ior io.Reader) (err error) {", exposedName)
	writeLine(w, "\ter := iohelp.NewErrorReader(ior)")
	if u.OpCode != 0 {
		writeLine(w, "\tiohelp.ReadUint32(er)")
	}
	writeLine(w, "\tbodyLen := iohelp.ReadUint32(er)")
	// why read the entire body upfront? Because we're allowed
	// to exit early, and if we do exit early and this message
	// is a field of another record we need that record to resume
	// reading at the byte after this entire body.
	writeLine(w, "\tbody := make([]byte, bodyLen)")
	writeLine(w, "\ter.Read(body)")
	writeLine(w, "\tr := iohelp.NewErrorReader(bytes.NewReader(body))")
	writeLine(w, "\tfor {")
	writeLine(w, "\t\tswitch iohelp.ReadByte(r) {")
	for _, fd := range fields {
		writeLine(w, "\t\tcase %d:", fd.num)
		name := exposeName(fd.Name)
		writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
		writeMessageFieldUnmarshaller("bbp."+name, fd.FieldType, w, settings, 3)
		writeLine(w, "\t\t\treturn er.Err")
	}
	// ref: https://github.com/RainwayApp/bebop/wiki/Wire-format#messages, final paragraph
	// for some reason we're allowed to skip parsing all remaining fields if we see one
	// that we don't know about.
	writeLine(w, "\t\tdefault:")
	writeLine(w, "\t\t\treturn er.Err")
	writeLine(w, "\t\t}")
	writeLine(w, "\t}")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func (bbp %s) Size() int {", exposedName)
	// size at front (4)
	writeLine(w, "\tbodyLen := 4")
	if u.OpCode != 0 {
		writeLine(w, "\tbodyLen += 4")
	}
	for _, fd := range fields {
		name := exposeName(fd.Name)
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "bodyLen += 1", 2)
		writeMessageFieldBodyCount(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "return bodyLen", 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn bodyLen")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func Make%[1]s(r iohelp.ErrorReader) (%[1]s, error) {", exposedName)
	writeLine(w, "\tv := %s{}", exposedName)
	writeLine(w, "\terr := v.DecodeBebop(r)")
	writeLine(w, "\treturn v, err")
	writeLine(w, "}")
	writeLine(w, "")
	writeLine(w, "func Make%[1]sFromBytes(buf []byte) (%[1]s, error) {", exposedName)
	writeLine(w, "\tv := %s{}", exposedName)
	writeLine(w, "\terr := v.UnmarshalBebop(buf)")
	writeLine(w, "\treturn v, err")
	writeLine(w, "}")
	writeLine(w, "")
	if settings.GenerateUnsafeMethods {
		writeLine(w, "func MustMake%[1]sFromBytes(buf []byte) %[1]s {", exposedName)
		writeLine(w, "\tv := %s{}", exposedName)
		writeLine(w, "\tv.MustUnmarshalBebop(buf)")
		writeLine(w, "\treturn v")
		writeLine(w, "}")
		writeLine(w, "")
	}
}

func (con Const) impossibleGoConst() bool {
	// unique floating point values (inf, nan) cannot be represented as consts in go,
	// at least not intuitively. This probagbly doesn't matter-- its rare to
	// rely on inf or nan floating points anyway, and even if you could get these as
	// consts you would need to use math.IsInf or math.IsNaN for many use cases.
	if con.FloatValue != nil {
		switch {
		case math.IsInf(*con.FloatValue, 1):
			return true
		case math.IsInf(*con.FloatValue, -1):
			return true
		case math.IsNaN(*con.FloatValue):
			return true
		}
	}
	return false
}

func (con Const) Generate(w io.Writer, settings GenerateSettings) {
	writeComment(w, 0, con.Comment)
	var val interface{}
	switch {
	case con.BoolValue != nil:
		val = *con.BoolValue
	case con.FloatValue != nil:
		switch {
		case math.IsInf(*con.FloatValue, 1):
			val = "math.Inf(1)"
		case math.IsInf(*con.FloatValue, -1):
			val = "math.Inf(-1)"
		case math.IsNaN(*con.FloatValue):
			val = "math.NaN()"
		default:
			val = strconv.FormatFloat(*con.FloatValue, 'g', -1, 64)
		}
	case con.IntValue != nil:
		val = *con.IntValue
	case con.UIntValue != nil:
		val = *con.UIntValue
	case con.StringValue != nil:
		val = fmt.Sprintf("%q", *con.StringValue)
	}
	writeLine(w, "\t%s = %v", exposeName(con.Name), val)
}

func writeFieldDefinition(fd Field, w io.Writer, readOnly bool, message bool, settings GenerateSettings) {
	writeComment(w, 1, fd.Comment)
	if fd.Deprecated {
		writeLine(w, "\t// Deprecated: %s", fd.DeprecatedMessage)
	}

	name := exposeName(fd.Name)
	if readOnly {
		name = unexposeName(fd.Name)
	}
	typ := fd.FieldType.goString(settings)
	if message {
		typ = "*" + typ
	}
	writeLine(w, "\t%s %s", name, typ)
}

func depthName(name string, depth int) string {
	return name + strconv.Itoa(depth)
}

var lengthInc = 0

func lengthName() string {
	lengthInc++
	return "ln" + strconv.Itoa(lengthInc)
}

func writeFieldByter(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	if typ.Array != nil {
		writeLineWithTabs(w, "iohelp.WriteUint32Bytes(buf[at:], uint32(len(%ASGN)))", depth, name)
		writeLineWithTabs(w, "at += 4", depth)
		if typ.Array.Simple == "byte" || typ.Array.Simple == "uint8" {
			writeLineWithTabs(w, "copy(buf[at:at+len(%ASGN)], %ASGN)", depth, name)
			writeLineWithTabs(w, "at += len(%ASGN)", depth, name)
			return
		}
		writeLineWithTabs(w, "for _, %VNAME := range %ASGN {", depth, name)
		writeFieldByter(depthName("v", depth), *typ.Array, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else if typ.Map != nil {
		writeLineWithTabs(w, "iohelp.WriteUint32Bytes(buf[at:], uint32(len(%ASGN)))", depth, name)
		writeLineWithTabs(w, "at += 4", depth)
		writeLineWithTabs(w, "for %KNAME, %VNAME := range %ASGN {", depth, name)
		writeLineWithTabs(w, settings.typeByters[typ.Map.Key], depth+1, depthName("k", depth))
		writeFieldByter(depthName("v", depth), typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeByters[simpleTyp], depth, name, typ.goString(settings))
	}
}

func writeLengthCheck(w io.Writer, ln string, depth int, args ...string) {
	writeLineWithTabs(w, "if len(buf[at:]) < "+ln+" {", depth, args...)
	writeLineWithTabs(w, "\t return iohelp.ErrTooShort", depth, args...)
	writeLineWithTabs(w, "}", depth, args...)
}

func writeFieldReadByter(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int, safe bool) {
	if typ.Array != nil {
		if safe {
			writeLengthCheck(w, "4", depth)
		}

		writeLineWithTabs(w, "%ASGN = make([]%TYPE, iohelp.ReadUint32Bytes(buf[at:]))", depth, name, typ.Array.goString(settings))
		writeLineWithTabs(w, "at += 4", depth)
		if safe {
			if sz, ok := fixedSizeTypes[typ.Array.Simple]; ok {
				writeLengthCheck(w, "len(%ASGN)*"+strconv.Itoa(int(sz)), depth, name)
				safe = false
			}
		}
		if typ.Array.Simple == "byte" || typ.Array.Simple == "uint8" {
			writeLineWithTabs(w, "copy(%ASGN, buf[at:at+len(%ASGN)])", depth, name)
			writeLineWithTabs(w, "at += len(%ASGN)", depth, name)
			return
		}
		iName := depthName("i", depth)
		writeLineWithTabs(w, "for "+iName+" := range %ASGN {", depth, name)
		writeFieldReadByter("("+name+")["+iName+"]", *typ.Array, w, settings, depth+1, safe)
		writeLineWithTabs(w, "}", depth)
	} else if typ.Map != nil {
		lnName := lengthName()
		writeLineWithTabs(w, lnName+" := iohelp.ReadUint32Bytes(buf[at:])", depth)
		writeLineWithTabs(w, "at += 4", depth)
		writeLineWithTabs(w, "%ASGN = make(%TYPE,"+lnName+")", depth, name, typ.Map.goString(settings))
		writeLineWithTabs(w, "for i := uint32(0); i < "+lnName+"; i++ {", depth, name)
		var ln string
		if format, ok := settings.typeByteReaders[typ.Map.Key+"&safe"]; ok && safe {
			ln = getLineWithTabs(format, depth+1, depthName("k", depth), simpleGoString(typ.Map.Key, settings))
		} else {
			if sz, ok := fixedSizeTypes[typ.Map.Key]; ok && safe {
				writeLengthCheck(w, strconv.Itoa(int(sz)), depth+1, depthName("k", depth))
			}
			ln = getLineWithTabs(settings.typeByteReaders[typ.Map.Key], depth+1, depthName("k", depth), typ.goString(settings))
		}
		w.Write([]byte(strings.Replace(ln, "=", ":=", 1)))
		writeFieldReadByter("("+name+")["+depthName("k", depth)+"]", typ.Map.Value, w, settings, depth+1, safe)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		if format, ok := settings.typeByteReaders[simpleTyp+"&safe"]; ok && safe {
			writeLineWithTabs(w, format, depth, name, typ.goString(settings))
		} else {
			if sz, ok := fixedSizeTypes[simpleTyp]; ok && safe {
				writeLengthCheck(w, strconv.Itoa(int(sz)), depth, name)
			}
			writeLineWithTabs(w, settings.typeByteReaders[simpleTyp], depth, name, typ.goString(settings))
		}
	}
}

func writeFieldMarshaller(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	if typ.Array != nil {
		writeLineWithTabs(w, "iohelp.WriteUint32(w, uint32(len(%ASGN)))", depth, name)
		writeLineWithTabs(w, "for _, elem := range %ASGN {", depth, name)
		writeFieldMarshaller("elem", *typ.Array, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else if typ.Map != nil {
		writeLineWithTabs(w, "iohelp.WriteUint32(w, uint32(len(%ASGN)))", depth, name)
		writeLineWithTabs(w, "for %KNAME, %VNAME := range %ASGN {", depth, name)
		writeLineWithTabs(w, settings.typeMarshallers[typ.Map.Key], depth+1, depthName("k", depth))
		writeFieldMarshaller(depthName("v", depth), typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeMarshallers[simpleTyp], depth, name, typ.goString(settings))
	}
}

var isFirstTopLength = true

func writeStructFieldUnmarshaller(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	iName := "i" + strconv.Itoa(depth)
	if typ.Array != nil {
		writeLineWithTabs(w, "%RECV = make([]%TYPE, iohelp.ReadUint32(r))", depth, name, typ.Array.goString(settings))
		writeLineWithTabs(w, "for "+iName+" := range %RECV {", depth, name)
		if name[0] == '&' {
			name = "&(" + name[1:] + "[" + iName + "])"
		} else {
			name = "(" + name + ")[" + iName + "]"
		}
		writeStructFieldUnmarshaller(name, *typ.Array, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else if typ.Map != nil {
		lnName := "ln" + strconv.Itoa(depth)
		if isFirstTopLength && depth == 1 {
			writeLineWithTabs(w, lnName+" := iohelp.ReadUint32(r)", depth)
			isFirstTopLength = false
		} else if depth == 1 {
			writeLineWithTabs(w, lnName+" = iohelp.ReadUint32(r)", depth)
		} else {
			writeLineWithTabs(w, lnName+" := iohelp.ReadUint32(r)", depth)
		}
		writeLineWithTabs(w, "%RECV = make(%TYPE, "+lnName+")", depth, name, typ.Map.goString(settings))
		writeLineWithTabs(w, "for "+iName+" := uint32(0); "+iName+" < "+lnName+"; "+iName+"++ {", depth, name)
		ln := getLineWithTabs(settings.typeUnmarshallers[typ.Map.Key], depth+1, "&"+depthName("k", depth))
		w.Write([]byte(strings.Replace(ln, "=", ":=", 1)))
		if name[0] == '&' {
			name = "&(" + name[1:] + "[" + depthName("k", depth) + "])"
		} else {
			name = "(" + name + ")[" + depthName("k", depth) + "]"
		}
		writeStructFieldUnmarshaller(name, typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeUnmarshallers[simpleTyp], depth, name, typ.goString(settings))
	}
}

func writeMessageFieldUnmarshaller(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	if typ.Array != nil {
		writeLineWithTabs(w, "%RECV = make([]%TYPE, iohelp.ReadUint32(r))", depth, name, typ.Array.goString(settings))
		writeLineWithTabs(w, "for i := range %RECV {", depth, name)
		writeMessageFieldUnmarshaller("("+name+")[i]", *typ.Array, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else if typ.Map != nil {
		lnName := depthName("ln", depth)
		writeLineWithTabs(w, lnName+" := iohelp.ReadUint32(r)", depth)
		writeLineWithTabs(w, "%RECV = make("+typ.Map.goString(settings)+")", depth, name)
		writeLineWithTabs(w, "for i := uint32(0); i < "+lnName+"; i++ {", depth, name)
		ln := getLineWithTabs(settings.typeUnmarshallers[typ.Map.Key], depth+1, "&"+depthName("k", depth))
		w.Write([]byte(strings.Replace(ln, "=", ":=", 1)))
		writeMessageFieldUnmarshaller("("+name+")["+depthName("k", depth)+"]", typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeUnmarshallers[simpleTyp], depth, name, typ.goString(settings))
	}
}

func typeNeedsElem(typ string, settings GenerateSettings) bool {
	switch typ {
	case "":
		return true
	case "string":
		return true
	}
	if _, ok := primitiveTypes[typ]; ok {
		return false
	}
	_, ok := settings.customRecordTypes[typ]
	return ok
}

func writeMessageFieldBodyCount(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	if typ.Array != nil {
		writeLineWithTabs(w, "bodyLen += 4", depth)
		if sz, ok := fixedSizeTypes[typ.Array.Simple]; ok {
			// short circuit-- write length times elem size
			writeLineWithTabs(w, "bodyLen += len(%ASGN) * "+strconv.Itoa(int(sz)), depth, name)
			return
		}
		if typeNeedsElem(typ.Array.Simple, settings) {
			writeLineWithTabs(w, "for _, elem := range %ASGN {", depth, name)
		} else {
			writeLineWithTabs(w, "for range %ASGN {", depth, name)
		}
		writeMessageFieldBodyCount("elem", *typ.Array, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else if typ.Map != nil {
		writeLineWithTabs(w, "bodyLen += 4", depth, name)
		useV := typeNeedsElem(typ.Map.Value.Simple, settings)
		useK := typ.Map.Key == "string"
		if useV && useK {
			writeLineWithTabs(w, "for %KNAME, %VNAME := range %ASGN {", depth, name)
		} else if useV {
			writeLineWithTabs(w, "for _, %VNAME := range %ASGN {", depth, name)
		} else if useK {
			writeLineWithTabs(w, "for %KNAME := range %ASGN {", depth, name)
		} else {
			writeLineWithTabs(w, "for range %ASGN {", depth, name)
		}
		writeLineWithTabs(w, settings.typeLengthers[typ.Map.Key], depth+1, depthName("k", depth))
		writeMessageFieldBodyCount(depthName("v", depth), typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeLengthers[simpleTyp], depth, name, typ.goString(settings))
	}
}

func exposeName(name string) string {
	if len(name) > 1 {
		return strings.ToUpper(string(name[0])) + name[1:]
	}
	if len(name) > 0 {
		return strings.ToUpper(string(name[0]))
	}
	return ""
}

func unexposeName(name string) string {
	if len(name) > 1 {
		return strings.ToLower(string(name[0])) + name[1:]
	}
	if len(name) > 0 {
		return strings.ToLower(string(name[0]))
	}
	return ""
}

var fixedSizeTypes = map[string]uint8{
	"bool":    1,
	"byte":    1,
	"uint8":   1,
	"uint16":  2,
	"int16":   2,
	"uint32":  4,
	"int32":   4,
	"uint64":  8,
	"int64":   8,
	"float32": 4,
	"float64": 8,
	"guid":    16,
	"date":    8,
}

func (f File) typeUnmarshallers() map[string]string {
	out := make(map[string]string)
	for typ := range fixedSizeTypes {
		out[typ] = "%RECV = iohelp.Read" + strings.Title(typ) + "(r)"
	}
	out["string"] = "%RECV = iohelp.ReadString(r)"
	out["guid"] = "%RECV = iohelp.ReadGUID(r)"
	for _, en := range f.Enums {
		out[en.Name] = "%RECV = %TYPE(iohelp.ReadUint32(r))"
	}
	for _, st := range f.Structs {
		var format string
		if st.Namespace != "" {
			format = "(%RECV), err = %NAMESPACE.Make%BARETYPE(r)\n" +
				"if err != nil {\n" +
				"\treturn err\n" +
				"}"
		} else {
			format = "(%RECV), err = make%TYPE(r)\n" +
				"if err != nil {\n" +
				"\treturn err\n" +
				"}"
		}
		out[st.Name] = format
	}
	for _, msg := range f.Messages {
		var format string
		if msg.Namespace != "" {
			format = "(%RECV), err = %NAMEPSACE.Make%BARETYPE(r)\n" +
				"if err != nil {\n" +
				"\treturn err\n" +
				"}"
		} else {
			format = "(%RECV), err = Make%TYPE(r)\n" +
				"if err != nil {\n" +
				"\treturn err\n" +
				"}"
		}
		out[msg.Name] = format
	}
	for _, union := range f.Unions {
		uout := union.typeUnmarshallers()
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeUnmarshallers() map[string]string {
	out := make(map[string]string)
	format := "(%RECV), err = make%TYPE(r)\n" +
		"if err != nil {\n" +
		"\treturn err\n" +
		"}"
	out[u.Name] = format
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			format := "(%RECV), err = make%TYPE(r)\n" +
				"if err != nil {\n" +
				"\treturn err\n" +
				"}"
			out[ufd.Struct.Name] = format
		}
		if ufd.Message != nil {
			format := "(%RECV), err = make%TYPE(r)\n" +
				"if err != nil {\n" +
				"\treturn err\n" +
				"}"
			out[ufd.Message.Name] = format
		}
	}
	return out
}

func (f File) typeMarshallers() map[string]string {
	out := make(map[string]string)
	for typ := range fixedSizeTypes {
		out[typ] = "iohelp.Write" + strings.Title(typ) + "(w, %ASGN)"
	}
	out["string"] = "iohelp.WriteUint32(w, uint32(len(%ASGN)))\n" +
		"w.Write([]byte(%ASGN))"
	out["guid"] = "iohelp.WriteGUID(w, %ASGN)"
	out["date"] = "if %ASGN != (time.Time{}) {\n" +
		"\tiohelp.WriteInt64(w, ((%ASGN).UnixNano()/100))\n" +
		"} else {\n" +
		"\tiohelp.WriteInt64(w, 0)\n" +
		"}"
	for _, en := range f.Enums {
		out[en.Name] = "iohelp.WriteUint32(w, uint32(%ASGN))"
	}
	for _, st := range f.Structs {
		format := "err = (%ASGN).EncodeBebop(w)\n" +
			"if err != nil {\n" +
			"\treturn err\n" +
			"}"
		out[st.Name] = format
	}
	for _, msg := range f.Messages {
		format := "err = (%ASGN).EncodeBebop(w)\n" +
			"if err != nil {\n" +
			"\treturn err\n" +
			"}"
		out[msg.Name] = format
	}
	for _, union := range f.Unions {
		uout := union.typeMarshallers()
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeMarshallers() map[string]string {
	out := make(map[string]string)
	format := "err = (%ASGN).EncodeBebop(w)\n" +
		"if err != nil {\n" +
		"\treturn err\n" +
		"}"
	out[u.Name] = format
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			format := "err = (%ASGN).EncodeBebop(w)\n" +
				"if err != nil {\n" +
				"\treturn err\n" +
				"}"
			out[ufd.Struct.Name] = format
		}
		if ufd.Message != nil {
			format := "err = (%ASGN).EncodeBebop(w)\n" +
				"if err != nil {\n" +
				"\treturn err\n" +
				"}"
			out[ufd.Message.Name] = format
		}
	}
	return out
}

func (f File) typeLengthers() map[string]string {
	out := make(map[string]string)
	for typ, sz := range fixedSizeTypes {
		out[typ] = "bodyLen += " + strconv.Itoa(int(sz))
	}
	out["string"] = "bodyLen += 4\n" + "bodyLen += len(%ASGN)"
	for _, en := range f.Enums {
		out[en.Name] = "bodyLen += 4"
	}
	for _, st := range f.Structs {
		out[st.Name] = "bodyLen += (%ASGN).Size()"
	}
	for _, msg := range f.Messages {
		out[msg.Name] = "bodyLen += (%ASGN).Size()"
	}
	for _, union := range f.Unions {
		uout := union.typeLengthers()
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeLengthers() map[string]string {
	out := make(map[string]string)
	out[u.Name] = "bodyLen += (%ASGN).Size()"
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			out[ufd.Struct.Name] = "bodyLen += (%ASGN).Size()"
		}
		if ufd.Message != nil {
			out[ufd.Message.Name] = "bodyLen += (%ASGN).Size()"
		}
	}
	return out
}

func (f File) typeByters() map[string]string {
	out := make(map[string]string)
	for typ, sz := range fixedSizeTypes {
		out[typ] = "iohelp.Write" + strings.Title(typ) + "Bytes(buf[at:], %ASGN)\n" +
			"at += " + strconv.Itoa(int(sz))
	}
	out["string"] = "iohelp.WriteUint32Bytes(buf[at:], uint32(len(%ASGN)))\n" +
		"at += 4\n" +
		"copy(buf[at:at+len(%ASGN)], []byte(%ASGN))\n" +
		"at += len(%ASGN)"

	out["guid"] = "iohelp.WriteGUIDBytes(buf[at:], %ASGN)\n" +
		"at += 16"
	out["date"] = "if %ASGN != (time.Time{}) {\n" +
		"\tiohelp.WriteInt64Bytes(buf[at:], ((%ASGN).UnixNano()/100))\n" +
		"} else {\n" +
		"\tiohelp.WriteInt64Bytes(buf[at:], 0)\n" +
		"}\n" +
		"at += 8"
	for _, en := range f.Enums {
		out[en.Name] = "iohelp.WriteUint32Bytes(buf[at:], uint32(%ASGN))\n" +
			"at += 4\n"
	}
	for _, st := range f.Structs {
		out[st.Name] = "(%ASGN).MarshalBebopTo(buf[at:])\n" +
			"at += (%ASGN).Size()"
	}
	for _, msg := range f.Messages {
		out[msg.Name] = "(%ASGN).MarshalBebopTo(buf[at:])\n" +
			"at += (%ASGN).Size()"
	}
	for _, union := range f.Unions {
		uout := union.typeByters()
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeByters() map[string]string {
	out := map[string]string{}
	out[u.Name] = "(%ASGN).MarshalBebopTo(buf[at:])\n" +
		"at += (%ASGN).Size()"
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			format := "(%ASGN).MarshalBebopTo(buf[at:])\n" +
				"at += (%ASGN).Size()"
			out[ufd.Struct.Name] = format
		}
		if ufd.Message != nil {
			format := "(%ASGN).MarshalBebopTo(buf[at:])\n" +
				"at += (%ASGN).Size()"
			out[ufd.Message.Name] = format
		}
	}
	return out
}

func (f File) typeByteReaders(gs GenerateSettings) map[string]string {
	out := make(map[string]string)
	for typ, sz := range fixedSizeTypes {
		out[typ] = "%ASGN = iohelp.Read" + strings.Title(typ) + "Bytes(buf[at:])\n" +
			"at += " + strconv.Itoa(int(sz))
	}
	out["guid"] = "%ASGN = iohelp.ReadGUIDBytes(buf[at:])\n" +
		"at += 16"

	stringRead := "ReadStringBytes(buf[at:])"
	if gs.SharedMemoryStrings {
		stringRead = "ReadStringBytesSharedMemory(buf[at:])"
	}

	out["string"] = "%ASGN =  iohelp.Must" + stringRead + "\n" +
		"at += 4+len(%ASGN)"

	out["string&safe"] = "%ASGN, err = iohelp." + stringRead + "\n" +
		"if err != nil {\n" +
		"\t return err\n" +
		"}\n" +
		"at += 4 + len(%ASGN)"

	for _, en := range f.Enums {
		out[en.Name] = "%ASGN = %TYPE(iohelp.ReadUint32Bytes(buf[at:]))\n" +
			"at += 4\n"
	}
	for _, st := range f.Structs {
		if st.Namespace != "" {
			out[st.Name] = "%ASGN = %NAMESPACE.MustMake%BARETYPEFromBytes(buf[at:])\n" +
				"at += (%ASGN).Size()"
			out[st.Name+"&safe"] = "%ASGN, err = %NAMESPACE.Make%BARETYPEFromBytes(buf[at:])\n" +
				"if err != nil {\n" +
				"\t return err\n" +
				"}\n" +
				"at += (%ASGN).Size()"
		} else {
			out[st.Name] = "%ASGN = MustMake%TYPEFromBytes(buf[at:])\n" +
				"at += (%ASGN).Size()"
			out[st.Name+"&safe"] = "%ASGN, err = Make%TYPEFromBytes(buf[at:])\n" +
				"if err != nil {\n" +
				"\t return err\n" +
				"}\n" +
				"at += (%ASGN).Size()"
		}
	}
	for _, msg := range f.Messages {
		if msg.Namespace != "" {
			out[msg.Name] = "%ASGN = %NAMESPACE.MustMake%BARETYPEFromBytes(buf[at:])\n" +
				"at += (%ASGN).Size()"
			out[msg.Name+"&safe"] = "%ASGN, err = %NAMESPACE.Make%BARETYPEFromBytes(buf[at:])\n" +
				"if err != nil {\n" +
				"\t return err\n" +
				"}\n" +
				"at += (%ASGN).Size()"
		} else {
			out[msg.Name] = "%ASGN = MustMake%TYPEFromBytes(buf[at:])\n" +
				"at += (%ASGN).Size()"
			out[msg.Name+"&safe"] = "%ASGN, err = Make%TYPEFromBytes(buf[at:])\n" +
				"if err != nil {\n" +
				"\t return err\n" +
				"}\n" +
				"at += (%ASGN).Size()"
		}
	}
	for _, union := range f.Unions {
		uout := union.typeByteReaders()
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeByteReaders() map[string]string {
	out := map[string]string{}
	if u.Namespace != "" {
		out[u.Name] = "%ASGN = %NAMESPACE.MustMake%BARETYPEFromBytes(buf[at:])\n" +
			"at += (%ASGN).Size()"
		out[u.Name+"&safe"] = "%ASGN, err = %NAMESPACE.Make%BARETYPEFromBytes(buf[at:])\n" +
			"if err != nil {\n" +
			"\t return err\n" +
			"}\n" +
			"at += (%ASGN).Size()"
	} else {
		out[u.Name] = "%ASGN = %NAMESPACE.MustMake%TYPEFromBytes(buf[at:])\n" +
			"at += (%ASGN).Size()"
		out[u.Name+"&safe"] = "%ASGN, err = %NAMESPACE.Make%TYPEFromBytes(buf[at:])\n" +
			"if err != nil {\n" +
			"\t return err\n" +
			"}\n" +
			"at += (%ASGN).Size()"
	}
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			st := ufd.Struct
			if st.Namespace != "" {
				out[st.Name] = "%ASGN = %NAMESPACE.MustMake%BARETYPEFromBytes(buf[at:])\n" +
					"at += (%ASGN).Size()"
				out[st.Name+"&safe"] = "%ASGN, err = %NAMESPACE.Make%BARETYPEFromBytes(buf[at:])\n" +
					"if err != nil {\n" +
					"\t return err\n" +
					"}\n" +
					"at += (%ASGN).Size()"
			} else {
				out[st.Name] = "%ASGN = MustMake%TYPEFromBytes(buf[at:])\n" +
					"at += (%ASGN).Size()"
				out[st.Name+"&safe"] = "%ASGN, err = Make%TYPEFromBytes(buf[at:])\n" +
					"if err != nil {\n" +
					"\t return err\n" +
					"}\n" +
					"at += (%ASGN).Size()"
			}
		}
		if ufd.Message != nil {
			msg := ufd.Message
			if msg.Namespace != "" {
				out[msg.Name] = "%ASGN = %NAMESPACE.MustMake%BARETYPEFromBytes(buf[at:])\n" +
					"at += (%ASGN).Size()"
				out[msg.Name+"&safe"] = "%ASGN, err = %NAMESPACE.Make%BARETYPEFromBytes(buf[at:])\n" +
					"if err != nil {\n" +
					"\t return err\n" +
					"}\n" +
					"at += (%ASGN).Size()"
			} else {
				out[msg.Name] = "%ASGN = MustMake%TYPEFromBytes(buf[at:])\n" +
					"at += (%ASGN).Size()"
				out[msg.Name+"&safe"] = "%ASGN, err = Make%TYPEFromBytes(buf[at:])\n" +
					"if err != nil {\n" +
					"\t return err\n" +
					"}\n" +
					"at += (%ASGN).Size()"
			}
		}
	}
	return out
}

func (f File) customRecordTypes() map[string]struct{} {
	out := make(map[string]struct{})
	for _, st := range f.Structs {
		out[st.Name] = struct{}{}
	}
	for _, msg := range f.Messages {
		out[msg.Name] = struct{}{}
	}
	for _, union := range f.Unions {
		out[union.Name] = struct{}{}
		for _, ufd := range union.Fields {
			if ufd.Struct != nil {
				out[ufd.Struct.Name] = struct{}{}
			}
			if ufd.Message != nil {
				out[ufd.Message.Name] = struct{}{}
			}
		}
	}
	return out
}

func (f File) usedTypes() map[string]bool {
	out := make(map[string]bool)
	for _, st := range f.Structs {
		stOut := st.usedTypes()
		for k, v := range stOut {
			out[k] = v
		}
	}
	for _, msg := range f.Messages {
		msgOut := msg.usedTypes()
		for k, v := range msgOut {
			out[k] = v
		}
	}
	for _, union := range f.Unions {
		unionOut := union.usedTypes()
		for k, v := range unionOut {
			out[k] = v
		}
	}
	return out
}

func (st Struct) usedTypes() map[string]bool {
	out := make(map[string]bool)
	for _, fd := range st.Fields {
		fdTypes := fd.usedTypes()
		for k, v := range fdTypes {
			out[k] = v
		}
	}
	return out
}

func (msg Message) usedTypes() map[string]bool {
	out := make(map[string]bool)
	for _, fd := range msg.Fields {
		fdTypes := fd.usedTypes()
		for k, v := range fdTypes {
			out[k] = v
		}
	}
	return out
}

func (u Union) usedTypes() map[string]bool {
	out := make(map[string]bool)
	for _, ufd := range u.Fields {
		fdTypes := ufd.usedTypes()
		for k, v := range fdTypes {
			out[k] = v
		}
	}
	return out
}

func (ft FieldType) usedTypes() map[string]bool {
	if ft.Array != nil {
		return ft.Array.usedTypes()
	}
	if ft.Map != nil {
		valTypes := ft.Map.Value.usedTypes()
		valTypes[ft.Map.Key] = true
		return valTypes
	}
	return map[string]bool{ft.Simple: true}
}

func (ufd UnionField) usedTypes() map[string]bool {
	if ufd.Struct != nil {
		return ufd.Struct.usedTypes()
	}
	if ufd.Message != nil {
		return ufd.Message.usedTypes()
	}
	return map[string]bool{}
}
