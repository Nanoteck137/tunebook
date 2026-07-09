package query

type Plan struct {
	Filter  FilterNode
	OrderBy []Ordering
}

type FilterNode interface {
	filterNode()
}

type AndNode struct {
	Left  FilterNode
	Right FilterNode
}

type OrNode struct {
	Left  FilterNode
	Right FilterNode
}

type NotNode struct {
	Expr FilterNode
}

type ComparisonNode struct {
	Field    *Field
	Operator Operator
	Value    Value
}

type IsNullNode struct {
	Field *Field
	Not   bool
}

type ContainsNode struct {
	Field *Field
	Value Value
}

type InNode struct {
	Field  *Field
	Values []Value
	Not    bool
}

type HasNode struct {
	Field    *Field
	Operator Operator
	Value    Value
	Not      bool
}

func (*AndNode) filterNode()        {}
func (*OrNode) filterNode()         {}
func (*NotNode) filterNode()        {}
func (*ComparisonNode) filterNode() {}
func (*IsNullNode) filterNode()     {}
func (*ContainsNode) filterNode()   {}
func (*InNode) filterNode()         {}
func (*HasNode) filterNode()        {}

type Field struct {
	Name     string
	Type     Type
	Nullable bool
	Meta     map[string]any
	Relation *RelationConfig
}

type RelationConfig struct {
	JoinTable      string
	JoinForeignKey string
	JoinReference  string
	ValueType      Type
}

type Value struct {
	Type  Type
	Value any
}

type Type int

const (
	TypeString Type = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeTime
	TypeRelation
)

type Operator int

const (
	OpEqual Operator = iota
	OpNotEqual
	OpLike
	OpGreater
	OpGreaterEqual
	OpLess
	OpLessEqual
)

type Ordering interface {
	ordering()
}

type Direction int

const (
	DirAsc Direction = iota
	DirDesc
)

type NullOrdering int

const (
	NullOrderingDefault NullOrdering = iota
	NullOrderingFirst
	NullOrderingLast
)

type FieldOrdering struct {
	Field      *Field
	Dir        Direction
	NullOrder  NullOrdering
}

type RandomOrdering struct{}

type ShuffleOrdering struct {
	Seed int64
}

type ScoreOrdering struct{}

func (*FieldOrdering) ordering()  {}
func (*RandomOrdering) ordering() {}
func (*ShuffleOrdering) ordering() {}
func (*ScoreOrdering) ordering()  {}
