package confusablematcher

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	var inMap []KeyValue

	inMap = append(inMap, KeyValue{"N", "T"})
	inMap = append(inMap, KeyValue{"I", "E"})
	inMap = append(inMap, KeyValue{"C", "S"})
	inMap = append(inMap, KeyValue{"E", "T"})

	var matcher = InitConfusableMatcher(inMap, true)
	index, length := IndexOf(matcher, "TEST", "NICE", false, 0)
	assert.Equal(t, index, 0)
	assert.Equal(t, length, 4)
	FreeConfusableMatcher(matcher)
}

func Test2(t *testing.T) {
	var inMap []KeyValue

	inMap = append(inMap, KeyValue{"V", "VA"})
	inMap = append(inMap, KeyValue{"V", "VO"})

	var matcher = InitConfusableMatcher(inMap, true)
	index, length := IndexOf(matcher, "VV", "VAVOVAVO", false, 0)
	assert.Equal(t, index, -1)
	assert.Equal(t, length, -1)

	index, length = IndexOf(matcher, "VAVOVAVO", "VV", false, 0)
	assert.Equal(t, index, 0)
	assert.True(t, length == 3 || length == 4)
	index, length = IndexOf(matcher, "VAVOVAVO", "VV", false, 4)
	assert.Equal(t, index, 4)
	assert.True(t, length == 3 || length == 4)
	index, length = IndexOf(matcher, "VAVOVAVO", "VV", false, 2)
	assert.Equal(t, index, 2)
	assert.True(t, length == 3 || length == 4)
	index, length = IndexOf(matcher, "VAVOVAVO", "VV", false, 3)
	assert.Equal(t, index, 4)
	assert.True(t, length == 3 || length == 4)
	FreeConfusableMatcher(matcher)
}

func Test3(t *testing.T) {
	var inMap []KeyValue

	inMap = append(inMap, KeyValue{"A", "\x02\x03"})
	inMap = append(inMap, KeyValue{"B", "\xC3\xBA\xC3\xBF"})

	var matcher = InitConfusableMatcher(inMap, true)
	index, length := IndexOf(matcher, "\x02\x03\xC3\xBA\xC3\xBF", "AB", false, 0)
	assert.Equal(t, index, 0)
	assert.Equal(t, length, 6)

	FreeConfusableMatcher(matcher)
}

func Test4(t *testing.T) {
	var inMap []KeyValue

	inMap = append(inMap, KeyValue{"S", "$"})
	inMap = append(inMap, KeyValue{"D", "[)"})

	var matcher = InitConfusableMatcher(inMap, true)
	SetIgnoreList(&matcher, []string{"_", " "})
	index, length := IndexOf(matcher, "A__ _ $$$[)D", "ASD", true, 0)
	assert.Equal(t, index, 0)
	assert.Equal(t, length, 11)

	FreeConfusableMatcher(matcher)
}

func Test5(t *testing.T) {

	var inMap []KeyValue

	inMap = append(inMap, KeyValue{"N", "/\\/"})
	inMap = append(inMap, KeyValue{"N", "/\\"})
	inMap = append(inMap, KeyValue{"I", "/"})

	for x := 0; x < 1000; x++ {
		var matcher = InitConfusableMatcher(inMap, true)
		var index, length int
		index, length = IndexOf(matcher, "/\\/CE", "NICE", false, 0)

		assert.Equal(t, index, 0)
		assert.Equal(t, length, 5)

		FreeConfusableMatcher(matcher)
	}
}

func Test6(t *testing.T) {
	var inMap []KeyValue

	inMap = append(inMap, KeyValue{"N", "/\\/"})
	inMap = append(inMap, KeyValue{"V", "\\/"})
	inMap = append(inMap, KeyValue{"I", "/"})

	var matcher = InitConfusableMatcher(inMap, true)
	index, length := IndexOf(matcher, "I/\\/AM", "INAN", false, 0)
	assert.Equal(t, index, -1)
	assert.Equal(t, length, -1)
	index, length = IndexOf(matcher, "I/\\/AM", "INAM", false, 0)
	assert.Equal(t, index, 0)
	assert.Equal(t, length, 6)
	index, length = IndexOf(matcher, "I/\\/AM", "IIVAM", false, 0)
	assert.Equal(t, index, 0)
	assert.Equal(t, length, 6)

	FreeConfusableMatcher(matcher)
}

func getDefaultMap() []KeyValue {
	var ret []KeyValue

	ret = append(ret, KeyValue{"N", "/[()[]]/"})
	ret = append(ret, KeyValue{"N", "ñ"})
	ret = append(ret, KeyValue{"N", "|\\|"})
	ret = append(ret, KeyValue{"N", "Ʌ/"})
	ret = append(ret, KeyValue{"N", "/IJ"})
	ret = append(ret, KeyValue{"N", "/|/"})

	var ns = []string{"Ӆ", "Π", "И", "𝐧", "𝑛", "𝒏", "𝓃", "𝓷", "𝔫", "𝕟", "𝖓", "𝗇", "𝗻", "𝘯", "𝙣", "𝚗", "ո", "ռ", "Ｎ", "ℕ", "𝐍", "𝑁", "𝑵", "𝒩", "𝓝", "𝔑", "𝕹", "𝖭", "𝗡", "𝘕", "𝙉", "𝙽", "Ν", "𝚴", "𝛮", "𝜨", "𝝢", "𝞜", "ꓠ", "Ń", "Ņ", "Ň", "ŋ", "Ɲ", "Ǹ", "Ƞ", "Ν", "Ṅ", "Ṇ", "Ṉ", "Ṋ", "₦", "ἠ", "ἡ", "ἢ", "ἣ", "ἤ", "ἥ", "ἦ", "ἧ", "ὴ", "ή", "ᾐ", "ᾑ", "ᾒ", "ᾓ", "ᾔ", "ᾕ", "ᾖ", "ᾗ", "ῂ", "ῃ", "ῄ", "ῆ", "ῇ", "ñ", "ń", "ņ", "ň", "ŉ", "Ŋ", "ƞ", "ǹ", "ȵ", "ɲ", "ɳ", "ɴ", "ᵰ", "ᶇ", "ṅ", "ṇ", "ṉ", "ṋ"}
	var is = []string{"Ỉ", "y", "i", "1", "|", "l", "j", "!", "/", "\\", "ｉ", "¡", "ⅰ", "ℹ", "ⅈ", "𝐢", "𝑖", "𝒊", "𝒾", "𝓲", "𝔦", "𝕚", "𝖎", "𝗂", "𝗶", "𝘪", "𝙞", "𝚒", "ı", "𝚤", "ɪ", "ɩ", "ι", "ι", "ͺ", "𝛊", "𝜄", "𝜾", "𝝸", "𝞲", "і", "Ⓘ", "ꙇ", "ӏ", "ꭵ", "Ꭵ", "ɣ", "ᶌ", "ｙ", "𝐲", "𝑦", "𝒚", "𝓎", "𝔂", "𝔶", "𝕪", "𝖞", "𝗒", "𝘆", "𝘺", "𝙮", "𝚢", "ʏ", "ỿ", "ꭚ", "γ", "ℽ", "𝛄", "𝛾", "𝜸", "𝝲", "𝞬", "у", "ү", "ყ", "Ｙ", "𝐘", "𝑌", "𝒀", "𝒴", "𝓨", "𝔜", "𝕐", "𝖄", "𝖸", "𝗬", "𝘠", "𝙔", "𝚈", "Υ", "ϒ", "𝚼", "𝛶", "𝜰", "𝝪", "𝞤", "Ⲩ", "У", "Ү", "Ꭹ", "Ꮍ", "ꓬ", "Ŷ", "Ÿ", "Ƴ", "Ȳ", "Ɏ", "ʏ", "Ẏ", "Ỳ", "Ỵ", "Ỷ", "Ỹ", "Ｙ", "Ì", "Í", "Î", "Ï", "Ĩ", "Ī", "Ĭ", "Į", "İ", "Ɩ", "Ɨ", "Ǐ", "Ȉ", "Ȋ", "ɪ", "Ί", "ΐ", "Ι", "Ϊ", "І", "Ѝ", "И", "Й", "Ӣ", "Ӥ", "Ḭ", "Ḯ", "Ỉ", "Ị", "Ῐ", "Ῑ", "Ⅰ", "Ｉ", "ェ", "エ", "ｪ", "ｴ", "ì", "í", "î", "ï", "ĩ", "ī", "ĭ", "į", "ı", "ǐ", "ȉ", "ȋ", "ɨ", "ɩ", "ͥ", "ί", "ϊ", "и", "й", "і", "ѝ", "ӣ", "ӥ", "ḭ", "ḯ", "ỉ", "ị", "ἰ", "ἱ", "ἲ", "ἳ", "ἴ", "ἵ", "ἶ", "ἷ", "ὶ", "ί", "ι", "ῐ", "ῑ", "ῒ", "ΐ", "ῖ", "ῗ", "ｉ", "ᶅ", "ḷ", "ḹ", "ḻ", "ḽ", "ý", "ÿ", "ŷ", "ƴ", "ȳ", "ɏ", "ʎ", "ʸ", "ẏ", "ẙ", "ỳ", "ỵ", "ỷ", "ỹ", "ｙ"}
	var gs = []string{"k", "g", "q", "4", "6", "9", "ğ", "൫", "ｇ", "ℊ", "𝐠", "𝑔", "𝒈", "𝓰", "𝔤", "𝕘", "𝖌", "𝗀", "𝗴", "𝘨", "𝙜", "𝚐", "ɡ", "ᶃ", "ƍ", "ց", "𝐆", "𝐺", "𝑮", "𝒢", "𝓖", "𝔊", "𝔾", "𝕲", "𝖦", "𝗚", "么", "𝘎", "𝙂", "𝙶", "Ԍ", "Ꮐ", "Ᏻ", "ꓖ", "Ĝ", "Ğ", "Ġ", "Ģ", "Ɠ", "Ǥ", "Ǧ", "Ǵ", "ʛ", "Γ", "Г", "Ḡ", "Ｇ", "Ꮆ", "ĝ", "ğ", "ġ", "ģ", "ǥ", "ǧ", "ǵ", "ɠ", "ɡ", "ɢ", "@"}
	var es = []string{"Ề", "Σ", "Ξ", "e", "3", "u", "℮", "ｅ", "ℯ", "ⅇ", "𝐞", "𝑒", "𝒆", "𝓮", "𝔢", "𝕖", "𝖊", "𝖾", "𝗲", "𝘦", "𝙚", "𝚎", "ꬲ", "е", "ҽ", "⋿", "Ｅ", "ℰ", "𝐄", "𝐸", "𝑬", "𝓔", "𝔈", "𝔼", "𝕰", "𝖤", "𝗘", "𝘌", "𝙀", "𝙴", "Ε", "𝚬", "𝛦", "𝜠", "𝝚", "𝞔", "Е", "ⴹ", "Ꭼ", "ꓰ", "È", "É", "Ê", "Ë", "Ē", "Ĕ", "Ė", "Ę", "Ě", "Ǝ", "Ɛ", "Ȅ", "Ȇ", "Ȩ", "Ɇ", "Έ", "Э", "Ӭ", "Ḕ", "Ḗ", "Ḙ", "Ḛ", "Ḝ", "Ẹ", "Ẻ", "Ẽ", "Ế", "Ề", "Ể", "Ễ", "Ệ", "Ἐ", "Ἑ", "Ἒ", "Ἓ", "Ἔ", "Ἕ", "Ὲ", "Έ", "è", "é", "ê", "ë", "ē", "ĕ", "ė", "ę", "ě", "Ə", "ȅ", "ȇ", "ȩ", "ɇ", "ɘ", "ɛ", "ɜ", "ɝ", "ɞ", "ͤ", "έ", "ε", "е", "э", "ӭ", "ḕ", "ḗ", "ḙ", "ḛ", "ḝ", "ẹ", "ẻ", "ẽ", "ế", "ề", "ể", "ễ", "ệ", "ἐ", "ἑ", "ἒ", "ἓ", "ἔ", "ἕ", "ὲ", "έ"}
	var rs = []string{"Ѓ", "Я", "r", "𝐫", "𝑟", "𝒓", "𝓇", "𝓻", "𝔯", "𝕣", "𝖗", "𝗋", "𝗿", "𝘳", "𝙧", "𝚛", "ꭇ", "ꭈ", "ᴦ", "ⲅ", "г", "ꮁ", "ℛ", "ℜ", "ℝ", "𝐑", "𝑅", "𝑹", "𝓡", "𝕽", "𝖱", "𝗥", "𝘙", "𝙍", "𝚁", "Ʀ", "Ꭱ", "Ꮢ", "𐒴", "ᖇ", "ꓣ", "Ŕ", "Ŗ", "Ř", "Ȑ", "Ȓ", "Ɍ", "ʀ", "ʁ", "Ṙ", "Ṛ", "Ṝ", "Ṟ", "Ɽ", "Ｒ", "Ꭱ", "ŕ", "ŗ", "ř", "ȑ", "ȓ", "ɍ", "ɹ", "ɺ", "ɻ", "ɼ", "ɽ", "ᚱ", "ᡵ", "ᵲ", "ᵳ", "ᶉ", "ṙ", "ṛ", "ṝ", "ṟ", "ｒ"}

	var s = [][]string{ns, is, gs, es, rs}

	for x := 0; x < 5; x++ {
		for i := 0; i < len(s[x]); i++ {
			switch x {
			case 0:
				ret = append(ret, KeyValue{"N", s[x][i]})
				break
			case 1:
				ret = append(ret, KeyValue{"I", s[x][i]})
				break
			case 2:
				ret = append(ret, KeyValue{"G", s[x][i]})
				break
			case 3:
				ret = append(ret, KeyValue{"E", s[x][i]})
				break
			case 4:
				ret = append(ret, KeyValue{"R", s[x][i]})
				break
			}
		}
	}

	return ret
}

func Test7(t *testing.T) {
	var inMap = getDefaultMap()

	var inp = "AAAAAAAAASSAFSAFNFNFNISFNSIFSIFJSDFUDSHF ASUF/|/__/|/___%/|/%I%%/|//|/%%%%%NNNN/|/NN__/|/N__𝘪G___%____$__G__𝓰𝘦Ѓ"
	var matcher = InitConfusableMatcher(inMap, true)
	SetIgnoreList(&matcher, []string{"_", "%", "$"})
	index, length := IndexOf(matcher, inp, "NIGGER", true, 0)

	assert.True(t, (index == 64 && length == 57) || (index == 89 && length == 32))

	FreeConfusableMatcher(matcher)
}

func TestLidlNormalizer(t *testing.T) {
	var inMap = getDefaultMap()

	// Additional test data
	var keys = []string{
		"A", "A", "A", "A", "B", "U", "U", "O", "O", "A", "A",
		"A", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y",
		"Z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"U", "A", " ", "S", "M", "O", "L", "N", "A", "T", "I", "O", "N", "N", "I", "G", "N", "I", "FREE", "AE",
	}
	var vals = []string{
		"ą", "ꬱ", "ᵃ", "å", "⒝", "ü", "Ü", "ö", "Ö", "ä", "Ä",
		"⒜", "⒞", "⒟", "⒠", "⒡", "⒢", "⒣", "⒤", "⒥", "⒦", "⒧", "⒨", "⒩", "⒪", "⒫", "⒬", "⒭", "⒮", "⒯", "⒰", "⒱", "⒲", "⒳", "⒴",
		"Ⓩ", "ⓐ", "ⓑ", "ⓒ", "ⓓ", "ⓔ", "ⓕ", "ⓖ", "ⓗ", "ⓘ", "ⓙ", "ⓚ", "ⓛ", "ⓜ", "ⓝ", "ⓞ", "ⓟ", "ⓠ", "ⓡ", "ⓢ", "ⓣ", "ⓤ", "ⓥ", "ⓦ", "ⓧ", "ⓨ", "ⓩ", "⓪",
		"𝕒", "𝕓", "𝕔", "𝕕", "𝕖", "𝕗", "𝕘", "𝕙", "𝕚", "𝕛", "𝕜", "𝕝", "𝕞", "𝕟", "𝕠", "𝕡", "𝕢", "𝕣", "𝕤", "𝕥", "𝕦", "𝕧", "𝕨", "𝕩", "𝕪", "𝕫",
		"🄰", "🄱", "🄲", "🄳", "🄴", "🄵", "🄶", "🄷", "🄸", "🄹", "🄺", "🄻", "🄼", "🄽", "🄾", "🄿", "🅀", "🅁", "🅂", "🅃", "🅄", "🅅", "🅆", "🅇", "🅈", "🅉",
		"₳", "฿", "₵", "Đ", "Ɇ", "₣", "₲", "Ⱨ", "ł", "J", "₭", "Ⱡ", "₥", "₦", "Ø", "₱", "Q", "Ɽ", "₴", "₮", "Ʉ", "V", "₩", "Ӿ", "Ɏ", "Ⱬ",
		"𝖆", "𝖇", "𝖈", "𝖉", "𝖊", "𝖋", "𝖌", "𝖍", "𝖎", "𝖏", "𝖐", "𝖑", "𝖒", "𝖓", "𝖔", "𝖕", "𝖖", "𝖗", "𝖘", "𝖙", "𝖚", "𝖛", "𝖜", "𝖝", "𝖞", "𝖟",
		"🅰", "🅱", "🅲", "🅳", "🅴", "🅵", "🅶", "🅷", "🅸", "🅹", "🅺", "🅻", "🅼", "🅽", "🅾", "🅿", "🆀", "🆁", "🆂", "🆃", "🆄", "🆅", "🆆", "🆇", "🆈", "🆉",
		"🇺", "🇦", " ", "ˢ", "ᵐ", "ᵒ", "ˡ", "ⁿ", "ᵃ", "ᵗ", "ᶦ", "ᵒ", "ⁿ", "Н", "и", "г", "🇳", "🇮", "🆓", "ᴭ",
	}

	for x := 0; x < len(keys); x++ {
		inMap = append(inMap, KeyValue{keys[x], vals[x]})
	}

	var matcher = InitConfusableMatcher(inMap, true)

	var data = []KeyValue{
		KeyValue{"ą", "A"},
		KeyValue{"ꬱ", "A"},
		KeyValue{"ᵃ", "A"},
		KeyValue{"abc å def", "ABC A DEF"},
		KeyValue{"ˢᵐᵒˡ ⁿᵃᵗᶦᵒⁿ", "SMOL NATION"},
		KeyValue{"Ниг", "NIG"},
		KeyValue{"🇺🇦XD", "UAXD"},
		KeyValue{"🆓 ICE", "FREE ICE"},
		KeyValue{"chocolate 🇳🇮b", "CHOCOLATE NIB"},
		KeyValue{"🅱lueberry", "BLUEBERRY"},
		KeyValue{"⒝", "B"},
		KeyValue{"ü Ü ö Ö ä Ä", "U U O O A A"},
		KeyValue{"ᴭ", "AE"},
		KeyValue{"⒜ ⒝ ⒞ ⒟ ⒠ ⒡ ⒢ ⒣ ⒤ ⒥ ⒦ ⒧ ⒨ ⒩ ⒪ ⒫ ⒬ ⒭ ⒮ ⒯ ⒰ ⒱ ⒲ ⒳ ⒴", "A B C D E F G H I J K L M N O P Q R S T U V W X Y"},
		KeyValue{"Ⓩⓐⓑⓒⓓⓔⓕⓖⓗⓘⓙⓚⓛⓜⓝⓞⓟⓠⓡⓢⓣⓤⓥⓦⓧⓨⓩ⓪", "ZABCDEFGHIJKLMNOPQRSTUVWXYZ0"},
		KeyValue{"𝕒𝕓𝕔𝕕𝕖𝕗𝕘𝕙𝕚𝕛𝕜𝕝𝕞𝕟𝕠𝕡𝕢𝕣𝕤𝕥𝕦𝕧𝕨𝕩𝕪𝕫", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		KeyValue{"🄰🄱🄲🄳🄴🄵🄶🄷🄸🄹🄺🄻🄼🄽🄾🄿🅀🅁🅂🅃🅄🅅🅆🅇🅈🅉", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		KeyValue{"₳฿₵ĐɆ₣₲ⱧłJ₭Ⱡ₥₦Ø₱QⱤ₴₮ɄV₩ӾɎⱫ", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		KeyValue{"𝖆𝖇𝖈𝖉𝖊𝖋𝖌𝖍𝖎𝖏𝖐𝖑𝖒𝖓𝖔𝖕𝖖𝖗𝖘𝖙𝖚𝖛𝖜𝖝𝖞𝖟", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		KeyValue{"🅰🅱🅲🅳🅴🅵🅶🅷🅸🅹🅺🅻🅼🅽🅾🅿🆀🆁🆂🆃🆄🆅🆆🆇🆈🆉", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
	}

	var data2 = [][]int{
		[]int{0, 2},
		[]int{0, 3},
		[]int{0, 3},
		[]int{0, 10},
		[]int{0, 29},
		[]int{0, 6},
		[]int{0, 10},
		[]int{0, 8},
		[]int{0, 19},
		[]int{0, 12},
		[]int{0, 3},
		[]int{0, 17},
		[]int{0, 3},
		[]int{0, 99},
		[]int{0, 84},
		[]int{0, 104},
		[]int{0, 104},
		[]int{0, 65},
		[]int{0, 104},
		[]int{0, 104},
	}

	for x := 0; x < len(data); x++ {
		index, length := IndexOf(matcher, data[x].Key, data[x].Value, true, 0)
		assert.Equal(t, index, data2[x][0])
		assert.Equal(t, length, data2[x][1])
	}

	FreeConfusableMatcher(matcher)
}

func Test8(t *testing.T) {
	var inMap []KeyValue

	var matcher = InitConfusableMatcher(inMap, true)
	SetIgnoreList(&matcher, []string{"̲", "̅", "[", "]"})
	index, length := IndexOf(matcher,
		"[̲̅a̲̅][̲̅b̲̅][̲̅c̲̅][̲̅d̲̅][̲̅e̲̅][̲̅f̲̅][̲̅g̲̅][̲̅h̲̅][̲̅i̲̅][̲̅j̲̅][̲̅k̲̅][̲̅l̲̅][̲̅m̲̅][̲̅n̲̅][̲̅o̲̅][̲̅p̲̅][̲̅q̲̅][̲̅r̲̅][̲̅s̲̅][̲̅t̲̅][̲̅u̲̅][̲̅v̲̅][̲̅w̲̅][̲̅x̲̅][̲̅y̲̅][̲̅z̲̅][̲̅0̲̅][̲̅1̲̅][̲̅2̲̅][̲̅3̲̅][̲̅4̲̅][̲̅5̲̅][̲̅6̲̅][̲̅7̲̅][̲̅8̲̅][̲̅9̲̅][̲̅0̲̅]",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890",
		false,
		0)
	assert.Equal(t, 5, index)
	assert.Equal(t, 397, length)

	FreeConfusableMatcher(matcher)
}

func Test9(t *testing.T) {
	var inMap []KeyValue

	inMap = append(inMap, KeyValue{" ", " "})

	for x := 0; x < 1000; x++ {
		var matcher = InitConfusableMatcher(inMap, true)

		index, length := IndexOf(matcher,
			"NOT NICE",
			"VERY NICE",
			false,
			0)
		assert.Equal(t, -1, index)
		assert.Equal(t, -1, length)

		AddMapping(matcher, "VERY", "NOT", false)

		index, length = IndexOf(matcher,
			"NOT NICE",
			"VERY NICE",
			false,
			0)
		assert.Equal(t, 0, index)
		assert.Equal(t, 8, length)

		RemoveMapping(matcher, "VERY", "NOT")
		FreeConfusableMatcher(matcher)
	}
}

func Test10(t *testing.T) {
	var inMap []KeyValue

	inMap = append(inMap, KeyValue{"B", "A"})
	inMap = append(inMap, KeyValue{"B", "AB"})
	inMap = append(inMap, KeyValue{"B", "ABC"})
	inMap = append(inMap, KeyValue{"B", "ABCD"})
	inMap = append(inMap, KeyValue{"B", "ABCDE"})
	inMap = append(inMap, KeyValue{"B", "ABCDEF"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFG"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGH"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHI"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJ"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJK"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJKL"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJKLM"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJKLMN"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJKLMNO"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJKLMNOP"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJKLMNOPQ"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJKLMNOPQR"})
	inMap = append(inMap, KeyValue{"B", "ABCDEFGHIJKLMNOPQRS"})

	var matcher = InitConfusableMatcher(inMap, true)

	index, length := IndexOf(matcher,
		"ABCDEFGHIJKLMNOPQRS",
		"B",
		false,
		0)
	assert.Equal(t, 0, index)
	assert.True(t, length >= 0 && length == 1)

	RemoveMapping(matcher, "B", "ABCDEFGHIJKLMNOP")
	AddMapping(matcher, "B", "P", false)
	AddMapping(matcher, "B", "PQ", false)
	AddMapping(matcher, "B", "PQR", false)
	AddMapping(matcher, "B", "PQRS", false)
	AddMapping(matcher, "B", "PQRST", false)
	AddMapping(matcher, "B", "PQRSTU", false)
	AddMapping(matcher, "B", "PQRSTUV", false)
	AddMapping(matcher, "B", "PQRSTUVW", false)
	AddMapping(matcher, "B", "PQRSTUVWX", false)
	AddMapping(matcher, "B", "PQRSTUVWXY", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ0", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ01", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ012", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ0123", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ01234", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ012345", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ0123456", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ01234567", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ012345678", false)
	AddMapping(matcher, "B", "PQRSTUVWXYZ0123456789", false)

	index, length = IndexOf(matcher,
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		"BB",
		false,
		0)
	assert.Equal(t, 0, index)
	assert.Equal(t, 2, length)

	index, length = IndexOf(matcher,
		"PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789PQRSTUVWXYZ0123456789",
		"BBBBBBBBBBBBBBBBBBBBBBBBBBB",
		true,
		0)
	assert.Equal(t, 0, index)
	assert.True(t, length >= 0 && length == 547)

	FreeConfusableMatcher(matcher)
}

func Test11(t *testing.T) {
	var inMap []KeyValue
	var matcher = InitConfusableMatcher(inMap, true)
	index, length := IndexOf(matcher, ":)", "", true, 0)
	assert.Equal(t, 0, index)
	assert.Equal(t, 0, length)

	index, length = IndexOf(matcher, "", ":)", true, 0)
	assert.Equal(t, -1, index)
	assert.Equal(t, -1, length)

	FreeConfusableMatcher(matcher)
}

func Test12(t *testing.T) {
	var inMap []KeyValue

	var matcher = InitConfusableMatcher(inMap, true)

	AddMapping(matcher, "A", "A", false)
	AddMapping(matcher, "A", "A", false)
	AddMapping(matcher, "A", "A", false)
	AddMapping(matcher, "A", "A", false)

	index, length := IndexOf(matcher, "ABAAA", "ABAR", true, 0)
	assert.Equal(t, -1, index)
	assert.Equal(t, -1, length)

	FreeConfusableMatcher(matcher)
}

func Test13(t *testing.T) {
	var inMap []KeyValue

	var matcher = InitConfusableMatcher(inMap, true)
	index, length := IndexOf(matcher, "?", "?", true, 0)
	assert.Equal(t, -1, index)
	assert.Equal(t, -1, length)

	for x := 0; x < 1000; x++ {
		assert.Equal(t, RemoveMapping(matcher, "?", "?"), false)
		assert.Equal(t, AddMapping(matcher, "?", "?", false), Success)
		assert.Equal(t, RemoveMapping(matcher, "?_", "?_"), false)
		assert.Equal(t, RemoveMapping(matcher, "?", "_"), false)
		assert.Equal(t, RemoveMapping(matcher, "?", "?__"), false)

		assert.Equal(t, AddMapping(matcher, "?_", "?_", false), Success)
		assert.Equal(t, AddMapping(matcher, "?_", "_", false), Success)
		assert.Equal(t, AddMapping(matcher, "?_", "?___", false), Success)
		assert.Equal(t, RemoveMapping(matcher, "?_", "?___"), true)
		assert.Equal(t, RemoveMapping(matcher, "?_", "_"), true)
		assert.Equal(t, RemoveMapping(matcher, "?_", "?_"), true)
		assert.Equal(t, RemoveMapping(matcher, "?", "?"), true)
	}

	FreeConfusableMatcher(matcher)
}

func Test14(t *testing.T) {
	var inMap []KeyValue

	var matcher = InitConfusableMatcher(inMap, true)

	index, length := IndexOf(matcher, "A", "A", false, 0)
	assert.Equal(t, 0, index)
	assert.Equal(t, 1, length)

	var matcher2 = InitConfusableMatcher(inMap, false)
	index, length = IndexOf(matcher2, "A", "A", false, 0)
	assert.Equal(t, -1, index)
	assert.Equal(t, -1, length)

	FreeConfusableMatcher(matcher)
}

func Test15(t *testing.T) {
	var inMap []KeyValue

	var matcher = InitConfusableMatcher(inMap, true)
	assert.Equal(t, AddMapping(matcher, "", "?", false), EmptyKey)
	assert.Equal(t, AddMapping(matcher, "?", "", false), EmptyValue)
	assert.True(t, AddMapping(matcher, "", "", false) != Success)
	assert.True(t, AddMapping(matcher, "\x00\x01", "?", false) != Success)
	assert.True(t, AddMapping(matcher, "?", "\x00", false) != Success)
	assert.Equal(t, AddMapping(matcher, "\x01", "?", false), InvalidKey)
	assert.Equal(t, AddMapping(matcher, "?", "\x01", false), InvalidValue)
	assert.True(t, AddMapping(matcher, "\x00", "\x01", false) != Success)
	assert.True(t, AddMapping(matcher, "\x00", "\x00", false) != Success)
	assert.True(t, AddMapping(matcher, "\x01", "\x00", false) != Success)
	assert.True(t, AddMapping(matcher, "\x01", "\x01", false) != Success)
	assert.True(t, AddMapping(matcher, "\x01\x00", "\x00\x01", false) != Success)
	assert.True(t, AddMapping(matcher, "A\x00", "\x00A", false) != Success)
	assert.True(t, AddMapping(matcher, "\x01\x00", "\x00\x01", false) != Success)
	assert.Equal(t, AddMapping(matcher, "A\x00", "A\x01", false), Success)
	assert.Equal(t, AddMapping(matcher, "A\x01", "A\x00", false), Success)
	assert.Equal(t, AddMapping(matcher, "A\x00", "A\x00", false), Success)
	assert.Equal(t, AddMapping(matcher, "A\x01", "A\x01", false), Success)

	FreeConfusableMatcher(matcher)
}

func Test16(t *testing.T) {
	var inMap []KeyValue
	var running = true
	var matcher = InitConfusableMatcher(inMap, true)

	go func() {
		for running {
			IndexOf(matcher, "ASD", "ZXC", false, 0)
		}
	}()

	go func() {
		for running {
			AddMapping(matcher, "Z", "A", false)
			RemoveMapping(matcher, "Z", "A")
		}
	}()

	time.Sleep(time.Second * 10)

	running = false
}

func Test17(t *testing.T) {
	var inMap []KeyValue
	inMap = append(inMap, KeyValue{"N", "/\\/"})
	var ignoreList []string
	var matcher = InitConfusableMatcher(inMap, true)
	var running = true
	var lock sync.Mutex

	go func() {
		for running {
			lock.Lock()
			{
				IndexOf(matcher, "/\\/", "N", false, 0)
			}
			lock.Unlock()
		}
	}()

	go func() {
		for running {
			lock.Lock()
			{
				FreeConfusableMatcher(matcher)
				matcher = InitConfusableMatcher(inMap, true)
			}
			lock.Unlock()
			SetIgnoreList(&matcher, ignoreList)
		}
	}()

	time.Sleep(time.Second * 10)

	running = false
}
