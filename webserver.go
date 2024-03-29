package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func startWebServer(newblockchain Blockchain) {
	http.HandleFunc("/upload", uploadImage)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		http.ServeFile(w, r, "static/js/script.js")
	})

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "static/css/style.css")
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

	http.HandleFunc("/imagepath", imagePathHandler)

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

func uploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the image file from the form data
	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error retrieving image file:", err)
		http.Error(w, "Error retrieving image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file on the server to store the uploaded image
	f, err := os.Create("./uploaded_image.jpg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the image file data to the newly created file
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("Error copying file data:", err)
		http.Error(w, "Error copying file data", http.StatusInternalServerError)
		return
	}

	// Pass the uploaded image path to the Python script for processing
	err = processImage("./uploaded_image.jpg")
	if err != nil {
		fmt.Println("Error processing image:", err)
		http.Error(w, "Error processing image", http.StatusInternalServerError)
		return
	}

	fmt.Println("Image uploaded and processed successfully")
	w.Write([]byte("Image uploaded and processed successfully"))
}

func processImage(imagePath string) error {
	// Get the absolute path to the Python script
	scriptPath, err := filepath.Abs("process_image.py")
	if err != nil {
		return err
	}

	// Execute the Python script as a separate process, passing the image path as an argument
	cmd := exec.Command("python", scriptPath, imagePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing Python script:", err)
	}

	// Print the output
	fmt.Println("Output of Python script:")
	fmt.Println(string(output))

	triggerEndpoint()
	return nil
}

func getImageBlob() (string, error) {
	// Get the directory of the Go file
	dir, err := filepath.Abs(filepath.Dir("./"))
	if err != nil {
		return "", err
	}

	// Construct the path to the image file
	imagePath := filepath.Join(dir, "adversarial_image.jpg") // Replace "image.jpg" with the name of your image file

	// Read the image file
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		fmt.Printf("oh nooooo 0 " + imagePath)
		return "", err
	}

	// Convert the image data to base64 encoding
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	// Create a data URL for the image blob
	blobURL := fmt.Sprintf("data:image/jpeg;base64,%s", base64Data)

	return blobURL, nil
}

func imagePathHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the image blob URL
	blobURL, err := getImageBlob()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("oh nooooo 1")
		return
	}

	// Encode the blob URL as JSON
	jsonResponse, err := json.Marshal(map[string]string{"blobURL": blobURL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("oh nooooo")
		return
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.Write(jsonResponse)
}

func triggerEndpoint() error {
	// Send a GET request to the endpoint
	resp, err := http.Get("http://localhost:8080/imagepath")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
