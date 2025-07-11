package server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"slai.io/takehome/pkg/common"
	"slai.io/takehome/pkg/fileutils"
)

func HandleEcho(msg []byte, client *Client) error {
	log.Println("Received ECHO request.")

	var request common.EchoRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid echo request.")
	}

	response := &common.EchoResponse{
		BaseResponse: common.BaseResponse{
			RequestId:   request.RequestId,
			RequestType: request.RequestType,
		},
		Value: request.Value,
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}

	return nil
}

func HandleSync(msg []byte, client *Client, syncDir string) error {
	log.Println("Recieved SYNC request")

	var request common.SyncRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid SYNC request.")
	}
	errlog := ""
	fullPath := filepath.Join(syncDir, request.FileOp.FileName)
	// request should be of type CREATE, UPDATE, or DELETE
	switch request.FileOp.OpType {
		case common.CREATE, common.UPDATE: {
			parent := filepath.Dir(fullPath)

			// create necessary folders for file placement
			if err := os.MkdirAll(parent, 0755); err != nil {
            	return err
        	}

			log.Printf("File OpType: %v\n", request.FileOp.OpType)
			data, err := base64.RawStdEncoding.DecodeString(string(request.Value))

			if err != nil {
				errlog = "Cannot decode file"
				log.Println("Cannot decode file.")
			}

			err = os.WriteFile(fullPath, data, 0644)

			if err != nil {
				if errlog == "" {
					errlog = "Cannot write file"
				}
				//TODO: add functionality to create directories if necessary
				log.Printf("Cannot write file: %v\n", err)
			}
		}
		case common.DELETE: {
			err := os.Remove(fullPath)
			if err != nil {
				errlog = "Cannot delete file"
				log.Printf("Cannot delete file: %v\n", err)
			}
			go fileutils.CleanDirs(syncDir, filepath.Dir(fullPath))
		}
	}
	value := ""
	if errlog == "" {
		value = "success"
	} else {
		value = errlog
	}
	// send response
	response := &common.SyncResponse{
		BaseResponse: common.BaseResponse{
			RequestId: request.RequestId,
			RequestType: request.RequestType,
		},
		Value: value,
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}
	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}
	return nil
}