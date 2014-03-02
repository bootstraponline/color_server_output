package main

import "io/ioutil"
import "os"

const (
	//                                0   1   2   3   4
	cyan      = "\033[36m" // 54 : 0 [27, 91, 51, 54, 109]
	green     = "\033[32m" // 50 : 1 [27, 91, 51, 50, 109]
	white_0   = "\033[0m"  // 48 : 2 [27, 91, 48, 109]
	white_1   = "\033[39m" // 57 : 3 [27, 91, 51, 57, 109]
	grey      = "\033[90m" // 90 : 4 [27, 91, 57, 48, 109]
	no_result = byte(10)
)

var state = 0
var last = byte(0)

func match(byte byte) byte {
	// \e always means the start of a new color code.
	if byte == 27 {
		state = 1
		return no_result
	}

	switch state {
	case 1:
		if byte == 91 {
			state = 2
		} else {
			state = 0
		}
	case 2:
		if byte == 51 || byte == 48 || byte == 57 {
			state = 3
			last = byte
		} else {
			state = 0
		}
	case 3:
		if byte == 109 && last == 48 {
			state = 0
			return last
		} else if byte == 54 || byte == 50 || byte == 57 || byte == 48 {
			state = 4
			last = byte
		} else {
			state = 0
		}
	case 4:
		if byte == 109 {
			char := last
			if char == 48 {
				char = 90
			}
			// only 4 length 0 is grey
			return char
		} else {
			state = 0
		}
	default:
		state = 0
	}

	return no_result
}

func code_len(byte byte) int {
	switch byte {
	case 54:
		return len(cyan)
	case 50:
		return len(green)
	case 48:
		return len(white_0)
	case 57:
		return len(white_1)
	case 90:
		return len(grey)
	default:
		return 0
	}
}

func code_div(input byte) []byte {
	switch input {
	case 54:
		return []byte("</div><div class=\"cyan\">")
	case 50:
		return []byte("</div><div class=\"green\">")
	case 48:
		return []byte("</div><div>")
	case 57:
		return []byte("</div><div>")
	case 90:
		return []byte("</div><div class=\"grey\">")
	default:
		return []byte("</div><div>")
	}
}

func pm(input []byte) []byte {
	index := 0
	len_input := len(input)
	result := make([]byte, 0)
	// iterate through each byte
	for i := 0; i < len_input; i++ {
		code := match(input[i])
		if code != no_result {
			new_i := i + 1 // skip current 'm' byte
			end := new_i - code_len(code)

			result = append(result, append(input[index:end], code_div(code)...)...)
			index = new_i
		}
	}

	result = append(result, input[index:len_input]...)

	return result
}

func main() { /*
		pm(cyan)
		pm(green)
		pm(white_0)
		pm(white_1)
		pm(grey)*/

	// start and end with a div
	// pm("<div>Welcome to " + cyan + " the server " + green + " today!</div>")

	file := "/Users/user/Desktop/tmp/log_messed.txt"

	content, err := ioutil.ReadFile(file)

	if err == nil {
		colored := pm(content)

		// before + colored + after
		data := append(append(html_before(), colored...), html_after()...)

		ioutil.WriteFile("/Users/user/Desktop/tmp/log_messed.html", data, os.FileMode(0644))
	}
}

func html_before() []byte {
	return []byte(`
<html>
<head>
<style>
* {
  padding: 0;
  margin: 0;
  border: 0;
  width: 0;
  height: 0;
}
div { display: inline; }
body { background-color: #262626; }
#terminal {
  display: inherit;
  white-space: pre;
  font-family: "Monaco";
  font-size: 14px;
  color: #f4f4f4;
  padding-left: 18px;
}
div.cyan  { color: #00eee9; }
div.green { color: #00e800; }
div.grey  { color: #666666; }
</style>
</head>
<body>
<div id="terminal">
<div>
	`)
}

func html_after() []byte {
	return []byte(`
</div>
</div>
</body>
</html>
`)
}
