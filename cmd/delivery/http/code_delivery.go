package http

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yarikTri/network-channel-layer/cmd/code"
	"io"
	"log"
	"math/rand"
	"net/http"
)

const frameLoseProbability = 1
const frameErrorProbability = 7

const transferEndpoint = "http://localhost:8080/encoded-message/transfer"

// Code
// @Summary		Code network flow
// @Tags		Code
// @Description	Осуществляет кодировку сообщения в код Хэмминга [15, 11], внесение ошибки в каждый закодированный 15-битовый кадр с вероятностью 7%, исправление внесённых ошибок, раскодировку кадров в изначальное сообщение. Затем отправляет результат в Procuder-сервис транспортного уровня. Сообщение может быть потеряно с вероятностью 1%.
// @Accept		json
// @Param		request		body		[]byte		true	"Сообщение"
// @Success		200			{object}	nil					"Кодировка и отправка запущены"
// @Failure		400			{object}	nil					"Ошибка при чтении сообщения"
// @Router		/code [post]
func Code(c *gin.Context) {
	rawRequest, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Data(http.StatusBadRequest, c.ContentType(), []byte("Can't read request body"))
		return
	}

	go transfer(rawRequest)

	c.Data(http.StatusOK, c.ContentType(), []byte{})
}

func transfer(rawData []byte) {
	if rand.Intn(100) < frameLoseProbability {
		fmt.Println("[Info] Message lost")
		return
	}

	coder := code.Coder{}

	// Кодирование полученных данных
	encodedFrames, err := coder.Encode(rawData)
	if err != nil {
		log.Println("[Error] Encoding issue")
		return
	}

	// Внесение ошибок в закодированные кадры
	encodedFrames = coder.SetRandomErrors(encodedFrames, frameErrorProbability)

	// Исправление ошибок и декодирование кадров
	decodedFrames, err := coder.FixAndDecode(encodedFrames)
	if err != nil {
		log.Println("[Error] Decoding issue")
		return
	}

	// Валидация совпадения пришедших данных и выходных
	for ind, _byte := range rawData {
		if _byte != decodedFrames[ind] {
			log.Println("[Error] Frames inequality")
			return
		}
	}

	// Отправка данных в Producer
	req, err := http.NewRequest("POST", transferEndpoint, bytes.NewBuffer(decodedFrames))
	if err != nil {
		fmt.Printf("[Error] Can't create transfer request: %s\n", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[Error] Transfer request issue: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[Info] Unexpected status code while transferring %d\n", resp.StatusCode)
	}
}
