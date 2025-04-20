document.addEventListener('DOMContentLoaded', () => {
    const uploadForm = document.getElementById('uploadForm');
    const fileInput = document.getElementById('fileInput');
    const dropZone = document.getElementById('dropZone');
    const previewContainer = document.getElementById('previewContainer');
    const imagePreview = document.getElementById('imagePreview');
    const resultsSection = document.getElementById('resultsSection');
    const statusMessage = document.getElementById('statusMessage');
    const resultContainer = document.getElementById('resultContainer');
    const processedImage = document.getElementById('processedImage');
    const downloadLink = document.getElementById('downloadLink');

    // Handle file selection
    dropZone.addEventListener('click', () => {
        fileInput.click();
    });

    fileInput.addEventListener('change', handleFileSelect);

    // Handle drag and drop
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    ['dragenter', 'dragover'].forEach(eventName => {
        dropZone.addEventListener(eventName, highlight, false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, unhighlight, false);
    });

    function highlight() {
        dropZone.classList.add('dragover');
    }

    function unhighlight() {
        dropZone.classList.remove('dragover');
    }

    dropZone.addEventListener('drop', handleDrop, false);

    function handleDrop(e) {
        const dt = e.dataTransfer;
        const files = dt.files;
        handleFiles(files);
    }

    function handleFileSelect(e) {
        const files = e.target.files;
        handleFiles(files);
    }

    function handleFiles(files) {
        if (files.length > 0) {
            const file = files[0];
            if (file.size > 10 * 1024 * 1024) {
                showError('File size exceeds 10MB limit');
                return;
            }
            if (!file.type.startsWith('image/')) {
                showError('Please select an image file');
                return;
            }
            previewFile(file);
        }
    }

    function previewFile(file) {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onloadend = () => {
            imagePreview.src = reader.result;
            previewContainer.classList.remove('hidden');
        };
    }

    function showError(message) {
        dropZone.classList.add('error');
        statusMessage.textContent = message;
        statusMessage.classList.add('text-red-500');
        setTimeout(() => {
            dropZone.classList.remove('error');
            statusMessage.textContent = '';
            statusMessage.classList.remove('text-red-500');
        }, 3000);
    }

    // Handle form submission
    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const file = fileInput.files[0];
        if (!file) {
            showError('Please select a file first');
            return;
        }

        const formData = new FormData();
        formData.append('file', file);

        // Show processing status
        resultsSection.classList.remove('hidden');
        resultContainer.classList.add('hidden');
        statusMessage.innerHTML = '<div class="spinner"></div><p class="mt-2">Processing your image...</p>';
        statusMessage.classList.add('processing');

        try {
            // Upload the file
            const uploadResponse = await fetch('http://localhost:8080/upload', {
                method: 'POST',
                body: formData,
                mode: 'cors'
            });

            if (!uploadResponse.ok) {
                throw new Error('Upload failed: ' + uploadResponse.statusText);
            }

            const { photo_id } = await uploadResponse.json();
            console.log('Received photo_id:', photo_id);

            // Poll for the result
            let processed = false;
            let attempts = 0;
            const maxAttempts = 30; // 30 attempts with 2-second intervals = 1 minute timeout

            while (!processed && attempts < maxAttempts) {
                await new Promise(resolve => setTimeout(resolve, 2000));
                
                const resultResponse = await fetch(`http://localhost:8080/images/${photo_id}`, {
                    mode: 'cors'
                });
                
                if (resultResponse.ok) {
                    const { url } = await resultResponse.json();
                    console.log('Received URL:', url);
                    
                    if (url) {
                        // Ensure URL has http:// prefix
                        const fullUrl = url.startsWith('http://') ? url : `http://${url}`;
                        console.log('Using URL:', fullUrl);

                        // Test if the image is accessible
                        try {
                            const img = new Image();
                            img.crossOrigin = 'anonymous'; // Add CORS attribute
                            
                            // Add error handling for CORS
                            img.onerror = (e) => {
                                console.error('Image load error:', e);
                                throw new Error('Failed to load image due to CORS or other error');
                            };

                            // Add timeout for image loading
                            const timeout = setTimeout(() => {
                                img.src = ''; // Cancel loading
                                throw new Error('Image loading timed out');
                            }, 10000); // 10 seconds timeout

                            img.onload = () => {
                                clearTimeout(timeout);
                                processed = true;
                                processedImage.src = fullUrl;
                                downloadLink.href = fullUrl;
                                resultContainer.classList.remove('hidden');
                                statusMessage.classList.remove('processing');
                                statusMessage.innerHTML = '<p class="text-green-500">Image processed successfully!</p>';
                            };

                            img.src = fullUrl;
                        } catch (error) {
                            console.error('Error testing image URL:', error);
                            throw error;
                        }
                    }
                } else {
                    console.error('Error fetching result:', resultResponse.statusText);
                }
                attempts++;
            }

            if (!processed) {
                throw new Error('Processing timed out after 1 minute');
            }
        } catch (error) {
            console.error('Error:', error);
            statusMessage.classList.remove('processing');
            statusMessage.innerHTML = `<p class="text-red-500">Error: ${error.message}</p>`;
        }
    });
}); 