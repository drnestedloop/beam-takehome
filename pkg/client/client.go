package client

import (
	"encoding/json"
	"log"

	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"slai.io/takehome/pkg/common"
	"slai.io/takehome/pkg/fileutils"
)

const maxConnectionAttempts = 100
const hostURL = "ws://localhost:5555/"

func init() {
}

type Client struct {
	Directory string
	SessionId string
	ws        *websocket.Conn
	connected bool
	hostURL   string
	channels  map[string]chan []byte
}

func NewClient(directory string) (*Client, error) {
	var client *Client = &Client{
		Directory: directory,
		hostURL:   hostURL,
	}

	err := client.connect()
	if err != nil {
		return nil, err
	}

	client.connected = true
	client.channels = make(map[string]chan []byte)

	return client, nil
}

func (c *Client) connect() error {
	connected := false
	attempts := 0

	for {
		log.Println("Connection attempt: ", attempts)

		if attempts > maxConnectionAttempts {
			break
		}

		ws, _, err := websocket.DefaultDialer.Dial(c.hostURL, nil)
		c.ws = ws

		if err != nil {
			attempts++
			continue
		}

		connected = true
		break
	}

	// We weren't able to connect to the host, bail
	if !connected {
		return nil
	}

	// Start receiving messages
	go c.rx()

	return nil
}

func (c *Client) rx() {
	for {
		_, message, err := c.ws.ReadMessage()
		if ce, ok := err.(*websocket.CloseError); ok {

			switch ce.Code {
			case websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure:
				return
			}
		}

		var msg common.BaseResponse

		err = json.Unmarshal(message, &msg)
		if err != nil {
			continue
		} else {
			if _, ok := c.channels[msg.RequestId]; ok {
				c.channels[msg.RequestId] <- message
			} else {
				log.Println("channel not found")
			}
		}
	}
}

func (c *Client) tx(msg []byte) error {
	err := c.ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}

	return nil
}

// Request implementations
func (r *Client) Echo(value string) (string, error) {
	requestId := uuid.NewString()

	var request *common.EchoRequest = &common.EchoRequest{
		BaseRequest: common.BaseRequest{
			RequestId:   requestId,
			RequestType: string(common.Echo),
		},
		Value: value,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	r.channels[requestId] = make(chan []byte)

	err = r.tx(payload)
	if err != nil {
		return "", err
	}

	var response common.EchoResponse = common.EchoResponse{}

	msg := <-r.channels[requestId]
	err = json.Unmarshal(msg, &response)
	if err != nil {
		log.Println("Unable to handle echo response: ", err)
		return "", err
	}

	return response.Value, err
}

func (r *Client) Sync(fileop common.FileOperation) (string, error) {
	requestId := uuid.NewString()
	// define value based on what the file op is
	var value string
	switch fileop.OpType {
		case common.CREATE, common.UPDATE: {
			//read file to memory and assign to value
			v, err := fileutils.FileSerializer(filepath.Join(r.Directory, fileop.FileName))
			if err != nil {
				return "", err
			}
			value = v
		}
		case common.DELETE:
			value = ""
	}
	var request *common.SyncRequest = &common.SyncRequest{
		BaseRequest: common.BaseRequest{
			RequestId: requestId,
			RequestType: string(common.Sync),
		},
		FileOp: fileop,
		Value: value,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	r.channels[requestId] = make(chan []byte)

	err = r.tx(payload)
	if err != nil {
		return "", err
	}

	var response common.SyncResponse = common.SyncResponse{}

	msg := <-r.channels[requestId]
	err = json.Unmarshal(msg, &response)
	if err != nil {
		log.Println("Unable to handle SYNC response: ", err)
		return "", err
	}

	return response.Value, err
}