package ascii

const (
	DefaultCharacterSet = "8@$e*+!:."

	DetailedCharacterSet = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'."
)

var CharacterSetMap = map[string]string{
	"default":  DefaultCharacterSet,
	"detailed": DetailedCharacterSet,
}

func GetCharacterSet(name string) string {
	if charset, exists := CharacterSetMap[name]; exists {
		return charset
	}
	return DefaultCharacterSet
}

func GetAvailableCharacterSets() []string {
	sets := make([]string, 0, len(CharacterSetMap))
	for name := range CharacterSetMap {
		sets = append(sets, name)
	}
	return sets
}
