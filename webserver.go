package main

import (
	"fmt"
	"net/http"
)

func startWebServer(newblockchain Blockchain) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.ServeFile(w, r, "index.html")
		}
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			// Clear inputList from session
			clearInputList("session-name", w, r)

			// Pass blockchain data to the session
			session, err := store.Get(r, "session-name")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			inputText := r.Form.Get("inputText")
			newblockchain.AddBlock(inputText)
			session.Values["resultList"] = newblockchain.Blocks
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Redirect to result page
			http.Redirect(w, r, "/result", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve blockchain data from the session
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		inputListInterface := session.Values["resultList"]
		blockchain, ok := inputListInterface.([]*Block)
		if !ok {
			http.Error(w, "No input list found in session or wrong type", http.StatusInternalServerError)
			return
		}

		// Render result.html with blockchain data
		err = templates.ExecuteTemplate(w, "result.html", blockchain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server listening on port 8080...")
	fmt.Println("open http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func clearInputList(sessionName string, w http.ResponseWriter, r *http.Request) {
	// Clear inputList from session
	session, err := store.New(r, sessionName)
	if err != nil {
		fmt.Println("Error creating new session:", err)
		return
	}
	session.Values["inputList"] = nil
	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Error saving session:", err)
		return
	}
}
