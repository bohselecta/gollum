//go:build !darwin

package metal

func Init() bool                 { return false }
func DeviceName() string         { return "" }
func DecodeStep(batch, seqlen int) int { return 0 }
