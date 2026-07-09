package sort

import (
	"testing"

	"github.com/nanoteck137/tunebook/tools/query"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, sort *Sort)
	}{
		{
			name:  "empty input",
			input: "",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 0 {
					t.Errorf("expected 0 orderings, got %d", len(sort.Orderings))
				}
			},
		},
		{
			name:  "whitespace only",
			input: "   ",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 0 {
					t.Errorf("expected 0 orderings, got %d", len(sort.Orderings))
				}
			},
		},
		{
			name:  "random mode",
			input: "random",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				if _, ok := sort.Orderings[0].(*query.RandomOrdering); !ok {
					t.Errorf("expected RandomOrdering, got %T", sort.Orderings[0])
				}
			},
		},
		{
			name:  "random mode uppercase",
			input: "RANDOM",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				if _, ok := sort.Orderings[0].(*query.RandomOrdering); !ok {
					t.Errorf("expected RandomOrdering, got %T", sort.Orderings[0])
				}
			},
		},
		{
			name:  "score mode",
			input: "score",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				if _, ok := sort.Orderings[0].(*query.ScoreOrdering); !ok {
					t.Errorf("expected ScoreOrdering, got %T", sort.Orderings[0])
				}
			},
		},
		{
			name:  "recent mode",
			input: "recent",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				fo, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo.Field.Name != "created" {
					t.Errorf("expected field 'created', got %q", fo.Field.Name)
				}
				if fo.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo.Dir)
				}
			},
		},
		{
			name:  "shuffle with seed",
			input: "shuffle=1234",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				so, ok := sort.Orderings[0].(*query.ShuffleOrdering)
				if !ok {
					t.Fatalf("expected ShuffleOrdering, got %T", sort.Orderings[0])
				}
				if so.Seed != 1234 {
					t.Errorf("expected seed 1234, got %d", so.Seed)
				}
			},
		},
		{
			name:  "shuffle uppercase",
			input: "SHUFFLE=5678",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				so, ok := sort.Orderings[0].(*query.ShuffleOrdering)
				if !ok {
					t.Fatalf("expected ShuffleOrdering, got %T", sort.Orderings[0])
				}
				if so.Seed != 5678 {
					t.Errorf("expected seed 5678, got %d", so.Seed)
				}
			},
		},
		{
			name:    "shuffle with invalid seed",
			input:   "shuffle=abc",
			wantErr: true,
		},
		{
			name:  "single field ascending with +",
			input: "+name",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				fo, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo.Field.Name != "name" {
					t.Errorf("expected field 'name', got %q", fo.Field.Name)
				}
				if fo.Dir != query.DirAsc {
					t.Errorf("expected DirAsc, got %v", fo.Dir)
				}
			},
		},
		{
			name:  "single field descending with -",
			input: "-year",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				fo, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo.Field.Name != "year" {
					t.Errorf("expected field 'year', got %q", fo.Field.Name)
				}
				if fo.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo.Dir)
				}
			},
		},
		{
			name:  "multiple fields with +/- syntax",
			input: "+artist,-year",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 2 {
					t.Fatalf("expected 2 orderings, got %d", len(sort.Orderings))
				}
				
				fo1, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo1.Field.Name != "artist" {
					t.Errorf("expected field 'artist', got %q", fo1.Field.Name)
				}
				if fo1.Dir != query.DirAsc {
					t.Errorf("expected DirAsc, got %v", fo1.Dir)
				}
				
				fo2, ok := sort.Orderings[1].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[1])
				}
				if fo2.Field.Name != "year" {
					t.Errorf("expected field 'year', got %q", fo2.Field.Name)
				}
				if fo2.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo2.Dir)
				}
			},
		},
		{
			name:  "field with asc suffix",
			input: "name asc",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				fo, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo.Field.Name != "name" {
					t.Errorf("expected field 'name', got %q", fo.Field.Name)
				}
				if fo.Dir != query.DirAsc {
					t.Errorf("expected DirAsc, got %v", fo.Dir)
				}
			},
		},
		{
			name:  "field with desc suffix",
			input: "year desc",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				fo, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo.Field.Name != "year" {
					t.Errorf("expected field 'year', got %q", fo.Field.Name)
				}
				if fo.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo.Dir)
				}
			},
		},
		{
			name:  "multiple fields with asc/desc syntax",
			input: "artist asc, year desc",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 2 {
					t.Fatalf("expected 2 orderings, got %d", len(sort.Orderings))
				}
				
				fo1, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo1.Field.Name != "artist" {
					t.Errorf("expected field 'artist', got %q", fo1.Field.Name)
				}
				if fo1.Dir != query.DirAsc {
					t.Errorf("expected DirAsc, got %v", fo1.Dir)
				}
				
				fo2, ok := sort.Orderings[1].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[1])
				}
				if fo2.Field.Name != "year" {
					t.Errorf("expected field 'year', got %q", fo2.Field.Name)
				}
				if fo2.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo2.Dir)
				}
			},
		},
		{
			name:  "field without direction defaults to asc",
			input: "name",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				fo, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo.Field.Name != "name" {
					t.Errorf("expected field 'name', got %q", fo.Field.Name)
				}
				if fo.Dir != query.DirAsc {
					t.Errorf("expected DirAsc, got %v", fo.Dir)
				}
			},
		},
		{
			name:  "mixed syntax",
			input: "+artist, year desc, -duration",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 3 {
					t.Fatalf("expected 3 orderings, got %d", len(sort.Orderings))
				}
				
				fo1, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo1.Field.Name != "artist" {
					t.Errorf("expected field 'artist', got %q", fo1.Field.Name)
				}
				if fo1.Dir != query.DirAsc {
					t.Errorf("expected DirAsc, got %v", fo1.Dir)
				}
				
				fo2, ok := sort.Orderings[1].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[1])
				}
				if fo2.Field.Name != "year" {
					t.Errorf("expected field 'year', got %q", fo2.Field.Name)
				}
				if fo2.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo2.Dir)
				}
				
				fo3, ok := sort.Orderings[2].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[2])
				}
				if fo3.Field.Name != "duration" {
					t.Errorf("expected field 'duration', got %q", fo3.Field.Name)
				}
				if fo3.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo3.Dir)
				}
			},
		},
		{
			name:  "field with nulls first",
			input: "name nulls first",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				fo, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo.Field.Name != "name" {
					t.Errorf("expected field 'name', got %q", fo.Field.Name)
				}
				if fo.Dir != query.DirAsc {
					t.Errorf("expected DirAsc, got %v", fo.Dir)
				}
				if fo.NullOrder != query.NullOrderingFirst {
					t.Errorf("expected NullOrderingFirst, got %v", fo.NullOrder)
				}
			},
		},
		{
			name:  "field with nulls last",
			input: "year desc nulls last",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 1 {
					t.Fatalf("expected 1 ordering, got %d", len(sort.Orderings))
				}
				fo, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo.Field.Name != "year" {
					t.Errorf("expected field 'year', got %q", fo.Field.Name)
				}
				if fo.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo.Dir)
				}
				if fo.NullOrder != query.NullOrderingLast {
					t.Errorf("expected NullOrderingLast, got %v", fo.NullOrder)
				}
			},
		},
		{
			name:  "multiple fields with null ordering",
			input: "+artist nulls first, -year nulls last",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 2 {
					t.Fatalf("expected 2 orderings, got %d", len(sort.Orderings))
				}
				
				fo1, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo1.Field.Name != "artist" {
					t.Errorf("expected field 'artist', got %q", fo1.Field.Name)
				}
				if fo1.Dir != query.DirAsc {
					t.Errorf("expected DirAsc, got %v", fo1.Dir)
				}
				if fo1.NullOrder != query.NullOrderingFirst {
					t.Errorf("expected NullOrderingFirst, got %v", fo1.NullOrder)
				}
				
				fo2, ok := sort.Orderings[1].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[1])
				}
				if fo2.Field.Name != "year" {
					t.Errorf("expected field 'year', got %q", fo2.Field.Name)
				}
				if fo2.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo2.Dir)
				}
				if fo2.NullOrder != query.NullOrderingLast {
					t.Errorf("expected NullOrderingLast, got %v", fo2.NullOrder)
				}
			},
		},
		{
			name:  "mixed null ordering syntax",
			input: "name asc nulls first, year desc nulls last, duration nulls first",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 3 {
					t.Fatalf("expected 3 orderings, got %d", len(sort.Orderings))
				}
				
				fo1, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo1.NullOrder != query.NullOrderingFirst {
					t.Errorf("expected NullOrderingFirst for first field, got %v", fo1.NullOrder)
				}
				
				fo2, ok := sort.Orderings[1].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[1])
				}
				if fo2.NullOrder != query.NullOrderingLast {
					t.Errorf("expected NullOrderingLast for second field, got %v", fo2.NullOrder)
				}
				
				fo3, ok := sort.Orderings[2].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[2])
				}
				if fo3.NullOrder != query.NullOrderingFirst {
					t.Errorf("expected NullOrderingFirst for third field, got %v", fo3.NullOrder)
				}
			},
		},
		{
			name:  "null ordering case insensitive",
			input: "name NULLS FIRST, year NULLS LAST",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 2 {
					t.Fatalf("expected 2 orderings, got %d", len(sort.Orderings))
				}
				
				fo1, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo1.NullOrder != query.NullOrderingFirst {
					t.Errorf("expected NullOrderingFirst, got %v", fo1.NullOrder)
				}
				
				fo2, ok := sort.Orderings[1].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[1])
				}
				if fo2.NullOrder != query.NullOrderingLast {
					t.Errorf("expected NullOrderingLast, got %v", fo2.NullOrder)
				}
			},
		},
		{
			name:    "empty field after +",
			input:   "+",
			wantErr: true,
		},
		{
			name:    "empty field after -",
			input:   "-",
			wantErr: true,
		},
		{
			name:    "empty field before asc",
			input:   " asc",
			wantErr: true,
		},
		{
			name:    "empty field before desc",
			input:   " desc",
			wantErr: true,
		},
		{
			name:  "whitespace handling",
			input: "  +artist  ,  -year  ",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 2 {
					t.Fatalf("expected 2 orderings, got %d", len(sort.Orderings))
				}
				
				fo1, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo1.Field.Name != "artist" {
					t.Errorf("expected field 'artist', got %q", fo1.Field.Name)
				}
				
				fo2, ok := sort.Orderings[1].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[1])
				}
				if fo2.Field.Name != "year" {
					t.Errorf("expected field 'year', got %q", fo2.Field.Name)
				}
			},
		},
		{
			name:  "case insensitive asc/desc",
			input: "name ASC, year DESC",
			check: func(t *testing.T, sort *Sort) {
				if len(sort.Orderings) != 2 {
					t.Fatalf("expected 2 orderings, got %d", len(sort.Orderings))
				}
				
				fo1, ok := sort.Orderings[0].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[0])
				}
				if fo1.Field.Name != "name" {
					t.Errorf("expected field 'name', got %q", fo1.Field.Name)
				}
				if fo1.Dir != query.DirAsc {
					t.Errorf("expected DirAsc, got %v", fo1.Dir)
				}
				
				fo2, ok := sort.Orderings[1].(*query.FieldOrdering)
				if !ok {
					t.Fatalf("expected FieldOrdering, got %T", sort.Orderings[1])
				}
				if fo2.Field.Name != "year" {
					t.Errorf("expected field 'year', got %q", fo2.Field.Name)
				}
				if fo2.Dir != query.DirDesc {
					t.Errorf("expected DirDesc, got %v", fo2.Dir)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				tt.check(t, sort)
			}
		})
	}
}
