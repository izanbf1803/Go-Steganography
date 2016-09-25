# Go-Steganography
> By <a href="http://izanbf.es/">izanbf1803</a>

#### Simple Go Steganography program without external libraries and support with multiple formats (.png, .jpg...)
#### Password cypher algorithm: AES

<br>

Example: 
```
go-steganography.exe -e "img.png" -m "Steganography it's amazing!" -p "sd4tst69"
```

<br>

Actions:
* <strong>-e</strong> (Encode) 
* <strong>-d</strong> (Decode)

<br><br>

Encode parameters:

| Parameter     | Description   |
| ------------- | ------------- |
| <strong>-e</strong> | File to encode |
| <strong>-m</strong> | Message to encode |
| <strong>-p</strong> | Password to encrypt your message |
| <strong>-o</strong> | Output filename |

<br>

Decode parameters:

| Parameter     | Description   |
| ------------- | ------------- |
| <strong>-d</strong> | File to decode |
| <strong>-p</strong> | Password to decrypt the message |
