package handler

import (
	"strconv"
	//	"bytes"
	"codegenerator"
	"conf"
	//	"encoding/base64"
	"facade/vo"
	"imagegenerator"
	"net/url"
	"redisclusteradapter"
	"strings"
	"syscommon"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gitstliu/log4go"
	"github.com/snluu/uuid"
)

var redisCommonKey = "VCODE:POSITINO:"

type VCodeFacade struct {
}

func (this *VCodeFacade) GetCode(w rest.ResponseWriter, r *rest.Request) {
	log4go.Debug("Get Code Start!!")

	fg := "FFFFFFFF"
	bg := "000000FF"
	scale := 1.0

	fgMeta, fgOK := r.URL.Query()["fg"]
	bgMeta, bgOK := r.URL.Query()["bg"]
	scaleMeta, scaleOK := r.URL.Query()["scale"]

	if fgOK {
		log4go.Debug("fgMeta[0] = %v", fgMeta[0])
		fg = fgMeta[0]
	}

	if bgOK {
		log4go.Debug("bgMeta[0] = %v", bgMeta[0])
		bg = bgMeta[0]
	}

	if scaleOK {
		log4go.Debug("scaleMeta[0] = %v", scaleMeta[0])
		scaleValue, scaleErr := strconv.ParseFloat(scaleMeta[0], 10)
		if scaleErr == nil {
			scale = scaleValue
		} else {
			log4go.Error(scaleErr)
			return
		}

	}

	key, noice := codegenerator.GetRandWords()
	log4go.Debug("key = %v, noice = %v", key, noice)

	positon, buffer, bufferErr := codegenerator.CreateNoiceConfusionCodeImage(key, noice, 2, 4, fg, bg, scale)

	if bufferErr != nil {
		log4go.Error(bufferErr)
		return
	}

	w.Header()["VCodeString"] = []string{url.QueryEscape(key)}
	id := uuid.Rand()
	idMeta := id.Hex()
	w.Header()["VCodeID"] = []string{idMeta}
	positionJson, positionJsonErr := syscommon.ObjectToJson(positon)

	if positionJsonErr != nil {
		log4go.Error(positionJsonErr)
		return
	}

	_, setPosiionErr := redisclusteradapter.GetAdapter().SETEX(getRedisKey(idMeta), []byte(positionJson), conf.GetConfigure().VCodeTimeOut)
	if setPosiionErr != nil {
		log4go.Error(setPosiionErr)
		return
	}

	//	encodeString := base64.StdEncoding.EncodeToString(buffer)
	//	log4go.Debug(encodeString)
	//	log4go.Debug(buffer)

	w.WriteBytes(buffer)
}

func (this *VCodeFacade) CheckCode(w rest.ResponseWriter, r *rest.Request) {

	log4go.Debug("Check Code Start!!")
	response := syscommon.CommonResponse{Code: syscommon.Success, Message: "Success"}
	//r.Header
	vcodeIDMeta, vcodeExist := r.Header["Vcodeid"]

	log4go.Debug("r.Header() = %v", r.Header)
	if !vcodeExist {
		response.Code = syscommon.Fail
		response.Message = "VCodeID do not exist"
		w.WriteJson(&response)
		return
	}

	if len(vcodeIDMeta) == 0 {
		response.Code = syscommon.Fail
		response.Message = "VCodeID length is invalid"
		w.WriteJson(response)
		return
	}
	vcodeID := vcodeIDMeta[0]

	log4go.Debug("vcodeID = %v", vcodeID)

	positionJson, getPositionErr := redisclusteradapter.GetAdapter().GET(getRedisKey(vcodeID))

	log4go.Debug("positionJson = %v", positionJson)
	if getPositionErr != nil {
		log4go.Debug(getPositionErr)
		response.Code = syscommon.Fail
		response.Message = "VCodeID invalid"
		w.WriteJson(response)
		return
	}

	if strings.TrimSpace(positionJson) == "" {
		response.Code = syscommon.Fail
		response.Message = "VCodeID length invalid"
		w.WriteJson(response)
		return
	}

	position := &imagegenerator.NoiceRunePosition{}
	positionErr := syscommon.JsonToObject(positionJson, position)
	if positionErr != nil {
		log4go.Debug(positionErr)
		response.Code = syscommon.Fail
		response.Message = "VCodeID decode invalid"
		w.WriteJson(response)
		return
	}

	currClickInfo := &vo.ClickInfo{}
	decodePayloadErr := r.DecodeJsonPayload(currClickInfo)

	log4go.Debug("currClickInfo = %v", currClickInfo)

	if decodePayloadErr != nil {
		log4go.Debug(decodePayloadErr)
		response.Code = syscommon.Fail
		response.Message = "ClickInfo structure is invalid!!"
		w.WriteJson(response)
		return
	}

	pointLen := len(currClickInfo.Points)
	positionLen := len(position.WordIndex)

	if pointLen != positionLen || pointLen == 0 {
		response.Code = syscommon.Fail
		response.Message = "ClickInfo Check Failed!!"
		w.WriteJson(response)
		return
	}

	log4go.Debug("position.LineCount = %v", position.LineCount)
	log4go.Debug("position.ColumnCount = %v", position.ColumnCount)
	log4go.Debug("position.Position = %v", position.Position)

	stepLine := currClickInfo.TotalImageHeight / position.LineCount
	stepColumn := currClickInfo.TotalImageWidth / position.ColumnCount

	log4go.Debug("stepLine = %v", stepLine)
	log4go.Debug("stepColumn = %v", stepColumn)

	checkResult := true

	log4go.Debug("position.WordIndex = %v", position.WordIndex)
	//	position.WordIndex

	for index, currKeyPos := range position.WordIndex {
		//		log4go.Debug("currKeyPos = %v", currKeyPos)
		currPoint := position.Position[currKeyPos]
		log4go.Debug("position.ColumnCount = %v", position.ColumnCount)
		//		actX := currKeyPos / position.ColumnCount
		//		actY := currKeyPos % position.ColumnCount
		actLine := currClickInfo.Points[index].Y / stepColumn
		actColumn := currClickInfo.Points[index].X / stepLine

		log4go.Debug("actLine = %v", actLine)
		log4go.Debug("actColumn = %v", actColumn)
		log4go.Debug("currPoint.X = %v", currPoint.X)
		log4go.Debug("currPoint.Y = %v", currPoint.Y)

		if currPoint.X != actLine || currPoint.Y != actColumn {
			checkResult = false
			break
		}
	}

	log4go.Debug("checkResult = %v", checkResult)
	if !checkResult {
		response.Code = syscommon.Fail
		response.Message = "ClickInfo Check Failed!!"
		w.WriteJson(response)
		return
	}
	w.WriteJson(response)
	return
}

func getRedisKey(key string) string {
	return strings.Join([]string{redisCommonKey, key}, "")
}
