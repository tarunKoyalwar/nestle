package nestle

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

/*
A Wrapper around regexp that supports nested groups.
A hacky but reliable way to find nested groups using regex.
Used with syntax `matcher_before_begin[start_delim:nested:end_delim]`
ex :  json:{xx} can be parsed from  json{a:b,b:{d:e}} using json[{:nested:}]
*/
type Nestle struct {
	// Does not support `[(:nested:)]` i.e cases with no pre or post match
	StartDelim rune           // Single char delimeter
	EndDelim   rune           // Single Char delimeter
	Prematch   string         // Match string before start delim
	PostMatch  string         // Match String after start delim
	pre_re     *regexp.Regexp // Pre Match regex
	post_re    *regexp.Regexp // Post Match regex

}

// MustCompile : Syntax to use nested matching
// `matcher_before_begin[start_delim:nested:end_delim]`
func MustCompile(regex string) (*Nestle, error) {
	if !strings.Contains(regex, ":nested:") {
		return nil, fmt.Errorf("regex does not contain `:nested:` matching. Try regex.Regexp ")
	}

	n := Nestle{}

	darr := strings.Split(regex, ":nested:")

	if len(darr) != 2 {
		return nil, fmt.Errorf("messed up syntax")
	}

	prefix := darr[0]
	suffix := darr[1]

	var begin string

	for i := len(prefix) - 1; i >= 0; i-- {
		if prefix[i] != '[' {
			// Everything until [ in reverse is delimeter
			begin = string(prefix[i]) + begin
		} else if prefix[i] == '[' {
			// Everything before [ is prematch regex
			n.Prematch = prefix[:i]
			break
		}
	}

	n.StartDelim = rune(begin[0])

	var end strings.Builder

	for i := 0; i < len(suffix); i++ {
		if suffix[i] != ']' {
			end.WriteByte(suffix[i])
		} else if suffix[i] == ']' {
			n.PostMatch = suffix[i+1:]
			break
		}
	}

	n.EndDelim = rune(end.String()[0])
	if n.StartDelim == 0 {
		return nil, fmt.Errorf("delim cannot be empty")
	}

	if n.EndDelim == 0 {
		return nil, fmt.Errorf("delim cannot be empty")
	}

	// create compile groups
	n.compile()

	return &n, nil

}

// compile : compile regex using regex.MustCompile
func (n *Nestle) compile() {

	n.pre_re = regexp.MustCompile(n.Prematch + string(n.StartDelim))

	n.post_re = regexp.MustCompile(string(n.EndDelim) + n.PostMatch)

}

// FindAllStringIndex
func (n *Nestle) FindAllStringIndex(data string) [][]int {
	/*
		Algorithm:
		1. Find indexes using prematch and save in a map
		2. Find indexes using postmatch and save in a map
		3. Using algo similar to Expression Evaluation using stack and check for consistency
		b/w prematch and postmatch indexes using stack
	*/

	// [][]int{}

	var start_indexes, end_indexes, final [][]int

	if n.Prematch != "" && n.PostMatch != "" {
		start_indexes = n.pre_re.FindAllStringIndex(data, -1)

		end_indexes = n.post_re.FindAllStringIndex(data, -1)
	} else {
		if n.Prematch == "" && n.PostMatch == "" {
			panic("finding nested groups {{{{}}}} with no pre/post match is bad idea and is not supported")
		} else {
			if n.Prematch == "" {
				end_indexes = n.post_re.FindAllStringIndex(data, -1)
				start_indexes = make([][]int, 0)

			} else if n.PostMatch == "" {
				start_indexes = n.pre_re.FindAllStringIndex(data, -1)
				end_indexes = make([][]int, 0)
			}
		}
	}

	// condition if both prematch and postmatch are not empty

	// contains end indexes of prematch group
	start_arr := []int{}

	for _, v := range start_indexes {
		start_arr = append(start_arr, v[0])
	}

	// contains start indexes of postmatch group
	end_arr := []int{}

	for _, v := range end_indexes {
		end_arr = append(end_arr, v[1])
	}

	// As a precaution sort in ascending order
	sort.Ints(start_arr)

	sort.Sort(sort.Reverse(sort.IntSlice(end_arr)))

	var estack stack = []int{}

	estack = append(estack, end_arr...)

	if len(start_arr) != 0 {
		final = n.forwardPropogration(data, estack, start_arr)
	} else {
		final = n.backPropogation(data, estack, end_indexes)
	}

	return final

}

// FindAllString
func (n *Nestle) FindAllString(data string) []string {
	res := []string{}

	arr := n.FindAllStringIndex(data)

	if len(arr) == 0 {
		return res
	}

	for _, v := range arr {
		res = append(res, data[v[0]:v[1]])
	}

	return res

}

// forwardPropogration for cases where regex has prematch
func (n *Nestle) forwardPropogration(data string, end stack, begin []int) [][]int {
	final := [][]int{}
	for _, v := range begin {

		if len(end) == 0 {
			// No postmatch
			x, last := n.validate(data, v, -1)
			if x {
				final = append(final, []int{v, last + 1})
			}
		} else {
			ok := true

			for ok {
				var val int
				// start for each prematch
				val, ok = end.GetLastElement()
				if ok {
					if val <= v {
						// remove element if it is smaller than start
						end.Pop()
					} else {
						// validate/check for continuity
						if x, _ := n.validate(data, v, val); x {
							final = append(final, []int{v, val + 1})
							ok = false
						}
					}
				}
			}
		}

	}

	return final
}

// validate forwardProgogration conditions
func (n *Nestle) validate(data string, start int, end int) (bool, int) {

	// cosidering delimeters are of 1 character

	OpenCount := 0
	OpenAtPos := -1

	finalindex := end
	if end == -1 {
		finalindex = len(data)
	}

	for i := start; i <= finalindex; i++ {
		if data[i] == byte(n.StartDelim) {
			// If symbol/delimeter is not open and is first
			if OpenCount == 0 && OpenAtPos == -1 {
				// open and increment
				OpenAtPos = i
				OpenCount += 1
			} else if OpenCount != 0 && OpenAtPos != -1 {
				// symbol is already  open
				// just increment
				OpenCount += 1
			} else {
				// i.e opencount is zero but OpenAtpos is not -1
				// malformed but not explored
				return false, -1
			}
		} else if data[i] == byte(n.EndDelim) {
			// fmt.Println(OpenCount)
			// If end delim is encountered
			if OpenCount == 0 {
				// this means {}}
				// such cases cannot be balanced
				return false, -1
			} else if OpenCount == 1 {
				// this means {}{
				return true, i
			} else {
				// lot of them are open
				OpenCount -= 1
			}
		} else {
			// skip until condition
			continue
		}
	}

	return false, -1
}

// backPropogation for cases where regex has postmatch with no prematch
func (n *Nestle) backPropogation(data string, end stack, end_arr [][]int) [][]int {
	final := [][]int{}

	lookup := map[int]int{}

	for _, v := range end_arr {
		lookup[v[1]] = v[0]
	}

	// endstack is leftmost match value
	// lookup is rightmost postmatch value

	for i := 0; i < len(end); i++ {

		var index int

		if i != len(end)-1 {
			index = n.reverseValidate(data, end[i], end[i+1])

		} else {
			index = n.reverseValidate(data, end[i], 0)
		}

		if index != -1 {
			final = append(final, []int{index, end[i] + 1})
		}
	}

	return final

}

// reverseValidate backPropogation conditions
func (n *Nestle) reverseValidate(data string, end int, start int) int {

	totalClosed := 0

	closedAt := -1

	for i := end; i >= start; i-- {
		//reverse
		if data[i] == byte(n.EndDelim) {

			if totalClosed == 0 && closedAt == -1 {
				totalClosed += 1
				closedAt = i
			} else {
				totalClosed += 1
			}

		} else if data[i] == byte(n.StartDelim) {
			if totalClosed == 0 {
				// fmt.Printf("malformed conditon got { before } while parsing in reverse")
				return -1
			} else if totalClosed == 1 && totalClosed != -1 {
				return i
			} else {
				totalClosed -= 1
			}
		}
	}

	return 0
}
