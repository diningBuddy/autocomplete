package hangul

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	hangulCHO    = []string{"ㄱ", "ㄲ", "ㄴ", "ㄷ", "ㄸ", "ㄹ", "ㅁ", "ㅂ", "ㅃ", "ㅅ", "ㅆ", "ㅇ", "ㅈ", "ㅉ", "ㅊ", "ㅋ", "ㅌ", "ㅍ", "ㅎ"}
	hangulJUN    = []string{"ㅏ", "ㅐ", "ㅑ", "ㅒ", "ㅓ", "ㅔ", "ㅕ", "ㅖ", "ㅗ", "ㅘ", "ㅙ", "ㅚ", "ㅛ", "ㅜ", "ㅝ", "ㅞ", "ㅟ", "ㅠ", "ㅡ", "ㅢ", "ㅣ"}
	hangulJON    = []string{"", "ㄱ", "ㄲ", "ㄳ", "ㄴ", "ㄵ", "ㄶ", "ㄷ", "ㄹ", "ㄺ", "ㄻ", "ㄼ", "ㄽ", "ㄾ", "ㄿ", "ㅀ", "ㅁ", "ㅂ", "ㅄ", "ㅅ", "ㅆ", "ㅇ", "ㅈ", "ㅊ", "ㅋ", "ㅌ", "ㅍ", "ㅎ"}
	doubleJonMap = map[string]string{"ㄳ": "ㄱㅅ", "ㄵ": "ㄴㅈ", "ㄶ": "ㄴㅎ", "ㄺ": "ㄹㄱ", "ㄻ": "ㄹㅁ",
		"ㄼ": "ㄹㅂ", "ㄽ": "ㄹㅅ", "ㄾ": "ㄹㅌ", "ㄿ": "ㄹㅍ", "ㅀ": "ㄹㅎ", "ㅄ": "ㅂㅅ",
	} // 뒤의 경우는 한번에 입력가능하므로 적용하지 않음 "ㄲ": "ㄱㄱ", "ㅆ": "ㅅㅅ",
	doubleJunMap = map[string]string{"ㅘ": "ㅗㅏ", "ㅙ": "ㅗㅐ", "ㅚ": "ㅗㅣ", "ㅝ": "ㅜㅓ", "ㅞ": "ㅜㅔ", "ㅟ": "ㅜㅣ", "ㅢ": "ㅡㅣ"}

	engKorMap = map[string]string{"q": "ㅂ", "w": "ㅈ", "e": "ㄷ", "r": "ㄱ", "t": "ㅅ", "y": "ㅛ", "u": "ㅕ", "i": "ㅑ", "o": "ㅐ", "p": "ㅔ",
		"a": "ㅁ", "s": "ㄴ", "d": "ㅇ", "f": "ㄹ", "g": "ㅎ", "h": "ㅗ", "j": "ㅓ", "k": "ㅏ", "l": "ㅣ", "z": "ㅋ", "x": "ㅌ", "c": "ㅊ", "v": "ㅍ", "b": "ㅠ", "n": "ㅜ", "m": "ㅡ",
		"Q": "ㅃ", "W": "ㅉ", "E": "ㄸ", "R": "ㄲ", "T": "ㅆ", "Y": "ㅛ", "U": "ㅕ", "I": "ㅑ", "O": "ㅒ", "P": "ㅖ",
		"A": "ㅁ", "S": "ㄴ", "D": "ㅇ", "F": "ㄹ", "G": "ㅎ", "H": "ㅗ", "J": "ㅓ", "K": "ㅏ", "L": "ㅣ", "Z": "ㅋ", "X": "ㅌ", "C": "ㅊ", "V": "ㅍ", "B": "ㅠ", "N": "ㅜ", "M": "ㅡ"}
)

const (
	// 한글의 시작은 '가' 끝은 '힣'
	hangulBASE = rune('가')
	hangulEND  = rune('힣')

	// 자음은 29개, 모음은 20개
	hangulJA = rune('ㄱ')
	hangulMO = rune('ㅏ')
)

// IsHangul checks if input string is all hangul
func IsHangul(s string) bool {
	if strings.TrimSpace(s) == "" {
		// empty string or whitespace
		return false
	}

	for _, c := range s {
		if c >= hangulJA && c <= (hangulJA+29) {
			// 자음 단독 입력
			continue
		} else if c >= hangulMO && c <= (hangulMO+20) {
			// 모음 단독 입력
			continue
		} else if c >= hangulBASE && c <= hangulEND {
			// 일반 한글 입력
			continue
		} else if unicode.IsSpace(c) {
			// white space
			continue
		} else {
			return false
		}
	}
	return true
}

// SplitJamoCharWithSplitDoubleJunJon takes korean string and returns a string of korean jamochars split and double jun, jon chars split
func SplitJamoCharWithSplitDoubleJunJon(word string) string {
	var sum string

	a := []rune(word)

	for i := 0; i < len(a); i++ {

		s := fmt.Sprintf("0x%x", a[i])
		s2, _ := strconv.ParseInt(s, 0, 64)

		if s2 >= int64(hangulBASE) && s2 <= int64(hangulEND) {

			cho := (((s2 - 0xAC00) - (s2-0xAC00)%28) / 28) / 21
			jung := (((s2 - 0xAC00) - (s2-0xAC00)%28) / 28) % 21
			jong := (s2 - 0xAC00) % 28

			var splitJon, splitJun string

			if c, ok := doubleJonMap[hangulJON[jong]]; ok {
				splitJon = c
			} else {
				splitJon = hangulJON[jong]
			}

			if c, ok := doubleJunMap[hangulJUN[jung]]; ok {
				splitJun = c
			} else {
				splitJun = hangulJUN[jung]
			}

			sum += hangulCHO[cho] + splitJun + splitJon

		} else {
			sum += string(a[i])
		}

	}

	return sum
}

// Eng2KorRaw maps english to korean keyboard positions without considering double jungsung, jongsung
func Eng2KorRaw(text string) string {
	var kor strings.Builder

	for _, v := range text {
		if h, ok := engKorMap[string(v)]; ok {
			kor.WriteString(h)
		} else {
			kor.WriteString(string(v))
		}
	}
	return kor.String()
}
