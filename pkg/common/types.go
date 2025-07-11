package common

type RequestType string
type FileOpType string

const (
	Echo RequestType = "ECHO"
	Sync RequestType = "SYNC"
)

const (
	UPDATE FileOpType = "UPDATE"
	CREATE FileOpType = "CREATE"
	DELETE FileOpType = "DELETE"
)

type FileOperation struct {
	OpType FileOpType `json:"op_type"`
	FileName   string `json:"file_name"`
}

type BaseRequest struct {
	RequestId   string `json:"request_id"`
	RequestType string `json:"request_type"`
}

type BaseResponse struct {
	RequestId   string `json:"request_id"`
	RequestType string `json:"request_type"`
}

type EchoRequest struct {
	BaseRequest
	Value string
}

type EchoResponse struct {
	BaseResponse
	Value string
}

// for SYNC requests

type SyncRequest struct {
	BaseRequest
	FileOp FileOperation `json:"file_op"`
	Value string
}

type SyncResponse struct {
	BaseResponse
	Value string
}