package functions

import (

	"fmt"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand"
)
const (
 	SMALLHAM  = 26
	BIGHAM  = 31
	PARITYBITS  = 5
	MMC = 13 // mmc (8, 26)
	BYTE = 8
)

//EncodeFile
//DecodeFile
//Hamming
//DeHamming
//IsPowerOfTwo
//WriteBufferInFile
//WriteSize
//ReadSize
//ExtractBytes
//MakeMistakes

var indexMap = []int{0, 1, 3, 7, 15} //powers of 2


//-------------------------------------------------------------
func EncodeFile (file *os.File, path string) error {

	if file == nil {
		return errors.New("No File")
	}
	
	lastDotIndex := strings.LastIndex(path, ".")
	if lastDotIndex != -1 {
		path = path[:lastDotIndex]
	} else {
		return errors.New("File type not identified")
	}

	//Enormous file where hamming will be dumped to 
	path = path + ".hamming"
	mid, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer mid.Close()
	
	info, err := file.Stat()
	if err != nil {
		return err
	}
	size := int(info.Size())
	err = WriteSize(mid, size)
	if err != nil {
		return err
	}

	var readBuffer [MMC]byte
	var byteBuffer []byte //for each bit a byte
	
	for {

		n, err := file.Read(readBuffer[:])
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			} else {
				return errors.New("Error reading the file:" + err.Error())
			}
		}

		//for each byte, turns its bits into bytes and adds to the new buffer
		for _, value := range readBuffer[:n] { 
			for pos := BYTE - 1; pos >= 0; pos-- {
				bit := (value >> pos) & 1
				byteBuffer = append(byteBuffer, byte(bit))
			}
		}
		//In case buffer is not full (MMC), adds the rest
		remainder := len(byteBuffer) % SMALLHAM
		if remainder != 0 {
			padding := SMALLHAM - remainder
			for i := 0; i < padding; i++ {
				byteBuffer = append(byteBuffer, 0)
			}
		}
		
		/*For each SMALLHAM size sequence, encode with hamming 
		and writes in the file*/
		start := 0
		for i := 0; i < len(byteBuffer) / SMALLHAM; i++ {
			
			bits := byteBuffer[start: start + SMALLHAM]
			ham, err := Hamming(bits)
			if err != nil {
				return err
			}
			start += SMALLHAM
			for _, b := range ham {
    			if b == 0 {
        			mid.Write([]byte{'0'})
    			} else {
        			mid.Write([]byte{'1'})
    			}	
			}
			mid.Write([]byte{byte(' ')})
		}

		//Flush in the buffer
		byteBuffer = byteBuffer[:0]
		for i := range readBuffer {
    		readBuffer[i] = 0
		}
	}
	return nil

}

//-------------------------------------------------------------
func DecodeFile ( file *os.File, path string) error {

	if file == nil {
		return errors.New("No file")
	}

	//Final file which will be equal to the first one
	name := strings.TrimSuffix(path, ".hamming")
	if name == path {
		return errors.New("Path with wrong extension")
	}
	name = name + ".dec"
	copyFile, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer copyFile.Close()

	size, err := ReadSize(file)
	if err != nil {
		return err
	}
	

	var bigSequence [32]byte
	var byteBuffer []byte
	var lastChunkSize = (size % MMC) * BYTE

	currentSize := 0
	for {
		_, err := file.Read(bigSequence[:])
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			} else {
				return errors.New("Error reading the file:" + err.Error())
			}
		}

		//Decode each sequence and adds it to the buffer
		smallSequence, err := ExtractBytes(bigSequence)
		if err != nil {
			return err
		}
		byteBuffer = append(byteBuffer, smallSequence...)
		currentSize += SMALLHAM

		if (currentSize / BYTE) >= size {
			byteBuffer = byteBuffer[:lastChunkSize] 
			err := WriteBufferInFile(copyFile, &byteBuffer)
			if err != nil {
				return err
			}	
		}
		//Write in copyFile
		if len(byteBuffer) == MMC * BYTE {
			err := WriteBufferInFile(copyFile, &byteBuffer)
			if err != nil {
				return err
			}
		} 
		
	} 

	return nil
}
//-------------------------------------------------------------
func Hamming(smallSequence []byte) ([]byte, error) {

	if len(smallSequence) != SMALLHAM {
		return nil, errors.New("Slice does not have 26 bits")
	}
	
	var bigSequence [BIGHAM]byte

	//Insert bits from smallSequence to bigSequence, skipping powers of 2
	bitIndex := 0
	for i := range bigSequence {
		if !IsPowerOfTwo(i + 1) {
			bigSequence[i] = smallSequence[bitIndex]
			bitIndex++
		}
	}
	//Calculates parity bits
	for i := range bigSequence {
		if !IsPowerOfTwo(i + 1) {
			for pos, value := range indexMap {
				if (i + 1) & (1 << pos) != 0 {
					bigSequence[value] ^= bigSequence[i]
				}
			}
		}
	}

	return bigSequence[:], nil
}

//-------------------------------------------------------------
func DeHamming (bigSequence []byte) ([]byte, error) {

	if len(bigSequence) != BIGHAM {
		return nil, errors.New("Slice does not have 31 bits")
	}
	
	var smallSequence [SMALLHAM]byte
	var parityBits [PARITYBITS]byte

	//Calculate parity bits
	for i := range bigSequence {
			
		if !IsPowerOfTwo(i+1){
			for pos := range parityBits {

				if (i + 1) & ( 1 << pos) != 0 {
					parityBits[pos] ^= bigSequence[i]
				}
			}
		}
	}

	//Compare and find error position if any
	errorPos := 0
	for i := range parityBits {
		if parityBits[i] != bigSequence[ indexMap[i] ] {
			errorPos += indexMap[i] + 1 //sums the power 
		}
	}

	if errorPos != 0 {
		bigSequence[errorPos - 1] ^= 1
	}

	//Insert bigSequence values in smallSequence
	smallIndex := 0
	for pos, value := range bigSequence {

		if !IsPowerOfTwo(pos + 1){
			smallSequence[smallIndex] = value
			smallIndex++
		}
	}
	return smallSequence[:], nil

}

//-------------------------------------------------------------
func IsPowerOfTwo (n int) bool {

	if n <= 0 {
		return false
	}
	return (n & (n - 1)) == 0
}

//-------------------------------------------------------------
func WriteBufferInFile(file *os.File, byteBuffer *[]byte) error {
		
		if file == nil {
			return errors.New("No File")

		}

		start := 0
		var abyte byte = 0
		//Write each BYTE sequence as a byte in the file
		for start < len(*byteBuffer) {
			bitSlice := (*byteBuffer)[start: start + BYTE]
			start += BYTE
				
			for i := 0; i < BYTE; i++ {
				abyte = abyte << 1 
				abyte = abyte | bitSlice[i]
			}
			file.Write([]byte{abyte})
		}
		*byteBuffer = (*byteBuffer)[:0]
		return nil
	

}

//-------------------------------------------------------------
func WriteSize(file *os.File, size int) error {

	if file == nil {
		return errors.New("No File")

	}

	//Transform int into a sequence of 26 bits
	var bits [SMALLHAM]byte
	for i := SMALLHAM - 1; i >= 0; i-- {
		bits[i] = byte((size >> (SMALLHAM - 1 - i)) & 1)
	}

	//Encode the sequence with Hamming
	ham, err := Hamming(bits[:])
	if err != nil {
		return err
	}
	//Write in the file
	for _, b := range ham {
    	if b == 0 {
        	file.Write([]byte{'0'})
    	} else {
        	file.Write([]byte{'1'})
    	}	
	}
	file.Write([]byte{byte(' ')})
	return nil


}

//-------------------------------------------------------------
func ReadSize (file *os.File) (int,error) {

	if file == nil {
		return 0, errors.New("No File")

	}

	//Read the first 32 byte sequence from the file
	var buffer [32]byte
	_ ,err := file.Read(buffer[:])
	if err != nil {
		if err.Error() == "EOF" {
			return 0, err
		}
		return 0, nil
	}

	//Get the 26 sequence
	smallSequence, err := ExtractBytes(buffer)
	if err != nil {
		return 0, err
	}

	//Convert the bit sequence into an int
	var binaryString string
	for _, bit := range smallSequence {
		binaryString += fmt.Sprintf("%d", bit)
	}
	num, err := strconv.ParseInt(binaryString, 2, 0)
	if err != nil {
		return 0, err
	}
	return int(num), nil

}

//-------------------------------------------------------------
func ExtractBytes( buffer [32]byte) ([]byte, error) {

	//For each 32 bits, ignore the last and 
	sqnceEnd := len(buffer) - 1
	smallSequence, err := DeHamming(buffer[:sqnceEnd])
	if err != nil {
		return nil, err
	}
	
	for pos, value := range smallSequence {
		if value == 48 {
			smallSequence[pos] = 0
		} else if value == 49 {
			smallSequence[pos] = 1
		} else {
			return nil, errors.New("File with unknown codification")
		}
	}
	return smallSequence, nil
}
//-------------------------------------------------------------
func MakeMistakes(filename string) error {

	
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	x := 20

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	fileSequences := int(fileSize) / 32

	if x > fileSequences {
		return fmt.Errorf("File does not have enough bytes for %d sequences of 32 bits", x)
	}


	//Pick x random sequences from the file
	rand.Seed(time.Now().UnixNano())
	sequences := make([]int,fileSequences)
	for i := range sequences {
		sequences[i] = i 
	}
	rand.Shuffle(len(sequences), func(i, j int) {
		sequences[i], sequences[j] = sequences[j], sequences[i]
	})

	choosen := make([]int, x)
	for i := 0; i < x; i++ {

		/*In each choosen sequence, picks a random bit position 
		and changes its bit */
		
		choosen[i] = sequences[i]
		randomInt := rand.Intn(30-1+1) + 1
		position := int64(choosen[i] * 32 + randomInt)
		fmt.Println("Sequence", choosen[i],"Position",randomInt)
		_, err = file.Seek(position, 0)
		if err != nil {
			return nil
		}

		var byteRead [1]byte 
		_, err = file.Read(byteRead[:])
		if err != nil {
			return nil
		}

		fmt.Println("Byte read:", byteRead[0]) 

		_, err = file.Seek(position, 0)
		if err != nil {
			return nil
		}

		if byteRead[0] == 48 {
			_, err = file.Write([]byte{'1'})
		} else if byteRead[0] == 49 {
			_, err = file.Write([]byte{'0'})
		} else {
			return errors.New("Byte is not 1 nor 0")
		}

		if err != nil {
			return nil
		}
	}
	return nil
		
}

