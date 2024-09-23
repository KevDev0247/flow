package blocks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockHandler(t *testing.T) {
	handler := NewBlockHandler()

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.blocks)
	assert.NotNil(t, handler.locksMap)
}

func TestBlockHandler_Process_ValidBlockAndVote(t *testing.T) {
	handler := NewBlockHandler()

	block := Block{
		Id:   "block-1",
		View: 1,
	}
	vote := Vote{
		BlockId: "block-1",
	}

	reqData := RequestData{
		Block: marshalJSON(t, block),
		Vote:  marshalJSON(t, vote),
	}

	body, _ := json.Marshal(reqData)
	req := httptest.NewRequest("POST", "/process", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Process(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Check if block was added
	blocks := handler.blocks["block-1"]
	assert.Len(t, blocks, 1)
	assert.Equal(t, block.Id, blocks[0].Id)
	assert.Equal(t, block.View, blocks[0].View)
	assert.True(t, blocks[0].accepted)
}

func TestBlockHandler_Process_InvalidBlock(t *testing.T) {
	handler := NewBlockHandler()

	vote := Vote{
		BlockId: "block-1",
	}

	reqData := RequestData{
		Block: json.RawMessage(`invalid`),
		Vote:  marshalJSON(t, vote),
	}

	body, _ := json.Marshal(reqData)
	req := httptest.NewRequest("POST", "/process", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Process(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestBlockHandler_Process_InvalidVote(t *testing.T) {
	handler := NewBlockHandler()

	block := Block{
		Id:   "block-1",
		View: 1,
	}

	reqData := RequestData{
		Block: marshalJSON(t, block),
		Vote:  json.RawMessage(`invalid`),
	}

	body, _ := json.Marshal(reqData)
	req := httptest.NewRequest("POST", "/process", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Process(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func marshalJSON(t *testing.T, v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	return data
}