package rag

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type VectorRecord struct {
	Id               string    `json:"id"`
	Prompt           string    `json:"prompt"`
	Embedding        []float64 `json:"embedding"`
	CosineSimilarity float64
}

type VectorStore interface {
	GetAll() ([]VectorRecord, error)
	Save(vectorRecord VectorRecord) (VectorRecord, error)
	SearchSimilarities(embeddingFromQuestion VectorRecord, limit float64) ([]VectorRecord, error)
	SearchTopNSimilarities(embeddingFromQuestion VectorRecord, limit float64, max int) ([]VectorRecord, error)
}

type MemoryVectorStore struct {
	Records map[string]VectorRecord
}

func (mvs *MemoryVectorStore) GetAll() ([]VectorRecord, error) {
	var records []VectorRecord
	for _, record := range mvs.Records {
		records = append(records, record)
	}
	return records, nil
}

// Save saves a vector record to the MemoryVectorStore.
// If the record does not have an ID, it generates a new UUID for it.
// It returns the saved vector record and an error if any occurred during the save operation.
// If the record already exists, it will be overwritten.
func (mvs *MemoryVectorStore) Save(vectorRecord VectorRecord) (VectorRecord, error) {
	if vectorRecord.Id == "" {
		vectorRecord.Id = uuid.New().String()
	}
	mvs.Records[vectorRecord.Id] = vectorRecord
	return vectorRecord, nil
}

// SearchSimilarities searches for vector records in the MemoryVectorStore that have a cosine distance similarity greater than or equal to the given limit.
//
// Parameters:
//   - embeddingFromQuestion: the vector record to compare similarities with.
//   - limit: the minimum cosine distance similarity threshold.
//
// Returns:
//   - []llm.VectorRecord: a slice of vector records that have a cosine distance similarity greater than or equal to the limit.
//   - error: an error if any occurred during the search.
func (mvs *MemoryVectorStore) SearchSimilarities(embeddingFromQuestion VectorRecord, limit float64) ([]VectorRecord, error) {

	var records []VectorRecord

	for _, v := range mvs.Records {
		distance := cosineSimilarity(embeddingFromQuestion.Embedding, v.Embedding)
		if distance >= limit {
			v.CosineSimilarity = distance
			records = append(records, v)
		}
	}
	return records, nil
}

// SearchTopNSimilarities searches for the top N similar vector records based on the given embedding from a question.
// It returns a slice of vector records and an error if any.
// The limit parameter specifies the minimum similarity score for a record to be considered similar.
// The max parameter specifies the maximum number of vector records to return.
func (mvs *MemoryVectorStore) SearchTopNSimilarities(embeddingFromQuestion VectorRecord, limit float64, max int) ([]VectorRecord, error) {
	records, err := mvs.SearchSimilarities(embeddingFromQuestion, limit)
	if err != nil {
		return nil, err
	}
	return getTopNVectorRecords(records, max), nil
}

func (mvs *MemoryVectorStore) Load(storeFilePath string) error {
	// Check if the store file exists
	if _, err := os.Stat(storeFilePath); os.IsNotExist(err) {
		return err
	}

	// Read the store file
	file, err := os.ReadFile(storeFilePath)
	if err != nil {
		return err
	}

	// Unmarshal the JSON into the vector store
	if err := json.Unmarshal(file, &mvs); err != nil {
		return err
	}

	return nil
}

func (mvs *MemoryVectorStore) Persist(storeFilePath string) error {
	// Marshal the store to JSON
	storeJSON, err := json.MarshalIndent(mvs, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON to a file

	err = os.WriteFile(storeFilePath, storeJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (mvs *MemoryVectorStore) ResetMemory() error {
	// Reset the vector store to a new empty MemoryVectorStore
	mvs.Records = make(map[string]VectorRecord)
	return nil
}
