package gotabulate

import "fmt"
import "bytes"
import "math"

// Basic Structure of TableFormat
type TableFormat struct {
	LineTop         Line
	LineBelowHeader Line
	LineBetweenRows Line
	LineBottom      Line
	HeaderRow       Row
	DataRow         Row
	Padding         int
	HeaderHide      bool
	FitScreen       bool
}

// Represents a Line
type Line struct {
	begin string
	hline string
	sep   string
	end   string
}

// Represents a Row
type Row struct {
	begin string
	sep   string
	end   string
}

// Table Formats that are available to the user
// The user can define his own format, just by addind an entry to this map
// and calling it with Render function e.g t.Render("customFormat")
var TableFormats = map[string]TableFormat{
	"simple": TableFormat{
		LineTop:         Line{"", "-", "  ", ""},
		LineBelowHeader: Line{"", "-", "  ", ""},
		LineBottom:      Line{"", "-", "  ", ""},
		HeaderRow:       Row{"", "  ", ""},
		DataRow:         Row{"", "  ", ""},
		Padding:         1,
	},
	"plain": TableFormat{
		HeaderRow: Row{"", "  ", ""},
		DataRow:   Row{"", "  ", ""},
		Padding:   1,
	},
	"grid": TableFormat{
		LineTop:         Line{"+", "-", "+", "+"},
		LineBelowHeader: Line{"+", "=", "+", "+"},
		LineBetweenRows: Line{"+", "-", "+", "+"},
		LineBottom:      Line{"+", "-", "+", "+"},
		HeaderRow:       Row{"|", "|", "|"},
		DataRow:         Row{"|", "|", "|"},
		Padding:         1,
	},
}

// Minimum padding that will be applied
var MIN_PADDING = 5

// Main Tabulate structure
type Tabulate struct {
	Data        []*TabulateRow
	Headers     []string
	FloatFormat byte
	TableFormat TableFormat
	Align       string
	EmptyVar    string
	HideLines   []string
	MaxSize     int
	WrapStrings bool
}

// Represents normalized tabulate Row
type TabulateRow struct {
	Elements  []string
	Continuos bool
	Extra     bool
}

// Add padding to each cell
func (t *Tabulate) padRow(arr []string, padding int) []string {
	if len(arr) < 1 {
		return arr
	}
	padded := make([]string, len(arr))
	for index, el := range arr {
		var buffer bytes.Buffer
		// Pad left
		for i := 0; i < padding; i++ {
			buffer.WriteString(" ")
		}

		buffer.WriteString(el)

		// Pad Right
		for i := 0; i < padding; i++ {
			buffer.WriteString(" ")
		}

		padded[index] = buffer.String()
	}
	return padded
}

// Align right (Add padding left)
func (t *Tabulate) padLeft(width int, str string) string {
	var buffer bytes.Buffer
	// Pad left
	padding := width - len(str)
	for i := 0; i < padding; i++ {
		buffer.WriteString(" ")
	}
	buffer.WriteString(str)
	return buffer.String()
}

// Align Left (Add padding right)
func (t *Tabulate) padRight(width int, str string) string {
	var buffer bytes.Buffer
	padding := width - len(str)

	buffer.WriteString(str)

	// Add Padding right
	for i := 0; i < padding; i++ {
		buffer.WriteString(" ")
	}
	return buffer.String()
}

// Center the element in the cell
func (t *Tabulate) padCenter(width int, str string) string {
	var buffer bytes.Buffer
	length := len(str)

	padding := int(math.Ceil(float64((width - length)) / 2.0))

	// Add padding left
	for i := 0; i < padding; i++ {
		buffer.WriteString(" ")
	}

	// Write string
	buffer.WriteString(str)

	// Calculate how much space is left
	current := (width - len(buffer.String()))

	// Add padding right
	for i := 0; i < current; i++ {
		buffer.WriteString(" ")
	}
	return buffer.String()
}

// Build Line based on padded_widths from t.GetWidths()
func (t *Tabulate) buildLine(padded_widths []int, padding []int, l Line) string {
	cells := make([]string, len(padded_widths))

	for i, _ := range cells {
		var buffer bytes.Buffer
		for j := 0; j < padding[i]+MIN_PADDING; j++ {
			buffer.WriteString(l.hline)
		}
		cells[i] = buffer.String()
	}
	var buffer bytes.Buffer

	// Print begin
	buffer.WriteString(l.begin)

	// Print contents
	for i := 0; i < len(cells); i++ {
		if i != len(cells)-1 {
			buffer.WriteString(cells[i] + l.sep)
		} else {
			buffer.WriteString(cells[i])
		}
	}

	// Print end
	buffer.WriteString(l.end)
	return buffer.String()
}

// Build Row based on padded_widths from t.GetWidths()
func (t *Tabulate) buildRow(elements []string, padded_widths []int, paddings []int, d Row) string {

	var buffer bytes.Buffer
	buffer.WriteString(d.begin)
	padFunc := t.getAlignFunc()
	// Print contents
	for i := 0; i < len(padded_widths); i++ {
		output := ""
		if len(elements) > i {
			output = padFunc(padded_widths[i], elements[i])
		} else {
			output = padFunc(padded_widths[i], t.EmptyVar)
		}
		buffer.WriteString(output)
		if i != len(padded_widths)-1 {
			buffer.WriteString(d.sep)
		}
	}
	// Print end
	buffer.WriteString(d.end)

	return buffer.String()
}

// Render the data table
func (t *Tabulate) Render(format ...interface{}) string {
	var lines []string

	// If headers are set use them, otherwise pop the first row
	if len(t.Headers) < 1 {
		t.Headers, t.Data = t.Data[0].Elements, t.Data[1:]
	}

	// Use the format that was passed as parameter, otherwise
	// use the format defined in the struct
	if len(format) > 0 {
		t.TableFormat = TableFormats[format[0].(string)]
	}

	// If Wrap Strings is set to True,then break up the string to multiple cells
	if t.WrapStrings {
		t.Data = t.wrapCellData()
	}

	// Check if Data is present
	if len(t.Data) < 1 {
		panic("No Data specified")
	}

	if len(t.Headers) < len(t.Data[0].Elements) {
		diff := len(t.Data[0].Elements) - len(t.Headers)
		padded_header := make([]string, diff)
		for _, e := range t.Headers {
			padded_header = append(padded_header, e)
		}
		t.Headers = padded_header
	}

	// Get Column widths for all columns
	cols := t.getWidths(t.Headers, t.Data)

	padded_widths := make([]int, len(cols))
	for i, _ := range padded_widths {
		padded_widths[i] = cols[i] + MIN_PADDING*t.TableFormat.Padding
	}

	// Start appending lines

	// Append top line if not hidden
	if !inSlice("top", t.HideLines) {
		lines = append(lines, t.buildLine(padded_widths, cols, t.TableFormat.LineTop))
	}

	// Add Header
	lines = append(lines, t.buildRow(t.padRow(t.Headers, t.TableFormat.Padding), padded_widths, cols, t.TableFormat.HeaderRow))

	// Add Line Below Header if not hidden
	if !inSlice("belowheader", t.HideLines) {
		lines = append(lines, t.buildLine(padded_widths, cols, t.TableFormat.LineBelowHeader))
	}

	// Add Data Rows
	for index, element := range t.Data {
		lines = append(lines, t.buildRow(t.padRow(element.Elements, t.TableFormat.Padding), padded_widths, cols, t.TableFormat.DataRow))
		if index < len(t.Data)-1 {
			if element.Continuos != true {
				lines = append(lines, t.buildLine(padded_widths, cols, t.TableFormat.LineBetweenRows))
			}
		}
	}

	if !inSlice("bottomLine", t.HideLines) {
		lines = append(lines, t.buildLine(padded_widths, cols, t.TableFormat.LineBottom))
	}

	// Join lines
	var buffer bytes.Buffer
	for _, line := range lines {
		buffer.WriteString(line + "\n")
	}

	return buffer.String()
}

// Calculate the max column width for each element
func (t *Tabulate) getWidths(headers []string, data []*TabulateRow) []int {
	widths := make([]int, len(headers))
	current_max := len(t.EmptyVar)
	for i := 0; i < len(headers); i++ {
		current_max = len(headers[i])
		for _, item := range data {
			if len(item.Elements) > i && len(widths) > i {
				element := item.Elements[i]
				if len(element) > current_max {
					widths[i] = len(element)
					current_max = len(element)
				} else {
					widths[i] = current_max
				}
			}
		}
	}

	return widths
}

// Set Headers of the table
// If Headers count is less than the data row count, the headers will be padded to the right
func (t *Tabulate) SetHeaders(headers []string) *Tabulate {
	t.Headers = headers
	return t
}

// Set Float Formatting
// will be used in strconv.FormatFloat(element, format, -1, 64)
func (t *Tabulate) SetFloatFormat(format byte) *Tabulate {
	t.FloatFormat = format
	return t
}

// Set Align Type, Available options: left, right, center
func (t *Tabulate) SetAlign(align string) {
	t.Align = align
}

// Select the padding function based on the align type
func (t *Tabulate) getAlignFunc() func(int, string) string {
	if len(t.Align) < 1 || t.Align == "right" {
		return t.padLeft
	} else if t.Align == "left" {
		return t.padRight
	} else {
		return t.padCenter
	}
}

// Set how an empty cell will be represented
func (t *Tabulate) SetEmptyString(empty string) {
	t.EmptyVar = empty + " "
}

// Set which lines to hide.
// Can be:
// top - Top line of the table,
// belowheader - Line below the header,
// bottom - Bottom line of the table
func (t *Tabulate) SetHideLines(hide []string) {
	t.HideLines = hide
}

func (t *Tabulate) SetWrapStrings(wrap bool) {
	t.WrapStrings = wrap
}

// Sets the maximum size of cell
// If WrapStrings is set to true, then the string inside
// the cell will be split up into multiple cell
func (t *Tabulate) SetMaxCellSize(max int) {
	t.MaxSize = max
}

// If string size is larger than t.MaxSize, then split it to multiple cells (downwards)
func (t *Tabulate) wrapCellData() []*TabulateRow {
	var arr []*TabulateRow
	next := t.Data[0]
	for index := 0; index <= len(t.Data); index++ {
		elements := next.Elements
		new_elements := make([]string, len(elements))

		for i, e := range elements {
			if len(e) > t.MaxSize {
				new_elements[i] = e[t.MaxSize:]
				elements[i] = e[:t.MaxSize]
				next.Continuos = true
			}
		}

		if next.Continuos {
			arr = append(arr, next)
			next = &TabulateRow{Elements: new_elements, Extra: true}
			index--
		} else if next.Extra && index+1 < len(t.Data) {
			arr = append(arr, next)
			next = t.Data[index+1]
		} else if index+1 < len(t.Data) {
			arr = append(arr, next)
			next = t.Data[index+1]
		} else if index >= len(t.Data) {
			arr = append(arr, next)
		}

	}
	return arr
}

// Create a new Tabulate Object
// Accepts 2D String Array, 2D Int Array, 2D Int64 Array,
// 2D Bool Array, 2D Float64 Array, 2D interface{} Array,
// Map map[strig]string, Map map[string]interface{},
func Create(data interface{}) *Tabulate {
	t := &Tabulate{FloatFormat: 'f', MaxSize: 30}

	switch v := data.(type) {
	case [][]string:
		t.Data = createFromString(data.([][]string))
	case [][]int32:
		t.Data = createFromInt32(data.([][]int32))
	case [][]int64:
		t.Data = createFromInt64(data.([][]int64))
	case [][]int:
		t.Data = createFromInt(data.([][]int))
	case [][]bool:
		t.Data = createFromBool(data.([][]bool))
	case [][]float64:
		t.Data = createFromFloat64(data.([][]float64), t.FloatFormat)
	case [][]interface{}:
		t.Data = createFromMixed(data.([][]interface{}), t.FloatFormat)
	case []string:
		t.Data = createFromString([][]string{data.([]string)})
	case []interface{}:
		t.Data = createFromMixed([][]interface{}{data.([]interface{})}, t.FloatFormat)
	case map[string][]interface{}:
		t.Headers, t.Data = createFromMapMixed(data.(map[string][]interface{}), t.FloatFormat)
	case map[string][]string:
		t.Headers, t.Data = createFromMapString(data.(map[string][]string))
	default:
		fmt.Println(v)
	}

	return t
}
