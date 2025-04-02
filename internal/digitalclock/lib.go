package digitalclock

import (
	"image"
	"image/color"
	"strconv"
	"time"
)

// Form png that contains time passed scaled by a factor of k
func formImage(k int, timeString string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, (6*width+2*w_colon)*k, height*k))
	left := 0
	line := 0
	for _, digit := range timeString {
		toPrint := convertCharToString(digit)
		wi := width
		if digit == ':' {
			wi = w_colon
		}
		linePos := 0
		for _, char := range toPrint {
			if colorSquare(img, k, linePos, left, char, line) {
				line++
				linePos = 0
			} else {
				linePos++
			}
		}
		left += wi * k
		line = 0
	}

	return img
}

func parseK(query string) (int, error) {
	var k int
	var err error

	if query == "" {
		k = 1
	} else {
		k, err = strconv.Atoi(query)
		if err != nil {
			return 0, err
		}
	}

	if k < 1 || k > 30 {
		return 0, err
	}

	return k, nil
}

func parseTime(query string) (string, error) {
	var t time.Time
	var err error

	if len(query) != 8 && len(query) != 0 {
		return "", err
	}

	if query == "" {
		t = time.Now()
	} else {
		t, err = time.Parse("15:04:05", query)
		if err != nil {
			return "", err
		}
	}

	return t.Format("15:04:05"), nil
}

func convertCharToString(c rune) string {
	var toPrint string
	switch c {
	case ':':
		toPrint = Colon
	case '0':
		toPrint = Zero
	case '1':
		toPrint = One
	case '2':
		toPrint = Two
	case '3':
		toPrint = Three
	case '4':
		toPrint = Four
	case '5':
		toPrint = Five
	case '6':
		toPrint = Six
	case '7':
		toPrint = Seven
	case '8':
		toPrint = Eight
	case '9':
		toPrint = Nine
	}

	return toPrint
}

func colorSquare(img *image.RGBA, k int, i int, left int, char rune, line int) bool {
	var c color.RGBA
	if char == '.' {
		c = color.RGBA{0xff, 0xff, 0xff, 0xff}
	} else if char == '1' {
		c = Cyan
	} else {
		return true
	}

	for x := left + i*k; x < left+(i+1)*k; x++ {
		for y := line * k; y < (line+1)*k; y++ {
			img.Set(x, y, c)
		}
	}

	return false
}
