package crossword

import "testing"

func TestSatisfiesAtPos(t *testing.T) {
	cases := []struct {
		name string
		expr string
		r    rune
		pos  int
		want bool
	}{
		{"simple pos 0", "ab", 'a', 0, true},
		{"simple pos 1", "ab", 'a', 1, false},
		{"or a", "a|b", 'a', 0, true},
		{"or b", "a|b", 'b', 0, true},
		{"brackets", "[abc]", 'c', 0, true},
		{"plus", "[abc]+d", 'd', 0, false},
		{"star true", "[abc]*d", 'd', 0, true},
		{"star false", "[abc]*d", 'e', 0, false},
		{"star exceeding length", "[a]*d", 'a', 3, true},
		{"parens", "(a)", 'a', 0, true},
		{"parens (a|b) a at 0", "(a|b)", 'a', 0, true},
		{"parens (a|b) a at 0", "(a|b)", 'a', 0, true},
		{"parens (a|b) a at 1", "(a|b)", 'a', 1, false},
		{"parens (a|b) b at 0", "(a|b)", 'b', 0, true},
		{"optional a?b: a at 0", "a?b", 'a', 0, true},
		{"optional a?b: b at 0", "a?b", 'b', 0, true},
		{"optional a?b: a at 1", "a?b", 'a', 1, false},

		{"(NA|FE|HE)[CV]: 'N' at 0", "(NA|FE|HE)[CV]", 'N', 0, true},
		{"(NA|FE|HE)[CV]: 'A' at 1", "(NA|FE|HE)[CV]", 'A', 1, true},
		{"(NA|FE|HE)[CV]: 'C' at 3", "(NA|FE|HE)[CV]", 'C', 2, true},
		{"(NA|FE|HE)[CV]: 'E' at 3", "(NA|FE|HE)[CV]", 'E', 2, false},

		{"EP|IP|EF: 'E' at 0", "EP|IP|EF", 'E', 0, true},
		{"HE|LL|O+: 'E' at 1", "HE|LL|O+", 'E', 1, true},
		{"[*]+: '*' at 0", "[*]+", '*', 0, true},
		{"[*]+: '*' at 1", "[*]+", '*', 1, true},
		{".?.+: '*' at 0", `.?.+`, '*', 0, true},
		{".?.+: '/' at 1", `.?.+`, '/', 1, true},
		{"[BORF].: 'L' at 1", `[BORF].`, 'L', 1, true},
	}
	for _, tt := range cases {
		got := satisfiesAtPos(tt.expr, tt.r, tt.pos)
		if got != tt.want {
			t.Errorf("%q case: got %t, want %t", tt.name, got, tt.want)
		}
	}
}

func TestSolve(t *testing.T) {
	cases := []struct {
		name string
		rows []string
		cols []string
		want string
	}{
		{
			name: "Beginner: Beatles",
			rows: []string{"HE|LL|O+", "[PLEASE]+"},
			cols: []string{"[^SPEAK]+", "EP|IP|EF"},
			want: "HELP",
		},
		{
			// this fails because Go does not support backreferences
			// in regex
			name: "Beginner Naughty",
			rows: []string{".*M?O.*", "(AN|FE|BE)"},
			cols: []string{`(A|B|C)\1`, `(AB|OE|SK)`},
			want: "",
		},
		{
			name: "Beginner: Symbolism",
			rows: []string{"[*]+", "/+"},
			cols: []string{`.?.+`, `.+`},
			want: "**//",
		},
		{
			name: "Beginner: Airstrip one",
			rows: []string{`18|19|20`, `[6789]\d`},
			cols: []string{`\d[2480]`, `56|94|73`},
			want: "1984",
		},
		{
			name: "Intermediate: Always remember",
			rows: []string{`[NOTAD]*`, `WEL|BAL|EAR`},
			cols: []string{`UB|IE|AW`, `[TUBE]*`, `[BORF].`},
			want: "ATOWEL",
		},
		{
			name: "Intermediate: Johnny",
			rows: []string{`[AWE]+`, `[ALP]+K`, `(PR|ER|EP)`},
			cols: []string{`[BQW](PR|LE)`, `[RANK]+`},
			want: "WALKER",
		},
		{
			// this fails because Go does not support backreferences
			// in regex
			name: "Intermediate: Earth",
			rows: []string{`CAT|FOR|FAT`, `RY|TY\-`, `[TOWEL]*`},
			cols: []string{`.(.)\1`, `.*[WAY]+`, `[RAM].[OH]`},
			want: "",
		},
		{
			name: "Intermediate: Encyclopedia",
			rows: []string{`[DEF][MNO]*`, `[^DJNU]P[ABC]`, `[ICAN]*`},
			cols: []string{`[JUNDT]*`, `APA|OPI|OLK`, `(NA|FE|HE)[CV]`},
			want: "DONTPANIC",
		},
		{
			// this fails because Go does not support backreferences
			// in regex
			name: "Intermediate: Technology",
			rows: []string{`[RUNT]*`, `O.*[HAT]`, `(.)*DO\1`},
			cols: []string{`[^NRU](NO|ON)`, `(D|FU|UF)+`, `(FO|A|R)*`, `(N|A)*`},
			want: "",
		},
		{
			name: "Palindromeda: Third",
			rows: []string{`(L|E|D|G|Y)*`, `(A|E|J)*Y.*`, `[FLEDG]*`},
			cols: []string{`(GE|PE)[AL]*`, `[EAF]+(YE|AB)*`, `(QR|LE)(G|I|M|P)`},
			want: "GELEYELEG",
		},
		{
			name: "Palindromeda: Horn",
			rows: []string{`[TRASH]*`, `(FA|AB)[TUP]*`, `(BA|TH|TU)*`, `.*A.*`},
			cols: []string{`(TS|RA|QA)*`, `(AB|UT|AR)*`, `(K|T)U.*(A|R)`, `(AR|FS|ST)+`},
			want: "RATSABUTTUBASTAR",
		},
		{
			name: "Palindromeda: Time walker",
			rows: []string{`(EP|ST)*`, `T[A-Z]*`, `.M.T`, `.*P.[S-X]+`},
			cols: []string{`.*E.*`, `[^P]I(IT|ME)`, `(EM|FE)(IT|IP)`, `(TS|PE|KE)*`},
			want: "STEPTIMEEMITPETS",
		},
	}
	for _, tt := range cases {
		got, err := Solve(tt.rows, tt.cols)
		if err != nil {
			t.Errorf("%q case error: %v", tt.name, err)
		}
		if got != tt.want {
			t.Errorf("%q case: got %q, want %q", tt.name, got, tt.want)
		}
	}
}
