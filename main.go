// By izanbf1803
// http://izanbf.es/

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"./err"
)

var _passwd = flag.String("p", "", "Password to encrypt your message")
var _retr = flag.String("d", "", "File to decode")
var _hide = flag.String("e", "", "File to encode")
var _out = flag.String("o", "", "Output filename")
var _msg = flag.String("m", "", "Message to encode")
var pass string
var key []byte

const offset = 64
const streamEnd = "1111111111111110"
const passChars = "d9z36awvYLbAPIr2DhxdbnCwNXRYmG1v"

func init() {
	flag.Parse()
	pass = *_passwd + passChars[len(*_passwd):]
	key = []byte(pass)
}

func main() {
	fmt.Printf("\n\tCreated by izanbf1803 - http://izanbf.es/\n\n")
	if *_hide != "" {
		if *_msg == "" {
			err.Exit("Please, set a message using -m")
		}
		encode(*_hide, *_msg, *_out)
	} else if *_retr != "" {
		res := retr(*_retr)
		fmt.Printf("TEXT RETRIVED:\n%s\n\n", res)
	} else {
		err.Exit("Please, set an action using -e/-d")
	}
}

func encrypt(text string) string {
	plaintext := []byte(text)
	block, e := aes.NewCipher(key)
	if e != nil {
		err.Exit("1 - %v", e)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, e := io.ReadFull(rand.Reader, iv); e != nil {
		err.Exit("2 - %v", e)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decrypt(cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, e := aes.NewCipher(key)
	if e != nil {
		err.Exit("3 - %v", e)
	}
	if len(ciphertext) < aes.BlockSize {
		err.Exit("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return fmt.Sprintf("%s", ciphertext)
}

func strRepeat(s string, times int) string {
	var res string
	for i := 0; i < times; i++ {
		res += s
	}
	return res
}

func rgb2hex(rgb [3]byte) string {
	return fmt.Sprintf("%02x%02x%02x", rgb[0], rgb[1], rgb[2])
}

func hex2rgb(hex string) [3]byte {
	var res [3]byte
	for i := 0; i < 3; i++ {
		num, _ := strconv.ParseInt(hex[0:2*(i+1)], 16, 32)
		res[i] = byte(num)
	}
	return res
}

func str2bin(txt string) string {
	var res, temp string
	for i := 0; i < len(txt); i++ {
		temp = fmt.Sprintf("%s%b", temp, txt[i])
		res += strRepeat("0", 8-len(temp)) + temp
		temp = ""
	}
	return res
}

func bin2str(bin string) string {
	var res string
	for i := 0; i < len(bin); i += 8 {
		num, _ := strconv.ParseInt(bin[i:i+8], 2, 32)
		res += string(byte(num))
	}
	return res
}

func readImage(fn string) [][3]byte {
	file, e := ioutil.ReadFile(fn)
	data := make([][3]byte, len(file)/3)
	if e != nil {
		err.Exit("Can't open \"%s\"", fn)
	}
	var index int
	for i := 0; i < len(file)-2; i += 3 {
		data[index] = [3]byte{file[i], file[i+1], file[i+2]}
		index++
	}
	return data
}

func writeImage(fn string, data [][3]byte) {
	buffer := make([]byte, len(data)*3, len(data)*3)
	var index int
	for i := range data {
		for j := 0; j < 3; j++ {
			buffer[index] = data[i][j]
			index++
		}
	}
	e := ioutil.WriteFile(fn, buffer, 0644)
	if e != nil {
		err.Exit("Can't write \"%s\"", fn)
	}
}

func setLSB(rgb string, digit byte) string {
	ld := rgb[len(rgb)-1:] // Last digit
	if ld == "0" || ld == "1" || ld == "2" || ld == "3" || ld == "4" || ld == "5" {
		return rgb[:len(rgb)-1] + string(digit)
	}
	return ""
}

func getLSB(rgb string) string {
	ld := rgb[len(rgb)-1:] // Last digit
	if ld == "0" || ld == "1" {
		return ld
	}
	return ""
}

func encode(fn string, msg string, out string) {
	var e error
	if pass != "" {
		msg = encrypt(msg)
		if e != nil {
			err.Exit("Can't crypt message")
		}
	}

	img := readImage(fn)
	if len(msg)*8 > len(img)-offset {
		err.Exit("Image too small")
	}

	binary := str2bin(msg) + streamEnd
	var index int

	for i := offset; i < len(img); i++ {
		if index < len(binary) {
			newByte := setLSB(rgb2hex(img[i]), binary[index])
			if newByte != "" {
				img[i] = hex2rgb(newByte)
				index++
			}
		}
	}

	if out == "" {
		out = "__crypted__" + fn
	}
	writeImage(out, img)

	fmt.Printf("\n\n\tMessage crypted and saved in \"%s\"\n\n", out)
}

func retr(fn string) string {
	var e error

	img := readImage(fn)

	var binary, digit, msg string

	for i := offset; i < len(img); i++ {
		digit = getLSB(rgb2hex(img[i]))
		if digit != "" {
			binary += digit
			if len(binary) > len(streamEnd) {
				if binary[len(binary)-len(streamEnd):] == streamEnd {
					msg = bin2str(binary[:len(binary)-len(streamEnd)])
					break
				}
			}
		}
	}

	if pass != "" {
		msg = decrypt(msg)
		if e != nil {
			err.Exit("Can't decrypt message: [%v]", e)
		}
	}

	return msg
}
