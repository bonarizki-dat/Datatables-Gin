package dto

// ========================
// Params → standard DataTables parameters
// ========================
type Params struct {
	Draw   int64
	Start  int
	Length int
	Search string
	Order  string
	Dir    string
}