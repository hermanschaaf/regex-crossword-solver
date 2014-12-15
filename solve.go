package crossword

import (
	"errors"
	"log"
	"regexp"
	"regexp/syntax"
)

type set []int

var flags = syntax.MatchNL | syntax.PerlX | syntax.UnicodeGroups

func (s *set) Add(v int) {
	// inefficient O(N), but fine for now
	for i := range *s {
		if (*s)[i] == v {
			return
		}
	}
	*s = append(*s, v)
}

func satisfiesAtPos(expr string, r rune, pos int) bool {
	re, err := syntax.Parse(expr, flags)
	if err != nil {
		panic(err)
	}

	re = re.Simplify()
	prog, err := syntax.Compile(re)
	if err != nil {
		panic(err)
	}
	// log.Println(prog)

	var el int
	s := 0                          // current step
	step := []int{}                 // states in this step
	queue := set([]int{prog.Start}) // queue of states to check

	for s <= pos {
		for len(queue) > 0 {
			el, queue = queue[0], queue[1:]
			inst := prog.Inst[el]
			switch inst.Op {
			case syntax.InstAlt:
				// alternative, so add both nodes to queue
				queue.Add(int(inst.Out))
				queue.Add(int(inst.Arg))
			case syntax.InstCapture:
				// capture group, just add next node to queue
				queue.Add(int(inst.Out))
			case syntax.InstRune, syntax.InstRune1, syntax.InstRuneAny:
				// rune, add instance to step
				step = append(step, el)
			}
		}
		// fmt.Println(s, step)
		for i := range step {
			inst := prog.Inst[step[i]]
			if s == pos && inst.MatchRune(r) {
				return true
			}
			// add rune's out to queue
			queue = append(queue, int(inst.Out))
		}
		step = []int{}
		s++
	}

	return false
}

func compileRegex(expr []string) ([]*regexp.Regexp, error) {
	res := []*regexp.Regexp{}
	// compile regex for rows
	for _, row := range expr {
		re, err := regexp.Compile("^" + row + "$")
		if err != nil {
			return res, err
		}
		res = append(res, re)
	}
	return res, nil
}

// Solve solves the given crossword.
func Solve(rows, cols []string) (string, error) {
	rowRe, err := compileRegex(rows)
	if err != nil {
		log.Println(err)
		return "", err
	}

	colRe, err := compileRegex(cols)
	if err != nil {
		return "", err
	}

	var start, end = 32, 90

	cells := len(cols) * len(rows)
	solution := make([]rune, cells)
	for i := range solution {
		solution[i] = rune(start)
	}

	for i := 0; i < cells; {
		r := i / len(cols)
		c := i % len(cols)

		solvedCell := false
		// find a rune that satisfies both regexes
		// for the current cell
	iterate:
		for u := int(solution[i]) + 1; u < end; u++ {
			rn := rune(u)
			// fmt.Println(string(solution[:i]) + string(rn))
			colOk := satisfiesAtPos(cols[c], rn, r)
			rowOk := satisfiesAtPos(rows[r], rn, c)
			if colOk && rowOk {
				solution[i] = rn
				solvedCell = true

				// if we're at the end of a row or col, check that
				// the regex is fully satisfied
				if c == len(cols)-1 {
					// end of a row
					row := solution[r*len(cols) : i+1]
					if rowRe[r].MatchString(string(row)) {
						solvedCell = true
					} else {
						solvedCell = false
						continue iterate
					}
				}
				if r == len(rows)-1 {
					// end of a col
					col := []rune{}
					for j := 0; j < len(rows); j++ {
						col = append(col, solution[j*len(cols)+c])
					}
					if colRe[c].MatchString(string(col)) {
						solvedCell = true
					} else {
						solvedCell = false
						continue iterate
					}
				}
			}
			if solvedCell {
				break iterate
			}
		}

		// fmt.Println(string(solution))
		if !solvedCell {
			// backtrack
			solution[i] = rune(start)
			i--
			if i < 0 {
				return "", errors.New("No solution")
			}
		} else {
			i++
		}
	}
	return string(solution), nil
}
