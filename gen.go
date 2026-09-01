//go:build ignore

// gen.go 从 swagger.json 生成 Ozon Performance API Go SDK 的模块代码。
//
// 用法: cd 项目根目录 && go run gen.go
//
// 输出（每个 API 域一个目录）:
//   - <dir>/types.go        请求/响应结构与枚举类型
//   - <dir>/service.go      API 方法封装
//   - <dir>/service_test.go 方法冒烟测试
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const modulePath = "github.com/QuoVadis86/go-ozon-performance"

var tagsToDir = map[string]string{
	"Campaign":     "campaign",
	"Statistics":   "statistics",
	"Ad":           "ad",
	"Product":      "product",
	"Search-Promo": "searchpromo",
	"Vendor":       "vendor",
}

// stripPrefixes 从 schema 名称剥离的服务端前缀，长前缀优先。
var stripPrefixes = []string{"extstatistics", "extorganisation", "extcampaign", "extstat", "camptype", "advcamptype"}

var abrMap = map[string]string{
	"id": "ID", "sku": "SKU", "url": "URL", "json": "JSON", "uuid": "UUID",
	"api": "API", "cpc": "CPC", "cpm": "CPM", "cpo": "CPO", "ctr": "CTR", "drr": "DRR",
}

type param struct {
	Name string
	Type string
	Fmt  string
	Enum []string
}

type rawSchema struct {
	Type        string                     `json:"type"`
	Properties  map[string]json.RawMessage `json:"properties"`
	Items       *rawSchema                 `json:"items"`
	Enum        json.RawMessage            `json:"enum"`
	Description string                     `json:"description"`
	Format      string                     `json:"format"`
	Default     string                     `json:"default"`
	Additional  *rawSchema                 `json:"additionalProperties"`
}

func main() {
	data, err := os.ReadFile("origin/swagger.json")
	if err != nil {
		panic(err)
	}
	var doc struct {
		Paths      map[string]map[string]json.RawMessage `json:"paths"`
		Components struct {
			Schemas       map[string]json.RawMessage `json:"schemas"`
			RequestBodies map[string]json.RawMessage `json:"requestBodies"`
		} `json:"components"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		panic(err)
	}

	schemas := map[string]*rawSchema{}
	for n, raw := range doc.Components.Schemas {
		var s rawSchema
		if err := json.Unmarshal(raw, &s); err != nil {
			panic("schema " + n + ": " + err.Error())
		}
		schemas[n] = &s
	}

	requestBodies := parseRequestBodies(doc.Components.RequestBodies)
	ops := collectOps(doc.Paths, requestBodies)

	// 将 query 参数与内联 body 物化为合成 schema，统一走类型生成管道。
	synthQuery := map[string][]param{}
	synthBody := map[string][]param{}
	for i := range ops {
		op := &ops[i]
		if op.reqSchema != "" || len(op.query) == 0 {
			continue
		}
		sname := op.Name + "Request"
		if len(op.reqInline) > 0 {
			synthBody[sname] = op.reqInline
		} else {
			synthQuery[sname] = op.query
		}
		op.reqSchema = sname
	}

	for name, params := range synthBody {
		schemas[name] = schemaFromParams(name, params, false)
	}
	for name, params := range synthQuery {
		schemas[name] = schemaFromParams(name, params, true)
	}

	g, err := newGenerator(schemas)
	if err != nil {
		panic(err)
	}
	g.querySchemas = synthQuery

	byDir := map[string][]opSpec{}
	for _, op := range ops {
		byDir[op.Dir] = append(byDir[op.Dir], op)
	}
	for dir := range byDir {
		sort.Slice(byDir[dir], func(i, j int) bool {
			if byDir[dir][i].Name == byDir[dir][j].Name {
				return byDir[dir][i].Path < byDir[dir][j].Path
			}
			return byDir[dir][i].Name < byDir[dir][j].Name
		})
	}
	dirs := make([]string, 0, len(byDir))
	for d := range byDir {
		dirs = append(dirs, d)
	}
	sort.Strings(dirs)

	for _, dir := range dirs {
		os.MkdirAll(dir, 0o755)
		cols := allTypesFor(byDir[dir], g)
		writeGo(dir, "types.go", g.genTypes(byDir[dir], cols))
		writeGo(dir, "service.go", g.genService(byDir[dir]))
		writeGo(dir, "service_test.go", genTests(byDir[dir], g))
		fmt.Printf("generated %s: %d methods, %d types\n", dir, len(byDir[dir]), len(cols))
	}
	fmt.Println("done")
}

// --- 请求体/操作解析 ---

func parseRequestBodies(rbs map[string]json.RawMessage) map[string]any {
	out := map[string]any{}
	for n, raw := range rbs {
		var v any
		_ = json.Unmarshal(raw, &v)
		out[n] = v
	}
	return out
}

func isHTTPMethod(s string) bool {
	switch s {
	case "get", "post", "put", "patch", "delete":
		return true
	}
	return false
}

type opSpec struct {
	Dir        string
	Method     string
	Path       string
	Name       string
	Summary    string
	pathArgs   []string
	pathFmt    []string
	query      []param
	reqSchema  string
	reqInline  []param
	respSchema string
	rawResp    bool
}

func collectOps(paths map[string]map[string]json.RawMessage, requestBodies map[string]any) []opSpec {
	var ops []opSpec
	for path, methods := range paths {
		for httpM, raw := range methods {
			if !isHTTPMethod(httpM) {
				continue
			}
			var item struct {
				Summary     string   `json:"summary"`
				OperationID string   `json:"operationId"`
				Tags        []string `json:"tags"`
				Parameters  []struct {
					Name   string         `json:"name"`
					In     string         `json:"in"`
					Schema map[string]any `json:"schema"`
				} `json:"parameters"`
				RequestBody json.RawMessage            `json:"requestBody"`
				Responses   map[string]json.RawMessage `json:"responses"`
			}
			if err := json.Unmarshal(raw, &item); err != nil {
				panic(err)
			}
			dir := ""
			for _, t := range item.Tags {
				if d, ok := tagsToDir[t]; ok {
					dir = d
					break
				}
			}
			if dir == "" {
				continue
			}
			op := opSpec{
				Dir:     dir,
				Method:  strings.ToUpper(httpM),
				Path:    path,
				Name:    opMethodName(item.OperationID),
				Summary: flatten(item.Summary),
			}
			for _, p := range item.Parameters {
				pv := parseParam(p.Name, p.Schema)
				if p.In == "path" {
					op.pathArgs = append(op.pathArgs, p.Name)
					op.pathFmt = append(op.pathFmt, fmt.Sprintf("%v", p.Schema["format"]))
				} else if p.In == "query" {
					op.query = append(op.query, pv)
				}
			}
			if len(item.RequestBody) > 0 {
				op.reqSchema, op.reqInline = resolveReqSchema(item.RequestBody, requestBodies)
			}
			op.respSchema, op.rawResp = resolveRespSchema(item.Responses)
			ops = append(ops, op)
		}
	}
	return ops
}

func parseParam(name string, m map[string]any) param {
	p := param{Name: name}
	if t, ok := m["type"].(string); ok {
		p.Type = t
	}
	if f, ok := m["format"].(string); ok {
		p.Fmt = f
	}
	if items, ok := m["items"].(map[string]any); ok {
		if t, ok := items["type"].(string); ok {
			p.Type = "array:" + t
		}
		if f, ok := items["format"].(string); ok {
			p.Fmt = f
		}
	}
	if en, ok := m["enum"].([]any); ok {
		for _, e := range en {
			p.Enum = append(p.Enum, fmt.Sprintf("%v", e))
		}
	} else if desc, ok := m["description"].(string); ok {
		if vals, _ := enumFromDesc(desc); vals != nil {
			p.Enum = vals
		}
	}
	return p
}

func resolveReqSchema(raw json.RawMessage, requestBodies map[string]any) (string, []param) {
	var body struct {
		Ref     string `json:"$ref"`
		Content map[string]struct {
			Schema map[string]any `json:"schema"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return "", nil
	}
	if body.Ref != "" {
		rname := strings.TrimPrefix(body.Ref, "#/components/requestBodies/")
		if inner, ok := requestBodies[rname]; ok {
			innerJSON, _ := json.Marshal(inner)
			return resolveReqSchema(innerJSON, requestBodies)
		}
		return rname, nil
	}
	for _, c := range body.Content {
		if ref, ok := c.Schema["$ref"].(string); ok {
			return strings.TrimPrefix(ref, "#/components/schemas/"), nil
		}
		props, _ := c.Schema["properties"].(map[string]any)
		keys := sortedAnyKeys(props)
		var params []param
		for _, k := range keys {
			pm, _ := props[k].(map[string]any)
			params = append(params, parseParam(k, pm))
		}
		return "", params
	}
	return "", nil
}

func resolveRespSchema(responses map[string]json.RawMessage) (string, bool) {
	r200, ok := responses["200"]
	if !ok {
		return "", false
	}
	var resp struct {
		Content map[string]struct {
			Schema map[string]any `json:"schema"`
		} `json:"content"`
	}
	if err := json.Unmarshal(r200, &resp); err != nil {
		return "", false
	}
	for ctype, c := range resp.Content {
		if ctype == "text/csv" {
			return "", true
		}
		if ref, ok := c.Schema["$ref"].(string); ok {
			sname := strings.TrimPrefix(ref, "#/components/schemas/")
			if sname == "extstatStatisticsReport" {
				return "", true
			}
			return sname, false
		}
	}
	return "", false
}

func opMethodName(opID string) string {
	parts := strings.Split(opID, "_")
	if len(parts) > 1 {
		return toGoName(strings.Join(parts[1:], "_"))
	}
	return toGoName(opID)
}

// schemaFromParams 从参数构造合成 schema。query 为真时字段使用 url tag。
func schemaFromParams(name string, params []param, query bool) *rawSchema {
	props := map[string]json.RawMessage{}
	for _, p := range params {
		m := map[string]any{}
		switch {
		case strings.HasPrefix(p.Type, "array:"):
			m["type"] = "array"
			elem := map[string]any{"type": strings.TrimPrefix(p.Type, "array:")}
			if p.Fmt == "uint64" {
				elem["format"] = "uint64"
			}
			m["items"] = elem
		case p.Fmt == "uint64":
			m["type"] = "string"
			m["format"] = "uint64"
		case len(p.Enum) > 0:
			m["type"] = "string"
			m["enum"] = p.Enum
		case p.Type == "boolean":
			m["type"] = "boolean"
		case strings.HasPrefix(p.Type, "array:"):
			m["type"] = "array"
			elem := map[string]any{"type": strings.TrimPrefix(p.Type, "array:")}
			if p.Fmt == "uint64" {
				elem["format"] = "uint64"
			}
			m["items"] = elem
		case p.Type == "integer":
			m["type"] = "integer"
			if p.Fmt != "" {
				m["format"] = p.Fmt
			}
		case p.Type == "number":
			m["type"] = "number"
		default:
			m["type"] = "string"
			if p.Fmt != "" {
				m["format"] = p.Fmt
			}
		}
		b, _ := json.Marshal(m)
		props[p.Name] = b
	}
	return &rawSchema{Type: "object", Properties: props}
}

// --- 命名工具 ---

func toCamel(s string) string {
	var b strings.Builder
	upperNext := true
	for _, r := range s {
		if r == '_' || r == '.' || r == '/' || r == '-' {
			upperNext = true
			continue
		}
		if upperNext {
			l := strings.ToLower(string(r))
			if v, ok := abrMap[l]; ok {
				b.WriteString(v)
			} else {
				b.WriteString(strings.ToUpper(string(r)))
			}
			upperNext = false
		} else {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "X"
	}
	return b.String()
}

func toGoName(s string) string {
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	var r []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		l := strings.ToLower(p)
		if v, ok := abrMap[l]; ok {
			r = append(r, v)
		} else {
			r = append(r, strings.ToUpper(p[:1])+p[1:])
		}
	}
	out := strings.Join(r, "")
	if out != "" && out[0] >= '0' && out[0] <= '9' {
		out = "X" + out
	}
	return out
}

func schemaGoName(name string) string {
	for _, p := range stripPrefixes {
		lower := strings.ToLower(name)
		if strings.HasPrefix(lower, p) && len(name) > len(p) {
			rest := name[len(p):]
			if idx := strings.Index(strings.ToLower(rest), p); idx > 0 {
				rest = rest[:idx]
			}
			return toCamel(rest)
		}
	}
	return toCamel(name)
}

func flatten(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

// enumFromDesc 解析描述行中的反引号枚举值。
func enumFromDesc(desc string) ([]string, []string) {
	if desc == "" {
		return nil, nil
	}
	boolSet := map[string]bool{"true": true, "false": true, "0": true, "1": true}
	seen := map[string]bool{}
	var vals, descs []string
	for _, line := range strings.Split(desc, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 支持 "- `VAL` — 描述" 与 "`VAL` — 描述" 两种格式
		rest := line
		if strings.HasPrefix(rest, "-") || strings.HasPrefix(rest, "*") || strings.HasPrefix(rest, "–") || strings.HasPrefix(rest, "—") {
			rest = strings.TrimLeft(rest, "-–—* \t")
		}
		if !strings.HasPrefix(rest, "`") {
			continue
		}
		rest = rest[1:]
		end := strings.Index(rest, "`")
		if end < 0 {
			continue
		}
		v := strings.TrimSpace(rest[:end])
		rest = strings.TrimSpace(rest[end+1:])
		if v == "" || boolSet[strings.ToLower(v)] || seen[v] {
			continue
		}
		seen[v] = true
		vals = append(vals, v)
		descs = append(descs, strings.TrimRight(strings.TrimLeft(rest, "—–- \t"), "。.;;；,， "))
	}
	if len(vals) < 2 {
		return nil, nil
	}
	return vals, descs
}

func valueToCamel(v string) string {
	clean := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, v)
	parts := strings.Split(strings.Trim(clean, "_"), "_")
	var out []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		if strings.ToUpper(p) == p && len(p) <= 3 {
			out = append(out, p)
		} else if strings.ToUpper(p) == p {
			out = append(out, strings.ToUpper(p[:1])+strings.ToLower(p[1:]))
		} else {
			out = append(out, strings.ToUpper(p[:1])+p[1:])
		}
	}
	if len(out) == 0 {
		return "X"
	}
	return strings.Join(out, "")
}

func enumTypeName(structName, fieldName string, used map[string]string) string {
	base := toGoName(fieldName)
	if len(base) < 3 {
		parts := splitCamel(structName)
		if len(parts) > 1 {
			base = parts[len(parts)-1] + base
		}
	}
	if _, taken := used[base]; !taken {
		return base
	}
	parts := splitCamel(structName)
	for i := len(parts) - 1; i >= 0; i-- {
		if _, taken := used[parts[i]+base]; !taken {
			return parts[i] + base
		}
	}
	return structName + base
}

func splitCamel(s string) []string {
	var parts []string
	var cur []rune
	for _, r := range s {
		if r >= 'A' && r <= 'Z' && len(cur) > 0 {
			parts = append(parts, string(cur))
			cur = []rune{r}
		} else {
			cur = append(cur, r)
		}
	}
	if len(cur) > 0 {
		parts = append(parts, string(cur))
	}
	return parts
}

func sortedAnyKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func writeGo(dir, file, src string) {
	formatted, err := format.Source([]byte(src))
	if err != nil {
		fmt.Fprintf(os.Stderr, "gofmt %s/%s failed: %v\n---\n%s\n---\n", dir, file, err, src)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(dir, file), formatted, 0o644); err != nil {
		panic(err)
	}
}

// --- 类型图 / 依赖闭包 ---

type generator struct {
	schemas      map[string]*rawSchema
	goName       map[string]string
	depOf        map[string]map[string]bool
	querySchemas map[string][]param
	emitted      map[string]bool
	usedEnum     map[string]string
	fieldEnum    map[string]string
	enumValSet   map[string]string
}

func newGenerator(schemas map[string]*rawSchema) (*generator, error) {
	g := &generator{
		schemas:    schemas,
		goName:     map[string]string{},
		depOf:      map[string]map[string]bool{},
		emitted:    map[string]bool{},
		usedEnum:   map[string]string{},
		fieldEnum:  map[string]string{},
		enumValSet: map[string]string{},
	}
	claimed := map[string]string{}
	names := make([]string, 0, len(schemas))
	for n := range schemas {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		short := schemaGoName(n)
		if prev, ok := claimed[short]; ok && prev != n {
			short = toCamel(n)
			if prev2, ok2 := claimed[short]; ok2 && prev2 != n {
				short = n
			}
		}
		claimed[short] = n
		g.goName[n] = short
	}
	for n, s := range schemas {
		g.depOf[n] = map[string]bool{}
		for _, raw := range s.Properties {
			p, _ := rawToProp(raw)
			if r, ok := p["$ref"].(string); ok {
				addDep(g.depOf[n], sanitizeRef(r))
			}
			if items, ok := p["items"].(map[string]any); ok {
				if r, ok := items["$ref"].(string); ok {
					addDep(g.depOf[n], sanitizeRef(r))
				}
			}
		}
	}
	return g, nil
}

func rawToProp(raw json.RawMessage) (map[string]any, error) {
	var m map[string]any
	err := json.Unmarshal(raw, &m)
	return m, err
}

func sanitizeRef(r string) string {
	return strings.TrimPrefix(r, "#/components/schemas/")
}

func addDep(deps map[string]bool, r string) {
	if r != "" {
		deps[r] = true
	}
}

func allTypesFor(ops []opSpec, g *generator) map[string]bool {
	set := map[string]bool{}
	for _, op := range ops {
		if op.reqSchema != "" && !strings.HasPrefix(op.reqSchema, "type:") {
			set[op.reqSchema] = true
		}
		if op.respSchema != "" {
			set[op.respSchema] = true
		}
	}
	visited := map[string]bool{}
	var visit func(string)
	visit = func(n string) {
		if visited[n] {
			return
		}
		visited[n] = true
		for dep := range g.depOf[n] {
			visit(dep)
		}
	}
	for n := range set {
		visit(n)
	}
	for n := range visited {
		if _, ok := g.schemas[n]; !ok {
			delete(visited, n)
		}
	}
	return visited
}

// --- types.go 生成 ---

func (g *generator) genTypes(ops []opSpec, cols map[string]bool) string {
	g.emitted = map[string]bool{}
	g.usedEnum = map[string]string{}
	g.fieldEnum = map[string]string{}
	g.enumValSet = map[string]string{}

	var body bytes.Buffer
	names := make([]string, 0, len(cols))
	for n := range cols {
		names = append(names, n)
	}
	sort.Strings(names)
	order := g.topoOrder(names)

	// 第一轮：先输出所有纯枚举 schema，注册到枚举池，供字段枚举按值集复用。
	for _, n := range order {
		if !cols[n] {
			continue
		}
		if g.isEnumSchema(n) {
			g.emitSchema(&body, n)
		}
	}
	// 第二轮：输出其余类型。
	for _, n := range order {
		if !cols[n] {
			continue
		}
		if !g.isEnumSchema(n) {
			g.emitSchema(&body, n)
		}
	}

	needsTransport := false
	for _, n := range order {
		if cols[n] && g.usesTransport(n) {
			needsTransport = true
		}
	}

	var b bytes.Buffer
	fmt.Fprintf(&b, "package %s\n", ops[0].Dir)
	if needsTransport {
		fmt.Fprintf(&b, "\nimport \"%s/transport\"\n", modulePath)
	}
	b.WriteString("\n" + body.String())
	return b.String()
}

func (g *generator) topoOrder(names []string) []string {
	visited := map[string]bool{}
	var order []string
	var visit func(string)
	visit = func(n string) {
		if visited[n] {
			return
		}
		visited[n] = true
		deps := make([]string, 0, len(g.depOf[n]))
		for d := range g.depOf[n] {
			deps = append(deps, d)
		}
		sort.Strings(deps)
		for _, dep := range deps {
			visit(dep)
		}
		order = append(order, n)
	}
	for _, n := range names {
		visit(n)
	}
	return order
}

func (g *generator) usesTransport(n string) bool {
	if _, ok := g.schemas[n]; !ok {
		return false
	}
	return hasUint64Props(g.schemas[n], g.querySchemas)
}

func hasUint64Props(s *rawSchema, querySchemas map[string][]param) bool {
	for _, raw := range s.Properties {
		p, _ := rawToProp(raw)
		if p["format"] == "uint64" && p["type"] == "string" {
			return true
		}
		if items, ok := p["items"].(map[string]any); ok {
			if items["format"] == "uint64" {
				return true
			}
		}
	}
	return false
}

func (g *generator) emitSchema(b *bytes.Buffer, n string) {
	if g.emitted[n] {
		return
	}
	g.emitted[n] = true
	s := g.schemas[n]
	gn := g.goName[n]

	// 枚举 schema（含 enum:null 但描述含枚举值）
	if vals, descs := g.extendedEnum(n); vals != nil {
		if gn == "State" {
		}
		// 注册到枚举池，供内联/查询枚举按值集复用
		valSet := strings.Join(vals, ",")
		if len(vals) >= 2 {
			g.usedEnum[gn] = valSet
			g.enumValSet[gn] = valSet
		}
		g.emitEnum(b, gn, vals, descs)
		return
	}
	if s.Type == "string" && s.Properties == nil {
		if gn == "State" {
		}
		if gn != "Empty" {
			if d := flatten(s.Description); d != "" {
				b.WriteString("// " + d + "\n")
			}
			b.WriteString("type " + gn + " string\n\n")
		}
		return
	}
	if s.Type == "integer" || s.Type == "number" || s.Type == "boolean" {
		goType := "int64"
		switch {
		case s.Type == "number":
			goType = "float64"
		case s.Type == "boolean":
			goType = "bool"
		case s.Format == "int32":
			goType = "int32"
		}
		b.WriteString("type " + gn + " " + goType + "\n\n")
		return
	}
	if s.Type == "array" {
		item := "any"
		if s.Items != nil {
			item = g.primitive(s.Items.Type, map[string]any{"format": s.Items.Format})
		}
		b.WriteString("type " + gn + " []" + item + "\n\n")
		return
	}

	if d := flatten(s.Description); d != "" && gn != "Empty" {
		b.WriteString("// " + d + "\n")
	}
	isQuery := g.querySchemas[n] != nil
	// 先输出内联枚举类型（独立于 struct 声明之外）
	for _, fn := range sortedKeys(s.Properties) {
		p, _ := rawToProp(s.Properties[fn])
		if en, ok := p["enum"].([]any); ok && len(en) > 0 {
			var vals []string
			for _, e := range en {
				vals = append(vals, fmt.Sprintf("%v", e))
			}
			valSet := strings.Join(vals, ",")
			// 值集与既有枚举一致时直接复用，避免输出冗余类型。
			if existing := g.existingEnum(valSet); existing != "" {
				g.fieldEnum[gn+"."+fn] = existing
				continue
			}
			etype := enumTypeName(gn, fn, g.usedEnum)
			if _, taken := g.usedEnum[etype]; !taken {
				g.usedEnum[etype] = valSet
				g.emitEnum(b, etype, vals, enumFieldDescs(p))
			}
			g.fieldEnum[gn+"."+fn] = etype
		}
	}
	b.WriteString("type " + gn + " struct {\n")
	for _, fn := range sortedKeys(s.Properties) {
		p, _ := rawToProp(s.Properties[fn])
		g.emitField(b, gn, fn, p, isQuery)
	}
	b.WriteString("}\n\n")
}

func enumFieldDescs(p map[string]any) []string {
	if d, ok := p["description"].(string); ok {
		if _, descs := enumFromDesc(d); descs != nil {
			return descs
		}
	}
	return nil
}

func (g *generator) emitEnum(b *bytes.Buffer, typeName string, vals, descs []string) {
	b.WriteString("// " + typeName + " values\n")
	b.WriteString("type " + typeName + " string\n\nconst (\n")
	used := map[string]bool{}
	descByVal := enumDescsByValue(vals, descs)
	for _, v := range vals {
		cname := typeName + valueToCamel(v)
		if used[cname] {
			suffix := 1
			for used[fmt.Sprintf("%s_%d", cname, suffix)] {
				suffix++
			}
			cname = fmt.Sprintf("%s_%d", cname, suffix)
		}
		used[cname] = true
		line := fmt.Sprintf("\t%s %s = %q", cname, typeName, v)
		if d := descByVal[v]; d != "" {
			line += " // " + d
		}
		b.WriteString(line + "\n")
	}
	b.WriteString(")\n\n")
}

// isEnumSchema 判断 schema 是否为纯枚举类型（值为集合，无属性）。
func (g *generator) isEnumSchema(n string) bool {
	s, ok := g.schemas[n]
	if !ok || s == nil || s.Properties != nil {
		return false
	}
	return g.enumValuesOf(n) != nil
}

// enumValuesOf 返回 schema 的枚举值集合；非枚举 schema 返回 nil。
func (g *generator) enumValuesOf(n string) []string {
	s, ok := g.schemas[n]
	if !ok || s == nil || s.Properties != nil {
		return nil
	}
	vals, _ := g.extendedEnum(n)
	return vals
}

// extendedEnum 返回 schema 的 (枚举值, 描述) 列表。
// 描述集合缺枚举值时，将 schema 的 default 值补入（其为合法状态值；
// 某些文档不在描述中列出 default 值）。
func (g *generator) extendedEnum(n string) ([]string, []string) {
	s, ok := g.schemas[n]
	if !ok || s == nil || s.Properties != nil {
		return nil, nil
	}
	var vals, descs []string
	hasReal := !bytes.Equal(bytes.TrimSpace(s.Enum), []byte("null")) && len(s.Enum) > 0
	if hasReal {
		_ = json.Unmarshal(s.Enum, &vals)
		descs = descsMatching(s.Description, vals)
	} else if s.Type == "string" {
		if pv, pd := enumFromDesc(s.Description); pv != nil {
			vals, descs = pv, pd
		}
	}
	if len(vals) == 0 {
		return nil, nil
	}
	if s.Default != "" {
		found := -1
		for i, v := range vals {
			if v == s.Default {
				found = i
				break
			}
		}
		if found < 0 {
			vals = append(vals, s.Default)
			descs = append(descs, "未定义（默认值）")
		} else if descs[found] == "" {
			descs[found] = "未定义（默认值）"
		}
	}
	return vals, descs
}

// 查询参数/内联枚举与顶层枚举 schema 常常声明同一组值，复用类型以避免重复并保留注释描述。
func (g *generator) existingEnum(valSet string) string {
	want := sortSet(valSet)
	names := make([]string, 0, len(g.usedEnum))
	for name := range g.usedEnum {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if g.enumValSet[name] != "" && sortSet(g.enumValSet[name]) == want {
			return name
		}
	}
	return ""
}

// sortSet 将逗号分隔的值集排序后重新拼接，用于无序比较。
func sortSet(s string) string {
	parts := strings.Split(s, ",")
	sort.Strings(parts)
	return strings.Join(parts, ",")
}
func enumDescsByValue(vals, descs []string) map[string]string {
	m := map[string]string{}
	for i, v := range vals {
		if i < len(descs) && descs[i] != "" {
			m[v] = flatten(descs[i])
		}
	}
	return m
}

// descsMatching 从 description 提取枚举值描述，返回与 vals 顺序一致的描述数组。
func descsMatching(desc string, vals []string) []string {
	parsedVals, descs := enumFromDesc(desc)
	if parsedVals == nil {
		return nil
	}
	byVal := enumDescsByValue(parsedVals, descs)
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		out = append(out, byVal[v])
	}
	return out
}

func (g *generator) emitField(b *bytes.Buffer, structName, fieldName string, p map[string]any, isQuery bool) {
	goName := toGoName(fieldName)
	tagKind := "json"
	tagName := fieldName
	if isQuery {
		tagKind = "url"
	}
	comment := ""
	if d, ok := p["description"].(string); ok && d != "" {
		comment = " // " + flatten(d)
	}

	if ref, ok := p["$ref"].(string); ok {
		b.WriteString(fmt.Sprintf("\t%s %s `%s:\"%s\"`%s\n", goName, g.goName[sanitizeRef(ref)], tagKind, tagName, comment))
		return
	}
	if en, ok := p["enum"].([]any); ok && len(en) > 0 {
		var vals []string
		for _, e := range en {
			vals = append(vals, fmt.Sprintf("%v", e))
		}
		valSet := strings.Join(vals, ",")
		// 值集去重：若已有枚举类型携带完全相同的值集，直接复用（避免重复类型并继承注释）。
		if existing := g.existingEnum(valSet); existing != "" {
			b.WriteString(fmt.Sprintf("\t%s %s `%s:\"%s\"`%s\n", goName, existing, tagKind, tagName, comment))
			return
		}
		etype := g.fieldEnum[structName+"."+fieldName]
		if etype == "" {
			etype = enumTypeName(structName, fieldName, g.usedEnum)
		}
		if g.enumValSet[etype] != valSet {
			g.enumValSet[etype] = valSet
			if _, taken := g.usedEnum[etype]; !taken {
				g.usedEnum[etype] = valSet
				g.emitEnum(b, etype, vals, enumFieldDescs(p))
			}
		}
		b.WriteString(fmt.Sprintf("\t%s %s `%s:\"%s\"`%s\n", goName, etype, tagKind, tagName, comment))
		return
	}

	goType := g.fieldType(p)
	b.WriteString(fmt.Sprintf("\t%s %s `%s:\"%s\"`%s\n", goName, goType, tagKind, tagName, comment))
}

func (g *generator) fieldType(p map[string]any) string {
	if items, ok := p["items"].(map[string]any); ok {
		if ref, ok := items["$ref"].(string); ok {
			return "[]" + g.goName[sanitizeRef(ref)]
		}
		t, _ := items["type"].(string)
		return "[]" + g.primitive(t, items)
	}
	switch p["type"] {
	case "string":
		if p["format"] == "uint64" {
			return "transport.Uint64"
		}
		return "string"
	case "integer":
		if p["format"] == "int32" {
			return "int32"
		}
		if p["format"] == "int64" {
			return "int64"
		}
		return "int64"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		return "[]any"
	case "object":
		if add, ok := p["additionalProperties"].(map[string]any); ok {
			return g.mapOf(add)
		}
		return "map[string]any"
	default:
		return "any"
	}
}

func (g *generator) mapOf(add map[string]any) string {
	if ref, ok := add["$ref"].(string); ok {
		return "map[string]" + g.goName[sanitizeRef(ref)]
	}
	t, _ := add["type"].(string)
	return "map[string]" + g.primitive(t, add)
}

func (g *generator) primitive(t string, m map[string]any) string {
	switch t {
	case "string":
		if m != nil && m["format"] == "uint64" {
			return "transport.Uint64"
		}
		return "string"
	case "integer":
		return "int64"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	default:
		return "any"
	}
}

// --- service.go 生成 ---

func (g *generator) genService(ops []opSpec) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "package %s\n\n", ops[0].Dir)
	b.WriteString("import (\n\t\"context\"\n")
	hasPathArgs := false
	for _, op := range ops {
		if len(op.pathArgs) > 0 {
			hasPathArgs = true
			break
		}
	}
	if hasPathArgs {
		b.WriteString("\t\"strings\"\n")
	}
	fmt.Fprintf(&b, "\n\t\"%s/transport\"\n)\n\n", modulePath)
	b.WriteString("type Service struct{ Client *transport.Client }\n")

	for _, op := range ops {
		comment := flatten(op.Summary)
		if comment == "" {
			comment = op.Name
		}
		b.WriteString("\n// " + comment + "\n")

		fn := fmt.Sprintf("func (s *Service) %s(ctx context.Context", op.Name)
		for i, p := range op.pathArgs {
			fn += fmt.Sprintf(", %s %s", argName(p), g.pathArgType(op.pathFmt[i]))
		}
		if op.reqSchema != "" {
			fn += ", req *" + g.goName[op.reqSchema]
		}
		fn += ")"

		switch {
		case op.rawResp:
			fn += " ([]byte, error)"
		case op.respSchema != "":
			fn += " (*" + g.goName[op.respSchema] + ", error)"
		default:
			fn += " error"
		}
		b.WriteString(fn + " {\n")
		b.WriteString(fmt.Sprintf("\tpath := %q\n", op.Path))
		for i, p := range op.pathArgs {
			if g.pathArgType(op.pathFmt[i]) == "string" {
				b.WriteString(fmt.Sprintf("\tpath = strings.Replace(path, \"{%s}\", %s, 1)\n", p, argName(p)))
			} else {
				b.WriteString(fmt.Sprintf("\tpath = strings.Replace(path, \"{%s}\", %s.String(), 1)\n", p, argName(p)))
			}
		}

		reqArg := "nil"
		if op.reqSchema != "" {
			reqArg = "req"
		}
		switch {
		case op.rawResp:
			b.WriteString(fmt.Sprintf("\treturn s.Client.%s(ctx, path, %s)\n", transportRawCall(op.Method), reqArg))
		case op.respSchema != "":
			b.WriteString(fmt.Sprintf("\tvar resp %s\n", g.goName[op.respSchema]))
			b.WriteString(fmt.Sprintf("\tif err := s.Client.%s(ctx, path, %s, &resp); err != nil {\n", transportCall(op.Method), reqArg))
			b.WriteString("\t\treturn nil, err\n\t}\n\treturn &resp, nil\n")
		default:
			b.WriteString(fmt.Sprintf("\treturn s.Client.%s(ctx, path, %s, nil)\n", transportCall(op.Method), reqArg))
		}
		b.WriteString("}\n")
	}
	return b.String()
}

func transportCall(m string) string {
	switch m {
	case "GET":
		return "Get"
	case "PUT":
		return "Put"
	case "PATCH":
		return "Patch"
	default:
		return "Post"
	}
}

func transportRawCall(m string) string {
	if m == "POST" {
		return "PostRaw"
	}
	return "GetRaw"
}

func argName(p string) string {
	n := toGoName(p)
	if n == "UUID" {
		return "uUID"
	}
	return strings.ToLower(n[:1]) + n[1:]
}

func (g *generator) pathArgType(f string) string {
	if f == "uint64" {
		return "transport.Uint64"
	}
	return "string"
}

func pathFmtFor(op opSpec, name string) string {
	for i, p := range op.pathArgs {
		if p == name {
			return op.pathFmt[i]
		}
	}
	return ""
}

// --- service_test.go 生成 ---

func genTests(ops []opSpec, g *generator) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "package %s\n\n", ops[0].Dir)
	b.WriteString("import (\n\t\"context\"\n\t\"net/http\"\n\t\"testing\"\n")
	fmt.Fprintf(&b, "\n\t\"%s/transport\"\n)\n\n", modulePath)
	b.WriteString("var testCtx = context.Background()\n")

	for _, op := range ops {
		b.WriteString("\nfunc Test" + op.Name + "(t *testing.T) {\n")
		b.WriteString("\tcl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {\n")
		b.WriteString("\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n")
		b.WriteString("\t\tw.Write([]byte(\"{}\"))\n")
		b.WriteString("\t}, \"test-token\")\n")
		b.WriteString("\tdefer srv.Close()\n")
		b.WriteString("\tsvc := &Service{Client: cl}\n")

		args := "testCtx"
		for _, p := range op.pathArgs {
			if g.pathArgType(pathFmtFor(op, p)) == "string" {
				args += ", \"\""
			} else {
				args += ", transport.Uint64(0)"
			}
		}
		if op.reqSchema != "" {
			args += ", &" + g.goName[op.reqSchema] + "{}"
		}

		if op.rawResp || op.respSchema != "" {
			b.WriteString(fmt.Sprintf("\t_, err := svc.%s(%s)\n", op.Name, args))
		} else {
			b.WriteString(fmt.Sprintf("\terr := svc.%s(%s)\n", op.Name, args))
		}
		b.WriteString("\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n")
		b.WriteString("}\n")
	}
	return b.String()
}
