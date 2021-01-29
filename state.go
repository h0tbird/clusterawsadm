package main

//-----------------------------------------------------------------------------
// Imports
//-----------------------------------------------------------------------------

import (

	// stdlib
	"bytes"
	"encoding/json"
	"io"
	"os"
)

//-----------------------------------------------------------------------------
// Types
//-----------------------------------------------------------------------------

type state struct {
	stateDir string
}

//-----------------------------------------------------------------------------
// State implementation
//-----------------------------------------------------------------------------

func (s *state) Init() (err error) {

	// Define the state directory
	s.stateDir = os.Getenv("HOME") + "/.clusterawsadm"

	// Create if not exists
	if _, err = os.Stat(s.stateDir); os.IsNotExist(err) {
		return os.MkdirAll(s.stateDir, os.ModePerm)
	}

	return err
}

func (s *state) Read(logicalID string, state interface{}) error {

	// Open a file handler
	f, err := os.Open(s.stateDir + "/" + logicalID + ".json")
	if err != nil {

		// No file means no state
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	// Unmarshal json
	return json.NewDecoder(f).Decode(state)
}

func (s *state) Write(logicalID string, state interface{}) error {

	// Open a file handler
	f, err := os.Create(s.stateDir + "/" + logicalID + ".json")
	if err != nil {
		return err
	}
	defer f.Close()

	// Marshal json
	jsonBytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	// Write to disk
	_, err = io.Copy(f, bytes.NewReader(jsonBytes))
	return err
}
