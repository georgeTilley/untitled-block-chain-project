document.addEventListener('DOMContentLoaded', function() {
    const imageUpload = document.getElementById('imageUpload');
    const originalImageElement = document.getElementById('originalImage');

    document.getElementById('uploadForm').addEventListener('submit', function(event) {
        event.preventDefault(); // Prevent default form submission
        if (imageUpload.files.length > 0) {
            const selectedImage = imageUpload.files[0];
    
            // Create URL for the selected image and display it
            const imageUrl = URL.createObjectURL(selectedImage);
            originalImageElement.src = imageUrl;
            generateAdversarialImage(imageUrl);
        }
    });
    
    function generateAdversarialImage(imagePath) {
        
        const formData = new FormData();
        formData.append('imagePath', imagePath);
    }
});

function uploadImage() {
    const fileInput = document.getElementById('fileInput');
    const file = fileInput.files[0];
    const originalImageElement = document.getElementById('originalImage');
    // Create URL for the selected image and display it
    const imageUrl = URL.createObjectURL(file);
    originalImageElement.src = imageUrl; 
    if (file) {
        const formData = new FormData();
        formData.append('image', file);

        fetch('/upload', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (response.ok) {
                console.log('Image uploaded successfully');
                // Fetch the generated images from the server
                // Make an AJAX request to retrieve the image path
                fetch('/imagepath')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    const imageURL = data.blobURL;
                    console.log("image " + imageURL)
                    const adversarialImageElement = document.getElementById('adversarialImage');
                    adversarialImageElement.src = imageURL;
                })
                .catch(error => console.error('Error fetching image path:', error));
                    } else {
                        console.error('Failed to upload image');
                    }
                })
        .catch(error => console.error('Error uploading image:', error));
    } else {
        console.error('No file selected');
    }
}
