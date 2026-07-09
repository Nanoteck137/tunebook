package schema

import (
	"github.com/nanoteck137/tunebook/tools/query"
)

type Schema struct {
	fields      map[string]*query.Field
	defaultSort []query.Ordering
}

func New() *Schema {
	return &Schema{
		fields: make(map[string]*query.Field),
	}
}

func (s *Schema) AddField(name string, typ query.Type, opts ...FieldOption) *Schema {
	f := &query.Field{
		Name: name,
		Type: typ,
		Meta: make(map[string]any),
	}
	for _, opt := range opts {
		opt(f)
	}
	s.fields[name] = f
	return s
}

func (s *Schema) SetDefaultSort(orderings ...query.Ordering) *Schema {
	s.defaultSort = orderings
	return s
}

func (s *Schema) GetDefaultSort() []query.Ordering {
	return s.defaultSort
}

func (s *Schema) Field(name string) (*query.Field, bool) {
	f, ok := s.fields[name]
	return f, ok
}

type FieldOption func(*query.Field)

func Nullable() FieldOption {
	return func(f *query.Field) {
		f.Nullable = true
	}
}

func Column(name string) FieldOption {
	return func(f *query.Field) {
		f.Meta["column"] = name
	}
}

type RelationOption func(*query.RelationConfig)

func Relation(joinTable, joinForeignKey, joinReference string, valueType query.Type) FieldOption {
	return func(f *query.Field) {
		f.Type = query.TypeRelation
		f.Relation = &query.RelationConfig{
			JoinTable:      joinTable,
			JoinForeignKey: joinForeignKey,
			JoinReference:  joinReference,
			ValueType:      valueType,
		}
	}
}
