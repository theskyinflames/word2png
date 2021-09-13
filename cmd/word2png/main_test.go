package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
