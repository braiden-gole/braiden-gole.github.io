// Configure manual scroll restoration to force the browser to the top on page refresh.
if ('scrollRestoration' in history) {
    history.scrollRestoration = 'manual'
}

/**
 * Main portfolio initialization function called explicitly after components load.
 */
window.initializePortfolio = async function() {
    // Force scroll position to the top ONLY on a clean load/refresh with no hash and no target parameter.
    const urlParams = new URLSearchParams(window.location.search)
    
    if (!window.location.hash && !urlParams.has('target')) {
        window.scrollTo(0, 0)
    }

    document.querySelectorAll(".notification .delete" || []).forEach((erase) => {
        const notification = erase.parentNode

        erase.addEventListener("click", () => {
            notification.parentNode.removeChild(notification)
        })
    })

    // Define configuration objects for each managed repository source.
    const repositoryConfigurations = [
        {
            manifestPath: "code/c_manifest.json",
            baseDirectory: "code/c",
            treeContainerId: "repository-tree-container",
            codeContainerId: "active-code-container",
            fileLabelId: "active-file-label",
            viewerWrapperId: "c-code-viewer",
            itemClass: "c-tree-file-item",
            folderThemeClass: "has-text-success",
            iconClass: "fas fa-file-code has-text-success"
        },
        {
            manifestPath: "code/python_manifest.json",
            baseDirectory: "code/python",
            treeContainerId: "python-repository-tree-container",
            codeContainerId: "python-active-code-container",
            fileLabelId: "python-active-file-label",
            viewerWrapperId: "python-code-viewer",
            itemClass: "python-tree-file-item",
            folderThemeClass: "has-text-warning",
            iconClass: "fab fa-python has-text-warning"
        },
        {
            manifestPath: "code/go_manifest.json",
            baseDirectory: "code/golang",
            treeContainerId: "golang-repository-tree-container",
            codeContainerId: "golang-active-code-container",
            fileLabelId: "golang-active-file-label",
            viewerWrapperId: "golang-code-viewer",
            itemClass: "golang-tree-file-item",
            folderThemeClass: "has-text-info",
            iconClass: "fa-brands fa-golang has-text-info"
        },
        {
            manifestPath: "code/under_the_hood_manifest.json",
            baseDirectory: "code/under-the-hood-ml",
            treeContainerId: "under-the-hood-repository-tree-container",
            codeContainerId: "under-the-hood-active-code-container",
            fileLabelId: "under-the-hood-active-file-label",
            viewerWrapperId: "under-the-hood-code-viewer",
            itemClass: "under-the-hood-tree-file-item",
            folderThemeClass: "has-text-primary",
            iconClass: "fas fa-gear has-text-primary"
        }
    ]

    // Loop through each configuration and initialize the respective repository viewer.
    for (const repositoryConfiguration of repositoryConfigurations) {
        await initializeRepository(repositoryConfiguration)
    }
}

/**
 * Determines the appropriate FontAwesome icon class based on the file extension.
 * @param {string} fileName - The name of the file.
 * @param {string} defaultIcon - The fallback icon class.
 * @returns {string} The CSS class string for the icon.
 */
function getFileIconForExtension(fileName, defaultIcon) {
    if (!fileName) {
        return defaultIcon
    }

    const lower = fileName.toLowerCase()

    if (lower.endsWith(".py") || lower.endsWith(".pyx") || lower.endsWith(".pyi") || lower.endsWith(".pxd")) {
        return "fab fa-python has-text-warning"
    }
    if (lower.endsWith(".go")) {
        return "fa-brands fa-golang has-text-info"
    }
    if (lower.endsWith(".cs")) {
        return "fas fa-code has-text-primary"
    }
    if (lower.endsWith(".java")) {
        return "fab fa-java has-text-danger"
    }
    if (lower.endsWith(".cpp") || lower.endsWith(".hpp") || lower.endsWith(".cc") || lower.endsWith(".c") || lower.endsWith(".h")) {
        return "fas fa-file-code has-text-success"
    }

    return defaultIcon
}

/**
 * Fetches a repository manifest, filters out ignored files/folders, and builds the UI file tree.
 * @param {Object} repositoryConfiguration - Configuration parameters for the specific repository.
 */
async function initializeRepository(repositoryConfiguration) {
    const treeContainerElement = document.getElementById(repositoryConfiguration.treeContainerId)
    const activeCodeContainer = document.getElementById(repositoryConfiguration.codeContainerId)
    const activeFileLabelElement = document.getElementById(repositoryConfiguration.fileLabelId)

    if (activeCodeContainer) {
        activeCodeContainer.textContent = ""
        activeCodeContainer.removeAttribute("data-highlighted")

        const parentWrapperElement = activeCodeContainer.closest(".code-viewer-wrapper")

        if (parentWrapperElement) {
            parentWrapperElement.style.display = "none"
        }
    }

    if (activeFileLabelElement) {
        activeFileLabelElement.innerHTML = `<i class="fas fa-file-code mr-2"></i> No file selected`
    }

    if (!treeContainerElement) {
        return
    }

    try {
        const manifestResponse = await fetch(repositoryConfiguration.manifestPath)

        if (!manifestResponse.ok) {
            throw new Error(`HTTP status ${manifestResponse.status}`)
        }

        const rawRepositoryFileEntries = await manifestResponse.json()

        if (!Array.isArray(rawRepositoryFileEntries) || rawRepositoryFileEntries.length === 0) {
            treeContainerElement.innerHTML = `<span class="has-text-grey">// No repository files found.</span>`
            return
        }

        let ignoredExtensions = [
            ".dll", ".exe", ".pdb", ".csproj", ".sln", ".slnx", ".user", ".suo",
            ".cache", ".json", ".o", ".obj", ".pyc", ".jar", "go.mod", "go.sum"
        ]

        if (repositoryConfiguration.baseDirectory === "code/under-the-hood-ml") {
            ignoredExtensions.push(".so")
        }

        const ignoredFolderPatterns = ["/obj/", "/bin/", ".vs/", "/lib/", ".git/"]

        const filteredRepositoryFileEntries = rawRepositoryFileEntries.filter(fileEntry => {
            const fileName = (fileEntry.fileName || "").toLowerCase()
            const filePath = (fileEntry.filePath || "").toLowerCase()
            const folderPath = (fileEntry.folderPath || "").toLowerCase()

            const hasIgnoredExtension = ignoredExtensions.some(ext => fileName.endsWith(ext))
            const isInIgnoredFolder = ignoredFolderPatterns.some(pat => filePath.includes(pat) || folderPath.includes(pat))

            let isIgnoredSpecificFile = false

            if (repositoryConfiguration.baseDirectory === "code/under-the-hood-ml") {
                if (fileName === "unity.h" || fileName === "unity.c" || fileName === "unity_internals.h") {
                    isIgnoredSpecificFile = true
                }
            }

            return !hasIgnoredExtension && !isInIgnoredFolder && !isIgnoredSpecificFile
        })

        const groupedDirectoryMap = {}

        filteredRepositoryFileEntries.forEach(fileEntry => {
            const folderName = fileEntry.folderPath || "Root"

            if (!groupedDirectoryMap[folderName]) {
                groupedDirectoryMap[folderName] = []
            }

            groupedDirectoryMap[folderName].push(fileEntry)
        })

        treeContainerElement.innerHTML = ""

        Object.keys(groupedDirectoryMap).forEach(folderPath => {
            const folderWrapperElement = document.createElement("div")
            folderWrapperElement.className = "mb-1"

            const folderHeaderElement = document.createElement("div")
            folderHeaderElement.className = `${repositoryConfiguration.folderThemeClass} is-flex is-align-items-center py-1 px-2`
            folderHeaderElement.innerHTML = `<i class="fas fa-folder mr-2"></i><span>${folderPath}</span>`
            folderWrapperElement.appendChild(folderHeaderElement)

            const filesListContainer = document.createElement("div")
            filesListContainer.className = "ml-3"
            filesListContainer.style.borderLeft = "1px dashed var(--border-color)"

            groupedDirectoryMap[folderPath].forEach(fileEntry => {
                const fileItemButton = document.createElement("a")
                const fileIconClass = getFileIconForExtension(fileEntry.fileName, repositoryConfiguration.iconClass)

                fileItemButton.className = `is-block py-1 px-2 tree-file-item ${repositoryConfiguration.itemClass}`
                fileItemButton.style.borderRadius = "0"
                fileItemButton.style.transition = "background 0.2s ease, color 0.2s ease"
                fileItemButton.style.color = "var(--text-body)"
                fileItemButton.innerHTML = `<i class="${fileIconClass} mr-2"></i>${fileEntry.fileName} <span class="is-size-8 has-text-grey">(${fileEntry.fileCategory || "Source"})</span>`

                fileItemButton.addEventListener("click", async () => {
                    document.querySelectorAll(`.${repositoryConfiguration.itemClass}`).forEach(element => {
                        element.classList.remove("is-active-file")
                        element.style.background = "transparent"
                        element.style.color = "var(--text-body)"
                    })

                    fileItemButton.classList.add("is-active-file")
                    fileItemButton.style.background = "rgba(16, 124, 16, 0.2)"
                    fileItemButton.style.color = "var(--text-heading)"

                    if (activeFileLabelElement) {
                        activeFileLabelElement.innerHTML = `<i class="${fileIconClass} mr-2"></i> ${fileEntry.fileName}`
                    }

                    if (activeCodeContainer) {
                        await fetchAndDisplaySourceFile(fileEntry, activeCodeContainer, repositoryConfiguration.baseDirectory)
                    }

                    const viewerWrapper = document.getElementById(repositoryConfiguration.viewerWrapperId)

                    if (viewerWrapper) {
                        viewerWrapper.scrollIntoView({ behavior: 'smooth', block: 'start' })
                    }
                })

                filesListContainer.appendChild(fileItemButton)
            })

            folderWrapperElement.appendChild(filesListContainer)
            treeContainerElement.appendChild(folderWrapperElement)
        })
    } catch (error) {
        treeContainerElement.innerHTML = `<span class="has-text-danger">// Error loading manifest: ${error.message}</span>`
    }
}

/**
 * Fetches the raw text content of a source file, determines its language class, and triggers syntax highlighting.
 */
async function fetchAndDisplaySourceFile(fileEntry, codeContainerElement, baseDirectory) {
    const parentWrapperElement = codeContainerElement.closest(".code-viewer-wrapper")

    if (parentWrapperElement) {
        parentWrapperElement.style.display = "block"
    }

    codeContainerElement.textContent = `// Loading ${fileEntry.fileName}...`
    let languageClass = "language-plaintext"

    if (fileEntry.fileName) {
        const lowercaseFileName = fileEntry.fileName.toLowerCase()

        if (lowercaseFileName.endsWith(".cpp") || lowercaseFileName.endsWith(".hpp") || lowercaseFileName.endsWith(".cc") || lowercaseFileName.endsWith(".cxx")) {
            languageClass = "language-cpp"
        } else if (lowercaseFileName.endsWith(".c") || lowercaseFileName.endsWith(".h")) {
            languageClass = "language-c"
        } else if (lowercaseFileName.endsWith(".cs")) {
            languageClass = "language-csharp"
        } else if (lowercaseFileName.endsWith(".py") || lowercaseFileName.endsWith(".pyx") || lowercaseFileName.endsWith(".pyi") || lowercaseFileName.endsWith(".pxd")) {
            languageClass = "language-python"
        } else if (lowercaseFileName.endsWith(".java")) {
            languageClass = "language-java"
        } else if (lowercaseFileName.endsWith(".go")) {
            languageClass = "language-go"
        }
    }

    try {
        if (!fileEntry.filePath) {
            throw new Error("Target file path is undefined or null.")
        }

        let resolvedFilePath = fileEntry.filePath.trim()

        if (resolvedFilePath.startsWith("/")) {
            resolvedFilePath = resolvedFilePath.slice(1)
        }

        if (resolvedFilePath.startsWith("code/darwin-systems/") || resolvedFilePath.startsWith("code/python-mac-os/")) {
        } else {
            const cleanBaseDirectory = baseDirectory ? (baseDirectory.endsWith("/") ? baseDirectory.slice(0, -1) : baseDirectory) : ""

            if (cleanBaseDirectory) {
                if (!resolvedFilePath.startsWith(cleanBaseDirectory + "/")) {
                    const baseName = cleanBaseDirectory.split("/").pop()

                    if (resolvedFilePath.startsWith("code/") && !resolvedFilePath.startsWith(`code/${baseName}/`)) {
                        resolvedFilePath = resolvedFilePath.replace(/^code\//, `${cleanBaseDirectory}/`)
                    } else {
                        resolvedFilePath = `${cleanBaseDirectory}/${resolvedFilePath}`
                    }
                }
            }
        }

        const encodedFilePath = encodeURI(resolvedFilePath)
        const fileResponse = await fetch(encodedFilePath)

        if (!fileResponse.ok) {
            throw new Error(`HTTP ${fileResponse.status} (${fileResponse.statusText}) attempting to fetch: ${resolvedFilePath}`)
        }

        const fileContentText = await fileResponse.text()

        codeContainerElement.textContent = fileContentText
        codeContainerElement.removeAttribute("data-highlighted")
        codeContainerElement.className = languageClass

        if (window.hljs && typeof window.hljs.highlightElement === "function") {
            window.hljs.highlightElement(codeContainerElement)
        }
    } catch (error) {
        console.error("Error fetching file content:", error)
        codeContainerElement.textContent = `// ERROR: Could not load file from path -> ${fileEntry.filePath}\n// Details: ${error.message}`
    }
}

/**
 * Resets and clears the C repository viewer state.
 */
function collapseAndClear() { 
    resetRepositoryViewerState("active-code-container", "active-file-label", "c-tree-file-item") 
}

/**
 * Resets and clears the Python repository viewer state.
 */
function collapseAndClearPython() { 
    resetRepositoryViewerState("python-active-code-container", "python-active-file-label", "python-tree-file-item") 
}

/**
 * Resets and clears the Golang repository viewer state.
 */
function collapseAndClearGolang() { 
    resetRepositoryViewerState("golang-active-code-container", "golang-active-file-label", "golang-tree-file-item") 
}

/**
 * Resets and clears the Under the Hood ML repository viewer state.
 */
function collapseAndClearUnderTheHood() { 
    resetRepositoryViewerState("under-the-hood-active-code-container", "under-the-hood-active-file-label", "under-the-hood-tree-file-item") 
}

/**
 * Resets the specified repository viewer state and clears its active selection.
 */
function resetRepositoryViewerState(codeContainerId, fileLabelId, itemClass) {
    const activeCodeContainer = document.getElementById(codeContainerId)
    const activeFileLabelElement = document.getElementById(fileLabelId)

    if (activeCodeContainer) {
        activeCodeContainer.textContent = ""
        activeCodeContainer.removeAttribute("data-highlighted")

        const parentWrapperElement = activeCodeContainer.closest(".code-viewer-wrapper")

        if (parentWrapperElement) {
            parentWrapperElement.style.display = "none"
        }
    }

    if (activeFileLabelElement) {
        activeFileLabelElement.innerHTML = `<i class="fas fa-file-code mr-2"></i> No file selected`
    }

    document.querySelectorAll(`.${itemClass}`).forEach(element => {
        element.classList.remove("is-active-file")
        element.style.background = "transparent"
        element.style.color = "var(--text-body)"
    })
}

// Expose functions globally for HTML inline handlers & module loaders.
window.collapseAndClear = collapseAndClear
window.collapseAndClearPython = collapseAndClearPython
window.collapseAndClearGolang = collapseAndClearGolang
window.collapseAndClearUnderTheHood = collapseAndClearUnderTheHood