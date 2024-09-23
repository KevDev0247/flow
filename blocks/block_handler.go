package blocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
)

type Block struct {
	Id       string `json:"id"`
	View     int    `json:"view"`
	accepted bool

	/* Future features */
	// Timestamp
	// Content
}

type Vote struct {
	BlockId string `json:"block_id"`

	/* Future features */
	// Timestamp
	// Votee
}

type RequestData struct {
	Block json.RawMessage `json:"block,omitempty"`
	Vote  json.RawMessage `json:"vote,omitempty"`
}

type BlockHandler struct {
	blocks       map[string][]Block
	locksMap     map[string]*sync.RWMutex
	locksMapLock sync.RWMutex
}

func NewBlockHandler() *BlockHandler {
	return &BlockHandler{
		blocks:   make(map[string][]Block),
		locksMap: make(map[string]*sync.RWMutex),
	}
}

func (h *BlockHandler) handleVote(vote Vote, block Block) {
	// Retrieve and lock block access
	blockLock := h.getOrCreateLockForBlock(block.Id)
	blockLock.Lock()
	defer blockLock.Unlock()

	voteObserved := vote.BlockId == block.Id
	if voteObserved {
		block.accepted = true
	}	
	if _, ok := h.blocks[vote.BlockId]; ok {
		allBlocksAccepted := false
		blocks := h.blocks[block.Id]
		for _, block := range blocks {
			allBlocksAccepted = block.accepted
		}

		h.blocks[block.Id] = append(h.blocks[block.Id], block)
		if voteObserved && allBlocksAccepted {
			// accept
			blocks = h.blocks[block.Id]
			sort.Slice(blocks, func(i, j int) bool {
				return blocks[i].View < blocks[j].View // Sort in ascending order
			})

			fmt.Printf("All accepted blocks:\n")
			for _, block := range blocks {
				fmt.Printf("Block %s, View %d\n", block.Id, block.View)
			}
		}
	} else {
		// Block doesn't exist
		h.blocks[block.Id] = append(h.blocks[block.Id], block)

		if voteObserved {
			blocks := h.blocks[block.Id]
			sort.Slice(blocks, func(i, j int) bool {
				return blocks[i].View < blocks[j].View // Sort in ascending order
			})

			fmt.Printf("All accepted blocks:\n")
			for _, block := range blocks {
				fmt.Printf("Block %s, View %d\n", block.Id, block.View)
			}
		}
	}
}

func (h *BlockHandler) Process(w http.ResponseWriter, r *http.Request) {
	var reqData RequestData

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Decoding block data
	var block Block
	err = json.Unmarshal(reqData.Block, &block)
	if err != nil {
		http.Error(w, "Error decoding Block", http.StatusBadRequest)
		return
	}
	block.accepted = false

	// Decoding vote data
	var vote Vote
	err = json.Unmarshal(reqData.Vote, &vote)
	if err != nil {
		http.Error(w, "Error decoding Vote", http.StatusBadRequest)
		return
	}

	h.handleVote(vote, block)
}

func (h *BlockHandler) getOrCreateLockForBlock(blockId string) *sync.RWMutex {
	// Use a mutex for locking access to locksMap
	h.locksMapLock.Lock()
	defer h.locksMapLock.Unlock()

	// Return lock if the lock for the block ID exists
	if mutex, exists := h.locksMap[blockId]; exists {
		return mutex
	}

	// Create a new lock and store it in the map
	newMutex := &sync.RWMutex{}
	h.locksMap[blockId] = newMutex
	return newMutex
}
