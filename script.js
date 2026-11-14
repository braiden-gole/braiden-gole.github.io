// Configure manual scroll restoration to prevent browser scroll retention on page refresh.
if ('scrollRestoration' in history) 
{
    history.scrollRestoration = 'manual'
}

/**
 * Main portfolio initialization function called explicitly after components load.
 */
window.initializePortfolio = async function() 
{
    // Force scroll position to the top only on a clean load or refresh with no hash and no target parameter.
    const currentUrlQueryParameters = new URLSearchParams(window.location.search)
    
    if (!window.location.hash && !currentUrlQueryParameters.has('target')) 
    {
        window.scrollTo(0, 0)
    }

    // Set up click handlers to dismiss any active notification cards present on the page.
    const notificationDeleteButtonElements = document.querySelectorAll(".notification .delete")
    
    if (notificationDeleteButtonElements) 
    {
        notificationDeleteButtonElements.forEach(function(deleteButtonElement) 
        {
            // Retrieve the parent notification card container element safely.
            const parentNotificationCardElement = deleteButtonElement.parentNode

            if (deleteButtonElement && parentNotificationCardElement) 
            {
                deleteButtonElement.addEventListener("click", function() 
                {
                    if (parentNotificationCardElement && parentNotificationCardElement.parentNode) 
                    {
                        parentNotificationCardElement.parentNode.removeChild(parentNotificationCardElement)
                    }
                })
            }
        })
    }

    // Define configuration objects for each managed repository source and destination viewer.
    const managedRepositoryConfigurationList = [
        {
            manifestPath: "code/c_manifest.json",
            baseDirectory: "code/darwin-systems",
            treeContainerId: "repository-tree-container",
            codeContainerId: "active-code-container",
            fileLabelId: "active-file-label",
            viewerWrapperId: "c-code-viewer",
            itemClass: "c-tree-file-item",
            folderThemeClass: "has-text-success",
            iconClass: "fas fa-file-code has-text-success",
            searchId: "c-tree-search"
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
            iconClass: "fab fa-python has-text-warning",
            searchId: "python-tree-search"
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
            iconClass: "fa-brands fa-golang has-text-info",
            searchId: "golang-tree-search"
        },
        {
            manifestPath: "code/under_the_hood_manifest.json",
            baseDirectory: "code/under-the-hood-ml",
            treeContainerId: "under-the-hood-repository-tree-container",
            codeContainerId: "under-the-hood-active-code-container",
            fileLabelId: "under-the-hood-active-file-label",
            viewerWrapperId: "under-the-hood-code-viewer",
            itemClass: "under-the-hood-tree-file-item",
            folderThemeClass: "has-text-purple",
            iconClass: "fas fa-gear has-text-purple",
            searchId: "under-the-hood-tree-search"
        }
    ]

    // Loop through each configuration and asynchronously initialize the respective repository viewer.
    for (let repositoryConfigurationIterationIndex = 0; repositoryConfigurationIterationIndex < managedRepositoryConfigurationList.length; repositoryConfigurationIterationIndex++) 
    {
        const currentRepositoryConfiguration = managedRepositoryConfigurationList[repositoryConfigurationIterationIndex]
        
        if (currentRepositoryConfiguration) 
        {
            await initializeRepository(currentRepositoryConfiguration)
        }
    }
}

/**
 * Determines the appropriate FontAwesome icon class based on the file extension.
 * @param {string} sourceFileName - The name of the file.
 * @param {string} fallbackDefaultIconClass - The fallback icon class.
 * @returns {string} The CSS class string for the icon.
 */
function getFileIconForExtension(sourceFileName, fallbackDefaultIconClass) 
{
    // Return default icon immediately if no file name is provided.
    if (!sourceFileName) 
    {
        return fallbackDefaultIconClass
    }

    const lowercaseNormalizedFileName = sourceFileName.toLowerCase()

    // Match extension strings to select correct language brand or file icon style.
    if (lowercaseNormalizedFileName.endsWith(".py") || lowercaseNormalizedFileName.endsWith(".pyx") || lowercaseNormalizedFileName.endsWith(".pyi") || lowercaseNormalizedFileName.endsWith(".pxd")) 
    {
        return "fab fa-python has-text-warning"
    }

    if (lowercaseNormalizedFileName.endsWith(".go")) 
    {
        return "fa-brands fa-golang has-text-info"
    }

    if (lowercaseNormalizedFileName.endsWith(".cs")) 
    {
        return "fas fa-code has-text-orange"
    }

    if (lowercaseNormalizedFileName.endsWith(".java")) 
    {
        return "fab fa-java has-text-danger"
    }

    if (lowercaseNormalizedFileName.endsWith(".cpp") || lowercaseNormalizedFileName.endsWith(".hpp") || lowercaseNormalizedFileName.endsWith(".cc") || lowercaseNormalizedFileName.endsWith(".c") || lowercaseNormalizedFileName.endsWith(".h")) 
    {
        return "fas fa-file-code has-text-success"
    }

    return fallbackDefaultIconClass
}

/**
 * Fetches a repository manifest, filters out ignored files and folders, and builds the UI file tree.
 * @param {Object} currentRepositoryConfiguration - Configuration parameters for the specific repository.
 */
async function initializeRepository(currentRepositoryConfiguration) 
{
    // Defensively guard against missing repository configuration objects.
    if (!currentRepositoryConfiguration) 
    {
        return
    }

    const repositoryTreeContainerElement = document.getElementById(currentRepositoryConfiguration.treeContainerId)
    const activeSourceCodeContainerElement = document.getElementById(currentRepositoryConfiguration.codeContainerId)
    const activeFileDescriptionLabelElement = document.getElementById(currentRepositoryConfiguration.fileLabelId)

    // Clear and hide the code viewing container when switching or initializing repositories.
    if (activeSourceCodeContainerElement) 
    {
        activeSourceCodeContainerElement.textContent = ""
        activeSourceCodeContainerElement.removeAttribute("data-highlighted")

        const codeViewerParentWrapperElement = activeSourceCodeContainerElement.closest(".code-viewer-wrapper")

        if (codeViewerParentWrapperElement) 
        {
            codeViewerParentWrapperElement.style.display = "none"
        }
    }

    // Reset the active file description label to its default placeholder.
    if (activeFileDescriptionLabelElement) 
    {
        activeFileDescriptionLabelElement.innerHTML = `<i class="fas fa-file-code mr-2"></i> No file selected`
    }

    // Exit early if the corresponding tree container element does not exist in the DOM.
    if (!repositoryTreeContainerElement) 
    {
        return
    }

    try 
    {
        // Fetch the JSON manifest describing the repository files.
        const repositoryManifestHttpResponse = await fetch(currentRepositoryConfiguration.manifestPath)

        if (!repositoryManifestHttpResponse.ok) 
        {
            throw new Error(`HTTP status ${repositoryManifestHttpResponse.status}`)
        }

        const rawRepositoryManifestFileEntries = await repositoryManifestHttpResponse.json()

        // Handle empty or invalid manifest entries gracefully.
        if (!Array.isArray(rawRepositoryManifestFileEntries) || rawRepositoryManifestFileEntries.length === 0) 
        {
            repositoryTreeContainerElement.innerHTML = `<span class="has-text-grey">// No repository files found.</span>`
            return
        }

        // Define file extensions that should be hidden from the portfolio tree.
        let blacklistedFileExtensionList = [
            ".dll", ".exe", ".pdb", ".csproj", ".sln", ".slnx", ".user", ".suo",
            ".cache", ".json", ".o", ".obj", ".pyc", ".jar", "go.mod", "go.sum"
        ]

        // Add specific excludes for machine learning sub-projects if applicable.
        if (currentRepositoryConfiguration.baseDirectory === "code/under-the-hood-ml") 
        {
            blacklistedFileExtensionList.push(".so")
        }

        // Define directory segments to ignore globally across repositories.
        const blacklistedFolderPatternList = ["/obj/", "/bin/", ".vs/", "/lib/", ".git/"]

        // Filter out binary builds, ignored folders, and specific framework headers.
        const filteredRepositoryFileEntryList = rawRepositoryManifestFileEntries.filter(repositoryFileEntry => 
        {
            if (!repositoryFileEntry) 
            {
                return false
            }

            const sourceFileName = (repositoryFileEntry.fileName || "").toLowerCase()
            const repositoryFileRelativePath = (repositoryFileEntry.filePath || "").toLowerCase()
            const repositoryFolderDirectoryPath = (repositoryFileEntry.folderPath || "").toLowerCase()

            const matchesBlacklistedExtension = blacklistedFileExtensionList.some(extensionString => sourceFileName.endsWith(extensionString))
            const matchesBlacklistedFolderPattern = blacklistedFolderPatternList.some(patternString => repositoryFileRelativePath.includes(patternString) || repositoryFolderDirectoryPath.includes(patternString))

            let isExplicitlyBlacklistedFile = false

            if (currentRepositoryConfiguration.baseDirectory === "code/under-the-hood-ml") 
            {
                if (sourceFileName === "unity.h" || sourceFileName === "unity.c" || sourceFileName === "unity_internals.h") 
                {
                    isExplicitlyBlacklistedFile = true
                }
            }

            return !matchesBlacklistedExtension && !matchesBlacklistedFolderPattern && !isExplicitlyBlacklistedFile
        })

        // Group filtered file entries by their destination folder path.
        const groupedDirectoryPathMapping = {}

        filteredRepositoryFileEntryList.forEach(repositoryFileEntry => 
        {
            if (repositoryFileEntry) 
            {
                const targetFolderName = repositoryFileEntry.folderPath || "Root"

                if (!groupedDirectoryPathMapping[targetFolderName]) 
                {
                    groupedDirectoryPathMapping[targetFolderName] = []
                }

                groupedDirectoryPathMapping[targetFolderName].push(repositoryFileEntry)
            }
        })

        // Clear existing tree container contents before rebuilding.
        repositoryTreeContainerElement.innerHTML = ""

        // Iterate through each grouped directory to construct DOM elements.
        Object.keys(groupedDirectoryPathMapping).forEach(directoryPathKey => 
        {
            const directoryFolderWrapperElement = document.createElement("div")

            // Include both classes so selectors in filter.js match reliably.
            directoryFolderWrapperElement.className = "mb-1 folder-item"

            const directoryFolderHeaderElement = document.createElement("div")

            directoryFolderHeaderElement.className = `${currentRepositoryConfiguration.folderThemeClass} is-flex is-align-items-center py-1 px-2`
            directoryFolderHeaderElement.innerHTML = `<i class="fas fa-folder mr-2"></i><span>${directoryPathKey}</span>`

            directoryFolderWrapperElement.appendChild(directoryFolderHeaderElement)

            const directoryFilesContainerElement = document.createElement("div")

            directoryFilesContainerElement.className = "ml-3"
            directoryFilesContainerElement.style.borderLeft = "1px dashed var(--border-color)"

            // Build individual selectable file buttons within each folder.
            groupedDirectoryPathMapping[directoryPathKey].forEach(repositoryFileEntry => 
            {
                if (repositoryFileEntry) 
                {
                    const selectableFileItemAnchorElement = document.createElement("a")
                    const resolvedFileIconCssClass = getFileIconForExtension(repositoryFileEntry.fileName, currentRepositoryConfiguration.iconClass)

                    selectableFileItemAnchorElement.className = `is-block py-1 px-2 tree-file-item file-item ${currentRepositoryConfiguration.itemClass}`
                    selectableFileItemAnchorElement.style.borderRadius = "0"
                    selectableFileItemAnchorElement.style.transition = "background 0.2s ease, color 0.2s ease"
                    selectableFileItemAnchorElement.style.color = "var(--text-body)"
                    selectableFileItemAnchorElement.innerHTML = `<i class="${resolvedFileIconCssClass} mr-2"></i>${repositoryFileEntry.fileName} <span class="is-size-8 has-text-grey">(${repositoryFileEntry.fileCategory || "Source"})</span>`

                    // Attach click handler to load and display the source code content.
                    selectableFileItemAnchorElement.addEventListener("click", async () => 
                    {
                        document.querySelectorAll(`.${currentRepositoryConfiguration.itemClass}`).forEach(treeItemElement => 
                        {
                            if (treeItemElement && treeItemElement.style) 
                            {
                                treeItemElement.classList.remove("is-active-file")
                                treeItemElement.style.background = "transparent"
                                treeItemElement.style.color = "var(--text-body)"
                            }
                        })

                        selectableFileItemAnchorElement.classList.add("is-active-file")
                        selectableFileItemAnchorElement.style.background = "rgba(16, 124, 16, 0.2)"
                        selectableFileItemAnchorElement.style.color = "var(--text-heading)"

                        if (activeFileDescriptionLabelElement) 
                        {
                            activeFileDescriptionLabelElement.innerHTML = `<i class="${resolvedFileIconCssClass} mr-2"></i> ${repositoryFileEntry.fileName}`
                        }

                        if (activeSourceCodeContainerElement) 
                        {
                            await fetchAndDisplaySourceFile(repositoryFileEntry, activeSourceCodeContainerElement, currentRepositoryConfiguration.baseDirectory)
                        }

                        const targetCodeViewerWrapperElement = document.getElementById(currentRepositoryConfiguration.viewerWrapperId)

                        // Smoothly scroll down to the code viewer container when a file is selected.
                        if (targetCodeViewerWrapperElement) 
                        {
                            targetCodeViewerWrapperElement.scrollIntoView({ behavior: 'smooth', block: 'start' })
                        }
                    })

                    directoryFilesContainerElement.appendChild(selectableFileItemAnchorElement)
                }
            })

            directoryFolderWrapperElement.appendChild(directoryFilesContainerElement)
            repositoryTreeContainerElement.appendChild(directoryFolderWrapperElement)
        })

        // Attach filter and search capabilities to the input elements after rendering the tree.
        if (typeof setupRepositorySearch === 'function') 
        {
            setupRepositorySearch(currentRepositoryConfiguration.searchId, repositoryTreeContainerElement)
        } 
        else 
        {
            console.warn("filter.js is not loaded or setupRepositorySearch is undefined.")
        }

    } 
    catch (errorException) 
    {
        if (repositoryTreeContainerElement) 
        {
            repositoryTreeContainerElement.innerHTML = `<span class="has-text-danger">// Error loading manifest: ${errorException.message}</span>`
        }
    }
}

/**
 * Fetches the raw text content of a source file, determines its language class, and triggers syntax highlighting.
 */
async function fetchAndDisplaySourceFile(repositoryFileEntry, targetSourceCodeContainerElement, repositoryBaseDirectoryPath) 
{
    // Defensively check for required parameters.
    if (!repositoryFileEntry || !targetSourceCodeContainerElement) 
    {
        return
    }

    const codeViewerParentWrapperElement = targetSourceCodeContainerElement.closest(".code-viewer-wrapper")

    // Show the code viewer wrapper box.
    if (codeViewerParentWrapperElement) 
    {
        codeViewerParentWrapperElement.style.display = "block"
    }

    targetSourceCodeContainerElement.textContent = `// Loading ${repositoryFileEntry.fileName}...`

    let highlightJsLanguageClass = "language-plaintext"

    // Map file extension types to highlight.js language classes.
    if (repositoryFileEntry.fileName) 
    {
        const lowercaseSourceFileName = repositoryFileEntry.fileName.toLowerCase()

        if (lowercaseSourceFileName.endsWith(".cpp") || lowercaseSourceFileName.endsWith(".hpp") || lowercaseSourceFileName.endsWith(".cc") || lowercaseSourceFileName.endsWith(".cxx")) 
        {
            highlightJsLanguageClass = "language-cpp"
        } 
        else if (lowercaseSourceFileName.endsWith(".c") || lowercaseSourceFileName.endsWith(".h")) 
        {
            highlightJsLanguageClass = "language-c"
        } 
        else if (lowercaseSourceFileName.endsWith(".cs")) 
        {
            highlightJsLanguageClass = "language-csharp"
        } 
        else if (lowercaseSourceFileName.endsWith(".py") || lowercaseSourceFileName.endsWith(".pyx") || lowercaseSourceFileName.endsWith(".pyi") || lowercaseSourceFileName.endsWith(".pxd")) 
        {
            highlightJsLanguageClass = "language-python"
        } 
        else if (lowercaseSourceFileName.endsWith(".java")) 
        {
            highlightJsLanguageClass = "language-java"
        } 
        else if (lowercaseSourceFileName.endsWith(".go")) 
        {
            highlightJsLanguageClass = "language-go"
        }
    }

    try 
    {
        if (!repositoryFileEntry.filePath) 
        {
            throw new Error("Target file path is undefined or null.")
        }

        let resolvedRepositoryFilePath = repositoryFileEntry.filePath.trim()

        if (resolvedRepositoryFilePath.startsWith("/")) 
        {
            resolvedRepositoryFilePath = resolvedRepositoryFilePath.slice(1)
        }

        // Construct the correct fetch path accounting for repository base directories.
        if (resolvedRepositoryFilePath.startsWith("code/darwin-systems/") || resolvedRepositoryFilePath.startsWith("code/python-mac-os/")) 
        {
            // Path is already fully qualified.
        } 
        else 
        {
            const sanitizedBaseDirectoryPath = repositoryBaseDirectoryPath ? (repositoryBaseDirectoryPath.endsWith("/") ? repositoryBaseDirectoryPath.slice(0, -1) : repositoryBaseDirectoryPath) : ""

            if (sanitizedBaseDirectoryPath) 
            {
                if (!resolvedRepositoryFilePath.startsWith(sanitizedBaseDirectoryPath + "/")) 
                {
                    const extractedBaseDirectoryName = sanitizedBaseDirectoryPath.split("/").pop()

                    if (resolvedRepositoryFilePath.startsWith("code/") && !resolvedRepositoryFilePath.startsWith(`code/${extractedBaseDirectoryName}/`)) 
                    {
                        resolvedRepositoryFilePath = resolvedRepositoryFilePath.replace(/^code\//, `${sanitizedBaseDirectoryPath}/`)
                    } 
                    else 
                    {
                        resolvedRepositoryFilePath = `${sanitizedBaseDirectoryPath}/${resolvedRepositoryFilePath}`
                    }
                }
            }
        }

        const urlEncodedRepositoryFilePath = encodeURI(resolvedRepositoryFilePath)
        const sourceFileHttpResponse = await fetch(urlEncodedRepositoryFilePath)

        if (!sourceFileHttpResponse.ok) 
        {
            throw new Error(`HTTP ${sourceFileHttpResponse.status} (${sourceFileHttpResponse.statusText}) attempting to fetch: ${resolvedRepositoryFilePath}`)
        }

        const fetchedSourceFileContentText = await sourceFileHttpResponse.text()

        // Populate code block contents and apply syntax highlighting.
        targetSourceCodeContainerElement.textContent = fetchedSourceFileContentText
        targetSourceCodeContainerElement.removeAttribute("data-highlighted")
        targetSourceCodeContainerElement.className = highlightJsLanguageClass

        if (window.hljs && typeof window.hljs.highlightElement === "function") 
        {
            window.hljs.highlightElement(targetSourceCodeContainerElement)
        }
    } 
    catch (errorException) 
    {
        console.error("Error fetching file content:", errorException)
        targetSourceCodeContainerElement.textContent = `// ERROR: Could not load file from path -> ${repositoryFileEntry.filePath}\n// Details: ${errorException.message}`
    }
}

/**
 * Resets and clears the C repository viewer state.
 */
function collapseAndClear() 
{ 
    resetRepositoryViewerState("active-code-container", "active-file-label", "c-tree-file-item", "c-repo-card") 
}

/**
 * Resets and clears the Python repository viewer state.
 */
function collapseAndClearPython() 
{ 
    resetRepositoryViewerState("python-active-code-container", "python-active-file-label", "python-tree-file-item", "python-repo-card") 
}

/**
 * Resets and clears the Golang repository viewer state.
 */
function collapseAndClearGolang() 
{ 
    resetRepositoryViewerState("golang-active-code-container", "golang-active-file-label", "golang-tree-file-item", "golang-repo-card") 
}

/**
 * Resets and clears the Under the Hood ML repository viewer state.
 */
function collapseAndClearUnderTheHood() 
{ 
    resetRepositoryViewerState("under-the-hood-active-code-container", "under-the-hood-active-file-label", "under-the-hood-tree-file-item", "under-the-hood-repo-card") 
}

/**
 * Resets the specified repository viewer state, scrolls to the repo card, and clears active selection.
 */
function resetRepositoryViewerState(codeContainerIdString, fileLabelIdString, treeItemCssClassNameString, repositoryCardIdString) 
{
    if (repositoryCardIdString) 
    {
        const repositoryCardContainerElement = document.getElementById(repositoryCardIdString)

        if (repositoryCardContainerElement) 
        {
            repositoryCardContainerElement.scrollIntoView({ behavior: 'smooth', block: 'start' })
        }
    }

    const activeSourceCodeContainerElement = document.getElementById(codeContainerIdString)
    const activeFileDescriptionLabelElement = document.getElementById(fileLabelIdString)

    if (activeSourceCodeContainerElement) 
    {
        activeSourceCodeContainerElement.textContent = ""
        activeSourceCodeContainerElement.removeAttribute("data-highlighted")

        const codeViewerParentWrapperElement = activeSourceCodeContainerElement.closest(".code-viewer-wrapper")

        if (codeViewerParentWrapperElement) 
        {
            codeViewerParentWrapperElement.style.display = "none"
        }
    }

    if (activeFileDescriptionLabelElement) 
    {
        activeFileDescriptionLabelElement.innerHTML = `<i class="fas fa-file-code mr-2"></i> No file selected`
    }

    // Clear active file highlight styling across tree items.
    document.querySelectorAll(`.${treeItemCssClassNameString}`).forEach(treeItemElement => 
    {
        if (treeItemElement && treeItemElement.style) 
        {
            treeItemElement.classList.remove("is-active-file")
            treeItemElement.style.background = "transparent"
            treeItemElement.style.color = "var(--text-body)"
        }
    })
}

// Expose reset functions globally for inline HTML event handlers.
window.collapseAndClear = collapseAndClear
window.collapseAndClearPython = collapseAndClearPython
window.collapseAndClearGolang = collapseAndClearGolang
window.collapseAndClearUnderTheHood = collapseAndClearUnderTheHood