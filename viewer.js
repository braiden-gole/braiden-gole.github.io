/**
 * Class representing an interactive SVG diagram viewer with pan, zoom, touch, and adaptive scaling capabilities.
 */
class DiagramViewer {
    /**
     * Initializes a new instance of the DiagramViewer class.
     * @param {string} viewportIdentifier - The DOM element ID for the viewport container.
     * @param {string} contentIdentifier - The DOM element ID for the SVG content container.
     */
    constructor(viewportIdentifier, contentIdentifier) {
        // Retrieve and assign the primary DOM element references for viewport and content.
        this.viewportElement = document.getElementById(viewportIdentifier)
        this.contentElement = document.getElementById(contentIdentifier)
        
        // Set initial scale and pan offsets for the viewport transform matrix.
        this.currentScale = 1
        this.panX = 0
        this.panY = 0
        
        // Define scale step rates and minimum/maximum boundaries for zoom operations.
        this.scaleStep = 0.015
        this.minScale = 0.05
        this.maxScale = 20

        // Initialize state variables for tracking active drag and pinch interactions.
        this.isDragging = false
        this.isPinching = false
        this.startX = 0
        this.startY = 0
        this.initialPinchDistance = 0
        this.initialScaleBeforePinch = 1

        // Execute initial application setup and event binding.
        this.initializeViewerApplication()
    }

    /**
     * Initializes application state by parsing URL parameters and setting up event bindings.
     * @returns {void}
     */
    initializeViewerApplication() {
        // Extract document source and title parameters from the browser query string.
        const urlParameters = new URLSearchParams(window.location.search)
        const svgSourceFile = urlParameters.get('src') || urlParameters.get('file')
        const diagramTitle = urlParameters.get('title')

        // Locate the target title element within the document interface.
        const titleElement = document.getElementById('diagram-title')

        // Update interface title elements based on provided query parameter values.
        if (titleElement) {
            if (diagramTitle) {
                const decodedTitle = decodeURIComponent(diagramTitle)

                titleElement.textContent = decodedTitle

                document.title = `${decodedTitle} | SVG Viewer`
            } else {
                titleElement.textContent = svgSourceFile ? svgSourceFile.split('/').pop() : 'Untitled Diagram'
            }
        }

        // Load the specified SVG file or update status text if none provided.
        if (svgSourceFile) {
            this.loadSvgDocument(svgSourceFile)
        } else if (titleElement) {
            titleElement.textContent = 'No SVG file specified'
        }

        // Attach event listeners for user interface interactions.
        this.attachInteractionEventListeners()
    }

    /**
     * Asynchronously fetches and processes an SVG document from a given URL.
     * @param {string} url - The relative or absolute path to the SVG asset.
     * @returns {Promise<void>}
     */
    async loadSvgDocument(url) {
        // Validate the existence of the content container before attempting load operations.
        if (!this.contentElement) {
            console.error("DOM element for SVG content was not found.")

            return
        }

        try {
            // Fetch raw SVG markup from the designated network resource path.
            const response = await fetch(url)

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`)
            }
            
            // Inject raw SVG content string into the target DOM container element.
            const svgText = await response.text()

            this.contentElement.innerHTML = svgText

            // Query and configure embedded SVG element attributes and styling.
            const svgElement = this.contentElement.querySelector('svg')

            if (svgElement) {
                // Remove Mermaid's hardcoded background rectangles to reveal transparent page styling.
                const bgElements = svgElement.querySelectorAll('rect.background, rect[width="100%"][height="100%"]')

                bgElements.forEach((rectangle) => rectangle.remove())

                // Set container background transparency.
                svgElement.style.backgroundColor = 'transparent'

                // Enforce hardware-accelerated vector rendering hints for high zoom crispness.
                svgElement.style.shapeRendering = 'geometricPrecision'
                svgElement.style.textRendering = 'geometricPrecision'
                svgElement.style.imageRendering = 'optimizeQuality'

                // Clear restricting dimensions to allow free fluid scaling.
                svgElement.style.maxWidth = 'none'
                svgElement.style.maxHeight = 'none'

                svgElement.removeAttribute('width')
                svgElement.removeAttribute('height')

                // Enforce uniform aspect ratio scaling relative to view bounds.
                svgElement.setAttribute('preserveAspectRatio', 'xMidYMid meet')

                // Center node labels by adjusting SVG text node alignment styles.
                const textElements = svgElement.querySelectorAll('text')

                textElements.forEach((textNode) => {
                    textNode.style.textAnchor = 'middle'
                    textNode.style.dominantBaseline = 'middle'
                })

                // Compute exact rendered dimensions using getBBox or fall back to bounding rect calculation.
                let naturalWidth = 800
                let naturalHeight = 600

                try {
                    const bbox = svgElement.getBBox()

                    if (bbox && bbox.width > 0 && bbox.height > 0) {
                        // Minimal padding around content to eliminate empty space margins.
                        const padding = 10

                        naturalWidth = Math.ceil(bbox.width + padding * 2)
                        naturalHeight = Math.ceil(bbox.height + padding * 2)

                        // Set the viewBox explicitly to fit the exact bounding box bounds.
                        svgElement.setAttribute(
                            'viewBox',
                            `${bbox.x - padding} ${bbox.y - padding} ${naturalWidth} ${naturalHeight}`
                        )
                    }
                } catch (error) {
                    // Fallback to viewBox parsing if getBBox is unavailable.
                    const viewBox = svgElement.getAttribute('viewBox')

                    if (viewBox) {
                        const parts = viewBox.split(/[\s,]+/).map(Number)

                        if (parts.length === 4 && parts[2] > 0 && parts[3] > 0) {
                            naturalWidth = parts[2]
                            naturalHeight = parts[3]
                        }
                    }
                }

                // HD Vector Density Factor: increases intrinsic render target resolution to avoid rasterization blur when scaling up.
                const HD_DENSITY_FACTOR = 2
                const renderWidth = naturalWidth * HD_DENSITY_FACTOR
                const renderHeight = naturalHeight * HD_DENSITY_FACTOR

                // Explicitly set calculated high-density dimensions on both elements.
                svgElement.style.width = `${renderWidth}px`
                svgElement.style.height = `${renderHeight}px`

                this.contentElement.style.width = `${renderWidth}px`
                this.contentElement.style.height = `${renderHeight}px`
            }

            // Schedule initial zoom and pan reset following DOM layout stabilization.
            setTimeout(() => {
                this.resetZoomAndPan()
            }, 50)

            console.log("Successfully loaded Mermaid SVG with adaptive height container sizing and HD vector density.")
        } catch (error) {
            // Handle and display resource load failures gracefully in the viewport interface.
            if (this.contentElement) {
                this.contentElement.innerHTML = `<p class="has-text-danger p-4">Error loading SVG file: ${error.message}</p>`
            }
        }
    }

    /**
     * Calculates the Euclidean distance between two touch points for multi-touch pinch zooming.
     * @param {TouchList} touches - The active touch points list from a TouchEvent.
     * @returns {number} The distance between the first two touch points in pixels.
     */
    getPinchDistance(touches) {
        // Calculate coordinate differences along X and Y axes.
        const deltaX = touches[0].clientX - touches[1].clientX
        const deltaY = touches[0].clientY - touches[1].clientY

        // Return calculated hypotenuse distance.
        return Math.sqrt(deltaX * deltaX + deltaY * deltaY)
    }

    /**
     * Registers DOM event listeners for mouse panning, touch dragging, zooming, navigation, and toolbar controls.
     * @returns {void}
     */
    attachInteractionEventListeners() {
        // Attach back button navigation listener to return to the designated section or next target ID.
        const backButton = document.getElementById('back-btn') || document.querySelector('.back-btn')

        if (backButton) {
            backButton.addEventListener('click', (event) => {
                event.preventDefault()
                
                const urlParameters = new URLSearchParams(window.location.search)
                const targetSelection = window.location.hash 
                    ? window.location.hash.substring(1) 
                    : (urlParameters.get("from") || "networking")

                window.location.href = `index.html#${targetSelection}`
            })
        }

        // Abort canvas interaction binding if viewport container is not available.
        if (!this.viewportElement) {
            return
        }

        // Bind mouse dragging listeners for interactive panning.
        this.viewportElement.addEventListener('mousedown', (event) => {
            if (event.button !== 0) {
                return
            }

            this.isDragging = true

            this.startX = event.clientX - this.panX
            this.startY = event.clientY - this.panY

            this.viewportElement.style.cursor = 'grabbing'
        })

        // Track global mouse motion when actively dragging the canvas.
        window.addEventListener('mousemove', (event) => {
            if (!this.isDragging) {
                return
            }

            event.preventDefault()

            this.panX = event.clientX - this.startX
            this.panY = event.clientY - this.startY

            this.updateTransformMatrix()
        })

        // Release mouse drag state on button up events.
        window.addEventListener('mouseup', () => {
            if (this.isDragging) {
                this.isDragging = false

                if (this.viewportElement) {
                    this.viewportElement.style.cursor = 'grab'
                }
            }
        })

        // Handle cursor state when leaving viewport boundaries.
        this.viewportElement.addEventListener('mouseleave', () => {
            if (this.isDragging) {
                this.isDragging = false
                this.viewportElement.style.cursor = 'grab'
            }
        })

        // Bind touch start events for mobile pointer dragging and pinch zooming.
        this.viewportElement.addEventListener('touchstart', (event) => {
            if (event.touches.length === 1) {
                this.isPinching = false
                this.isDragging = true
                this.startX = event.touches[0].clientX - this.panX
                this.startY = event.touches[0].clientY - this.panY
            } else if (event.touches.length === 2) {
                this.isDragging = false
                this.isPinching = true
                this.initialPinchDistance = this.getPinchDistance(event.touches)
                this.initialScaleBeforePinch = this.currentScale
            }
        }, { passive: true })

        // Process active touch moves across mobile screens (dragging or pinch-to-zoom).
        window.addEventListener('touchmove', (event) => {
            if (event.touches.length === 2 && this.isPinching) {
                event.preventDefault()

                const currentDistance = this.getPinchDistance(event.touches)

                if (this.initialPinchDistance > 0) {
                    const scaleFactor = currentDistance / this.initialPinchDistance

                    this.currentScale = Math.max(
                        this.minScale,
                        Math.min(this.maxScale, this.initialScaleBeforePinch * scaleFactor)
                    )

                    this.updateTransformMatrix()
                }
            } else if (this.isDragging && event.touches.length === 1) {
                this.panX = event.touches[0].clientX - this.startX
                this.panY = event.touches[0].clientY - this.startY

                this.updateTransformMatrix()
            }
        }, { passive: false })

        // Terminate touch tracking state when touch points detach.
        window.addEventListener('touchend', (event) => {
            if (event.touches.length < 2) {
                this.isPinching = false
            }

            if (event.touches.length === 0) {
                this.isDragging = false
            } else if (event.touches.length === 1) {
                this.isDragging = true
                this.startX = event.touches[0].clientX - this.panX
                this.startY = event.touches[0].clientY - this.panY
            }
        })

        // Intercept mouse wheel events to perform incremental zooming centered on pointer position.
        this.viewportElement.addEventListener('wheel', (event) => {
            event.preventDefault()

            const zoomDirection = event.deltaY < 0 ? 1 : -1
            const zoomFactor = 1 + (zoomDirection * this.scaleStep)
            
            this.currentScale = Math.max(this.minScale, Math.min(this.maxScale, this.currentScale * zoomFactor))

            this.updateTransformMatrix()
        }, { passive: false })

        // Attach action handlers to external UI toolbar controls.
        const zoomInButton = document.getElementById('zoom-in-btn')
        const zoomOutButton = document.getElementById('zoom-out-btn')
        const resetZoomButton = document.getElementById('reset-zoom-btn')

        if (zoomInButton) {
            zoomInButton.addEventListener('click', () => this.zoomIn())
        }

        if (zoomOutButton) {
            zoomOutButton.addEventListener('click', () => this.zoomOut())
        }

        if (resetZoomButton) {
            resetZoomButton.addEventListener('click', () => this.resetZoomAndPan())
        }

        // Recalculate zoom and alignment automatically on window resize events.
        window.addEventListener('resize', () => {
            this.resetZoomAndPan()
        })

        // Set default interaction grab cursor.
        this.viewportElement.style.cursor = 'grab'
    }

    /**
     * Increases the current scale factor incrementally and updates the canvas matrix.
     * @returns {void}
     */
    zoomIn() {
        // Compute new scale step bounded by max limit.
        this.currentScale = Math.min(this.maxScale, this.currentScale + 0.1)

        // Apply updated transformation matrix.
        this.updateTransformMatrix()
    }

    /**
     * Decreases the current scale factor incrementally and updates the canvas matrix.
     * @returns {void}
     */
    zoomOut() {
        // Compute new scale step bounded by min limit.
        this.currentScale = Math.max(this.minScale, this.currentScale - 0.1)

        // Apply updated transformation matrix.
        this.updateTransformMatrix()
    }

    /**
     * Resets pan offsets to center and computes optimal fitting scale for the SVG content.
     * @returns {void}
     */
    resetZoomAndPan() {
        // Reset translation coordinates to origin.
        this.panX = 0
        this.panY = 0

        if (!this.contentElement) {
            return
        }

        const svgElement = this.contentElement.querySelector('svg')

        if (svgElement && this.viewportElement) {
            // Retrieve container bounding dimensions.
            const viewportWidth = this.viewportElement.clientWidth || window.innerWidth
            const viewportHeight = this.viewportElement.clientHeight || window.innerHeight

            let svgWidth = parseFloat(svgElement.style.width) || 800
            let svgHeight = parseFloat(svgElement.style.height) || 600

            // Apply 0.85 padding ratio for optimal zoomed-out breathing room by default.
            const paddingRatio = 0.85
            const scaleX = (viewportWidth * paddingRatio) / svgWidth
            const scaleY = (viewportHeight * paddingRatio) / svgHeight

            // Calculate bounding scale factor.
            this.currentScale = Math.min(scaleX, scaleY)
            this.currentScale = Math.max(this.minScale, Math.min(this.maxScale, this.currentScale))
        } else {
            this.currentScale = 1
        }

        // Apply reset transformation matrix.
        this.updateTransformMatrix()
    }

    /**
     * Applies the computed pan offsets and scale factor to the CSS transform matrix of the container.
     * @returns {void}
     */
    updateTransformMatrix() {
        if (this.contentElement) {
            // Transform element matrix using standard centered translation and scale values.
            this.contentElement.style.transform = `translate(calc(-50% + ${this.panX}px), calc(-50% + ${this.panY}px)) scale(${this.currentScale})`
        }
    }
}

document.addEventListener("DOMContentLoaded", () => {
    window.diagramViewerInstance = new DiagramViewer("svg-viewport", "svg-content")
})