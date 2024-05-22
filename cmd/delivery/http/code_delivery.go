package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yarikTri/network-channel-layer/cmd/code"
	"log"
	"math/rand"
	"net/http"
)

const frameLoseProbability = 1
const frameErrorProbability = 7

const transferEndpoint = "http://localhost:8080/encoded-message/transfer"

type CodeRequest struct {
	Sender        string `json:"sender"`
	Timestamp     uint32 `json:"timestamp"`
	PartMessageID uint32 `json:"part_message_id"`
	Total         uint32 `json:"total"`
	Message       string `json:"message"`
}

type CodeTransferRequest struct {
	Sender        string `json:"sender"`
	Timestamp     uint32 `json:"timestamp"`
	PartMessageID uint32 `json:"part_message_id"`
	Total         uint32 `json:"total"`
	Message       string `json:"message"`
	FlagError     bool   `json:"flag_error"`
}

// Code
// @Summary		Code network flow
// @Tags		Code
// @Description	Осуществляет кодировку сообщения в код Хэмминга [15, 11], внесение ошибки в каждый закодированный 15-битовый кадр с вероятностью 7%, исправление внесённых ошибок, раскодировку кадров в изначальное сообщение. Затем отправляет результат в Procuder-сервис транспортного уровня. Сообщение может быть потеряно с вероятностью 1%.
// @Accept		json
// @Produce     json
// @Param		request		body		CodeRequest		true	"Информация о сегменте сообщения"
// @Success		200			{object}	nil					"Обработка и отправка запущены"
// @Failure		400			{object}	nil					"Ошибка при чтении сообщения"
// @Router		/code [post]
func Code(c *gin.Context) {
	var codeRequest CodeRequest
	if err := c.Bind(&codeRequest); err != nil {
		fmt.Println(err.Error())
		c.Data(http.StatusBadRequest, c.ContentType(), []byte("Can't read request body"))
		return
	}

	fmt.Printf("Got request: %v\n", codeRequest)

	go transfer(codeRequest)

	c.Data(http.StatusOK, c.ContentType(), []byte{})
}

func transfer(codeRequest CodeRequest) {
	if rand.Intn(100) < frameLoseProbability {
		fmt.Println("[Info] Message lost")
		return
	}

	processedMessage, hasErrors, err := processMessage([]byte(codeRequest.Message))
	if err != nil {
		fmt.Printf("[Error] Error while processing message: %s\n", err.Error())
		return
	}

	// Отправка данных в Producer
	transferReqBody, err := json.Marshal(
		CodeTransferRequest{
			Sender:        codeRequest.Sender,
			Timestamp:     codeRequest.Timestamp,
			PartMessageID: codeRequest.PartMessageID,
			Total:         codeRequest.Total,
			Message:       string(processedMessage),
			FlagError:     hasErrors,
		},
	)
	if err != nil {
		fmt.Printf("[Error] Can't create transfer request: %s\n", err.Error())
		return
	}

	req, err := http.NewRequest("POST", transferEndpoint, bytes.NewBuffer(transferReqBody))
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
		fmt.Printf("[Info] Unexpected status code [%d] while transferring: %s\n", resp.StatusCode, resp.Body)
		return
	}
	fmt.Println("[Info] Message successfully transferred")
}

// Возвращает декодированные биты (в случае успеха), флаг ошибки декодирования
// и ошибку исполнения при наличии
func processMessage(message []byte) ([]byte, bool, error) {
	coder := code.Coder{}

	// Кодирование полученных данных
	encodedFrames, err := coder.Encode(message)
	if err != nil {
		log.Println("[Error] Encoding issue")
		return nil, false, errors.New(fmt.Sprintf("Enocding issue: %s", err.Error()))
	}

	// Внесение ошибок в закодированные кадры
	encodedFrames = coder.SetRandomErrors(encodedFrames, frameErrorProbability)

	// Исправление ошибок и декодирование кадров
	decodedFrames, err := coder.FixAndDecode(encodedFrames)
	if err != nil {
		log.Println("[Error] Decoding issue")
		return nil, false, errors.New(fmt.Sprintf("Decoding issue: %s", err.Error()))
	}

	// Валидация совпадения пришедших данных и выходных
	hasError := false // Что-то вроде флага 500-ки при декодировании сообщения
	for ind, _byte := range message {
		if _byte != decodedFrames[ind] {
			log.Println("[Error] Frames inequality")
			hasError = true
			break
		}
	}

	return decodedFrames, hasError, nil
}
