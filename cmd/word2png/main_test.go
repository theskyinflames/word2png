package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadWordsFromFile(t *testing.T) {
	t.Parallel()

	t.Run("reads entire file content as a single word", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "words.txt")
		content := "word1\nword2\nword3"
		err := os.WriteFile(path, []byte(content), 0o644)
		require.NoError(t, err)

		words, err := readWordsFromFile(path)
		require.NoError(t, err)
		require.Equal(t, []string{content}, words)
	})

	t.Run("trims surrounding whitespace", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "words.txt")
		err := os.WriteFile(path, []byte("  hello world  \n"), 0o644)
		require.NoError(t, err)

		words, err := readWordsFromFile(path)
		require.NoError(t, err)
		require.Equal(t, []string{"hello world"}, words)
	})

	t.Run("returns error on non-existent file", func(t *testing.T) {
		_, err := readWordsFromFile("/nonexistent/path.txt")
		require.Error(t, err)
	})

	t.Run("returns error on oversized file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "big.txt")
		data := make([]byte, maxFileWordsSize+1)
		err := os.WriteFile(path, data, 0o644)
		require.NoError(t, err)

		_, err = readWordsFromFile(path)
		require.Error(t, err)
	})

	t.Run("returns empty slice for empty file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "empty.txt")
		err := os.WriteFile(path, []byte{}, 0o644)
		require.NoError(t, err)

		words, err := readWordsFromFile(path)
		require.NoError(t, err)
		require.Empty(t, words)
	})
}

func TestRemoveWordsByIdx(t *testing.T) {
	t.Parallel()
	t.Run(`Given a list of words and a negative index to be removed,
		   when it's called, 
		   then no words are removed`, func(t *testing.T) {
		words := []string{"w1", "w2", "w3"}
		rmIdxs := []int{-1}
		require.Equal(t, words, RemoveWordsByIdx(words, rmIdxs))
	})

	t.Run(`Given a list of words and an index to be removed bigger than list's size,
		   when it's called, 
		   then no words are removed`, func(t *testing.T) {
		words := []string{"w1", "w2", "w3"}
		rmIdxs := []int{4}
		require.Equal(t, words, RemoveWordsByIdx(words, rmIdxs))
	})

	t.Run(`Given a list of words and a valid index to be removed,
	when it's called, 
	then the word with that index has bee removed from the list`, func(t *testing.T) {
		words := []string{"w1", "w2", "w3"}
		expected := []string{"w1", "w3"}
		rmIdxs := []int{2}
		require.Equal(t, expected, RemoveWordsByIdx(words, rmIdxs))
	})
}
