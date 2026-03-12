package hashtable

type Flags struct {
	Collision bool
	Occupied  bool
	Terminal  bool
	Link      bool
	Deleted   bool
}

type Entry struct {
	Key   string
	Value string
	V     int
	Home  int
	Flags Flags
	Next  int
}

type Table struct {
	size  int
	base  int
	slots []Entry
}

type Record struct {
	Key   string
	Value string
}
