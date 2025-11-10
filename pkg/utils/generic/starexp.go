package generic

import "strings"

const Wildcard = "*"
const WildHelp = `
opensearch-cli supports wildcard expressions using '*' character.
The '*' recognized as any number of characters(0 .. N). Examples of supported expressions:
'a*'    - matches any string that starts with or equals to 'a'
'*b'    - matches any string that ends with or equals to 'b'
'a*b'   - matches any string that starts with 'a' and ends with 'b' with any number of characters in between
'*a*'   - matches any string contains 'a'
'*'     - matches any string`

func ContainsWildcard(s string) bool {
	return strings.Contains(s, Wildcard)
}

// GetMatchFunc - returns matching func for stringPattern,
// stringPattern - is either simple string or pattern,
// where the pattern is the string with the Wildcard occurrence
// supported expressions:
//   - 'aaaaa' - simple string
//   - 'a*' - matches any string that starts with or equals to 'a'
//   - '*b' - matches any string that ends with or equals to 'b'
//   - 'a*b' - matches any string that starts with 'a' and ends with 'b' with any number of characters in between
//   - *a*' - matches any string contains 'a'
//   - '*' - matches any string
func GetMatchFunc(stringPattern string) func(string) bool {
	if ContainsWildcard(stringPattern) {
		if len(stringPattern) == 1 || strings.Count(stringPattern, Wildcard) == len(stringPattern) {
			return func(other string) bool {
				return true
			}
		}
		cleanedPattern := strings.ReplaceAll(stringPattern, Wildcard, "")
		if strings.HasSuffix(stringPattern, Wildcard) && strings.HasPrefix(stringPattern, Wildcard) && len(stringPattern) > 1 && len(cleanedPattern) > 0 {
			return func(other string) bool {
				return strings.Contains(other, cleanedPattern)
			}
		} else if strings.HasPrefix(stringPattern, Wildcard) {
			return func(other string) bool {
				return strings.HasSuffix(other, stringPattern[1:])
			}
		} else if strings.HasSuffix(stringPattern, Wildcard) {
			return func(other string) bool {
				return strings.HasPrefix(other, stringPattern[0:len(stringPattern)-1])
			}
		} else {
			return func(other string) bool {
				groups := strings.Split(stringPattern, Wildcard)
				return strings.HasPrefix(other, groups[0]) && strings.HasSuffix(other, groups[1])
			}
		}
	}
	return func(other string) bool {
		return stringPattern == other
	}
}
