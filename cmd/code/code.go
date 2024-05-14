package code

import "math/rand"

/*
 Кодирование/декодирование кодом Хэмминга [15, 11]
*/
rgegrhreth
const ByteSize = 8

const RawFrameSize = 11     // 11 бит - кадр до кодирования
const EncodedFrameSize = 15 // 15 бит - кадр после кодирования

// Все возможные индексы контрольных битов кода Хэмминга [15, 11]
var controlBitsIndexes = 
}

type Coder struct {
	garbageTailBits int
}

// Контрольные биты:
// 1 - нечётное количество единиц в контролируемых битах
// 0 - чётное
func (c *Coder) controlBit(num uint16, index uint) uint16 {
	controlledBits := num & controlledBitsMasks[index]
	oddCounter := 0
	for ind := EncodedFrameSize - 1; ind >= 0; ind-- {
		if (controlledBits>>ind)&1 == 1 {
			oddCounter++
		}
	}

	return uint16(oddCounter % 2)
}

// SetRandomErrors инвертирует один случайный бит в каждом кадре
// с вероятностью `probabilityPerFrame`
func (c *Coder) SetRandomErrors(frames []uint16, probabilityPerFrame int) []uint16 {
	for frameInd, frame := range frames {
		if rand.Intn(100) < probabilityPerFrame {
			errorFrame := InvertedAtShift(frame, uint(rand.Intn(15)))
			frames[frameInd] = errorFrame
		}
	}

	return frames
}

// InvertedAtShift возвращает `frame` с инвертированным битом в `shift` разряде
func InvertedAtShift(frame uint16, shift uint) uint16 {
	bitValue := (frame >> shift) & 1

	if bitValue == 0 {
		return frame | (1 << shift)
	}
	return frame & ^(1 << shift)
}
