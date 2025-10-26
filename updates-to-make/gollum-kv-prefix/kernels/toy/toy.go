package toy
import ("math/rand"; "time"; "strings")
type ToyBackend struct{}
func NewToyBackend()*ToyBackend{ return &ToyBackend{} }
func (t *ToyBackend) NextToken(prefix string) string {
    rand.Seed(time.Now().UnixNano())
    w := []string{" llama"," on"," the"," high"," plain","."," gentle",","," wind"," hums"," softly"}
    if strings.HasSuffix(prefix, "."){ return " " }
    return w[rand.Intn(len(w))]
}
