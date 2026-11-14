/**
 * Immediately invoked function expression to encapsulate the repository search, filtering, and file loading module.
 */
(function () {
    // Map each repository search input identifier to its corresponding unique file tree container identifier.
    const searchInputMappingDictionary = {
        'c-tree-search': 'repository-tree-container',
        'python-tree-search': 'python-repository-tree-container',
        'golang-tree-search': 'golang-repository-tree-container',
        'under-the-hood-tree-search': 'under-the-hood-repository-tree-container'
    }

    /**
     * Filters repository tree files and folders so only matching files and relevant categories are displayed.
     * @param {HTMLElement} treeContainerElement - The container element holding the repository file tree.
     * @param {string} rawQueryString - The raw search input query string.
     */
    function filterRepositoryTree(treeContainerElement, rawQueryString) {
        // Guard against missing tree container elements to ensure defensive robustness.
        if (!treeContainerElement) {
            return
        }

        // Normalize the raw query string by converting all characters to lowercase and trimming outer whitespace.
        const normalizedQuery = (rawQueryString || '').toLowerCase().trim()
        
        // Clean the query string by removing all hyphens, underscores, spaces, and forward slashes.
        const cleanedQuery = normalizedQuery.replace(/[-_\s/]/g, '')
        
        // Query all folder wrapper elements within the file tree container.
        const folderWrapperElements = treeContainerElement.querySelectorAll('div.mb-1, .folder-item')

        // Guard against missing folder wrapper elements.
        if (!folderWrapperElements || folderWrapperElements.length === 0) {
            return
        }

        // If the query string is empty, reset and restore all folder wrappers and file items to their default visible display state.
        if (!normalizedQuery) {
            folderWrapperElements.forEach(function (folderWrapperElement) {
                if (folderWrapperElement && folderWrapperElement.style) {
                    folderWrapperElement.style.display = 'block'
                }

                const innerContainer = folderWrapperElement.querySelector('.ml-3')

                if (innerContainer && innerContainer.style) {
                    innerContainer.style.display = 'block'
                }

                folderWrapperElement.querySelectorAll('.tree-file-item, .file-item').forEach(function (fileItem) {
                    if (fileItem && fileItem.style) {
                        fileItem.style.display = 'block'
                    }
                })
            })

            return
        }

        // Iterate through each folder wrapper element to evaluate search matching conditions.
        folderWrapperElements.forEach(function (folderWrapperElement) {
            // Guard against missing folder wrapper elements.
            if (!folderWrapperElement) {
                return
            }

            // Extract the header text from the folder element and clean it for precise comparison.
            const folderHeaderElement = folderWrapperElement.querySelector('.folder-header, .folder-title, span, strong, div:first-child')
            const folderHeaderText = folderHeaderElement ? String(folderHeaderElement.textContent || '').toLowerCase().trim() : ''
            const cleanedHeader = folderHeaderText.replace(/[-_\s/]/g, '')

            // Query all individual file item elements contained within the current folder wrapper.
            const fileItemElements = folderWrapperElement.querySelectorAll('.tree-file-item, .file-item')

            let hasMatchingFile = false

            // Iterate over each file item to check if its text matches the search query.
            fileItemElements.forEach(function (fileItem) {
                const fileText = String(fileItem.textContent || '').toLowerCase()
                const cleanedFileText = fileText.replace(/[-_\s/]/g, '')

                if (cleanedFileText.includes(cleanedQuery)) {
                    hasMatchingFile = true
                }
            })

            // Determine if the folder header starts with the cleaned search query.
            const headerMatches = cleanedHeader.startsWith(cleanedQuery)
            
            // Final match flag evaluated from either header prefix match or matching child files.
            const isMatch = headerMatches || hasMatchingFile

            // Update display styles based on whether the folder matches the search criteria.
            if (isMatch) {
                folderWrapperElement.style.display = 'block'

                const innerContainer = folderWrapperElement.querySelector('.ml-3')

                if (innerContainer && innerContainer.style) {
                    innerContainer.style.display = 'block'
                }

                fileItemElements.forEach(function (fileItem) {
                    if (fileItem && fileItem.style) {
                        fileItem.style.display = 'block'
                    }
                })
            } else {
                folderWrapperElement.style.display = 'none'
            }
        })
    }

    /**
     * Safely fetches repository file content by targeting the exact code directory layout.
     * @param {string} rawFilePathString - The raw file path extracted from data-path or href.
     */
    async function loadRepositoryFileContent(rawFilePathString) {
        // Guard against missing file path strings to prevent runtime errors.
        if (!rawFilePathString) {
            console.error("Error: Cannot load file because the file path is empty or null.")
            return
        }

        let normalizedFilePathString = rawFilePathString.trim()

        // Ignore anchor hash links.
        if (normalizedFilePathString.startsWith('#')) {
            return
        }

        // Clean leading relative path markers and slashes.
        const cleanPathString = normalizedFilePathString
            .replace(/^(\.\/|\/+)+/, '')

        const candidateUrlList = []

        // If the path explicitly points to tests/, include/, or src/, try it directly and with relative fallback prefixes.
        if (cleanPathString.startsWith('tests/') || cleanPathString.startsWith('include/') || cleanPathString.startsWith('src/')) {
            candidateUrlList.push(cleanPathString)
            candidateUrlList.push(`/${cleanPathString}`)
            candidateUrlList.push(`../${cleanPathString}`)
            candidateUrlList.push(`../../${cleanPathString}`)
        } else {
            // Fallback candidate URLs for legacy darwin-systems path formatting.
            candidateUrlList.push(`code/darwin-systems/${cleanPathString}`)
            candidateUrlList.push(`/code/darwin-systems/${cleanPathString}`)
            candidateUrlList.push(cleanPathString)
            candidateUrlList.push(`/${cleanPathString}`)
        }

        // Deduplicate candidate URLs to prevent redundant network requests.
        const uniqueCandidateUrlList = [...new Set(candidateUrlList)]

        let successfulFileContentText = null
        let resolvedFilePath = normalizedFilePathString

        // Iterate through candidate URLs until a successful HTTP response is received.
        for (const candidateUrl of uniqueCandidateUrlList) {
            try {
                console.log(`Attempting to fetch file from candidate URL: ${candidateUrl}`)

                const response = await fetch(candidateUrl)

                if (response.ok) {
                    successfulFileContentText = await response.text()
                    console.log(`Successfully fetched file using URL: ${candidateUrl}`)
                    break
                }
            } catch (networkError) {
                // Continue silently to the next candidate URL upon network failure.
            }
        }

        // Render the fetched file content or log an error if all candidate URLs fail.
        if (successfulFileContentText !== null) {
            displayFileContentInViewer(resolvedFilePath, successfulFileContentText)
        } else {
            console.error("Error: Could not load file from path across all candidate URLs ->", normalizedFilePathString)
        }
    }

    /**
     * Updates the UI code viewer container with the fetched file content and file path.
     * @param {string} filePathString - The normalized file path.
     * @param {string} fileContentText - The raw text content of the file.
     */
    function displayFileContentInViewer(filePathString, fileContentText) {
        // Search for active code viewer containers across the document using supported portfolio container IDs and selectors.
        const codeViewerContainerElement = document.getElementById('active-code-container') ||
            document.getElementById('python-active-code-container') ||
            document.getElementById('golang-active-code-container') ||
            document.getElementById('under-the-hood-active-code-container') ||
            document.getElementById('code-viewer') ||
            document.getElementById('file-content') ||
            document.querySelector('pre code') ||
            document.querySelector('.code-container')

        if (codeViewerContainerElement) {
            codeViewerContainerElement.textContent = fileContentText
            codeViewerContainerElement.removeAttribute("data-highlighted")

            // Determine syntax highlighting language class based on file extension.
            let languageClass = "language-plaintext"

            const lowerCasePath = filePathString.toLowerCase()
            
            if (lowerCasePath.endsWith(".c") || lowerCasePath.endsWith(".h")) {
                languageClass = "language-c"
            } else if (lowerCasePath.endsWith(".cpp") || lowerCasePath.endsWith(".hpp")) {
                languageClass = "language-cpp"
            } else if (lowerCasePath.endsWith(".py")) {
                languageClass = "language-python"
            } else if (lowerCasePath.endsWith(".go")) {
                languageClass = "language-go"
            }

            codeViewerContainerElement.className = languageClass

            // Apply syntax highlighting if highlight.js is loaded in the window environment.
            if (window.hljs && typeof window.hljs.highlightElement === "function") {
                window.hljs.highlightElement(codeViewerContainerElement)
            }

            console.log("Successfully rendered file content into viewer for:", filePathString)
        } else {
            console.log("File content loaded successfully (no standard code viewer container found, logging content):", filePathString)
        }

        // Dispatch a custom window event in case external scripts listen for repository file loading events.
        const fileLoadedCustomEvent = new CustomEvent('repositoryFileLoaded', {
            detail: {
                filePath: filePathString,
                content: fileContentText
            }
        })

        window.dispatchEvent(fileLoadedCustomEvent)
    }

    /**
     * Sets up live input filtering and clear button functionality for a specific repository search container.
     * @param {string} searchInputIdString - The DOM element identifier for the search input field.
     * @param {HTMLElement} treeContainerElement - The container element holding the repository file tree.
     */
    window.setupRepositorySearch = function (searchInputIdString, treeContainerElement) {
        const searchInputElement = document.getElementById(searchInputIdString)

        // Guard against missing search input or tree container elements to prevent runtime errors.
        if (!searchInputElement || !treeContainerElement) {
            return
        }

        const controlWrapperElement = searchInputElement.closest('.control')
        const clearButtonElement = controlWrapperElement ? controlWrapperElement.querySelector('.clear-search-btn') : null

        // Prevent attaching duplicate event listeners if the initialization function is invoked multiple times.
        if (searchInputElement.dataset.searchInitialized === 'true') {
            return
        }

        searchInputElement.dataset.searchInitialized = 'true'

        // Define the handler function for live input changes.
        const handleInputChangeEvent = function () {
            const queryValueString = searchInputElement.value

            if (clearButtonElement && clearButtonElement.style) {
                clearButtonElement.style.display = queryValueString.trim().length > 0 ? 'flex' : 'none'
            }

            filterRepositoryTree(treeContainerElement, queryValueString)
        }

        // Attach the input event listener to trigger real-time filtering as the user types.
        searchInputElement.addEventListener('input', handleInputChangeEvent)

        // Attach a click event listener to the clear button to reset the search input field and restore the tree view.
        if (clearButtonElement) {
            clearButtonElement.addEventListener('click', function () {
                searchInputElement.value = ''

                if (clearButtonElement.style) {
                    clearButtonElement.style.display = 'none'
                }

                filterRepositoryTree(treeContainerElement, '')

                searchInputElement.focus()
            })
        }

        // Execute the handler function immediately if the search input already contains text upon initialization.
        if (searchInputElement.value.trim().length > 0) {
            handleInputChangeEvent()
        }
    }

    // Attach a global input listener as a bulletproof fallback mechanism for dynamically rendered search inputs.
    document.addEventListener('input', function (event) {
        const targetInputElement = event.target

        if (!targetInputElement || !searchInputMappingDictionary[targetInputElement.id]) {
            return
        }

        const containerIdString = searchInputMappingDictionary[targetInputElement.id]
        const containerElement = document.getElementById(containerIdString)

        if (containerElement) {
            filterRepositoryTree(containerElement, targetInputElement.value)
        }
    })

    // Attach a global click event listener to intercept file item clicks across repository tree containers.
    document.addEventListener('click', function (event) {
        const clickedElement = event.target

        // Find if the click originated from a file item or a link within a file tree container.
        const fileItemElement = clickedElement.closest('.tree-file-item, .file-item')
        const linkElement = clickedElement.tagName === 'A' ? clickedElement : clickedElement.closest('a')

        if (!fileItemElement && !linkElement) {
            return
        }

        // Ensure the click is within one of our designated repository tree containers.
        const parentContainerElement = (fileItemElement || linkElement).closest('#repository-tree-container, #python-repository-tree-container, #golang-repository-tree-container, #under-the-hood-repository-tree-container')

        if (!parentContainerElement) {
            return
        }

        // Extract the file path from the data-path attribute or href attribute.
        const targetLinkElement = linkElement || (fileItemElement ? fileItemElement.querySelector('a') : null)
        const rawFilePathString = targetLinkElement ? (targetLinkElement.getAttribute('data-path') || targetLinkElement.getAttribute('href')) : null

        if (rawFilePathString) {
            // Prevent default browser navigation to avoid 404 page routing errors.
            event.preventDefault()
            loadRepositoryFileContent(rawFilePathString)
        }
    })
})()